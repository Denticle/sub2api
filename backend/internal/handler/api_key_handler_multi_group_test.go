package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// 5.3：用户端密钥接口的多分组字段（add-api-key-multi-group-routing）。

type multiGroupKeyRepo struct {
	service.APIKeyRepository
	nextID int64
	keys   map[int64]*service.APIKey
}

func (r *multiGroupKeyRepo) Create(_ context.Context, key *service.APIKey) error {
	r.nextID++
	key.ID = r.nextID
	clone := *key
	r.keys[key.ID] = &clone
	return nil
}

func (r *multiGroupKeyRepo) GetByID(_ context.Context, id int64) (*service.APIKey, error) {
	key, ok := r.keys[id]
	if !ok {
		return nil, service.ErrAPIKeyNotFound
	}
	clone := *key
	return &clone, nil
}

func (r *multiGroupKeyRepo) Update(_ context.Context, key *service.APIKey, _ service.APIKeyUpdateFields) error {
	clone := *key
	r.keys[key.ID] = &clone
	return nil
}

type multiGroupUserRepo struct {
	service.UserRepository
	user *service.User
}

func (r *multiGroupUserRepo) GetByID(_ context.Context, id int64) (*service.User, error) {
	if r.user == nil || r.user.ID != id {
		return nil, service.ErrUserNotFound
	}
	clone := *r.user
	return &clone, nil
}

type multiGroupGroupRepo struct {
	service.GroupRepository
	groups map[int64]*service.Group
}

func (r *multiGroupGroupRepo) GetByID(_ context.Context, id int64) (*service.Group, error) {
	g, ok := r.groups[id]
	if !ok {
		return nil, service.ErrGroupNotFound
	}
	clone := *g
	return &clone, nil
}

type multiGroupKeyFixture struct {
	router  *gin.Engine
	keyRepo *multiGroupKeyRepo
}

const multiGroupKeyUserID = int64(7)

// 分组 1–3 是公开标准分组；分组 9 是专属分组，用户没有权限。
func newMultiGroupKeyFixture(t *testing.T) *multiGroupKeyFixture {
	t.Helper()
	gin.SetMode(gin.TestMode)

	groups := map[int64]*service.Group{
		1: {ID: 1, Name: "claude", Platform: service.PlatformAnthropic, Status: service.StatusActive, RateMultiplier: 1},
		2: {ID: 2, Name: "openai", Platform: service.PlatformOpenAI, Status: service.StatusActive, RateMultiplier: 0.8},
		3: {ID: 3, Name: "gemini", Platform: service.PlatformGemini, Status: service.StatusActive, RateMultiplier: 1.5},
		9: {ID: 9, Name: "vip", Platform: service.PlatformAnthropic, Status: service.StatusActive, IsExclusive: true, RateMultiplier: 1},
	}
	for id := int64(10); id <= 30; id++ {
		groups[id] = &service.Group{ID: id, Name: "g" + strconv.FormatInt(id, 10), Platform: service.PlatformAnthropic, Status: service.StatusActive, RateMultiplier: 1}
	}

	keyRepo := &multiGroupKeyRepo{keys: map[int64]*service.APIKey{}}
	userRepo := &multiGroupUserRepo{user: &service.User{ID: multiGroupKeyUserID, Status: service.StatusActive, Role: service.RoleUser}}
	groupRepo := &multiGroupGroupRepo{groups: groups}

	svc := service.NewAPIKeyService(keyRepo, userRepo, groupRepo, nil, nil, nil, &config.Config{RunMode: config.RunModeStandard})
	h := NewAPIKeyHandler(svc)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: multiGroupKeyUserID})
		c.Next()
	})
	router.POST("/keys", h.Create)
	router.PUT("/keys/:id", h.Update)
	return &multiGroupKeyFixture{router: router, keyRepo: keyRepo}
}

type multiGroupKeyResponse struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
	Data   struct {
		ID        int64   `json:"id"`
		GroupID   *int64  `json:"group_id"`
		GroupIDs  []int64 `json:"group_ids"`
		GroupMode string  `json:"group_mode"`
		Groups    []struct {
			ID             int64   `json:"id"`
			Platform       string  `json:"platform"`
			RateMultiplier float64 `json:"rate_multiplier"`
		} `json:"groups"`
	} `json:"data"`
}

