//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// 每个分组配置了哪些模型（第 3 条）。
type routerDiagnoserStub struct {
	models map[int64]map[string]bool
	calls  int
}

func (s *routerDiagnoserStub) DiagnoseModelAvailabilityForPlatform(_ context.Context, groupID *int64, model, _ string) ModelAvailabilityDiagnosis {
	s.calls++
	ok := s.models[*groupID][model]
	return ModelAvailabilityDiagnosis{HasAccountsInPool: true, HasModelSupport: ok}
}

// 每个分组的调度快照（第 5 条）。
type routerSchedulerStub struct {
	accounts map[int64][]Account
	err      error
}

func (s *routerSchedulerStub) ListSchedulableAccounts(_ context.Context, groupID *int64, _ string, _ bool) ([]Account, bool, error) {
	return s.accounts[*groupID], false, s.err
}

type routerSubscriptionStub struct {
	subs map[int64]*UserSubscription
}

func (s *routerSubscriptionStub) GetActiveSubscription(_ context.Context, _ int64, groupID int64) (*UserSubscription, error) {
	if sub, ok := s.subs[groupID]; ok {
		cp := *sub
		return &cp, nil
	}
	return nil, ErrSubscriptionNotFound
}

func (s *routerSubscriptionStub) ValidateAndCheckLimits(sub *UserSubscription, group *Group) (bool, error) {
	if !sub.CheckDailyLimit(group, 0) {
		return false, ErrDailyLimitExceeded
	}
	return false, nil
}

type routerGroupSourceStub struct {
	groups []Group
	calls  int
}

func (s *routerGroupSourceStub) GetAvailableGroups(context.Context, int64) ([]Group, error) {
	s.calls++
	out := make([]Group, len(s.groups))
	copy(out, s.groups)
	return out, nil
}

type routerCompositeStub struct {
	decisions map[string]CompositeRouteDecision
}

func (s *routerCompositeStub) Resolve(_ context.Context, _ int64, model, _ string) (CompositeRouteDecision, error) {
	return s.decisions[model], nil
}

func healthyAccount() Account {
	return Account{Status: StatusActive, Schedulable: true}
}

func rateLimitedAccount() Account {
	until := time.Now().Add(time.Hour)
	return Account{Status: StatusActive, Schedulable: true, RateLimitResetAt: &until}
}

func activeGroup(id int64, platform string) *Group {
	return &Group{ID: id, Platform: platform, Status: StatusActive, Hydrated: true, RateMultiplier: 1}
}

func multiGroupKey(groups ...*Group) *APIKey {
	ids := make([]int64, len(groups))
	for i, g := range groups {
		ids[i] = g.ID
	}
	return &APIKey{
		ID: 42, UserID: 7, Status: StatusActive,
		GroupID: &ids[0], GroupIDs: ids, Group: groups[0], Groups: groups,
		User: &User{ID: 7, Status: StatusActive},
	}
}

type routerFixture struct {
	router    *GroupRouter
	anthropic *routerDiagnoserStub
	openai    *routerDiagnoserStub
	scheduler *routerSchedulerStub
	subs      *routerSubscriptionStub
}

// 分组 1 = anthropic（claude-*），分组 2 = openai（gpt-*），都有健康账号。
func newRouterFixture() *routerFixture {
	f := &routerFixture{
		anthropic: &routerDiagnoserStub{models: map[int64]map[string]bool{
			1: {"claude-sonnet-4": true},
		}},
		openai: &routerDiagnoserStub{models: map[int64]map[string]bool{
			2: {"gpt-5": true, "gpt-image-1": true},
			3: {"gpt-5": true},
		}},
		scheduler: &routerSchedulerStub{accounts: map[int64][]Account{
			1: {healthyAccount()},
			2: {healthyAccount()},
			3: {healthyAccount()},
		}},
		subs: &routerSubscriptionStub{subs: map[int64]*UserSubscription{}},
	}
	f.router = NewGroupRouter(GroupRouterDeps{
		Anthropic:     f.anthropic,
		OpenAI:        f.openai,
		Scheduler:     f.scheduler,
		Subscriptions: f.subs,
	})
	return f
}

