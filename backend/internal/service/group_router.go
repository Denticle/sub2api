package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/domain"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	gocache "github.com/patrickmn/go-cache"
)

// 多分组密钥的请求时选组（add-api-key-multi-group-routing design D3）。
//
// 候选：手动模式按 api_key_groups 顺序；自动模式取用户全部可用分组，按
// 订阅分组优先 → 实际倍率升序 → sort_order → id 排序。
//
// 对每个候选依次判断，第一个全部通过的被选中：
//  1. 分组 active，用户有权绑定（订阅分组留给第 4 条）
//  2. 分组模型白名单放行
//  3. 分组里有账号配置了该模型（composite 先解析目标平台）
//  4. 订阅分组：有有效订阅且日/周/月未超限
//  5. 调度快照里有当前对该模型可调度的账号
//
// 依赖查询失败时一律按"通过"处理，与现有诊断的保守语义一致：宁可选中后由
// 原流程报错，也不因为查询抖动把请求错判成 404。

const (
	groupRouteResultTTL    = 5 * time.Second
	groupRouteCandidateTTL = 5 * time.Second
)

// GroupRouteOutcome 说明选组结果的来源。
type GroupRouteOutcome int

const (
	// GroupRouteSelected 候选通过了全部判断。
	GroupRouteSelected GroupRouteOutcome = iota
	// GroupRouteCapacityFallback 没有候选全部通过；返回第一个模型可用但卡在
	// 订阅额度或当前容量的分组，交给原流程返回它原本的错误（429 / 503）。
	GroupRouteCapacityFallback
	// GroupRoutePrimary 不需要选组（单分组、请求没有模型）或没有任何候选可用，
	// 退回主路由，由原流程按主路由返回错误（分组停用、无权等）。
	GroupRoutePrimary
)

// GroupRouteResult 是一次选组的结果。
type GroupRouteResult struct {
	// Group 是选中的分组；Outcome 为 GroupRoutePrimary 且密钥未绑分组时为 nil。
	Group   *Group
	Outcome GroupRouteOutcome
	// Candidates 是通过第 1 条（可用且有权）的候选，给模型列表等不带模型的请求用。
	Candidates []*Group
}

// ErrNoGroupSupportsModel 所有候选分组都不支持请求的模型。
var ErrNoGroupSupportsModel = infraerrors.NotFound("model_not_found", "no group bound to this api key supports the requested model")

func newNoGroupSupportsModelError(model string) error {
	return ErrNoGroupSupportsModel.WithMetadata(map[string]string{"model": model})
}

// GroupRouterGroupSource 提供自动模式的候选分组。
type GroupRouterGroupSource interface {
	GetAvailableGroups(ctx context.Context, userID int64) ([]Group, error)
}

// GroupRouterPlatformResolver 把 composite 分组按模型解析为具体平台。
type GroupRouterPlatformResolver interface {
	Resolve(ctx context.Context, groupID int64, model, endpoint string) (CompositeRouteDecision, error)
}

// GroupRouterSubscriptionChecker 判断订阅分组当前能否使用。
type GroupRouterSubscriptionChecker interface {
	GetActiveSubscription(ctx context.Context, userID, groupID int64) (*UserSubscription, error)
	ValidateAndCheckLimits(sub *UserSubscription, group *Group) (bool, error)
}

// GroupRouterSchedulableSource 读取调度快照。
type GroupRouterSchedulableSource interface {
	ListSchedulableAccounts(ctx context.Context, groupID *int64, platform string, hasForcePlatform bool) ([]Account, bool, error)
}

// GroupRouterRateSource 返回用户在分组的实际倍率（含专属倍率）。
type GroupRouterRateSource interface {
	Resolve(ctx context.Context, userID, groupID int64, groupDefaultMultiplier float64) float64
}

// GroupRouterVideoGroupSource 查视频任务创建时使用的分组（没有记录返回 0）。
type GroupRouterVideoGroupSource interface {
	ResolveVideoRequestGroup(ctx context.Context, requestID string, userID, apiKeyID int64) (int64, error)
}

// GroupRouterDeps 汇总选组需要的依赖；任意字段为 nil 时对应判断视为通过。
type GroupRouterDeps struct {
	Groups        GroupRouterGroupSource
	Composite     GroupRouterPlatformResolver
	Anthropic     ModelAvailabilityDiagnoser // 非 OpenAI 网关族平台
	OpenAI        ModelAvailabilityDiagnoser // OpenAI 网关族平台
	Subscriptions GroupRouterSubscriptionChecker
	Scheduler     GroupRouterSchedulableSource
	Rates         GroupRouterRateSource
	VideoGroups   GroupRouterVideoGroupSource
	Now           func() time.Time
}