func (f *multiGroupKeyFixture) do(t *testing.T, method, path string, body any) (int, multiGroupKeyResponse) {
	t.Helper()
	raw, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(method, path, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	f.router.ServeHTTP(rec, req)
	var resp multiGroupKeyResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp), rec.Body.String())
	return rec.Code, resp
}

func (f *multiGroupKeyFixture) seedKey(groupIDs ...int64) int64 {
	key := &service.APIKey{UserID: multiGroupKeyUserID, Key: "sk-seed", Name: "seed", Status: service.StatusActive, GroupMode: service.APIKeyGroupModeManual}
	if len(groupIDs) > 0 {
		primary := groupIDs[0]
		key.GroupID = &primary
		key.GroupIDs = append([]int64(nil), groupIDs...)
	}
	_ = f.keyRepo.Create(context.Background(), key)
	return key.ID
}

func TestAPIKeyCreate_OnlyGroupIDIsSingleGroupKey(t *testing.T) {
	f := newMultiGroupKeyFixture(t)

	code, resp := f.do(t, http.MethodPost, "/keys", map[string]any{"name": "k", "group_id": 2})

	require.Equal(t, http.StatusOK, code, "%+v", resp)
	require.NotNil(t, resp.Data.GroupID)
	require.Equal(t, int64(2), *resp.Data.GroupID)
	require.Equal(t, []int64{2}, resp.Data.GroupIDs)
	require.Equal(t, service.APIKeyGroupModeManual, resp.Data.GroupMode)
	require.Empty(t, resp.Data.Groups, "single-group keys do not return the groups list")
}

func TestAPIKeyCreate_GroupIDsKeepsOrderAndReturnsGroups(t *testing.T) {
	f := newMultiGroupKeyFixture(t)

	code, resp := f.do(t, http.MethodPost, "/keys", map[string]any{"name": "k", "group_ids": []int64{3, 1, 2}})

	require.Equal(t, http.StatusOK, code, "%+v", resp)
	require.Equal(t, int64(3), *resp.Data.GroupID, "first entry becomes the primary route")
	require.Equal(t, []int64{3, 1, 2}, resp.Data.GroupIDs)
	require.Len(t, resp.Data.Groups, 3)
	require.Equal(t, service.PlatformGemini, resp.Data.Groups[0].Platform)
	require.Equal(t, 1.5, resp.Data.Groups[0].RateMultiplier)
	require.Equal(t, 0.8, resp.Data.Groups[2].RateMultiplier)

	stored := f.keyRepo.keys[resp.Data.ID]
	require.Equal(t, []int64{3, 1, 2}, stored.GroupIDs, "order is persisted")
	require.Equal(t, int64(3), *stored.GroupID)
}

func TestAPIKeyCreate_GroupIDsWinsOverGroupID(t *testing.T) {
	f := newMultiGroupKeyFixture(t)

	code, resp := f.do(t, http.MethodPost, "/keys", map[string]any{"name": "k", "group_id": 1, "group_ids": []int64{2, 3}})

	require.Equal(t, http.StatusOK, code, "%+v", resp)
	require.Equal(t, []int64{2, 3}, resp.Data.GroupIDs)
	require.Equal(t, int64(2), *resp.Data.GroupID)
}

func TestAPIKeyCreate_UnbindableGroupIsForbidden(t *testing.T) {
	f := newMultiGroupKeyFixture(t)

	code, resp := f.do(t, http.MethodPost, "/keys", map[string]any{"name": "k", "group_ids": []int64{1, 9}})

	require.Equal(t, http.StatusForbidden, code, "%+v", resp)
	require.Equal(t, "GROUP_NOT_ALLOWED", resp.Reason)
	require.Empty(t, f.keyRepo.keys, "nothing is written when any group is not allowed")
}