func TestGroupRouter_SkipsGroupWithoutModel(t *testing.T) {
	f := newRouterFixture()
	key := multiGroupKey(activeGroup(1, PlatformAnthropic), activeGroup(2, PlatformOpenAI))

	res, err := f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, GroupRouteSelected, res.Outcome)
	require.Equal(t, int64(2), res.Group.ID)

	res, err = f.router.Select(context.Background(), key, []string{"claude-sonnet-4"}, "")
	require.NoError(t, err)
	require.Equal(t, int64(1), res.Group.ID)
}

func TestGroupRouter_AllowlistSkipsGroup(t *testing.T) {
	f := newRouterFixture()
	g2 := activeGroup(2, PlatformOpenAI)
	g2.ModelAllowlist = GroupModelAllowlist{Enabled: true, Models: []string{"gpt-image-1"}}
	key := multiGroupKey(g2, activeGroup(3, PlatformOpenAI))

	res, err := f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, int64(3), res.Group.ID)
}

func TestGroupRouter_RateLimitedGroupIsSkipped(t *testing.T) {
	f := newRouterFixture()
	f.scheduler.accounts[2] = []Account{rateLimitedAccount()}
	key := multiGroupKey(activeGroup(2, PlatformOpenAI), activeGroup(3, PlatformOpenAI))

	res, err := f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, GroupRouteSelected, res.Outcome)
	require.Equal(t, int64(3), res.Group.ID)
}

func TestGroupRouter_SubscriptionOverLimitIsSkipped(t *testing.T) {
	f := newRouterFixture()
	limit := 10.0
	g2 := activeGroup(2, PlatformOpenAI)
	g2.SubscriptionType = SubscriptionTypeSubscription
	g2.DailyLimitUSD = &limit
	f.subs.subs[2] = &UserSubscription{GroupID: 2, DailyUsageUSD: 11}
	key := multiGroupKey(g2, activeGroup(3, PlatformOpenAI))

	res, err := f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, int64(3), res.Group.ID)

	// 订阅未超限时优先用它
	f.subs.subs[2].DailyUsageUSD = 1
	f.router = NewGroupRouter(f.router.deps) // 清空结果缓存
	res, err = f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, int64(2), res.Group.ID)
}

// 全部容量不足：返回第一个模型可用的分组，让原流程报它原本的错误。
func TestGroupRouter_AllCapacityExhaustedFallsBackToFirstServingGroup(t *testing.T) {
	f := newRouterFixture()
	f.scheduler.accounts[2] = []Account{rateLimitedAccount()}
	f.scheduler.accounts[3] = nil
	key := multiGroupKey(activeGroup(1, PlatformAnthropic), activeGroup(2, PlatformOpenAI), activeGroup(3, PlatformOpenAI))

	res, err := f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, GroupRouteCapacityFallback, res.Outcome)
	require.Equal(t, int64(2), res.Group.ID)
}

func TestGroupRouter_NoGroupSupportsModelReturns404(t *testing.T) {
	f := newRouterFixture()
	key := multiGroupKey(activeGroup(1, PlatformAnthropic), activeGroup(2, PlatformOpenAI))

	_, err := f.router.Select(context.Background(), key, []string{"gemini-2.5-pro"}, "")
	require.ErrorIs(t, err, ErrNoGroupSupportsModel)
}

// 多个候选模型值时，分组必须对全部候选都合格。
func TestGroupRouter_MultipleModelCandidatesMustAllPass(t *testing.T) {
	f := newRouterFixture()
	key := multiGroupKey(activeGroup(3, PlatformOpenAI), activeGroup(2, PlatformOpenAI))

	res, err := f.router.Select(context.Background(), key, []string{"gpt-5", "gpt-image-1"}, "")
	require.NoError(t, err)
	require.Equal(t, int64(2), res.Group.ID, "group 3 lacks gpt-image-1")
}

func TestGroupRouter_DisabledOrForbiddenGroupsAreNotCandidates(t *testing.T) {
	f := newRouterFixture()
	disabled := activeGroup(2, PlatformOpenAI)
	disabled.Status = StatusDisabled
	exclusive := activeGroup(3, PlatformOpenAI)
	exclusive.IsExclusive = true
	key := multiGroupKey(activeGroup(1, PlatformAnthropic), disabled, exclusive)

	require.Equal(t, []*Group{key.Groups[0]}, f.router.EligibleCandidates(context.Background(), key))
	_, err := f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.ErrorIs(t, err, ErrNoGroupSupportsModel)

	key.User.AllowedGroups = []int64{3}
	res, err := NewGroupRouter(f.router.deps).Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, int64(3), res.Group.ID)
}