// GroupRouter 为多分组密钥选出本次请求使用的分组。
type GroupRouter struct {
	deps       GroupRouterDeps
	results    *gocache.Cache
	candidates *gocache.Cache
}

func NewGroupRouter(deps GroupRouterDeps) *GroupRouter {
	return &GroupRouter{
		deps:       deps,
		results:    gocache.New(groupRouteResultTTL, time.Minute),
		candidates: gocache.New(groupRouteCandidateTTL, time.Minute),
	}
}

// Candidates 返回密钥的候选分组（已按优先级排序，未做资格过滤）。
func (r *GroupRouter) Candidates(ctx context.Context, apiKey *APIKey) []*Group {
	if apiKey == nil {
		return nil
	}
	if apiKey.IsAutoGroupMode() {
		return r.autoCandidates(ctx, apiKey.UserID)
	}
	if len(apiKey.Groups) > 1 {
		return apiKey.Groups
	}
	if apiKey.Group != nil {
		return []*Group{apiKey.Group}
	}
	return nil
}

// EligibleCandidates 返回通过第 1 条（分组可用、用户有权）的候选，保持优先级顺序。
func (r *GroupRouter) EligibleCandidates(ctx context.Context, apiKey *APIKey) []*Group {
	all := r.Candidates(ctx, apiKey)
	out := make([]*Group, 0, len(all))
	for _, g := range all {
		if groupRouteUsable(apiKey, g) {
			out = append(out, g)
		}
	}
	return out
}

func groupRouteUsable(apiKey *APIKey, g *Group) bool {
	if g == nil || !g.IsActive() {
		return false
	}
	// 订阅分组的权限就是"有有效订阅"，由第 4 条判断，与 validateAPIKeyGroupAllowed 一致。
	if g.IsSubscriptionType() || apiKey == nil || apiKey.User == nil {
		return true
	}
	return apiKey.User.CanBindGroup(g.ID, g.IsExclusive)
}

func (r *GroupRouter) autoCandidates(ctx context.Context, userID int64) []*Group {
	key := strconv.FormatInt(userID, 10)
	if cached, ok := r.candidates.Get(key); ok {
		if groups, ok := cached.([]*Group); ok {
			return groups
		}
	}
	if r.deps.Groups == nil {
		return nil
	}
	available, err := r.deps.Groups.GetAvailableGroups(ctx, userID)
	if err != nil {
		// 不缓存失败结果；本次请求没有候选，交给主路由处理。
		return nil
	}
	groups := orderAutoGroupCandidates(ctx, userID, available, r.deps.Rates)
	r.candidates.Set(key, groups, gocache.DefaultExpiration)
	return groups
}

// orderAutoGroupCandidates 按自动模式的优先级排序候选分组（D3）。
// 建密钥时写 group_id 与请求时选组共用这一个排序，保证两边看到的「第一组」一致。
func orderAutoGroupCandidates(ctx context.Context, userID int64, available []Group, rateSource GroupRouterRateSource) []*Group {
	groups := make([]*Group, 0, len(available))
	rates := make(map[int64]float64, len(available))
	for i := range available {
		g := &available[i]
		groups = append(groups, g)
		rates[g.ID] = g.RateMultiplier
		if rateSource != nil {
			rates[g.ID] = rateSource.Resolve(ctx, userID, g.ID, g.RateMultiplier)
		}
	}
	sortAutoGroupCandidates(groups, rates)
	return groups
}

// sortAutoGroupCandidates 订阅分组优先（GetAvailableGroups 只返回有有效订阅的订阅分组），
// 其余按实际倍率升序，再按 sort_order、id。
func sortAutoGroupCandidates(groups []*Group, rates map[int64]float64) {
	sort.SliceStable(groups, func(i, j int) bool {
		a, b := groups[i], groups[j]
		if as, bs := a.IsSubscriptionType(), b.IsSubscriptionType(); as != bs {
			return as
		}
		if ra, rb := rates[a.ID], rates[b.ID]; ra != rb {
			return ra < rb
		}
		if a.SortOrder != b.SortOrder {
			return a.SortOrder < b.SortOrder
		}
		return a.ID < b.ID
	})
}

type imageGenerationRouteCtxKey struct{}

// WithImageGenerationRoute 标记本次请求会生成图片：选组时只考虑开启了生图的分组。
// 生图权限是分组级开关（Group.AllowImageGeneration），不随模型配置变化，必须在选组时就排除，
// 否则会选中一个能接模型、但在 handler 里被 403 拒绝的分组。
func WithImageGenerationRoute(ctx context.Context) context.Context {
	return context.WithValue(ctx, imageGenerationRouteCtxKey{}, true)
}

