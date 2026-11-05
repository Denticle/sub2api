## Context

已核实的事实（只读分析影策 `dc035d5` 与本仓库）：

**登录**

- 影策的 Linux.do 登录是通用 OAuth 2.0 客户端：授权地址、Token 地址、userinfo 地址、Client ID / Secret、回调、scope、客户端鉴权方式（`client_secret_post` / `client_secret_basic`）与 userinfo 字段映射（支持点路径）全部后台可配（`auth/linuxdo.go`）。
- 固定使用 PKCE S256；Token 请求带 `code_verifier`，只读 `access_token`；userinfo 用 Bearer 拉取。
- 三个地址必须是 HTTPS，出站校验默认拒绝私网地址。
- 发起登录的入口是 `GET /api/auth/linuxdo/start?next=…&acceptedTerms=true`，会直接 302 到授权地址——**这就是「自动登录」的入口**。
- 首次登录且未绑定时新建用户，要求影策「开放注册」打开、且未开启「短信+邮箱双重验证」；身份按 `(provider, subject)` 唯一绑定。

**计费**

- `FeatureCredits` 关闭时，`taskBillingOrder` 与 `ReserveProxyBillingWithBody` 直接返回 `nil`，不扣积分，生成照常。
- 系统渠道密钥只在 `app/provider.go` `resolveProviderConfig` 中一处注入（`config.APIKey = channel.APIKey`）；调用方都在任务上下文里，任务记录带 `UserID`。

**品牌与外观（全部后台可配，无需改源码）**

- 外观设置 `AppearanceSetting`：`brandName`、`brandSlug`、浅色 / 深色 logo（上传）、`logoFrameEnabled`、登录页标题与描述、登录页背景视频与封面（上传）、`seoTitle` / `seoDescription` / `seoKeywords`、`footerCopyright`、ICP 备案、兑换购买链接。
- 配色：`skinThemes` 是完整的设计令牌（画布、表面、文本、边框、控件、主色、选中、图标、成功 / 警告 / 危险等约 40 项，浅色 / 深色两套），每项校验为 hex 颜色。可以新建一套 forapi 皮肤并设为当前皮肤。

**功能开关（后台可关）**

- `shortDrama`、`taskCenter`、`credits`、`customChannels`、`frontendModels`、`pluginCenter`、`systemPluginsVisibleToUsers`、`inspirationSourcesVisible`、`timelineTranscription`、`imageWorkbench`、`videoWorkbench`。

**第三方内容盘点**

| 位置 | 内容 | 处理 |
| --- | --- | --- |
| 前端默认外观 `use-appearance-store.ts` | 登录页背景视频和封面引用 B 站 CDN（`biliapi.net` / `hdslb.com`） | 外观设置上传 forapi 自己的素材；补丁把默认值改为空，防止配置被清空时回落到外链 |
| 欢迎页 `pages/welcome/` | 贡献者卡片、上游故事、赞助 | 补丁：`/welcome` 重定向到 `/create`，不渲染欢迎页 |
| 渠道请求头 `channel-headers-editor.tsx` | 默认 UA `Yingce/1.0 (+github.com/ddcat-ai/open-ai-canvas)` | 补丁：改为 `ForAPI-Studio/1.0` |
| 服务商目录 `providercatalog` | 远程拉 `ai.yingce.tv/v1/catalog.json` | 环境变量 `CANVAS_PROVIDER_CATALOG_DISABLED=1` |
| 在线更新 `hostupdate` | 查询 GitHub releases、下载产物 | 不部署 host-updater，不设 `CANVAS_UPDATER_TOKEN`；补丁隐藏 `settings/system-update` 菜单 |
| 第三方登录 | Linux.do 原文案、Gmail 注册 | Linux.do 按钮改文案为「用 forapi 账号登录」；其余登录方式隐藏 |
| 画布导入（**保留**） | TapNow 分享链接导入、LibTV 项目导入：把对方画布的图片 / 视频节点连同提示词和参数导入当前画布 | 属于视频创作功能，保留。TapNow 无需配置；LibTV 需要管理员在 `settings/libtv` 填一个站点级 LibTV Token，该菜单保留 |
| 灵感大厅 `lib/inspirations/catalog.ts` | 外部提示词仓库来源说明与链接 | 功能开关 `inspirationSourcesVisible=false` 隐藏来源；条目保留（属于内容，不是第三方服务） |
| 3D 预演默认模型 `previs-mannequin.ts` | 从 jsDelivr 加载 three.js 示例模型 | 补丁：模型文件打进前端产物，改为本地路径 |
| 内置技能 / Agent 文案 | 若干处出现「影策」 | 补丁：统一替换为品牌名读取或「视频工作台」 |
| 管理后台 | 支付插件、短信插件、兑换购买链接 | 后台不启用；积分关闭后钱包与支付入口自动隐藏 |
| 仓库自带预编译二进制 `backend/server` | 65 MB，来源不可审计 | 不使用，镜像一律从源码构建 |

