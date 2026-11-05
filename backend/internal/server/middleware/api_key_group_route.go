package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/httputil"
	"github.com/Wei-Shaw/sub2api/internal/pkg/requestmodel"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// 多分组密钥在鉴权阶段选组（add-api-key-multi-group-routing design D2）。
//
// 选中的分组被装进 APIKey 副本（GroupID/Group 指向它），后续的分组校验、订阅、
// 额度检查、setGroupContext 与下游 handler 都按它工作，不需要感知多分组。
//
// 单分组密钥直接返回原对象：不读请求体、不查任何依赖。

// apiKeyGroupRouteErrorWriter 写出"没有分组支持该模型"的错误。
type apiKeyGroupRouteErrorWriter func(c *gin.Context, status int, message string)

// routeAPIKeyGroup 为多分组密钥选组。返回 false 表示已写出错误响应并 Abort。
func routeAPIKeyGroup(c *gin.Context, apiKeyService *service.APIKeyService, apiKey *service.APIKey, writeErr apiKeyGroupRouteErrorWriter) (*service.APIKey, bool) {
	if apiKey == nil || !apiKey.HasMultipleGroups() {
		return apiKey, true
	}
	router := apiKeyService.GroupRouter()
	if router == nil {
		return apiKey, true
	}

	if requestID, ok := videoLookupRequestID(c); ok {
		return applyAPIKeyGroupRoute(c, apiKey, router.SelectForVideoLookup(c.Request.Context(), apiKey, requestID)), true
	}

	var models []string
	var body []byte
	if apiKeyGroupRouteNeedsModel(c) {
		var done bool
		models, body, done = apiKeyGroupRouteModels(c)
		if !done {
			return nil, false
		}
	}

	endpoint := service.CompositeRouteEndpointForPath(c.Request.URL.Path)
	ctx := c.Request.Context()
	if apiKeyGroupRouteIsImageGeneration(c, models, body) {
		ctx = service.WithImageGenerationRoute(ctx)
	}
	result, err := router.Select(ctx, apiKey, models, endpoint)
	if err != nil {
		if errors.Is(err, service.ErrNoGroupSupportsModel) {
			service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalModelConfiguration)
			MarkIngressRejected(c, IngressRejectModelNotAllowed)
			writeErr(c, http.StatusNotFound, fmt.Sprintf("Model %q is not available for any group bound to this API key", firstNonEmpty(models)))
			c.Abort()
			return nil, false
		}
		// Select 只会返回上面这一种错误；其余情况保守地退回主路由。
		return apiKey, true
	}

	return applyAPIKeyGroupRoute(c, apiKey, result), true
}

func applyAPIKeyGroupRoute(c *gin.Context, apiKey *service.APIKey, result service.GroupRouteResult) *service.APIKey {
	selected := service.CloneAPIKeyForGroup(apiKey, result.Group)
	c.Request = c.Request.WithContext(service.WithAPIKeyGroupRoute(c.Request.Context(), service.APIKeyGroupRoute{
		All:      result.Candidates,
		Selected: result.Group,
		Outcome:  result.Outcome,
	}))
	SetOpsFallbackAPIKey(c, selected)
	return selected
}

// videoLookupRequestID 识别视频任务的状态/内容查询（Grok /videos/... 与 Seedance tasks），
// 返回路径里的任务 ID。这类请求不带模型，按创建任务时的分组处理。
func videoLookupRequestID(c *gin.Context) (string, bool) {
	if c.Request == nil {
		return "", false
	}
	switch c.Request.Method {
	case http.MethodGet, http.MethodDelete:
	default:
		return "", false
	}
	route := c.FullPath()
	if strings.Contains(route, "/videos/") {
		if id := strings.TrimSpace(c.Param("request_id")); id != "" {
			return id, true
		}
	}
	if strings.Contains(route, "/contents/generations/tasks/") {
		if id := strings.TrimSpace(c.Param("task_id")); id != "" {
			// 创建时按 result.ResponseID 记录分组，Seedance 的 ResponseID 带 "seedance:" 前缀。
			return service.SeedanceTaskKey(id), true
		}
	}
	return "", false
}

