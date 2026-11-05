package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler"
	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// 多分组密钥跨平台请求落到正确 handler（tasks 4.6）：一把密钥绑定
// Anthropic（claude-*）、OpenAI（gpt-*、gpt-image-*）、Grok（grok-*）三个分组。

type multiGroupRouteDiagnoser map[int64]map[string]bool

func (d multiGroupRouteDiagnoser) DiagnoseModelAvailabilityForPlatform(_ context.Context, groupID *int64, model, _ string) service.ModelAvailabilityDiagnosis {
	return service.ModelAvailabilityDiagnosis{HasAccountsInPool: true, HasModelSupport: d[*groupID][model]}
}

// multiGroupRouteAPIKeyRepo 在 keyBillingRouteAPIKeyRepo 基础上实现标准模式鉴权会调用的 TouchLastUsed。
type multiGroupRouteAPIKeyRepo struct {
	keyBillingRouteAPIKeyRepo
}

func (r *multiGroupRouteAPIKeyRepo) UpdateLastUsed(context.Context, int64, time.Time) error {
	return nil
}

type multiGroupRouteVideoGroups map[string]int64

func (v multiGroupRouteVideoGroups) ResolveVideoRequestGroup(_ context.Context, requestID string, _, _ int64) (int64, error) {
	return v[requestID], nil
}

// newMultiGroupRouteTestRouter 返回真实鉴权 + 真实网关路由；selected 记录每次请求
// 鉴权后交给下游的分组平台（即 getGroupPlatform 的分派依据）。
func newMultiGroupRouteTestRouter(t *testing.T) (*gin.Engine, *string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	mk := func(id int64, platform string) *service.Group {
		return &service.Group{ID: id, Name: platform, Platform: platform, Status: service.StatusActive, Hydrated: true, RateMultiplier: 1}
	}
	claudeGroup, openaiGroup, grokGroup := mk(1, service.PlatformAnthropic), mk(2, service.PlatformOpenAI), mk(3, service.PlatformGrok)
	user := &service.User{ID: 7, Role: service.RoleUser, Status: service.StatusActive, Balance: 10}
	apiKey := &service.APIKey{
		ID: 100, UserID: user.ID, Key: "multi-group-route-key", Status: service.StatusActive, User: user,
		GroupID: &claudeGroup.ID, Group: claudeGroup,
		GroupIDs: []int64{1, 2, 3}, Groups: []*service.Group{claudeGroup, openaiGroup, grokGroup},
	}
	cfg := &config.Config{
		RunMode: config.RunModeStandard,
		Gateway: config.GatewayConfig{MaxBodySize: 1 << 20, TextMaxBodySize: 1 << 20},
	}
	apiKeyService := service.NewAPIKeyService(&multiGroupRouteAPIKeyRepo{keyBillingRouteAPIKeyRepo{apiKey: apiKey}}, nil, nil, nil, nil, nil, cfg)
	diag := multiGroupRouteDiagnoser{
		1: {"claude-sonnet-4": true},
		2: {"gpt-5": true, "gpt-image-1": true},
		3: {"grok-4.3": true, "grok-imagine-video": true},
	}
	apiKeyService.SetGroupRouter(service.NewGroupRouter(service.GroupRouterDeps{
		Anthropic: diag, OpenAI: diag,
		VideoGroups: multiGroupRouteVideoGroups{"vid-grok": 3},
	}))

	selected := new(string)
	realAuth := servermiddleware.NewAPIKeyAuthMiddleware(apiKeyService, nil, cfg)
	auth := servermiddleware.APIKeyAuthMiddleware(func(c *gin.Context) {
		*selected = ""
		// 鉴权内部会调用 c.Next()；下游 handler 可能 panic，用 defer 保证记录。
		defer func() {
			if k, ok := servermiddleware.GetAPIKeyFromContext(c); ok && k != nil && k.Group != nil {
				*selected = k.Group.Platform
			}
		}()
		gin.HandlerFunc(realAuth)(c)
	})

	router := gin.New()
	// 测试用的空 handler 依赖在真正转发前会 panic；能走到这里说明鉴权、白名单、
	// composite 与平台门都已放行，按 500 返回即可。
	router.Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.AbortWithStatus(http.StatusInternalServerError)
	}))
	RegisterGatewayRoutes(
		router,
		&handler.Handlers{
			Gateway:       &handler.GatewayHandler{},
			OpenAIGateway: &handler.OpenAIGatewayHandler{},
			AsyncImage:    handler.NewAsyncImageHandler(nil, nil),
		},
		auth, apiKeyService, nil, nil, nil, nil, cfg,
	)
	return router, selected
}

