//go:build unit

package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// groupRouteDiagnoserStub 按分组声明支持的模型。
type groupRouteDiagnoserStub map[int64]map[string]bool

func (s groupRouteDiagnoserStub) DiagnoseModelAvailabilityForPlatform(_ context.Context, groupID *int64, model, _ string) service.ModelAvailabilityDiagnosis {
	return service.ModelAvailabilityDiagnosis{HasAccountsInPool: true, HasModelSupport: s[*groupID][model]}
}

func routeTestGroup(id int64, platform string) *service.Group {
	return &service.Group{ID: id, Name: platform, Platform: platform, Status: service.StatusActive, Hydrated: true, RateMultiplier: 1}
}

// newGroupRouteTestService 返回绑定 groups（按顺序）的密钥与挂好选组服务的 APIKeyService。
// 分组 1 只支持 claude-sonnet-4，分组 2 只支持 gpt-5。
func newGroupRouteTestService(t *testing.T, groups ...*service.Group) (*service.APIKeyService, *service.APIKey) {
	t.Helper()
	ids := make([]int64, len(groups))
	for i, g := range groups {
		ids[i] = g.ID
	}
	apiKey := &service.APIKey{
		ID: 100, UserID: 7, Key: "multi-key", Status: service.StatusActive,
		GroupID: &ids[0], GroupIDs: ids, Group: groups[0],
		User: &service.User{ID: 7, Role: service.RoleUser, Status: service.StatusActive, Balance: 10, Concurrency: 3},
	}
	if len(groups) > 1 {
		apiKey.Groups = groups
	}
	repo := &stubApiKeyRepo{getByKey: func(_ context.Context, key string) (*service.APIKey, error) {
		if key != apiKey.Key {
			return nil, service.ErrAPIKeyNotFound
		}
		clone := *apiKey
		return &clone, nil
	}}
	cfg := &config.Config{RunMode: config.RunModeStandard}
	svc := service.NewAPIKeyService(repo, nil, nil, nil, nil, nil, cfg)
	diag := groupRouteDiagnoserStub{
		1: {"claude-sonnet-4": true},
		2: {"gpt-5": true},
	}
	svc.SetGroupRouter(service.NewGroupRouter(service.GroupRouterDeps{Anthropic: diag, OpenAI: diag}))
	return svc, apiKey
}

type groupRouteSeen struct {
	keyGroupID int64
	ctxGroupID int64
	body       string
	route      service.APIKeyGroupRoute
	hasRoute   bool
}

func groupRouteRecorder(seen *groupRouteSeen) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, _ := GetAPIKeyFromContext(c)
		if apiKey != nil && apiKey.GroupID != nil {
			seen.keyGroupID = *apiKey.GroupID
		}
		if g, ok := c.Request.Context().Value(ctxkey.Group).(*service.Group); ok && g != nil {
			seen.ctxGroupID = g.ID
		}
		if c.Request.Body != nil {
			b, _ := io.ReadAll(c.Request.Body)
			seen.body = string(b)
		}
		seen.route, seen.hasRoute = service.APIKeyGroupRouteFromContext(c.Request.Context())
		c.Status(http.StatusOK)
	}
}

func serveGroupRoute(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("x-api-key", "multi-key")
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func TestAPIKeyAuth_MultiGroupKeySelectsGroupByModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc, _ := newGroupRouteTestService(t, routeTestGroup(1, service.PlatformAnthropic), routeTestGroup(2, service.PlatformOpenAI))
	cfg := &config.Config{RunMode: config.RunModeStandard}

	var seen groupRouteSeen
	r := gin.New()
	r.Use(gin.HandlerFunc(NewAPIKeyAuthMiddleware(svc, nil, cfg)))
	r.POST("/v1/messages", groupRouteRecorder(&seen))

	body := `{"model":"gpt-5","messages":[]}`
	w := serveGroupRoute(r, http.MethodPost, "/v1/messages", body)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.Equal(t, int64(2), seen.keyGroupID, "downstream apiKey uses the selected group")
	require.Equal(t, int64(2), seen.ctxGroupID, "group context follows the selected group")
	require.Equal(t, body, seen.body, "request body is restored for the handler")
	require.True(t, seen.hasRoute)
	require.Len(t, seen.route.All, 2)
	require.Equal(t, service.GroupRouteSelected, seen.route.Outcome)

	seen = groupRouteSeen{}
	w = serveGroupRoute(r, http.MethodPost, "/v1/messages", `{"model":"claude-sonnet-4"}`)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, int64(1), seen.keyGroupID)
}

