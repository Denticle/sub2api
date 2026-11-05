## Context

已核实的现状（本仓库 `theme` 分支）：

- `api_keys.group_id` 单列；`service.APIKey` 带 `GroupID *int64` 和预加载的 `Group *Group`（`internal/service/api_key.go`）。
- 认证在 `middleware.APIKeyAuth`（`internal/server/middleware/api_key_auth.go`）：取 key → 校验分组可用/有权 → 按 `apiKey.Group` 取订阅 → 余额/订阅额度检查 → `c.Set(ContextKeyAPIKey, apiKey)` + `setGroupContext(c, apiKey.Group)`。Google 变体（`api_key_auth_google.go`）结构相同。
- 鉴权之后的网关链：`groupModelAllowlist`（读请求体取模型）→ `compositeTarget`（composite 分组按模型解析目标平台）→ 各平台 handler。handler 全部从 context 取 `apiKey`，用 `apiKey.GroupID` 选号、`apiKey.Group` 计费。仓库里 `apiKey.Group` / `.GroupID` 在 handler 包有约 295 处引用。
- 已有「换组」先例：`cloneAPIKeyWithGroup`（`gateway_handler.go`）把 `FallbackGroupIDOnInvalidRequest` 换进来重试，下游照常工作。这证明「把选中的分组装进 APIKey 副本交给下游」是可行的。
- 认证缓存：L1 进程内 + L2 Redis，快照 `APIKeyAuthSnapshot` 里存分组的完整字段；`InvalidateAuthCacheByGroupID` 用 `ListKeysByGroupID`（按 `api_keys.group_id`）找要失效的 key。
- 判断「分组里有没有账号配置了这个模型」已有现成实现：`DiagnoseModelAvailabilityForPlatform(ctx, groupID, model, platform)`，只看持久配置、不看限流，失败时保守返回「有」。
- composite 分组按模型解析目标平台：`CompositeRouteResolver.Resolve`。
- Grok 视频状态/内容查询靠「requestID → 账号」绑定，缓存 key 含 groupID（`ResolveGrokMediaVideoRequestAccount`）。
- 工作台密钥由 `oauthprovider.ensureStudioKey` 用 `APIKeyService.Create` 创建，分组取 `Config.KeyGroupID`；模型目录 `StudioModels` 只看这一个分组。

## Goals / Non-Goals

**Goals**

- 一把密钥按顺序挂多个分组，可跨平台；请求自动落到第一个能服务该模型的分组。
- 单分组密钥的行为、性能、响应完全不变（快速路径不读请求体、不多查库）。
- 选中分组后，下游计费、订阅、白名单、调度、用量日志全部按选中分组，不逐个改 handler。
- 工作台密钥不用管理员维护：随用户权限自动包含全部可用分组。

**Non-Goals**

- 上游已发出请求后失败再换组重试（流式已写出无法撤销；各平台重试循环差异大）。组内已有的账号级 failover 不变。
- 同一请求按组拆分或负载均衡（永远只用一个组）。
- 管理员侧的分组管理、composite 路由规则不改。
- 不做「按成本自动选最便宜的组」，顺序由用户（手动）或固定规则（自动）决定。

## Decisions

### D1 数据模型

```sql
CREATE TABLE api_key_groups (
  api_key_id BIGINT NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
  group_id   BIGINT NOT NULL REFERENCES groups(id)   ON DELETE CASCADE,
  position   INT    NOT NULL,
  PRIMARY KEY (api_key_id, group_id),
  UNIQUE (api_key_id, position)
);
CREATE INDEX idx_api_key_groups_group ON api_key_groups (group_id);
ALTER TABLE api_keys ADD COLUMN group_mode VARCHAR(16) NOT NULL DEFAULT 'manual';
-- 回填：每把有 group_id 的密钥写一行 position 0
```

- `api_keys.group_id` 保留，始终等于 position 0（主路由）。所有只认 `group_id` 的旧代码（管理员筛选、统计、`ListKeysByGroupID`）继续成立。
- 手动模式最多 20 个分组（截图那家是 18），防止误配导致每个请求遍历过多分组。
- 自动模式不写 `api_key_groups`，顺序在请求时计算（D3），`group_id` 写当时计算出的第一组，仅用于展示和旧代码兼容。
- 用 database/sql 手写 repository（与 `oauthprovider` 一致），不改 ent schema，避免重新生成 ent 代码。

**备选**：把分组列表存成 `api_keys.group_ids JSONB`。否决：分组删除时无法级联，`InvalidateAuthCacheByGroupID` 也没法按组反查。

### D2 选组发生在鉴权中间件里

在 `APIKeyAuth`「校验分组可用/有权」那一步改为调用 `GroupRouter.Select(ctx, apiKey, request)`，得到选中的 `*Group` 后用它构造 APIKey 副本（同 `cloneAPIKeyWithGroup`），后续订阅读取、额度检查、`setGroupContext` 都用副本。下游中间件和 handler 拿到的就是选中分组，不需要感知多分组。