func TestGroupRouter_CompositeResolvesTargetPlatform(t *testing.T) {
	f := newRouterFixture()
	f.openai.models[5] = map[string]bool{"gpt-5-upstream": true}
	f.scheduler.accounts[5] = []Account{healthyAccount()}
	f.router.deps.Composite = &routerCompositeStub{decisions: map[string]CompositeRouteDecision{
		"gpt-5": {Matched: true, TargetPlatform: PlatformOpenAI, UpstreamModel: "gpt-5-upstream"},
	}}
	key := multiGroupKey(activeGroup(1, PlatformAnthropic), activeGroup(5, PlatformComposite))

	res, err := f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, int64(5), res.Group.ID)

	// composite 解析不出来的模型视为该组不支持
	_, err = f.router.Select(context.Background(), key, []string{"unknown-model"}, "")
	require.ErrorIs(t, err, ErrNoGroupSupportsModel)
}

// 依赖查询失败时按可调度处理，不把请求错判成别的分组或 404。
func TestGroupRouter_SchedulerErrorTreatedAsSchedulable(t *testing.T) {
	f := newRouterFixture()
	f.scheduler.err = errors.New("redis down")
	f.scheduler.accounts[2] = nil
	key := multiGroupKey(activeGroup(2, PlatformOpenAI), activeGroup(3, PlatformOpenAI))

	res, err := f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, GroupRouteSelected, res.Outcome)
	require.Equal(t, int64(2), res.Group.ID)
}

func TestGroupRouter_SingleGroupOrNoModelUsesPrimaryWithoutLookups(t *testing.T) {
	f := newRouterFixture()
	single := &APIKey{ID: 1, UserID: 7, GroupID: int64Ptr(2), Group: activeGroup(2, PlatformOpenAI), User: &User{ID: 7}}

	res, err := f.router.Select(context.Background(), single, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, GroupRoutePrimary, res.Outcome)
	require.Equal(t, int64(2), res.Group.ID)

	multi := multiGroupKey(activeGroup(1, PlatformAnthropic), activeGroup(2, PlatformOpenAI))
	res, err = f.router.Select(context.Background(), multi, nil, "")
	require.NoError(t, err)
	require.Equal(t, GroupRoutePrimary, res.Outcome)
	require.Equal(t, int64(1), res.Group.ID)
	require.Len(t, res.Candidates, 2)

	require.Zero(t, f.anthropic.calls+f.openai.calls)
}

func TestGroupRouter_ResultIsCached(t *testing.T) {
	f := newRouterFixture()
	key := multiGroupKey(activeGroup(1, PlatformAnthropic), activeGroup(2, PlatformOpenAI))

	_, err := f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	calls := f.anthropic.calls + f.openai.calls

	res, err := f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, int64(2), res.Group.ID)
	require.Equal(t, calls, f.anthropic.calls+f.openai.calls, "second call must hit the cache")

	// 缓存命中的分组被停用后重新判断，不再返回它
	key.Groups[1].Status = StatusDisabled
	_, err = f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.ErrorIs(t, err, ErrNoGroupSupportsModel)
}

func TestGroupRouter_AutoModeOrdering(t *testing.T) {
	sub := Group{ID: 10, Status: StatusActive, Platform: PlatformOpenAI, SubscriptionType: SubscriptionTypeSubscription, RateMultiplier: 2}
	cheap := Group{ID: 11, Status: StatusActive, Platform: PlatformOpenAI, RateMultiplier: 0.5, SortOrder: 9}
	mid := Group{ID: 12, Status: StatusActive, Platform: PlatformOpenAI, RateMultiplier: 1, SortOrder: 2}
	midFirst := Group{ID: 13, Status: StatusActive, Platform: PlatformOpenAI, RateMultiplier: 1, SortOrder: 1}
	source := &routerGroupSourceStub{groups: []Group{mid, cheap, sub, midFirst}}
	router := NewGroupRouter(GroupRouterDeps{Groups: source})
	key := &APIKey{ID: 1, UserID: 7, GroupMode: APIKeyGroupModeAuto, User: &User{ID: 7}}

	got := router.Candidates(context.Background(), key)
	ids := make([]int64, len(got))
	for i, g := range got {
		ids[i] = g.ID
	}
	require.Equal(t, []int64{10, 11, 13, 12}, ids)

	router.Candidates(context.Background(), key)
	require.Equal(t, 1, source.calls, "auto candidates are cached per user")

	res, err := router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, int64(10), res.Group.ID)
}