func isImageGenerationRoute(ctx context.Context) bool {
	v, _ := ctx.Value(imageGenerationRouteCtxKey{}).(bool)
	return v
}

type cachedGroupRoute struct {
	groupID  int64
	outcome  GroupRouteOutcome
	notFound bool
}

// Select 为本次请求选组。models 是从请求中解析出的模型候选值（可能多个，
// 要求分组对全部候选都合格）；endpoint 用于 composite 分组解析。
//
// 返回 ErrNoGroupSupportsModel 表示所有可用候选都不支持该模型。
func (r *GroupRouter) Select(ctx context.Context, apiKey *APIKey, models []string, endpoint string) (GroupRouteResult, error) {
	if apiKey == nil {
		return GroupRouteResult{Outcome: GroupRoutePrimary}, nil
	}
	eligible := r.EligibleCandidates(ctx, apiKey)
	imageRoute := isImageGenerationRoute(ctx)
	if imageRoute && apiKey.HasMultipleGroups() {
		// 生图请求只在开启了生图的分组里选；一个都没有时保留原候选，由 handler 返回原有的 403。
		if allowed := filterImageGenerationGroups(eligible); len(allowed) > 0 {
			eligible = allowed
		}
	}
	primary := apiKey.Group
	if apiKey.IsAutoGroupMode() && len(eligible) > 0 {
		primary = eligible[0]
	} else if imageRoute && apiKey.HasMultipleGroups() && !GroupAllowsImageGeneration(primary) && len(eligible) > 0 {
		primary = eligible[0]
	}
	models = normalizeGroupRouteModels(models)
	if !apiKey.HasMultipleGroups() || len(models) == 0 || len(eligible) == 0 {
		return GroupRouteResult{Group: primary, Outcome: GroupRoutePrimary, Candidates: eligible}, nil
	}

	routeEndpoint := endpoint
	if imageRoute {
		routeEndpoint += "|image"
	}
	cacheKey := groupRouteCacheKey(apiKey.ID, routeEndpoint, models)
	if cached, ok := r.results.Get(cacheKey); ok {
		if route, ok := cached.(cachedGroupRoute); ok {
			if route.notFound {
				return GroupRouteResult{Candidates: eligible}, newNoGroupSupportsModelError(models[0])
			}
			for _, g := range eligible {
				if g.ID == route.groupID {
					return GroupRouteResult{Group: g, Outcome: route.outcome, Candidates: eligible}, nil
				}
			}
			// 缓存里的分组已不在候选中（停用、权限变化），重新判断。
		}
	}

	var capacityFallback *Group
	for _, g := range eligible {
		targets, ok := r.groupServesModels(ctx, g, models, endpoint)
		if !ok {
			continue
		}
		subState := r.subscriptionState(ctx, apiKey.UserID, g)
		// 兜底只取用户确实能用、只是额度或容量不足的分组；没有订阅的订阅分组不能兜底，
		// 否则会把别的分组的 503 换成 403 SUBSCRIPTION_NOT_FOUND。
		if capacityFallback == nil && subState != subscriptionMissing {
			capacityFallback = g
		}
		if subState == subscriptionOK && r.hasSchedulableAccount(ctx, g, targets) {
			r.results.Set(cacheKey, cachedGroupRoute{groupID: g.ID, outcome: GroupRouteSelected}, gocache.DefaultExpiration)
			return GroupRouteResult{Group: g, Outcome: GroupRouteSelected, Candidates: eligible}, nil
		}
	}
	if capacityFallback != nil {
		r.results.Set(cacheKey, cachedGroupRoute{groupID: capacityFallback.ID, outcome: GroupRouteCapacityFallback}, gocache.DefaultExpiration)
		return GroupRouteResult{Group: capacityFallback, Outcome: GroupRouteCapacityFallback, Candidates: eligible}, nil
	}
	r.results.Set(cacheKey, cachedGroupRoute{notFound: true}, gocache.DefaultExpiration)
	return GroupRouteResult{Candidates: eligible}, newNoGroupSupportsModelError(models[0])
}

