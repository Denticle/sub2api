//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type multiGroupRepoStub struct {
	groupRepoStub
	groups map[int64]*Group
}

func (s *multiGroupRepoStub) GetByID(_ context.Context, id int64) (*Group, error) {
	if g, ok := s.groups[id]; ok {
		return g, nil
	}
	return nil, ErrGroupNotFound
}

func newMultiGroupService(key *APIKey, user *User) (*APIKeyService, *updateFieldsAPIKeyRepoStub) {
	repo := &updateFieldsAPIKeyRepoStub{key: key}
	groups := &multiGroupRepoStub{groups: map[int64]*Group{
		1: {ID: 1, Name: "claude"},
		2: {ID: 2, Name: "openai"},
		3: {ID: 3, Name: "image"},
		9: {ID: 9, Name: "vip", IsExclusive: true},
	}}
	return &APIKeyService{
		apiKeyRepo: repo,
		userRepo:   &userRepoStub{user: user},
		groupRepo:  groups,
	}, repo
}

func TestRequestedGroupIDs(t *testing.T) {
	ids, err := requestedGroupIDs(nil, nil)
	require.NoError(t, err)
	require.Nil(t, ids)

	ids, err = requestedGroupIDs(int64Ptr(5), nil)
	require.NoError(t, err)
	require.Equal(t, []int64{5}, ids)

	// group_ids 优先于 group_id
	ids, err = requestedGroupIDs(int64Ptr(5), []int64{2, 1})
	require.NoError(t, err)
	require.Equal(t, []int64{2, 1}, ids)

	_, err = requestedGroupIDs(nil, []int64{})
	require.ErrorIs(t, err, ErrAPIKeyGroupsEmpty)

	_, err = requestedGroupIDs(nil, []int64{1, 2, 1})
	require.ErrorIs(t, err, ErrAPIKeyGroupsDuplicate)

	tooMany := make([]int64, MaxAPIKeyGroups+1)
	for i := range tooMany {
		tooMany[i] = int64(i + 1)
	}
	_, err = requestedGroupIDs(nil, tooMany)
	require.ErrorIs(t, err, ErrAPIKeyGroupsTooMany)
}

func TestAPIKeyUpdate_GroupIDsReplaceListAndPrimary(t *testing.T) {
	svc, repo := newMultiGroupService(
		&APIKey{ID: 1, UserID: 7, Key: "sk-test", Status: StatusActive, GroupID: int64Ptr(1)},
		&User{ID: 7},
	)

	got, err := svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{GroupIDs: []int64{2, 1, 3}})
	require.NoError(t, err)
	require.Equal(t, []APIKeyUpdateFields{{GroupID: true}}, repo.updateFields)
	require.Equal(t, int64(2), *got.GroupID)
	require.Equal(t, []int64{2, 1, 3}, got.OrderedGroupIDs())
	require.Equal(t, int64(2), got.Group.ID)
	require.Len(t, got.Groups, 3)
}

// 旧版前端编辑名称时会原样回传 group_id，不能因此清掉后面的分组。
func TestAPIKeyUpdate_SamePrimaryGroupIDKeepsList(t *testing.T) {
	svc, repo := newMultiGroupService(
		&APIKey{ID: 1, UserID: 7, Key: "sk-test", Status: StatusActive, GroupID: int64Ptr(1), GroupIDs: []int64{1, 2}},
		&User{ID: 7},
	)
	name := "renamed"

	got, err := svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{Name: &name, GroupID: int64Ptr(1)})
	require.NoError(t, err)
	require.Equal(t, []APIKeyUpdateFields{{Name: true}}, repo.updateFields)
	require.Equal(t, []int64{1, 2}, got.OrderedGroupIDs())
}

// 只传一个不同的 group_id 等价于 group_ids=[group_id]。
func TestAPIKeyUpdate_DifferentGroupIDCollapsesToSingle(t *testing.T) {
	svc, repo := newMultiGroupService(
		&APIKey{ID: 1, UserID: 7, Key: "sk-test", Status: StatusActive, GroupID: int64Ptr(1), GroupIDs: []int64{1, 2}},
		&User{ID: 7},
	)

	got, err := svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{GroupID: int64Ptr(3)})
	require.NoError(t, err)
	require.Equal(t, []APIKeyUpdateFields{{GroupID: true}}, repo.updateFields)
	require.Equal(t, []int64{3}, got.OrderedGroupIDs())
	require.Nil(t, got.Groups)
}

// 列表里任何一个分组无权绑定都要拒绝，而不只是首位。
func TestAPIKeyUpdate_RejectsUnbindableGroupAnywhereInList(t *testing.T) {
	svc, repo := newMultiGroupService(
		&APIKey{ID: 1, UserID: 7, Key: "sk-test", Status: StatusActive, GroupID: int64Ptr(1)},
		&User{ID: 7},
	)

	_, err := svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{GroupIDs: []int64{1, 9}})
	require.ErrorIs(t, err, ErrGroupNotAllowed)
	require.Empty(t, repo.updateFields)
}