type routerRateStub map[int64]float64

func (s routerRateStub) Resolve(_ context.Context, _ int64, groupID int64, def float64) float64 {
	if v, ok := s[groupID]; ok {
		return v
	}
	return def
}

// 用户专属倍率参与自动模式排序。
func TestGroupRouter_AutoModeUsesUserRate(t *testing.T) {
	a := Group{ID: 1, Status: StatusActive, RateMultiplier: 0.5}
	b := Group{ID: 2, Status: StatusActive, RateMultiplier: 1}
	router := NewGroupRouter(GroupRouterDeps{
		Groups: &routerGroupSourceStub{groups: []Group{a, b}},
		Rates:  routerRateStub{2: 0.1},
	})
	got := router.Candidates(context.Background(), &APIKey{UserID: 7, GroupMode: APIKeyGroupModeAuto})
	require.Equal(t, int64(2), got[0].ID)
}

type routerVideoGroupStub map[string]int64

func (s routerVideoGroupStub) ResolveVideoRequestGroup(_ context.Context, requestID string, _, _ int64) (int64, error) {
	return s[requestID], nil
}

func TestGroupRouter_VideoLookupReturnsCreateTimeGroup(t *testing.T) {
	g1, g2 := activeGroup(1, PlatformAnthropic), activeGroup(2, PlatformGrok)
	router := NewGroupRouter(GroupRouterDeps{VideoGroups: routerVideoGroupStub{"req-1": 2, "req-gone": 99}})
	key := multiGroupKey(g1, g2)

	res := router.SelectForVideoLookup(context.Background(), key, "req-1")
	require.Equal(t, GroupRouteSelected, res.Outcome)
	require.Equal(t, int64(2), res.Group.ID)

	// 没有记录（升级前的任务）或记录的分组已不在候选里：退回主路由
	for _, id := range []string{"req-unknown", "req-gone", ""} {
		res = router.SelectForVideoLookup(context.Background(), key, id)
		require.Equal(t, GroupRoutePrimary, res.Outcome, id)
		require.Equal(t, int64(1), res.Group.ID, id)
	}

	// 分组停用后不再路由过去
	g2.Status = StatusDisabled
	res = router.SelectForVideoLookup(context.Background(), key, "req-1")
	require.Equal(t, GroupRoutePrimary, res.Outcome)
}

func TestGrokVideoRequestGroupCacheKeyIsScopedAndDistinct(t *testing.T) {
	a := grokVideoRequestGroupCacheKey("req", 1, 2)
	require.NotEmpty(t, a)
	require.NotEqual(t, a, grokVideoRequestGroupCacheKey("req", 1, 3), "scoped by api key")
	require.NotEqual(t, a, grokVideoRequestGroupCacheKey("req", 9, 2), "scoped by user")
	require.NotEqual(t, a, "openai:"+GrokMediaVideoRequestSessionHash("req", 1, 2), "must not collide with the account binding key")
	require.Empty(t, grokVideoRequestGroupCacheKey(" ", 1, 2))
}

// 建自动模式密钥时用户还没有分组（group_id 为空）：之后有了权限，请求时直接按当前候选选组。
func TestGroupRouter_AutoModeKeyWithoutStoredGroupGetsPrimary(t *testing.T) {
	openai := Group{ID: 5, Status: StatusActive, Platform: PlatformOpenAI, RateMultiplier: 1}
	router := NewGroupRouter(GroupRouterDeps{Groups: &routerGroupSourceStub{groups: []Group{openai}}})
	key := &APIKey{ID: 1, UserID: 7, GroupMode: APIKeyGroupModeAuto, User: &User{ID: 7}}

	res, err := router.Select(context.Background(), key, nil, "")
	require.NoError(t, err)
	require.Equal(t, GroupRoutePrimary, res.Outcome)
	require.NotNil(t, res.Group)
	require.Equal(t, int64(5), res.Group.ID)

	routed := CloneAPIKeyForGroup(key, res.Group)
	require.NotNil(t, routed.GroupID, "routed key carries a group so RequireGroupAssignment lets it through")
	require.Equal(t, int64(5), *routed.GroupID)
}

