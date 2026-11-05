## ADDED Requirements

### Requirement: 密钥可以绑定多个有序分组
系统 SHALL 允许一把 API Key 绑定 1 到 20 个有序分组。第一个分组 MUST 与 `api_keys.group_id` 保持一致。只绑一个分组的密钥 MUST 与变更前行为完全一致。

#### Scenario: 旧密钥升级后行为不变
- **WHEN** 迁移前只有 `group_id` 的密钥发起请求
- **THEN** 系统 MUST 使用该分组处理请求，响应、计费、日志与变更前一致
- **THEN** 系统 MUST NOT 为该请求额外读取请求体

#### Scenario: 删除主路由分组
- **WHEN** 管理员删除某密钥的第一个分组
- **THEN** 系统 MUST 删除该关联，并把原第二个分组设为新的 `group_id`

### Requirement: 网关在选号前按顺序选出一个分组
对绑定多个分组的密钥，系统 SHALL 在账号选择、订阅读取和计费检查之前，按候选顺序选出第一个合格分组，整条请求 MUST 只使用这一个分组。分组合格的条件：分组可用且用户有权使用；分组白名单允许请求模型；分组内至少一个账号配置了该模型；订阅分组有有效订阅且未超日/周/月额度；分组当前至少有一个可调度账号。

#### Scenario: 跳过不支持该模型的分组
- **WHEN** 密钥分组顺序为 [Gemini 组, Claude 组]，请求模型为 `claude-sonnet-5-5`
- **THEN** 系统 MUST 选中 Claude 组
- **THEN** 用量日志的 `group_id` MUST 为 Claude 组，扣费倍率 MUST 为 Claude 组倍率

#### Scenario: 当前组无可调度账号时用下一组
- **WHEN** 第一个分组支持该模型但所有账号都处于限流
- **THEN** 系统 MUST 选中下一个支持该模型且有可调度账号的分组

#### Scenario: 订阅额度用完跳到下一组
- **WHEN** 第一个分组是订阅分组且当日额度已用完，第二个是余额分组
- **THEN** 系统 MUST 选中第二个分组并按余额计费

#### Scenario: 只有额度或容量不足
- **WHEN** 所有支持该模型的分组都因订阅额度或无可调度账号不合格
- **THEN** 系统 MUST 用第一个支持该模型的分组走原流程，返回该分组原本的错误（429 额度或 503 无可用账号）

#### Scenario: 没有任何分组支持该模型
- **WHEN** 所有分组都不支持请求模型
- **THEN** 系统 MUST 返回 404，错误类型为 `model_not_found`

#### Scenario: 请求体里出现多个候选模型值
- **WHEN** 请求体被不同解析器解析出多个模型值
- **THEN** 只有对全部候选值都合格的分组才能被选中

#### Scenario: 上游失败不换组
- **WHEN** 分组已被选中、请求已发往上游，且组内账号全部失败
- **THEN** 系统 MUST 返回该分组的错误，MUST NOT 换到下一个分组重试

### Requirement: 自动模式按用户权限实时计算分组顺序
`group_mode=auto` 的密钥 SHALL 在每次请求时以用户当前可用的全部分组为候选。顺序 MUST 为：有效订阅的订阅分组在前；其余按用户实际倍率（存在专属倍率时用专属倍率）从低到高；倍率相同按分组排序值，再按分组 ID。

#### Scenario: 用户新开订阅后立即生效
- **WHEN** 用户新购了订阅分组 S，且 S 支持请求模型
- **THEN** 该用户的自动模式密钥在认证缓存失效后的下一次请求 MUST 优先选中 S

#### Scenario: 专属分组被收回
- **WHEN** 管理员收回用户的专属分组 E
- **THEN** 自动模式密钥 MUST NOT 再选中 E

#### Scenario: 用户端不能创建自动模式密钥
- **WHEN** 用户通过 `POST /api/v1/keys` 提交 `group_mode=auto`
- **THEN** 系统 MUST 返回 400

### Requirement: 不带模型的请求有确定的分组
系统 SHALL 对取不到模型的请求使用确定规则：模型列表接口合并全部可用候选分组的可见模型；Grok 视频状态与内容查询使用创建该任务时选中的分组；其余请求使用主路由。

#### Scenario: 模型列表取并集
- **WHEN** 多分组密钥请求 `GET /v1/models`
- **THEN** 响应 MUST 包含各候选分组（已应用各自白名单）可见模型的去重并集

#### Scenario: 视频状态回到创建时的分组
- **WHEN** 视频任务创建时选中了分组 B，随后查询该任务状态
- **THEN** 系统 MUST 用分组 B 查询任务所属账号

#### Scenario: 升级前创建的视频任务
- **WHEN** 查询的视频任务没有分组记录
- **THEN** 系统 MUST 使用主路由分组

### Requirement: 分组相关变更会失效认证缓存
系统 SHALL 在分组配置变更、分组删除、用户订阅变化或专属分组变化时，失效所有引用该分组（任一位置）或属于该用户的密钥认证缓存。

#### Scenario: 非主路由分组被停用
- **WHEN** 管理员停用某密钥第三个位置的分组
- **THEN** 该密钥的认证缓存 MUST 被失效，下一次请求 MUST NOT 选中该分组

### Requirement: 工作台专用密钥使用自动模式
工作台（OAuth 客户端）为用户创建的专用密钥 SHALL 使用 `group_mode=auto`。工作台模型目录 SHALL 返回全部 active 分组可见模型的并集。

#### Scenario: 新用户授权工作台
- **WHEN** 用户首次授权工作台
- **THEN** 系统 MUST 创建 `group_mode=auto` 的「工作台」密钥，且不要求管理员配置专用密钥分组

#### Scenario: 已有工作台密钥
- **WHEN** 迁移执行完成
- **THEN** 已有的工作台密钥 MUST 变为 `group_mode=auto`
