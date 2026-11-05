package service

import (
	"context"
	"fmt"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var (
	ErrAPIKeyGroupsEmpty      = infraerrors.BadRequest("API_KEY_GROUPS_EMPTY", "group_ids must contain at least one group")
	ErrAPIKeyGroupsTooMany    = infraerrors.BadRequest("API_KEY_GROUPS_TOO_MANY", fmt.Sprintf("an api key can bind at most %d groups", MaxAPIKeyGroups))
	ErrAPIKeyGroupsDuplicate  = infraerrors.BadRequest("API_KEY_GROUPS_DUPLICATE", "group_ids must not contain duplicates")
	ErrAPIKeyGroupModeInvalid = infraerrors.BadRequest("API_KEY_GROUP_MODE_INVALID", "group_mode must be manual or auto")
	ErrAPIKeyAutoModeGroups   = infraerrors.BadRequest("API_KEY_AUTO_MODE_GROUPS", "group_ids is not accepted in auto group mode")
)

// normalizeGroupMode 校验并规范化分组模式；空字符串视为 manual。
func normalizeGroupMode(mode string) (string, error) {
	switch mode {
	case "", APIKeyGroupModeManual:
		return APIKeyGroupModeManual, nil
	case APIKeyGroupModeAuto:
		return APIKeyGroupModeAuto, nil
	default:
		return "", ErrAPIKeyGroupModeInvalid
	}
}

// requestedGroupIDs 把请求里的 group_ids / group_id 合成有序列表。
// group_ids 非 nil 时以它为准；否则只传 group_id 等价于 [group_id]；都没传返回 nil。
func requestedGroupIDs(groupID *int64, groupIDs []int64) ([]int64, error) {
	if groupIDs == nil {
		if groupID == nil {
			return nil, nil
		}
		return []int64{*groupID}, nil
	}
	if len(groupIDs) == 0 {
		return nil, ErrAPIKeyGroupsEmpty
	}
	if len(groupIDs) > MaxAPIKeyGroups {
		return nil, ErrAPIKeyGroupsTooMany
	}
	seen := make(map[int64]struct{}, len(groupIDs))
	out := make([]int64, 0, len(groupIDs))
	for _, id := range groupIDs {
		if _, dup := seen[id]; dup {
			return nil, ErrAPIKeyGroupsDuplicate
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out, nil
}

// loadBindableGroups 按顺序加载分组并逐个校验用户是否可以绑定。
func (s *APIKeyService) loadBindableGroups(ctx context.Context, user *User, groupIDs []int64) ([]*Group, error) {
	groups := make([]*Group, 0, len(groupIDs))
	for _, id := range groupIDs {
		group, err := s.groupRepo.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("get group: %w", err)
		}
		if !s.canUserBindGroup(ctx, user, group) {
			return nil, ErrGroupNotAllowed
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// applyGroups 把校验过的分组写到密钥上（GroupID = 首位）。
func applyGroups(apiKey *APIKey, groupIDs []int64, groups []*Group) {
	if len(groupIDs) == 0 {
		apiKey.GroupID = nil
		apiKey.GroupIDs = nil
		apiKey.Groups = nil
		apiKey.Group = nil
		return
	}
	primary := groupIDs[0]
	apiKey.GroupID = &primary
	apiKey.GroupIDs = groupIDs
	apiKey.Group = groups[0]
	if len(groups) > 1 {
		apiKey.Groups = groups
	} else {
		apiKey.Groups = nil
	}
}

func sameGroupIDs(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// applyGroupUpdate 处理 Update 请求里的 group_id / group_ids / group_mode。
//
// 只传 group_id 且等于当前主路由时保留现有的多分组列表：旧版前端编辑名称等字段时
// 会原样回传 group_id，不能因此把后面的分组清掉。
func (s *APIKeyService) applyGroupUpdate(ctx context.Context, userID int64, apiKey *APIKey, req UpdateAPIKeyRequest, fields *APIKeyUpdateFields) error {
	currentMode, _ := normalizeGroupMode(apiKey.GroupMode)
	mode := currentMode
	if req.GroupMode != nil {
		m, err := normalizeGroupMode(*req.GroupMode)
		if err != nil {
			return err
		}
		mode = m
	}
	if mode == APIKeyGroupModeAuto && req.GroupIDs != nil {
		return ErrAPIKeyAutoModeGroups
	}
	// 自动模式的分组由工作台按用户权限决定，用户端不能改（含批量编辑只传 group_id 的情况）。
	if mode == APIKeyGroupModeAuto && currentMode == APIKeyGroupModeAuto && req.GroupID != nil {
		return ErrAPIKeyAutoModeGroups
	}
	if req.GroupID == nil && req.GroupIDs == nil {
		if mode != currentMode {
			apiKey.GroupMode = mode
			fields.GroupID = true
		}
		return nil
	}

	current := apiKey.OrderedGroupIDs()
	target, err := requestedGroupIDs(req.GroupID, req.GroupIDs)
	if err != nil {
		return err
	}
	keepCurrent := req.GroupIDs == nil && len(current) > 0 && target[0] == current[0]
	// 原样回传当前列表（编辑名称等其他字段时）不重新校验权限：否则某个备用分组的权限被收回后，
	// 用户连改名都保存不了。收回的分组在请求时由选组服务跳过（D3 第 1 条）。
	if !keepCurrent && req.GroupIDs != nil && sameGroupIDs(target, current) && mode == currentMode {
		return nil
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	groups, err := s.loadBindableGroups(ctx, user, target)
	if err != nil {
		return err
	}

	apiKey.GroupMode = mode
	if keepCurrent {
		if mode != currentMode {
			fields.GroupID = true
		}
		return nil
	}
	applyGroups(apiKey, target, groups)
	if mode != currentMode || !sameGroupIDs(target, current) {
		fields.GroupID = true
	}
	return nil
}

// autoModePrimaryGroup 返回自动模式密钥建立时的主路由：按 D3 排序后的第一个可用分组。
// 用户当前没有任何可用分组时返回空列表（密钥先不绑分组，权限到位后请求时再选）。
func (s *APIKeyService) autoModePrimaryGroup(ctx context.Context, userID int64) ([]int64, []*Group, error) {
	available, err := s.GetAvailableGroups(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	var rates GroupRouterRateSource
	if s.groupRouter != nil {
		rates = s.groupRouter.deps.Rates
	}
	ordered := orderAutoGroupCandidates(ctx, userID, available, rates)
	if len(ordered) == 0 {
		return nil, nil, nil
	}
	return []int64{ordered[0].ID}, ordered[:1], nil
}

// replacePrimaryGroupID 把有序分组的首位换成 primary，其余保持顺序并去掉与 primary 重复的项。
func replacePrimaryGroupID(ids []int64, primary int64) []int64 {
	out := []int64{primary}
	for i, id := range ids {
		if i == 0 || id == primary {
			continue
		}
		out = append(out, id)
	}
	return out
}
