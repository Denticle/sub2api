package oauthprovider

import (
	"context"
	"sort"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

// StudioModel 是工作台需要的最小模型描述：调用名、所属平台、能力类别。
// 价格不下发——工作台积分已关闭，费用由 forapi 按用户自己的密钥结算。
type StudioModel struct {
	ID         string `json:"id"`
	Platform   string `json:"platform"`
	Capability string `json:"capability"` // text / image / video
}

// groupLister 只取模型广场的分组口径：分组关联了哪些渠道、渠道里有哪些模型。
type groupLister interface {
	ListGroups(ctx context.Context) ([]service.PlazaGroup, error)
}

// StudioModels 返回工作台的模型目录：全部 active 分组模型的并集。
//
// 目录是全站一份（工作台定时同步，不区分用户），专用密钥又是自动模式、按用户权限
// 选组，所以这里列出任何用户可能用到的模型。用户选了自己没有权限的模型时，网关
// 按选组规则返回 404 model_not_found，由工作台展示。
// 与 /oauth/studio-key 使用同一组客户端凭据。
func (s *Service) StudioModels(ctx context.Context, clientID, clientSecret string) ([]StudioModel, error) {
	if _, err := s.authenticateClient(ctx, clientID, clientSecret); err != nil {
		return nil, err
	}
	if s.groups == nil {
		return []StudioModel{}, nil
	}
	groups, err := s.groups.ListGroups(ctx)
	if err != nil {
		return nil, err
	}
	out := []StudioModel{}
	seen := map[string]bool{}
	for _, group := range groups {
		for _, m := range group.Models {
			id := strings.TrimSpace(m.Name)
			if id == "" || seen[id] {
				continue
			}
			seen[id] = true
			out = append(out, StudioModel{ID: id, Platform: m.Platform, Capability: modelCapability(m)})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// modelCapability 优先用渠道上配置的计费方式判断类别；没有定价时按模型名兜底。
func modelCapability(m service.PlazaModel) string {
	if m.Pricing != nil {
		switch m.Pricing.BillingMode {
		case service.BillingModeImage:
			return "image"
		case service.BillingModeVideo:
			return "video"
		}
	}
	name := strings.ToLower(m.Name)
	for _, marker := range []string{"video", "veo", "seedance", "sora", "kling", "hailuo"} {
		if strings.Contains(name, marker) {
			return "video"
		}
	}
	for _, marker := range []string{"image", "dall-e", "imagen", "seedream", "flux"} {
		if strings.Contains(name, marker) {
			return "image"
		}
	}
	return "text"
}
