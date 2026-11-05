## 1. 数据与仓储

- [x] 1.1 迁移 `250_api_key_multi_group.sql`：`api_key_groups` 表、`api_keys.group_mode` 列、回填 position 0、工作台密钥改 `auto`；迁移测试覆盖重复执行
- [x] 1.2 `service.APIKey` 增加 `GroupIDs []int64`、`Groups []*Group`、`GroupMode`；常量 `APIKeyGroupModeManual/Auto`
- [x] 1.3 `api_key_groups` 读写放在 `apiKeyRepository` 内（`api_key_group_repo.go`，原生 SQL）：`Create`/`Update` 与密钥行同事务重写；`ListKeysByGroupID` 合并两种绑定
- [x] 1.4 `apiKeyRepository.GetByKeyForAuth` 等读取路径加载有序分组（手动模式）；`ListKeysByGroupID` 合并 `api_key_groups` 结果
- [x] 1.5 删除分组时重排 position 并更新 `group_id`；集成测试 `TestMultiGroup_*`（无 Docker 时自动跳过，已在本地 Postgres 上跑通）

## 2. 认证缓存

- [x] 2.1 快照加 `GroupMode`、`Groups`；版本号 +1；`snapshotToAPIKey` 还原
- [x] 2.2 失效：分组变更按 `api_key_groups` 反查；订阅变化、专属分组变化按用户失效（检查现有调用点是否覆盖）
- [x] 2.3 单测：快照往返、版本升级、非主路由分组停用后失效

## 3. 选组服务

- [x] 3.1 `GroupRouter.Candidates(ctx, apiKey)`：手动按 position；自动按 D3 排序（订阅优先、实际倍率升序、sort_order、id）
- [x] 3.2 `GroupRouter.Select(ctx, apiKey, models []string)`：D3 五项资格判断，composite 先 Resolve；兜底错误语义
- [x] 3.3 可调度判断复用调度快照；查询失败视为可调度
- [x] 3.4 `(apiKeyID, model)` 5 秒进程内结果缓存
- [x] 3.5 单测：跳过不支持模型、限流跳组、订阅超额跳组、全部容量不足返回首组原错误、全部不支持返回 404、多候选模型值、自动模式排序

## 4. 网关接入

- [x] 4.1 `APIKeyAuth`：多分组时取模型（与白名单中间件同一解析，回填请求体），选组后构造 APIKey 副本，再做订阅/额度检查与上下文设置
- [x] 4.2 `APIKeyAuthGoogle` 同样接入
- [x] 4.3 context 存全部候选分组；单分组快速路径不读请求体（加测试断言）
- [x] 4.4 模型列表并集：`Gateway.Models`、`/models/:model`、Gemini、Antigravity、Codex 模型列表
- [x] 4.5 Grok 视频创建记录 requestID→分组；状态/内容查询取回；无记录退回主路由
- [x] 4.6 路由级测试：多分组密钥跨平台请求落到正确 handler（Claude→Anthropic、gpt→OpenAI、gpt-image→图片、grok 视频）

## 5. 用户端接口

- [x] 5.1 `CreateAPIKeyRequest`/`UpdateAPIKeyRequest` 加 `group_ids`、`group_mode`；校验去重、1–20、`CanBindGroup`、用户端拒绝 `auto`
- [x] 5.2 DTO 返回 `group_ids`、`groups`（含平台、倍率）、`group_mode`
- [x] 5.3 handler 单测：只传 `group_id` 兼容、传 `group_ids` 写入顺序、无权分组 403、超限 400

## 6. 前端

- [x] 6.1 `api/keys.ts` 类型与请求字段
- [x] 6.2 `KeysView.vue` 创建/编辑弹窗：多分组列表、主路由/路由 N、上移下移删除、添加分组下拉（只列有权分组）、倍率显示
- [x] 6.3 列表分组列：按顺序显示多个分组徽标，超出折叠
- [x] 6.4 自动模式只读展示与说明文案；i18n 中英文
- [x] 6.5 组件测试；`vue-tsc` 与 vitest 全量通过

## 7. 工作台

- [x] 7.1 `ensureStudioKey` 创建 `auto` 密钥；`validateConfig` 不再要求 `KeyGroupID`；后台隐藏该配置项
- [x] 7.2 `StudioModels` 改为全部 active 分组模型并集；测试更新
- [x] 7.3 同步更新 `add-yingce-studio-bridge` 的 design.md D4/D8 描述

## 8. 验证与交付

- [x] 8.1 `go build ./...`、`go vet`、相关包 `go test` 全通过；前端类型检查与测试全通过
- [ ] 8.2 本地起服务，用一把多分组密钥分别调 Claude、GPT、图片模型，确认用量日志分组与倍率正确
- [x] 8.3 独立 code review（另一个上下文）
- [ ] 8.4 部署：备份生产库 → 迁移 → 切镜像 → 线上用工作台密钥验证跨平台调用；记录回滚标签