func serveMultiGroup(router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("x-api-key", "multi-group-route-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestMultiGroupKeyRoutesEachModelToItsPlatform(t *testing.T) {
	router, selected := newMultiGroupRouteTestRouter(t)

	for _, tc := range []struct {
		name, method, path, body, want string
	}{
		{"claude on messages", http.MethodPost, "/v1/messages", `{"model":"claude-sonnet-4","messages":[]}`, service.PlatformAnthropic},
		{"gpt on messages", http.MethodPost, "/v1/messages", `{"model":"gpt-5","messages":[]}`, service.PlatformOpenAI},
		{"gpt on chat completions", http.MethodPost, "/v1/chat/completions", `{"model":"gpt-5","messages":[]}`, service.PlatformOpenAI},
		{"gpt on responses", http.MethodPost, "/v1/responses", `{"model":"gpt-5","input":"hi"}`, service.PlatformOpenAI},
		{"grok on chat completions", http.MethodPost, "/v1/chat/completions", `{"model":"grok-4.3","messages":[]}`, service.PlatformGrok},
		{"root alias", http.MethodPost, "/chat/completions", `{"model":"grok-4.3","messages":[]}`, service.PlatformGrok},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := serveMultiGroup(router, tc.method, tc.path, tc.body)
			require.NotEqual(t, http.StatusNotFound, w.Code, w.Body.String())
			require.Equal(t, tc.want, *selected)
		})
	}
}

// 平台门按选中的分组判断：主路由是 Anthropic，但 gpt-image / grok 视频请求能通过。
func TestMultiGroupKeyPassesPlatformGatesOfSelectedGroup(t *testing.T) {
	router, selected := newMultiGroupRouteTestRouter(t)

	w := serveMultiGroup(router, http.MethodPost, "/v1/images/generations", `{"model":"gpt-image-1","prompt":"cat"}`)
	require.NotContains(t, w.Body.String(), "not supported for this platform")
	require.Equal(t, service.PlatformOpenAI, *selected)

	w = serveMultiGroup(router, http.MethodPost, "/v1/videos/generations", `{"model":"grok-imagine-video","prompt":"cat"}`)
	require.NotContains(t, w.Body.String(), "not supported for this platform")
	require.Equal(t, service.PlatformGrok, *selected)

	// 视频状态查询不带模型：回到创建时的 Grok 分组
	w = serveMultiGroup(router, http.MethodGet, "/v1/videos/vid-grok", "")
	require.NotContains(t, w.Body.String(), "not supported for this platform")
	require.Equal(t, service.PlatformGrok, *selected)

	// 没有记录的任务走主路由（Anthropic），被平台门拒绝，与单分组行为一致
	w = serveMultiGroup(router, http.MethodGet, "/v1/videos/unknown", "")
	require.Contains(t, w.Body.String(), "not supported for this platform")
	require.Equal(t, service.PlatformAnthropic, *selected)
}

func TestMultiGroupKeyUnsupportedModelReturns404BeforeHandlers(t *testing.T) {
	router, selected := newMultiGroupRouteTestRouter(t)

	w := serveMultiGroup(router, http.MethodPost, "/v1/messages", `{"model":"gemini-2.5-pro","messages":[]}`)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "not available for any group bound to this API key")
	require.Empty(t, *selected, "request never reaches a handler")
}