func TestAPIKeyAuth_MultiGroupKeyUnknownModelReturns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc, _ := newGroupRouteTestService(t, routeTestGroup(1, service.PlatformAnthropic), routeTestGroup(2, service.PlatformOpenAI))
	cfg := &config.Config{RunMode: config.RunModeStandard}

	var reached atomic.Bool
	r := gin.New()
	r.Use(gin.HandlerFunc(NewAPIKeyAuthMiddleware(svc, nil, cfg)))
	r.POST("/v1/messages", func(c *gin.Context) { reached.Store(true) })
	r.POST("/v1/chat/completions", func(c *gin.Context) { reached.Store(true) })

	w := serveGroupRoute(r, http.MethodPost, "/v1/messages", `{"model":"gemini-2.5-pro"}`)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), `"type":"error"`, "Anthropic error format on /messages")
	require.Contains(t, w.Body.String(), "gemini-2.5-pro")

	w = serveGroupRoute(r, http.MethodPost, "/v1/chat/completions", `{"model":"gemini-2.5-pro"}`)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), `"error"`)
	require.False(t, reached.Load())
}

// countingBody 记录请求体被读了多少次。
type countingBody struct {
	io.Reader
	reads *atomic.Int32
}

func (b countingBody) Read(p []byte) (int, error) {
	b.reads.Add(1)
	return b.Reader.Read(p)
}

func (b countingBody) Close() error { return nil }

// 单分组密钥走快速路径：鉴权阶段不读请求体（tasks 4.3）。
func TestAPIKeyAuth_SingleGroupKeyDoesNotReadBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc, _ := newGroupRouteTestService(t, routeTestGroup(2, service.PlatformOpenAI))
	cfg := &config.Config{RunMode: config.RunModeStandard}

	var reads atomic.Int32
	var readsAtHandler int32 = -1
	var seen groupRouteSeen
	r := gin.New()
	r.Use(gin.HandlerFunc(NewAPIKeyAuthMiddleware(svc, nil, cfg)))
	r.POST("/v1/chat/completions", func(c *gin.Context) {
		readsAtHandler = reads.Load()
		groupRouteRecorder(&seen)(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	req.Body = countingBody{Reader: bytes.NewBufferString(`{"model":"gpt-5"}`), reads: &reads}
	req.Header.Set("x-api-key", "multi-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.Equal(t, int32(0), readsAtHandler, "auth must not read the body for single-group keys")
	require.Equal(t, int64(2), seen.keyGroupID)
	require.False(t, seen.hasRoute, "single-group keys carry no route context")
}

// 不带模型的请求（模型列表等）走主路由，但 context 里带上全部候选（design D4）。
func TestAPIKeyAuth_MultiGroupKeyWithoutModelUsesPrimary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc, _ := newGroupRouteTestService(t, routeTestGroup(1, service.PlatformAnthropic), routeTestGroup(2, service.PlatformOpenAI))
	cfg := &config.Config{RunMode: config.RunModeStandard}

	var seen groupRouteSeen
	r := gin.New()
	r.Use(gin.HandlerFunc(NewAPIKeyAuthMiddleware(svc, nil, cfg)))
	r.GET("/v1/models", groupRouteRecorder(&seen))

	w := serveGroupRoute(r, http.MethodGet, "/v1/models", "")
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, int64(1), seen.keyGroupID)
	require.True(t, seen.hasRoute)
	require.Equal(t, service.GroupRoutePrimary, seen.route.Outcome)
	require.Len(t, seen.route.All, 2)
}

func TestAPIKeyAuthGoogle_MultiGroupKeySelectsGroupByPathModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gemini := routeTestGroup(3, service.PlatformGemini)
	svc, _ := newGroupRouteTestService(t, routeTestGroup(1, service.PlatformAnthropic), gemini)
	svc.SetGroupRouter(service.NewGroupRouter(service.GroupRouterDeps{
		Anthropic: groupRouteDiagnoserStub{1: {"claude-sonnet-4": true}, 3: {"gemini-2.5-pro": true}},
	}))
	cfg := &config.Config{RunMode: config.RunModeStandard}

	var seen groupRouteSeen
	r := gin.New()
	r.Use(APIKeyAuthWithSubscriptionGoogle(svc, nil, cfg))
	r.POST("/v1beta/models/*modelAction", groupRouteRecorder(&seen))

	w := serveGroupRoute(r, http.MethodPost, "/v1beta/models/gemini-2.5-pro:generateContent", `{"contents":[]}`)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.Equal(t, int64(3), seen.keyGroupID)

	w = serveGroupRoute(r, http.MethodPost, "/v1beta/models/gpt-5:generateContent", `{}`)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), `"status":"NOT_FOUND"`, "Google error format")
}

type groupRouteVideoStub map[string]int64

