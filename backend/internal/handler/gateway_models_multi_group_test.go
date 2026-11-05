package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func mappedAccount(id int64, platform string, models ...string) service.Account {
	mapping := map[string]any{}
	for _, m := range models {
		mapping[m] = m
	}
	return service.Account{
		ID: id, Platform: platform, Type: service.AccountTypeAPIKey,
		Status: service.StatusActive, Schedulable: true,
		Credentials: map[string]any{"model_mapping": mapping},
	}
}

func requestMultiGroupModels(h *GatewayHandler, groups []*service.Group, modelID string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/models", nil)
	if modelID != "" {
		c.Params = gin.Params{{Key: "model", Value: modelID}}
	}
	c.Set(string(middleware2.ContextKeyAPIKey), &service.APIKey{GroupID: &groups[0].ID, Group: groups[0]})
	c.Request = c.Request.WithContext(service.WithAPIKeyGroupRoute(context.Background(), service.APIKeyGroupRoute{
		All: groups, Selected: groups[0], Outcome: service.GroupRoutePrimary,
	}))
	h.Models(c)
	return rec
}

func modelIDsFromList(t *testing.T, rec *httptest.ResponseRecorder) []string {
	t.Helper()
	var catalog struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &catalog))
	ids := make([]string, 0, len(catalog.Data))
	for _, m := range catalog.Data {
		ids = append(ids, m.ID)
	}
	return ids
}

func TestMultiGroupModelsMergesCandidateCatalogs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	claudeGroup := &service.Group{ID: 81, Platform: service.PlatformAnthropic}
	openaiGroup := &service.Group{ID: 82, Platform: service.PlatformOpenAI}
	h := newGatewayModelsHandlerForTest(&gatewayModelsAccountRepoStub{byGroup: map[int64][]service.Account{
		81: {mappedAccount(1, service.PlatformAnthropic, "claude-only")},
		82: {mappedAccount(2, service.PlatformOpenAI, "gpt-only", "gpt-hidden")},
	}})
	// 各组白名单在合并前生效
	openaiGroup.ModelAllowlist = service.GroupModelAllowlist{Enabled: true, Models: []string{"gpt-only"}}

	list := requestMultiGroupModels(h, []*service.Group{claudeGroup, openaiGroup}, "")
	require.Equal(t, http.StatusOK, list.Code, list.Body.String())
	ids := modelIDsFromList(t, list)
	require.Contains(t, ids, "claude-only")
	require.Contains(t, ids, "gpt-only")
	require.NotContains(t, ids, "gpt-hidden")

	// /models/:model 在任一组可见即返回
	got := requestMultiGroupModels(h, []*service.Group{claudeGroup, openaiGroup}, "gpt-only")
	require.Equal(t, http.StatusOK, got.Code, got.Body.String())
	require.Contains(t, got.Body.String(), `"gpt-only"`)
	require.Equal(t, http.StatusNotFound, requestMultiGroupModels(h, []*service.Group{claudeGroup, openaiGroup}, "gpt-hidden").Code)
}

// 同平台多分组保留该平台的列表格式（OpenAI 形状带 object/owned_by）。
func TestMultiGroupModelsSamePlatformKeepsShape(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := &service.Group{ID: 91, Platform: service.PlatformOpenAI}
	b := &service.Group{ID: 92, Platform: service.PlatformOpenAI}
	h := newGatewayModelsHandlerForTest(&gatewayModelsAccountRepoStub{byGroup: map[int64][]service.Account{
		91: {mappedAccount(1, service.PlatformOpenAI, "gpt-a")},
		92: {mappedAccount(2, service.PlatformOpenAI, "gpt-b", "gpt-a")},
	}})

	list := requestMultiGroupModels(h, []*service.Group{a, b}, "")
	require.Equal(t, http.StatusOK, list.Code)
	require.Equal(t, []string{"gpt-a", "gpt-b"}, modelIDsFromList(t, list), "merged in candidate order without duplicates")
	require.Contains(t, list.Body.String(), `"owned_by":"openai"`)
}

func TestAntigravityModelsAllowedByAnyCandidate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	blocked := &service.Group{ID: 1, ModelAllowlist: service.GroupModelAllowlist{Enabled: true, Models: []string{"none"}}}
	open := &service.Group{ID: 2}

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/antigravity/models", nil)
	c.Set(string(middleware2.ContextKeyAPIKey), &service.APIKey{GroupID: &blocked.ID, Group: blocked})
	(&GatewayHandler{}).AntigravityModels(c)
	require.Empty(t, modelIDsFromList(t, rec), "single blocked group lists nothing")

	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/antigravity/models", nil)
	c.Set(string(middleware2.ContextKeyAPIKey), &service.APIKey{GroupID: &blocked.ID, Group: blocked})
	c.Request = c.Request.WithContext(service.WithAPIKeyGroupRoute(context.Background(), service.APIKeyGroupRoute{
		All: []*service.Group{blocked, open},
	}))
	(&GatewayHandler{}).AntigravityModels(c)
	require.NotEmpty(t, modelIDsFromList(t, rec), "a permissive candidate exposes the catalog")
}