## Goals / Non-Goals

**Goals**

- 主站顶部菜单点「视频工作台」，新标签页打开，**一次跳转内**完成登录并落在创作页，用户无感。
- 外观与主站一致：logo、品牌名、主色、字体气质；用户看不到影策品牌和作者信息。
- 每次模型调用用该用户自己的 forapi 密钥；费用、限速、用量记录与 API 调用一致。
- 上游发布新版本后，按固定流程重建即可；补丁锚点变化时构建失败并指出位置。
- 不影响 new-api、sub2api 现有服务与数据。

**Non-Goals**

- 不内嵌到 forapi 主站（不用 iframe，不共享主站外壳）。
- 不做通用 OAuth 平台（无动态注册、无 refresh token、无第三方应用管理）。
- 不改影策的短剧、画布、Agent、时间线业务逻辑。
- 不追求像素级与主站一致：影策用 React + antd，主站是 Vue；统一到品牌、配色和字体层。

## Decisions

### D1 入口：顶部菜单直达登录入口

- 主站 `AppHeader` 顶栏和首页 `lp-nav` 各加一项「视频工作台」，`target="_blank"`，链接是：

  `https://<子域>/api/auth/linuxdo/start?next=/create&acceptedTerms=true`

- 影策收到后直接 302 到 forapi 授权页；forapi 判断用户已登录、且已授权过 → 立即带 `code` 回跳；影策完成登录后进入 `/create`。整个过程是几次 302，用户看到的就是新标签页打开了工作台。
- 首次使用时 forapi 显示一次授权确认（应用名、会创建专用密钥、费用记在余额），确认后记住。
- 用户在影策里已有有效会话时，再走一次这个入口也无害：会刷新一次身份并回到 `/create`。
- `acceptedTerms=true` 的含义是用户已在 forapi 同意过服务条款；影策的注册协议页由补丁隐藏。
- 菜单项只在 forapi 后台开关打开、且非 simple / 后端模式时显示；未登录点击时 forapi 授权页会先带去登录。

### D2 授权服务器协议

| 端点 | 方法 | 说明 |
| --- | --- | --- |
| `/oauth/authorize` | GET（前端页面） | 校验 `client_id`、`redirect_uri` 精确匹配、`response_type=code`、`code_challenge_method=S256`。未登录先登录再回来；已授权过直接签发；否则显示确认页 |
| `/api/v1/oauth/authorize/approve` | POST（JWT） | 生成授权码，返回带 `code` 与原样 `state` 的回调地址 |
| `/oauth/token` | POST（表单） | 校验客户端密钥（basic / post）、授权码未用未过期、`redirect_uri` 一致、`code_verifier` 的 S256 等于 `code_challenge`；返回 `{access_token, token_type: "Bearer", expires_in}` |
| `/oauth/userinfo` | GET（Bearer） | `{sub, username, name, email, picture}`；`sub` 为 forapi 用户 ID 的字符串，永不变化 |
| `/oauth/studio-key` | GET（客户端密钥 + 访问令牌） | 返回该用户的视频工作台专用密钥，见 D4 |