func (s groupRouteVideoStub) ResolveVideoRequestGroup(_ context.Context, requestID string, _, _ int64) (int64, error) {
	return s[requestID], nil
}

// 视频状态查询不带模型：回到创建任务时的分组（tasks 4.5）。
func TestAPIKeyAuth_VideoLookupUsesCreateTimeGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc, _ := newGroupRouteTestService(t, routeTestGroup(1, service.PlatformAnthropic), routeTestGroup(2, service.PlatformGrok))
	svc.SetGroupRouter(service.NewGroupRouter(service.GroupRouterDeps{VideoGroups: groupRouteVideoStub{"vid-1": 2}}))
	cfg := &config.Config{RunMode: config.RunModeStandard}

	var seen groupRouteSeen
	r := gin.New()
	r.Use(gin.HandlerFunc(NewAPIKeyAuthMiddleware(svc, nil, cfg)))
	r.GET("/v1/videos/:request_id", groupRouteRecorder(&seen))
	r.GET("/v1/videos/:request_id/content", groupRouteRecorder(&seen))

	w := serveGroupRoute(r, http.MethodGet, "/v1/videos/vid-1", "")
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, int64(2), seen.keyGroupID)

	seen = groupRouteSeen{}
	w = serveGroupRoute(r, http.MethodGet, "/v1/videos/vid-1/content", "")
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, int64(2), seen.keyGroupID)

	seen = groupRouteSeen{}
	w = serveGroupRoute(r, http.MethodGet, "/v1/videos/legacy", "")
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, int64(1), seen.keyGroupID, "unknown task falls back to the primary group")
}

// Seedance 任务创建时按 ResponseID（带 "seedance:" 前缀）记录分组；查询 / 删除要用同一个 key。
func TestAPIKeyAuth_SeedanceTaskLookupUsesPrefixedTaskKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc, _ := newGroupRouteTestService(t, routeTestGroup(1, service.PlatformAnthropic), routeTestGroup(2, service.PlatformOpenAI))
	svc.SetGroupRouter(service.NewGroupRouter(service.GroupRouterDeps{
		VideoGroups: groupRouteVideoStub{service.SeedanceTaskKey("task-9"): 2},
	}))
	cfg := &config.Config{RunMode: config.RunModeStandard}

	var seen groupRouteSeen
	r := gin.New()
	r.Use(gin.HandlerFunc(NewAPIKeyAuthMiddleware(svc, nil, cfg)))
	r.GET("/api/v3/contents/generations/tasks/:task_id", groupRouteRecorder(&seen))
	r.DELETE("/api/v3/contents/generations/tasks/:task_id", groupRouteRecorder(&seen))

	for _, method := range []string{http.MethodGet, http.MethodDelete} {
		seen = groupRouteSeen{}
		w := serveGroupRoute(r, method, "/api/v3/contents/generations/tasks/task-9", "")
		require.Equal(t, http.StatusOK, w.Code, method)
		require.Equal(t, int64(2), seen.keyGroupID, "%s must return to the create-time group", method)
	}
}

// 线上复现：工作台用 multipart 调 /v1/images/edits，请求体里没有 model 字段；
// 选组必须避开没开生图的分组，而不是退回倍率最低的主路由。
func TestAPIKeyAuth_ImageEditsAvoidsGroupWithoutImageGeneration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	noImage := routeTestGroup(44, service.PlatformOpenAI)
	withImage := routeTestGroup(43, service.PlatformOpenAI)
	withImage.AllowImageGeneration = true
	svc, _ := newGroupRouteTestService(t, noImage, withImage)
	cfg := &config.Config{RunMode: config.RunModeStandard}

	var seen groupRouteSeen
	r := gin.New()
	r.Use(gin.HandlerFunc(NewAPIKeyAuthMiddleware(svc, nil, cfg)))
	r.POST("/v1/images/edits", groupRouteRecorder(&seen))
	r.POST("/v1/chat/completions", groupRouteRecorder(&seen))

	body := "--b\r\nContent-Disposition: form-data; name=\"prompt\"\r\n\r\ncat\r\n--b--\r\n"
	req := httptest.NewRequest(http.MethodPost, "/v1/images/edits", strings.NewReader(body))
	req.Header.Set("x-api-key", "multi-key")
	req.Header.Set("Content-Type", "multipart/form-data; boundary=b")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.Equal(t, int64(43), seen.keyGroupID)
	require.Equal(t, body, seen.body, "multipart body is restored for the handler")

	// 非生图请求不受影响，仍走主路由
	seen = groupRouteSeen{}
	w = serveGroupRoute(r, http.MethodPost, "/v1/chat/completions", `{"messages":[]}`)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, int64(44), seen.keyGroupID)
}