- **快速路径**：候选分组只有 1 个时，直接走现有逻辑，不读请求体。
- **取模型**：多于 1 个时用 `requestmodel.FromBodyCandidates`（与 `groupModelAllowlist` 同一套解析，读完回填请求体），路由参数和查询参数也按白名单中间件的规则取。取到多个候选模型值时，要求分组对**所有**候选都合格（与白名单中间件「任一不过即拒」一致，防止解析器分歧绕过）。
- **原分组信息保留**：context 里另存 `ctxkey.APIKeyGroupRoute{All []*Group, Selected int}`，给模型列表和视频查询用（D4）。

**备选**：在每个 handler 的选号循环里遍历分组。否决：要改 6 类 handler 的重试循环，而选号前判断（用户已选）在中间件一处即可完成。

### D3 候选分组与资格判断

候选列表：
- 手动：`api_key_groups` 按 position。
- 自动：用户可用的全部分组（`GetAvailableGroups` 同口径：公开分组 + 被授予的专属分组 + 有有效订阅的订阅分组），排序 = 有效订阅的订阅分组在前，其余按 `rate_multiplier` 升序，相同按 `sort_order`、再按 id。用户专属倍率（`user_group_rates`）存在时按专属倍率排。

对每个候选依次判断，第一个全部通过的被选中：

1. 分组 active、未删除，用户 `CanBindGroup`（与现有 `validateAPIKeyGroupAllowed` 同语义）。
2. 白名单：`group.ModelAllowlist.Allows(model)`（未开启视为通过）。
3. 配置支持：composite 分组先 `Resolve` 出目标平台，再用 `DiagnoseModelAvailabilityForPlatform(group.ID, model, platform)`；`HasModelSupport=false` 跳过。诊断失败时视为支持（与现有保守语义一致）。
4. 订阅分组：有有效订阅，且 `ValidateAndCheckLimits` 不报日/周/月超限。
5. 当前可调度：调度快照里该组存在至少一个对该模型可调度、未被限流的账号。查询失败时视为可调度。

全部候选都不通过时：
- 有分组通过 1–3 但卡在 4/5：按**第一个**这样的分组走现有流程，让它返回原本的错误（额度用完 429、无可用账号 503），客户端看到的错误语义与单分组一致。
- 都卡在 1–3：返回 404 `model_not_found`，提示「该密钥的分组都不支持模型 X」。
- 没有任何候选通过第 1 条（全部停用或无权）：退回主路由，由现有校验返回「分组已停用 / 无权」，不报 404，避免把权限问题说成模型不存在。

第 5 条只看账号运行时状态（限流、过载、临时不可调度、模型级限流）。账号是否配置了该模型已由第 3 条按分组确认，这里不再逐账号匹配映射，避免 Antigravity 混合调度等映射规则不一致时误判为「无容量」。

自动模式的候选列表按用户进程内缓存 5 秒（`GetAvailableGroups` 本身不带缓存，每次 3 次查库）。

判断 5 用的是调度快照，与真正选号之间存在时间差；被选中后若真正选号失败，按该组现有错误返回，不再换组（Non-Goal）。

**结果缓存**：`(apiKeyID, model)` → 选中分组 ID，进程内缓存 5 秒，减少高并发下重复判断；分组/订阅/账号变更不主动失效，5 秒内自然过期。

### D4 不带模型的请求

| 请求 | 规则 |
| --- | --- |
| `GET /v1/models`、`/models/:model`（含根路径别名） | 合并全部候选分组（只取通过 D3 第 1 条的）各自可见的模型（已应用各组白名单），按候选顺序去重；`/models/:model` 在任一组可见即返回。同平台沿用该平台的列表格式，跨平台用 composite 格式。开启固定账号 manifest 的 OpenAI 分组在合并时按账号映射列出（固定 manifest 是上游原始响应体，无法合并） |
| Antigravity 模型列表 | 任一候选分组的白名单放行即列出 |
| Gemini 模型列表 / 详情（`/v1beta/models`） | 这两个接口透传某个上游账号的模型目录，无法跨组合并：主路由不是 Gemini 时改用第一个 Gemini 候选分组，没有则保持主路由（返回原有的平台错误） |
| Codex 模型 manifest（`?client_version=`、`/backend-api/codex/models`） | 同理按单组生成：主路由不是 OpenAI/Composite 时改用第一个 OpenAI/Composite 候选分组 |
| Grok / Seedance 视频创建 | 选组照常；创建成功后额外写 requestID → 分组 ID（存在粘性会话缓存里，key 按 user、apiKey、requestID 派生，与账号绑定 key 不同；TTL 与现有账号绑定一致） |
| Grok 视频状态 / 内容查询、Seedance 任务查询 / 删除 | 按上面的记录取回分组（分组须仍通过 D3 第 1 条）；没有记录时退回主路由（兼容升级前创建的任务） |
| 异步图片任务查询 `/v1/images/tasks/*` | 现有逻辑只按密钥归属校验，不依赖分组，走主路由 |
| `/v1/usage`、`/v1/sub2api/billing` | 走主路由（这两个接口展示的是密钥层面的额度，主路由用于显示订阅信息） |
| 其他取不到模型的请求 | 走主路由 |