// SelectForVideoLookup 为视频任务的状态/内容查询选组：回到创建任务时的分组。
//
// 没有记录（升级前创建的任务、单分组创建）或记录的分组已不可用时退回主路由，
// 由下游按主路由的账号绑定查找，查不到时返回原有的 404。
func (r *GroupRouter) SelectForVideoLookup(ctx context.Context, apiKey *APIKey, requestID string) GroupRouteResult {
	if apiKey == nil {
		return GroupRouteResult{Outcome: GroupRoutePrimary}
	}
	eligible := r.EligibleCandidates(ctx, apiKey)
	primary := apiKey.Group
	if apiKey.IsAutoGroupMode() && len(eligible) > 0 {
		primary = eligible[0]
	}
	fallback := GroupRouteResult{Group: primary, Outcome: GroupRoutePrimary, Candidates: eligible}
	if r.deps.VideoGroups == nil || strings.TrimSpace(requestID) == "" {
		return fallback
	}
	groupID, err := r.deps.VideoGroups.ResolveVideoRequestGroup(ctx, requestID, apiKey.UserID, apiKey.ID)
	if err != nil || groupID <= 0 {
		return fallback
	}
	for _, g := range eligible {
		if g.ID == groupID {
			return GroupRouteResult{Group: g, Outcome: GroupRouteSelected, Candidates: eligible}
		}
	}
	return fallback
}

func filterImageGenerationGroups(groups []*Group) []*Group {
	out := make([]*Group, 0, len(groups))
	for _, g := range groups {
		if GroupAllowsImageGeneration(g) {
			out = append(out, g)
		}
	}
	return out
}

func normalizeGroupRouteModels(models []string) []string {
	out := make([]string, 0, len(models))
	seen := make(map[string]struct{}, len(models))
	for _, m := range models {
		m = strings.TrimSpace(m)
		if m == "" {
			continue
		}
		if _, dup := seen[m]; dup {
			continue
		}
		seen[m] = struct{}{}
		out = append(out, m)
	}
	return out
}

func groupRouteCacheKey(apiKeyID int64, endpoint string, models []string) string {
	return strconv.FormatInt(apiKeyID, 10) + "|" + endpoint + "|" + strings.Join(models, "\x00")
}

// groupRouteTarget 是某个模型在分组里实际落到的平台与上游模型名。
type groupRouteTarget struct {
	platform string
	model    string
}

// groupServesModels 判断第 2、3 条：白名单放行，且分组有账号配置了每个模型。
func (r *GroupRouter) groupServesModels(ctx context.Context, g *Group, models []string, endpoint string) ([]groupRouteTarget, bool) {
	targets := make([]groupRouteTarget, 0, len(models))
	for _, model := range models {
		if !g.ModelAllowlist.Allows(model) {
			return nil, false
		}
		target, ok := r.resolveTarget(ctx, g, model, endpoint)
		if !ok {
			return nil, false
		}
		if target.platform != "" && !r.modelConfigured(ctx, g, target) {
			return nil, false
		}
		targets = append(targets, target)
	}
	return targets, true
}

// resolveTarget 返回模型在分组里的目标平台；平台为空表示无法判断（按支持处理）。
func (r *GroupRouter) resolveTarget(ctx context.Context, g *Group, model, endpoint string) (groupRouteTarget, bool) {
	if g.Platform != PlatformComposite {
		return groupRouteTarget{platform: g.Platform, model: model}, true
	}
	if r.deps.Composite == nil {
		return groupRouteTarget{model: model}, true
	}
	decision, err := r.deps.Composite.Resolve(ctx, g.ID, model, endpoint)
	if err != nil {
		return groupRouteTarget{model: model}, true
	}
	if !decision.Matched {
		// 未解析（含 OAuth 歧义）的模型在该 composite 分组里无法可靠服务。
		return groupRouteTarget{}, false
	}
	upstream := strings.TrimSpace(decision.UpstreamModel)
	if upstream == "" {
		upstream = model
	}
	return groupRouteTarget{platform: decision.TargetPlatform, model: upstream}, true
}

func (r *GroupRouter) modelConfigured(ctx context.Context, g *Group, target groupRouteTarget) bool {
	diagnoser := r.deps.Anthropic
	if domain.UsesOpenAIGateway(target.platform) {
		diagnoser = r.deps.OpenAI
	}
	if diagnoser == nil {
		return true
	}
	groupID := g.ID
	return diagnoser.DiagnoseModelAvailabilityForPlatform(ctx, &groupID, target.model, target.platform).HasModelSupport
}

type subscriptionRouteState int

const (
	subscriptionOK      subscriptionRouteState = iota // 非订阅分组，或订阅有效且未超限
	subscriptionLimited                               // 有订阅但日/周/月额度用完（或状态异常）
	subscriptionMissing                               // 订阅分组但用户没有有效订阅
)