func TestAPIKeyCreate_GroupListValidation(t *testing.T) {
	tooMany := make([]int64, 0, service.MaxAPIKeyGroups+1)
	for id := int64(10); id < 10+service.MaxAPIKeyGroups+1; id++ {
		tooMany = append(tooMany, id)
	}
	cases := []struct {
		name   string
		body   map[string]any
		reason string
	}{
		{"too many", map[string]any{"name": "k", "group_ids": tooMany}, "API_KEY_GROUPS_TOO_MANY"},
		{"empty", map[string]any{"name": "k", "group_ids": []int64{}}, "API_KEY_GROUPS_EMPTY"},
		{"duplicate", map[string]any{"name": "k", "group_ids": []int64{1, 2, 1}}, "API_KEY_GROUPS_DUPLICATE"},
		{"auto mode", map[string]any{"name": "k", "group_mode": "auto"}, "API_KEY_GROUP_MODE_NOT_ALLOWED"},
		{"unknown mode", map[string]any{"name": "k", "group_mode": "random", "group_id": 1}, "API_KEY_GROUP_MODE_INVALID"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newMultiGroupKeyFixture(t)
			code, resp := f.do(t, http.MethodPost, "/keys", tc.body)
			require.Equal(t, http.StatusBadRequest, code, "%+v", resp)
			require.Equal(t, tc.reason, resp.Reason)
			require.Empty(t, f.keyRepo.keys)
		})
	}
}

func TestAPIKeyCreate_MaxGroupsIsAccepted(t *testing.T) {
	f := newMultiGroupKeyFixture(t)
	ids := make([]int64, 0, service.MaxAPIKeyGroups)
	for id := int64(10); id < 10+service.MaxAPIKeyGroups; id++ {
		ids = append(ids, id)
	}

	code, resp := f.do(t, http.MethodPost, "/keys", map[string]any{"name": "k", "group_ids": ids})

	require.Equal(t, http.StatusOK, code, "%+v", resp)
	require.Len(t, resp.Data.GroupIDs, service.MaxAPIKeyGroups)
}

func TestAPIKeyUpdate_GroupIDsReplacesList(t *testing.T) {
	f := newMultiGroupKeyFixture(t)
	id := f.seedKey(1, 2)

	code, resp := f.do(t, http.MethodPut, "/keys/"+strconv.FormatInt(id, 10), map[string]any{"group_ids": []int64{2, 3, 1}})

	require.Equal(t, http.StatusOK, code, "%+v", resp)
	require.Equal(t, []int64{2, 3, 1}, resp.Data.GroupIDs)
	require.Equal(t, int64(2), *f.keyRepo.keys[id].GroupID)
}

func TestAPIKeyUpdate_OnlyNameKeepsGroups(t *testing.T) {
	f := newMultiGroupKeyFixture(t)
	id := f.seedKey(1, 2)

	// 旧版前端编辑名称时会原样回传 group_id。
	code, resp := f.do(t, http.MethodPut, "/keys/"+strconv.FormatInt(id, 10), map[string]any{"name": "renamed", "group_id": 1})

	require.Equal(t, http.StatusOK, code, "%+v", resp)
	require.Equal(t, []int64{1, 2}, f.keyRepo.keys[id].GroupIDs)
}

func TestAPIKeyUpdate_Rejections(t *testing.T) {
	cases := []struct {
		name   string
		body   map[string]any
		status int
		reason string
	}{
		{"unbindable", map[string]any{"group_ids": []int64{2, 9}}, http.StatusForbidden, "GROUP_NOT_ALLOWED"},
		{"duplicate", map[string]any{"group_ids": []int64{2, 2}}, http.StatusBadRequest, "API_KEY_GROUPS_DUPLICATE"},
		{"auto mode", map[string]any{"group_mode": "auto"}, http.StatusBadRequest, "API_KEY_GROUP_MODE_NOT_ALLOWED"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newMultiGroupKeyFixture(t)
			id := f.seedKey(1, 2)
			code, resp := f.do(t, http.MethodPut, "/keys/"+strconv.FormatInt(id, 10), tc.body)
			require.Equal(t, tc.status, code, "%+v", resp)
			require.Equal(t, tc.reason, resp.Reason)
			require.Equal(t, []int64{1, 2}, f.keyRepo.keys[id].GroupIDs, "stored groups are unchanged")
		})
	}
}
