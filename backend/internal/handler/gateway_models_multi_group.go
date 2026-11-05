package handler

import (
	"context"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// 多分组密钥的模型列表（add-api-key-multi-group-routing design D4）：
// 合并全部可用候选分组各自「可见」的模型（已应用各组白名单），按候选顺序去重。
// /models/:model 在合并后的列表里查找，因此任一组可见即返回。

// writeMultiGroupModels 处理多分组密钥的 /v1/models；返回 false 表示不是多分组请求，
// 调用方继续走单分组逻辑。
func (h *GatewayHandler) writeMultiGroupModels(c *gin.Context) bool {
	if _, forced := middleware2.GetForcePlatformFromContext(c); forced {
		return false
	}
	candidates := middleware2.APIKeyGroupCandidates(c)
	if len(candidates) == 0 {
		return false
	}

	ctx := c.Request.Context()
	var merged []string
	shape := ""
	for i, group := range candidates {
		ids, platform := h.listedModelIDsForGroup(ctx, group)
		merged = mergeModelIDs(merged, ids)
		if i == 0 {
			shape = platform
		} else if shape != platform {
			// 跨平台合并：沿用 composite 分组的列表格式（Anthropic 形状，含 id）。
			shape = service.PlatformComposite
		}
	}
	writeModelsList(c, shape, merged)
	return true
}

// listedModelIDsForGroup 返回单个分组在 /v1/models 里会列出的模型 ID 与列表格式平台，
// 与 Models 的单分组分支保持同一套规则（固定账号 manifest 分支除外：它直接透传上游
// 响应体，无法与其他分组合并，这里退回按账号映射列出）。
func (h *GatewayHandler) listedModelIDsForGroup(ctx context.Context, group *service.Group) ([]string, string) {
	if group == nil {
		return nil, ""
	}
	groupID := &group.ID
	platform := group.Platform

	if platform == service.PlatformComposite {
		models := h.compositeAvailableModels(ctx, groupID, "", true)
		if len(models) == 0 {
			models = defaultModelIDsForPlatform(service.PlatformComposite)
		}
		if group.ModelAllowlistEnabled() {
			models = group.ModelAllowlist.FilterForListing(models)
		}
		return models, platform
	}

	available := h.gatewayService.GetAvailableModels(ctx, groupID, platform)
	if group.ModelAllowlistEnabled() {
		source := modelListingSource(platform, available, defaultModelIDsForPlatform(platform))
		return group.ModelAllowlist.FilterForListing(source), platform
	}
	if len(available) > 0 || (platform == service.PlatformOpenAI && available != nil) {
		return available, platform
	}
	return defaultModelIDsForPlatform(platform), platform
}
