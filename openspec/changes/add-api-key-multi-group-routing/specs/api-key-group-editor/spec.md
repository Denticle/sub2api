## ADDED Requirements

### Requirement: 用户可以为密钥多选分组并排序
「API 密钥」创建与编辑弹窗 SHALL 允许用户添加多个分组、上移下移、删除，第一项标为「主路由」，其余标为「路由 2、路由 3…」，每项显示平台图标、分组名和用户实际倍率。只能添加用户有权使用的分组。

#### Scenario: 添加第二个分组
- **WHEN** 用户在编辑弹窗添加分组 B 并保存
- **THEN** 前端 MUST 提交 `group_ids=[原分组, B]`
- **THEN** 列表 MUST 显示两个分组及其顺序

#### Scenario: 调整顺序
- **WHEN** 用户把路由 2 上移并保存
- **THEN** 该分组 MUST 成为主路由，`group_id` MUST 随之更新

#### Scenario: 选择无权使用的分组
- **WHEN** 请求的 `group_ids` 含用户无权使用的分组
- **THEN** 后端 MUST 返回 403 `GROUP_NOT_ALLOWED`，不修改任何数据

#### Scenario: 超过上限
- **WHEN** 请求的 `group_ids` 超过 20 个或含重复
- **THEN** 后端 MUST 返回 400

### Requirement: 接口兼容只传单个分组的客户端
`POST/PUT /api/v1/keys` SHALL 继续接受只传 `group_id` 的请求，等价于 `group_ids=[group_id]`。响应 SHALL 新增 `group_ids`、`groups`、`group_mode` 字段，原有字段不变。

#### Scenario: 旧客户端改分组
- **WHEN** 客户端只提交 `group_id=C`
- **THEN** 密钥 MUST 变为只绑定分组 C

### Requirement: 自动模式密钥只读展示
自动模式密钥（工作台密钥）在列表与编辑弹窗中 SHALL 显示「自动（按权限）」和当前计算出的顺序，用户 MUST NOT 能编辑其分组。

#### Scenario: 打开工作台密钥编辑弹窗
- **WHEN** 用户编辑「工作台」密钥
- **THEN** 分组区域 MUST 只读，并说明顺序规则