func TestAPIKeyUpdate_AutoModeRejectsGroupIDs(t *testing.T) {
	svc, _ := newMultiGroupService(
		&APIKey{ID: 1, UserID: 7, Key: "sk-test", Status: StatusActive, GroupID: int64Ptr(1)},
		&User{ID: 7},
	)
	auto := APIKeyGroupModeAuto

	_, err := svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{GroupMode: &auto, GroupIDs: []int64{1}})
	require.ErrorIs(t, err, ErrAPIKeyAutoModeGroups)

	bad := "random"
	_, err = svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{GroupMode: &bad})
	require.ErrorIs(t, err, ErrAPIKeyGroupModeInvalid)
}

func TestAPIKeyCreate_GroupIDsSetsOrderedGroups(t *testing.T) {
	repo, cache := newCreateLimitStubs()
	svc := newCreateLimitService(repo, cache, 0, 0)
	svc.userRepo = &userRepoStub{user: &User{ID: 7, AllowedGroups: []int64{9}}}
	svc.groupRepo = &multiGroupRepoStub{groups: map[int64]*Group{
		1: {ID: 1}, 9: {ID: 9, IsExclusive: true},
	}}

	key, err := svc.Create(context.Background(), 7, CreateAPIKeyRequest{Name: "multi", GroupIDs: []int64{9, 1}})
	require.NoError(t, err)
	require.Len(t, repo.created, 1)
	require.Equal(t, int64(9), *key.GroupID)
	require.Equal(t, []int64{9, 1}, key.OrderedGroupIDs())
	require.Equal(t, APIKeyGroupModeManual, key.GroupMode)

	_, err = svc.Create(context.Background(), 7, CreateAPIKeyRequest{Name: "dup", GroupIDs: []int64{1, 1}})
	require.ErrorIs(t, err, ErrAPIKeyGroupsDuplicate)
	require.Len(t, repo.created, 1)
}

func TestOrderedGroupIDs_StaleListFallsBackToGroupID(t *testing.T) {
	// 只改了 GroupID 的旧调用方（如管理员换分组）不能留下过期列表。
	k := &APIKey{GroupID: int64Ptr(3), GroupIDs: []int64{1, 2}}
	require.Equal(t, []int64{3}, k.OrderedGroupIDs())

	k = &APIKey{GroupIDs: []int64{1, 2}}
	require.Nil(t, k.OrderedGroupIDs())
}

func TestAuthSnapshot_RoundTripsMultipleGroups(t *testing.T) {
	svc := &APIKeyService{}
	primary := &Group{ID: 1, Name: "claude", Platform: PlatformAnthropic, RateMultiplier: 1, Hydrated: true}
	in := &APIKey{
		ID: 5, UserID: 7, Key: "sk-test", Status: StatusActive,
		GroupID:   int64Ptr(1),
		GroupMode: APIKeyGroupModeManual,
		GroupIDs:  []int64{1, 2},
		Group:     primary,
		Groups:    []*Group{primary, {ID: 2, Name: "openai", Platform: PlatformOpenAI, RateMultiplier: 0.5}},
		User:      &User{ID: 7, Status: StatusActive},
	}

	snap := svc.snapshotFromAPIKey(context.Background(), in)
	require.Equal(t, apiKeyAuthSnapshotVersion, snap.Version)
	require.Equal(t, []int64{1, 2}, snap.GroupIDs)
	require.Len(t, snap.Groups, 2)

	out := svc.snapshotToAPIKey("sk-test", snap)
	require.Equal(t, APIKeyGroupModeManual, out.GroupMode)
	require.Equal(t, []int64{1, 2}, out.OrderedGroupIDs())
	require.Len(t, out.Groups, 2)
	require.Same(t, out.Group, out.Groups[0])
	require.Equal(t, PlatformOpenAI, out.Groups[1].Platform)
	require.Equal(t, 0.5, out.Groups[1].RateMultiplier)
	require.True(t, out.Groups[1].Hydrated)
	require.True(t, out.HasMultipleGroups())
}

func TestAuthSnapshot_SingleGroupHasNoGroupsList(t *testing.T) {
	svc := &APIKeyService{}
	in := &APIKey{
		ID: 5, UserID: 7, Key: "sk-test", Status: StatusActive,
		GroupID: int64Ptr(1), GroupIDs: []int64{1},
		Group: &Group{ID: 1}, User: &User{ID: 7},
	}
	snap := svc.snapshotFromAPIKey(context.Background(), in)
	require.Nil(t, snap.Groups)
	out := svc.snapshotToAPIKey("sk-test", snap)
	require.Nil(t, out.Groups)
	require.False(t, out.HasMultipleGroups())
}