// subscriptionState 判断第 4 条。查询失败按可用处理（与其他依赖的保守语义一致）。
func (r *GroupRouter) subscriptionState(ctx context.Context, userID int64, g *Group) subscriptionRouteState {
	if !g.IsSubscriptionType() || r.deps.Subscriptions == nil {
		return subscriptionOK
	}
	sub, err := r.deps.Subscriptions.GetActiveSubscription(ctx, userID, g.ID)
	if err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			return subscriptionMissing
		}
		return subscriptionOK
	}
	if sub == nil {
		return subscriptionMissing
	}
	if _, err := r.deps.Subscriptions.ValidateAndCheckLimits(sub, g); err != nil {
		return subscriptionLimited
	}
	return subscriptionOK
}

// hasSchedulableAccount 判断第 5 条：每个目标都有当前对该模型可调度的账号。
// 模型是否配置已由第 3 条确认，这里只看限流、过载等运行时状态。
func (r *GroupRouter) hasSchedulableAccount(ctx context.Context, g *Group, targets []groupRouteTarget) bool {
	if r.deps.Scheduler == nil {
		return true
	}
	groupID := g.ID
	for _, target := range targets {
		if target.platform == "" {
			continue
		}
		accounts, _, err := r.deps.Scheduler.ListSchedulableAccounts(ctx, &groupID, target.platform, false)
		if err != nil {
			continue
		}
		found := false
		for i := range accounts {
			if accounts[i].IsSchedulableForModelWithContext(ctx, target.model) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// String 便于日志输出。
func (o GroupRouteOutcome) String() string {
	switch o {
	case GroupRouteSelected:
		return "selected"
	case GroupRouteCapacityFallback:
		return "capacity_fallback"
	case GroupRoutePrimary:
		return "primary"
	default:
		return fmt.Sprintf("outcome(%d)", int(o))
	}
}

var (
	_ GroupRouterGroupSource         = (*APIKeyService)(nil)
	_ GroupRouterPlatformResolver    = (*CompositeRouteResolver)(nil)
	_ ModelAvailabilityDiagnoser     = (*GatewayService)(nil)
	_ ModelAvailabilityDiagnoser     = (*OpenAIGatewayService)(nil)
	_ GroupRouterSubscriptionChecker = (*SubscriptionService)(nil)
	_ GroupRouterSchedulableSource   = (*SchedulerSnapshotService)(nil)
	_ GroupRouterRateSource          = (*userGroupRateResolver)(nil)
	_ GroupRouterVideoGroupSource    = (*OpenAIGatewayService)(nil)
)

// APIKeyGroupRoute 记录本次请求的选组结果，给模型列表、视频查询等需要看全部候选的路径用。
type APIKeyGroupRoute struct {
	// All 是通过第 1 条的候选分组（优先级顺序）。
	All []*Group
	// Selected 是本次请求使用的分组。
	Selected *Group
	Outcome  GroupRouteOutcome
}

type apiKeyGroupRouteCtxKey struct{}

// WithAPIKeyGroupRoute 把选组结果写入 context。
func WithAPIKeyGroupRoute(ctx context.Context, route APIKeyGroupRoute) context.Context {
	return context.WithValue(ctx, apiKeyGroupRouteCtxKey{}, route)
}

// APIKeyGroupRouteFromContext 读取选组结果；单分组密钥没有该值。
func APIKeyGroupRouteFromContext(ctx context.Context) (APIKeyGroupRoute, bool) {
	if ctx == nil {
		return APIKeyGroupRoute{}, false
	}
	route, ok := ctx.Value(apiKeyGroupRouteCtxKey{}).(APIKeyGroupRoute)
	return route, ok
}

// CloneAPIKeyForGroup 返回把 group 装作主路由的密钥副本，下游 handler 按它计费、调度。
//
// 认证快照里的 UserGroupRPMOverride 属于原主路由分组；换组后必须清掉，
// 由计费层按新分组回源查询，否则会套用错误分组的 RPM 上限。
func CloneAPIKeyForGroup(apiKey *APIKey, group *Group) *APIKey {
	if apiKey == nil || group == nil {
		return apiKey
	}
	if apiKey.GroupID != nil && *apiKey.GroupID == group.ID && apiKey.Group == group {
		return apiKey
	}
	cloned := *apiKey
	groupID := group.ID
	cloned.GroupID = &groupID
	cloned.Group = group
	if apiKey.User != nil && (apiKey.GroupID == nil || *apiKey.GroupID != group.ID) {
		user := *apiKey.User
		user.UserGroupRPMOverride = nil
		cloned.User = &user
	}
	return &cloned
}