// apiKeyGroupRouteNeedsModel 判断本次请求是否按模型选组；其余走主路由（design D4）。
func apiKeyGroupRouteNeedsModel(c *gin.Context) bool {
	if c.Request == nil || c.Request.URL == nil {
		return false
	}
	path := c.Request.URL.Path
	if path == "/v1/usage" || path == "/v1/sub2api/billing" || isAsyncImageTaskRead(c.Request.Method, path) {
		return false
	}
	// Responses WebSocket 的模型在升级后的帧里，鉴权阶段拿不到。
	if isResponsesWebSocketRoute(c) {
		return false
	}
	return true
}

// apiKeyGroupRouteModels 按白名单中间件同一套规则取模型：路由参数 → 请求体 → 查询参数。
// 返回 false 表示读请求体失败，已写出响应。
func apiKeyGroupRouteModels(c *gin.Context) ([]string, []byte, bool) {
	if model := groupModelAllowlistModelFromParams(c); model != "" {
		return []string{model}, nil, true
	}
	var models []string
	var body []byte
	switch c.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		var err error
		body, err = httputil.ReadRequestBodyWithPrealloc(c.Request)
		if err != nil {
			status := http.StatusBadRequest
			message := "Failed to read request body"
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				status = http.StatusRequestEntityTooLarge
				message = "Request body is too large"
			}
			c.JSON(status, gin.H{"error": gin.H{"type": "invalid_request_error", "message": message}})
			c.Abort()
			return nil, nil, false
		}
		requestmodel.ResetRequestBody(c.Request, body)
		models = requestmodel.FromBodyCandidates(c.FullPath(), c.GetHeader("Content-Type"), body)
	}
	if len(models) == 0 {
		if model := strings.TrimSpace(c.Query("model")); model != "" {
			models = []string{model}
		}
	}
	return models, body, true
}

// apiKeyGroupRouteIsImageGeneration 判断请求会生成图片，规则与各 handler 的生图权限检查一致：
// 生图端点、图片模型，或请求体里显式的 image_generation 工具（不含 Codex 的被动 namespace 声明）。
func apiKeyGroupRouteIsImageGeneration(c *gin.Context, models []string, body []byte) bool {
	path := c.Request.URL.Path
	if strings.Contains(path, "/images/") {
		return true
	}
	model := ""
	if len(models) > 0 {
		model = models[0]
	}
	return service.IsExplicitImageGenerationIntent(path, model, body)
}

func firstNonEmpty(values []string) string {
	for _, v := range values {
		if v = strings.TrimSpace(v); v != "" {
			return v
		}
	}
	return ""
}

// APIKeyGroupCandidates 返回多分组密钥本次请求的全部可用候选分组；单分组密钥返回 nil。
func APIKeyGroupCandidates(c *gin.Context) []*service.Group {
	if c == nil || c.Request == nil {
		return nil
	}
	route, ok := service.APIKeyGroupRouteFromContext(c.Request.Context())
	if !ok || len(route.All) < 2 {
		return nil
	}
	return route.All
}

// FirstAPIKeyGroupCandidate 返回第一个满足 accept 的候选分组；单分组密钥或没有命中时返回 nil。
func FirstAPIKeyGroupCandidate(c *gin.Context, accept func(*service.Group) bool) *service.Group {
	for _, g := range APIKeyGroupCandidates(c) {
		if accept(g) {
			return g
		}
	}
	return nil
}

// UseAPIKeyGroup 把上下文里的密钥换成以 group 为主路由的副本，并同步分组上下文。
//
// 只给不计费的只读请求（模型列表、Codex manifest）按平台挑组用：鉴权阶段的订阅与
// 额度检查是按原分组做的，计费请求必须在鉴权里选组。
func UseAPIKeyGroup(c *gin.Context, group *service.Group) *service.APIKey {
	apiKey, ok := GetAPIKeyFromContext(c)
	if !ok || apiKey == nil || group == nil {
		return apiKey
	}
	cloned := service.CloneAPIKeyForGroup(apiKey, group)
	c.Set(string(ContextKeyAPIKey), cloned)
	setGroupContext(c, group)
	return cloned
}