// 自动模式（工作台）密钥：group_id 写当前排序的第一组，不写有序列表（tasks 7.1）。
func TestAPIKeyCreate_AutoModeUsesFirstOrderedGroup(t *testing.T) {
	now := time.Now()
	repo, cache := newCreateLimitStubs()
	svc := newCreateLimitService(repo, cache, 0, 0)
	svc.userRepo = &visibilityUserRepo{user: &User{ID: 7}}
	svc.groupRepo = &visibilityGroupRepo{groups: []Group{
		{ID: 1, Status: StatusActive, RateMultiplier: 1},
		{ID: 2, Status: StatusActive, RateMultiplier: 0.5},
		{ID: 3, Status: StatusActive, RateMultiplier: 2, SubscriptionType: SubscriptionTypeSubscription},
	}}
	svc.userSubRepo = &visibilitySubRepo{subscriptions: []UserSubscription{
		{UserID: 7, GroupID: 3, Status: SubscriptionStatusActive, ExpiresAt: now.Add(time.Hour)},
	}}

	key, err := svc.Create(context.Background(), 7, CreateAPIKeyRequest{Name: "studio", GroupMode: APIKeyGroupModeAuto})
	require.NoError(t, err)
	require.Equal(t, APIKeyGroupModeAuto, key.GroupMode)
	require.NotNil(t, key.GroupID)
	require.Equal(t, int64(3), *key.GroupID, "subscription group first, matching the request-time router order")
	require.Equal(t, []int64{3}, key.OrderedGroupIDs(), "auto keys never store a manual group list")

	// 没有订阅时取倍率最低的
	svc.userSubRepo = &visibilitySubRepo{}
	key, err = svc.Create(context.Background(), 7, CreateAPIKeyRequest{Name: "studio-2", GroupMode: APIKeyGroupModeAuto})
	require.NoError(t, err)
	require.Equal(t, int64(2), *key.GroupID)
}

// 用户当前没有任何可用分组时仍能建自动模式密钥，权限到位后请求时再选组。
func TestAPIKeyCreate_AutoModeWithoutGroups(t *testing.T) {
	repo, cache := newCreateLimitStubs()
	svc := newCreateLimitService(repo, cache, 0, 0)
	svc.userRepo = &visibilityUserRepo{user: &User{ID: 7}}
	svc.groupRepo = &visibilityGroupRepo{}
	svc.userSubRepo = &visibilitySubRepo{}

	key, err := svc.Create(context.Background(), 7, CreateAPIKeyRequest{Name: "studio", GroupMode: APIKeyGroupModeAuto})
	require.NoError(t, err)
	require.Nil(t, key.GroupID)
	require.Equal(t, APIKeyGroupModeAuto, key.GroupMode)
}

// 用户端（含批量编辑）不能改自动模式密钥的分组，只传 group_id 也不行。
func TestAPIKeyUpdate_AutoModeKeyRejectsGroupIDChange(t *testing.T) {
	svc, repo := newMultiGroupService(
		&APIKey{ID: 1, UserID: 7, Key: "sk-studio", Status: StatusActive, GroupID: int64Ptr(1), GroupMode: APIKeyGroupModeAuto},
		&User{ID: 7},
	)
	_, err := svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{GroupID: int64Ptr(2)})
	require.ErrorIs(t, err, ErrAPIKeyAutoModeGroups)
	require.Empty(t, repo.updateFields)

	name := "renamed"
	_, err = svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{Name: &name})
	require.NoError(t, err, "other fields of an auto key are still editable")
}

// 原样回传当前列表时不重新校验权限：某个备用分组权限被收回后，用户仍能改名（审查发现的问题）。
func TestAPIKeyUpdate_UnchangedListSkipsPermissionRecheck(t *testing.T) {
	svc, repo := newMultiGroupService(
		&APIKey{ID: 1, UserID: 7, Key: "sk-test", Status: StatusActive, GroupID: int64Ptr(1), GroupIDs: []int64{1, 9}},
		&User{ID: 7}, // 用户已没有专属分组 9 的权限
	)
	name := "renamed"
	_, err := svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{Name: &name, GroupIDs: []int64{1, 9}})
	require.NoError(t, err)
	require.Equal(t, []APIKeyUpdateFields{{Name: true}}, repo.updateFields, "groups are not rewritten")

	// 真正改了列表时仍要校验每个分组
	_, err = svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{GroupIDs: []int64{9, 1}})
	require.ErrorIs(t, err, ErrGroupNotAllowed)
}

func TestReplacePrimaryGroupID(t *testing.T) {
	require.Equal(t, []int64{4, 2, 3}, replacePrimaryGroupID([]int64{1, 2, 3}, 4))
	require.Equal(t, []int64{3, 2}, replacePrimaryGroupID([]int64{1, 2, 3}, 3))
	require.Equal(t, []int64{5}, replacePrimaryGroupID(nil, 5))
}