// 没有订阅的订阅分组不能当容量兜底：否则别的分组的 503 会变成 403 SUBSCRIPTION_NOT_FOUND（审查发现的问题）。
func TestGroupRouter_CapacityFallbackSkipsMissingSubscription(t *testing.T) {
	f := newRouterFixture()
	sub := activeGroup(2, PlatformOpenAI)
	sub.SubscriptionType = SubscriptionTypeSubscription // 用户没有该订阅
	f.scheduler.accounts[3] = nil                       // 分组 3 暂时没有可调度账号
	key := multiGroupKey(sub, activeGroup(3, PlatformOpenAI))

	res, err := f.router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, GroupRouteCapacityFallback, res.Outcome)
	require.Equal(t, int64(3), res.Group.ID)

	// 有订阅但额度用完：可以兜底（返回它原本的 429）
	limit := 1.0
	sub.DailyLimitUSD = &limit
	f.subs.subs[2] = &UserSubscription{GroupID: 2, DailyUsageUSD: 5}
	res, err = NewGroupRouter(f.router.deps).Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, GroupRouteCapacityFallback, res.Outcome)
	require.Equal(t, int64(2), res.Group.ID)
}

// 线上问题：工作台（自动模式）密钥选了倍率最低的分组，但该分组没开生图，生图请求被 403。
// 生图请求必须只在开启了生图的分组里选。
func TestGroupRouter_ImageRequestSkipsGroupsWithoutImageGeneration(t *testing.T) {
	cheapNoImage := Group{ID: 44, Status: StatusActive, Platform: PlatformOpenAI, RateMultiplier: 0.08, AllowImageGeneration: false}
	imageGroup := Group{ID: 43, Status: StatusActive, Platform: PlatformOpenAI, RateMultiplier: 1, AllowImageGeneration: true}
	router := NewGroupRouter(GroupRouterDeps{Groups: &routerGroupSourceStub{groups: []Group{imageGroup, cheapNoImage}}})
	key := &APIKey{ID: 350, UserID: 13, GroupMode: APIKeyGroupModeAuto, User: &User{ID: 13}}

	// 普通请求仍按倍率选最便宜的
	res, err := router.Select(context.Background(), key, []string{"gpt-5"}, "")
	require.NoError(t, err)
	require.Equal(t, int64(44), res.Group.ID)

	// 生图请求（/v1/images/edits 的 multipart 里取不到模型时也一样）跳过没开生图的分组
	imageCtx := WithImageGenerationRoute(context.Background())
	for _, models := range [][]string{{"gpt-image-2"}, nil} {
		res, err = router.Select(imageCtx, key, models, CompositeRouteEndpointImages)
		require.NoError(t, err)
		require.Equal(t, int64(43), res.Group.ID, "models=%v", models)
	}

	// 一个开了生图的分组都没有：不改变选择，由 handler 返回原有的 403
	onlyNoImage := NewGroupRouter(GroupRouterDeps{Groups: &routerGroupSourceStub{groups: []Group{cheapNoImage}}})
	res, err = onlyNoImage.Select(imageCtx, key, []string{"gpt-image-2"}, CompositeRouteEndpointImages)
	require.NoError(t, err)
	require.Equal(t, int64(44), res.Group.ID)
}

// 手动多分组密钥：主路由没开生图时，生图请求换到列表里第一个开了生图的分组。
func TestGroupRouter_ImageRequestManualKeyUsesFirstImageGroup(t *testing.T) {
	f := newRouterFixture()
	f.openai.models[2]["gpt-image-1"] = true
	noImage := activeGroup(3, PlatformOpenAI)
	withImage := activeGroup(2, PlatformOpenAI)
	withImage.AllowImageGeneration = true
	key := multiGroupKey(noImage, withImage)

	res, err := f.router.Select(WithImageGenerationRoute(context.Background()), key, nil, CompositeRouteEndpointImages)
	require.NoError(t, err)
	require.Equal(t, int64(2), res.Group.ID)
}