- 授权码 32 字节随机，只存 SHA-256，60 秒有效，单次使用；重放时作废其签发的令牌。
- 访问令牌 32 字节随机，只存哈希，10 分钟有效，只用于登录瞬间。
- 唯一客户端配置在 forapi 设置里，Client Secret 用 SecretEncryptor 加密。
- `/oauth/token`、`/oauth/studio-key` 按客户端 + IP 限流。

### D3 计费：关闭影策积分，按用户密钥直扣 forapi

- 影策 `credits=false`。真正扣费发生在 forapi 网关，按密钥归属用户；现有费率、分组、限速、用量记录全部生效。
- 代价：影策「创作历史」里积分消耗一列为空；用户在 forapi「使用记录」查看花费。可以在后续补丁里把这一列改成「去 forapi 查看」链接。

### D4 每用户密钥

- 首次授权时为用户创建或复用一把 API Key：名称「视频工作台」，用途标记 `studio_bridge`，额度取管理员配置。
- 密钥为**自动分组模式**（`group_mode=auto`，见 `add-api-key-multi-group-routing` D3/D8）：不绑定固定分组，
  每个请求按用户当前权限选组（订阅分组优先，其次实际倍率最低）。管理员不再配置「专用密钥分组」，
  `Config.KeyGroupID` 字段保留但不参与建密钥与启用校验；已有的工作台密钥由迁移 250 改为自动模式。
- 「API 密钥」页显示「视频工作台专用」标记，可停用、可重置。
- 影策登录回调拉完 userinfo 后，用同一访问令牌 + 客户端密钥调 `/oauth/studio-key`，取到明文后用影策自己的 `EncryptSecret` 加密存入新表 `user_bridge_credentials`。
- `resolveProviderConfig` 增加 `userID`；对标记为桥接的系统渠道，用该用户的密钥替换 `channel.APIKey`。取不到即失败，**不回落**到任何共享密钥。
- forapi 返回 401（密钥被停用或重置）时，影策作废该用户会话；用户下次从顶部菜单进入即刷新密钥。

### D5 品牌与外观：配置优先，补丁兜底

- **配置**（`video-studio-integration/provision/` 下的脚本，用管理员会话调用影策后台 API，可重复执行）：
  - 品牌名「ForAPI 视频工作台」，`brandSlug`、SEO、页脚版权取 forapi。
  - 上传 forapi 的浅色 / 深色 logo（与主站同一份 `logo-concept-a-256.png` 及深色版）。
  - 新建一套 forapi 皮肤：主色 `#16a34a` / 深色 `#22c55e`，表面与文本取主站 `tokens-forapi.css` 的 `n-*` / `d-*` 阶梯；设为当前皮肤并禁止用户切换。
  - 登录页背景：上传本地素材或留空。
- **补丁**只处理配置覆盖不到的：字体（Inter，与主站一致）、默认外链素材、欢迎页、硬编码的品牌字样、不需要的菜单。
- 配置脚本也负责关闭功能开关（积分、插件中心、系统插件可见、灵感来源可见）和 LibTV / TapNow。

### D6 上游同步

- `video-studio-integration/` 结构同 `canvas-integration/`：
  - `upstream.lock`：记录当前对齐的上游提交。
  - `patches/*.py`：每个补丁只做一件事，所有替换都校验锚点；锚点找不到即构建失败并打印文件与锚点。
  - `overlay/`：整文件覆盖的新增文件（如取密钥客户端、本地化的 3D 模型）。
  - `provision/`：后台配置脚本，幂等。
  - `build.sh`：`git fetch` → checkout 到 `upstream.lock` → 还原工作区 → 依次打补丁 → 构建 backend / web / agent 镜像 → 打 tag。
  - `UPGRADE.md`：升级步骤：改 `upstream.lock` → `build.sh` → 修失败的补丁 → 跑数据库迁移（上游自带 migrate 服务）→ 替换容器 → `provision` 再执行一次。
- 补丁面控制在约 10 处，集中在：登录页与菜单、`resolveProviderConfig`、登录回调、欢迎页路由、默认外观、少量品牌字符串。
- 不采用 fork 后长期维护：fork 会让每次上游更新都变成合并冲突。

### D7 网络与部署