### D5 认证缓存

- 快照 `APIKeyAuthSnapshot` 增加 `GroupMode` 和 `Groups []APIKeyAuthGroupSnapshot`（手动模式的完整列表）；缓存版本号 +1。
- 自动模式不把候选列表写进认证快照（它依赖用户订阅和专属分组）。`GetAvailableGroups` 本身没有缓存，
  选组服务按用户在进程内缓存排好序的候选 5 秒（见 D3），权限变化最多 5 秒后生效。
- `InvalidateAuthCacheByGroupID` 改为同时查 `api_key_groups.group_id` 与 `api_keys.group_id`。
- 失效覆盖核查（tasks 2.2）：
  - 分组变更（编辑、删除、渠道 / 模型配置变化）走 `InvalidateAuthCacheByGroupID`，经上一条覆盖非主路由分组。
  - 专属分组授权变化（管理员改 `allowed_groups`、替换专属分组）走 `InvalidateAuthCacheByUserID`，
    按用户列出全部密钥，多分组密钥同样被失效。
  - 订阅不在认证快照里：鉴权中间件与选组服务都通过 `SubscriptionService.GetActiveSubscription` 实时读取，
    订阅的分配、撤销、续期已失效该缓存，认证缓存不需要随订阅失效。
- 删除分组：外键级联删掉 `api_key_groups` 行；若删掉的是主路由，把 position 重排，`group_id` 改为新的第一组（在删除分组的同一事务里做）。

### D6 计费与日志

不改计费代码。选中分组后 `apiKey.Group` 就是它，倍率、图片/视频独立倍率、订阅扣减、用量日志 `group_id` 自然按实际命中的分组。用量日志不加新字段，已有的 `group_id` 记录的就是这次请求实际用的分组。

### D7 用户端接口与校验

- `POST /api/v1/keys`、`PUT /api/v1/keys/:id` 新增可选 `group_ids []int64`、`group_mode`。传了 `group_ids` 时以它为准并把第一个写进 `group_id`；只传 `group_id` 时等价于 `group_ids=[group_id]`。
- 校验：去重、数量 1–20、每个分组用户都 `CanBindGroup`；自动模式不接受 `group_ids`。
- 用户端只能选手动模式；自动模式只给工作台密钥（`oauthprovider`）使用，避免用户密钥的分组在不知情时变化。
- 管理员「用户密钥」页暂只读展示多分组（不提供编辑），后续需要再加。

### D8 工作台

- `ensureStudioKey` 创建密钥时固定 `group_mode=auto`。后台「视频工作台接入」的「专用密钥分组」配置项隐藏，`Config.KeyGroupID` 字段保留但不再参与建密钥，启用校验不再要求它。
- 已有的工作台密钥（线上 2 把）通过迁移改为 `auto`。
- `StudioModels`：工作台的模型目录是全站一份（工作台每 10 分钟同步一次，不区分用户），所以目录取**全部 active 分组**模型的并集，即任何用户可能用到的模型。某个用户选了自己没有权限的模型时，forapi 按 D3 返回 404 `model_not_found`，工作台把这条错误原样展示。

## Risks / Trade-offs

| 风险 | 处理 |
| --- | --- |
| 多分组密钥每个请求多一次读请求体和 N 次资格判断 | 单分组快速路径不受影响；结果缓存 5 秒；最多 20 组 |
| 调度快照判断与实际选号之间状态变化 | 只影响边界情况，按选中组的原错误返回；不做二次换组 |
| 自动模式下用户权限变化导致同一模型在不同时间落到不同分组，价格不同 | 这是用户要求的语义（订阅优先、倍率低优先）；密钥页展示当前顺序与倍率 |
| 模型列表并集后，客户端看到的模型在某组被白名单禁止 | 并集只取各组「自身可见」的模型（已应用各组白名单） |
| 工作台全站模型目录比单个用户能用的多 | 调用时 forapi 明确报错；后续可加按用户过滤 |
| 认证缓存版本升级导致短时间回源增多 | 版本号变更只让旧缓存失效一次，已有 singleflight |

## Migration Plan

1. 部署前备份生产库（同 249 的做法）。
2. 迁移 `250_api_key_multi_group.sql`：建表、加列、回填 position 0、把 `oauth_client_api_keys` 里的密钥改为 `auto`。全部 `IF NOT EXISTS`，可重复执行。
3. 新代码在 `api_key_groups` 为空时退回 `group_id`，所以先上迁移、后上代码，或反过来，都不会中断。
4. 回滚：镜像切回上一个标签即可；新表和新列保留不影响旧代码（旧代码只读 `group_id`）。