- 子域独立站点，`enterprise-new-api-nginx` 新增 server 块反代到 `127.0.0.1` 上的影策 web 容器；证书走现有 certbot webroot 与续期定时器。
- 影策调用 forapi 的 OAuth 与模型接口走公网 `https://forapi.cloud`（OAuth 地址必须是 HTTPS 且非私网）。生图已有心跳，视频是异步任务。
- 环境：`CANVAS_REGISTRATION_ENABLED=true`（OAuth 首登建号需要；密码注册由补丁拒绝）、`CANVAS_CORS_ORIGINS` 与 `CANVAS_PUBLIC_BASE_URL` 指向子域、`CANVAS_ALLOW_PRIVATE_UPSTREAMS=false`、`CANVAS_PROVIDER_CATALOG_DISABLED=1`、不设 `CANVAS_UPDATER_TOKEN`。
- 容器：上游 `docker-compose.server.yml` 去掉 host-updater；端口只绑 `127.0.0.1`；设内存上限。

### D8 模型目录：forapi 为准，工作台自动同步

管理员只在 forapi 维护一份模型：各分组关联的渠道里的模型和计费方式（token / 按张 / 按秒）
就是工作台能用的模型。工作台不再手工建模型。

- forapi 新增 `GET /oauth/studio-models`（客户端 Basic 鉴权，与 `/oauth/studio-key` 同一组凭据），
  返回**全部 active 分组**在「模型广场」口径下的模型并集：`{id, platform, capability}`。
  工作台的模型目录是全站一份、不区分用户，而专用密钥是自动模式，所以目录列出任何用户可能用到的模型；
  用户选了自己没有权限的模型时，forapi 返回 404 `model_not_found`，工作台原样展示。
  capability 优先取渠道计费方式（token→text、image→image、video→video），没有定价时按模型名兜底。
- 工作台后端每 10 分钟（及启动时）拉取一次，按 capability 选协议写入「ForAPI」系统渠道：
  text→chat-completion（forapi 组合分组会把 Claude 等转到对应平台），image→openai-image，
  video：grok→xai-video、其余（Seedance）→volcengine-ark-video。
- 写入走工作台自己的 `SaveAdminChannelModel`，能力参数用上游默认值，价格档为 0（积分已关，
  真实扣费在 forapi 按用户密钥完成）。
- 只管理自己同步出来的模型（记录在 system_settings `forapi_model_sync`）：上游新增→创建并启用，
  上游消失→停用（不删除，保留历史任务引用）。管理员在工作台手工加的模型不受影响；
  管理员手工改过的同步模型（显示名、能力参数）在后续同步中保留，只补齐启用状态。

## Risks / Trade-offs

| 风险 | 缓解 |
| --- | --- |
| 影策持有用户的真实密钥 | 专用密钥、只能使用用户本人有权限的分组、独立额度、用户可停用；影策库内加密；授权页说明 |
| 全站模型目录比单个用户能用的多 | 调用时 forapi 明确返回 404 `model_not_found`；后续可加按用户过滤 |
| 上游改动认证、`resolveProviderConfig` 或外观结构 | 锚点校验失败即构建失败；`UPGRADE.md` 说明每处补丁的意图，便于重写 |
| 上游新增第三方内容 | `build.sh` 末尾跑一次外部域名扫描（与本次盘点同法），出现新域名时构建报警 |
| 外观与主站不完全一致 | 统一到 logo、品牌名、配色、字体；组件结构保持 antd，不追求像素一致 |
| 用户在影策直接注册密码账号 | 前端隐藏，后端拒绝（管理员除外） |
| 公网链路慢 | 生图心跳已开；必要时把 forapi 解析为内网地址并精确放行 |

## Migration Plan

1. forapi：上线 OAuth 提供方与专用密钥，后台开关默认关闭，顶部菜单不显示。
2. 部署视频工作台，运行 `provision`，管理员内部验证登录、品牌和一次真实生成。
3. 打开 forapi 开关，顶部菜单出现。

**回滚**：关闭 forapi 开关即隐藏菜单并拒绝授权；视频工作台栈可整体停止，不影响主站。
