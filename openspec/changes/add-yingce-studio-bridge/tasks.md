## 1. forapi：OAuth 授权服务器

- [ ] 1.1 新增 `internal/oauthprovider/`：授权码表、访问令牌表（只存 SHA-256），迁移脚本
- [ ] 1.2 设置项：开关、Client ID、Client Secret（SecretEncryptor 加密）、回调地址（精确匹配）、专用密钥额度；管理后台新增分区（专用密钥为自动分组模式，不再配置分组，见 add-api-key-multi-group-routing D8）
- [ ] 1.3 `POST /api/v1/oauth/authorize/approve`（JWT）：参数校验、PKCE S256、60 秒单次授权码
- [ ] 1.4 前端 `/oauth/authorize`：未登录跳登录带回；首次显示确认页；已授权直接回跳
- [ ] 1.5 `POST /oauth/token`：basic / post 客户端鉴权、单次消费、重放作废、`code_verifier` 校验、10 分钟访问令牌
- [ ] 1.6 `GET /oauth/userinfo`：`sub`、`username`、`name`、`email`、`picture`
- [ ] 1.7 限流：`/oauth/token`、`/oauth/studio-key`
- [ ] 1.8 单测：PKCE 错误、redirect_uri 不一致、过期、重放、客户端密钥错误、开关关闭、用户被封禁

## 2. forapi：视频工作台专用密钥

- [ ] 2.1 API Key 增加用途标记 `studio_bridge`
- [ ] 2.2 授权成功时创建或复用专用密钥（自动分组模式，额度来自 1.2）
- [ ] 2.3 `GET /oauth/studio-key`：客户端密钥 + 有效访问令牌，返回明文
- [ ] 2.4 「API 密钥」页：「视频工作台专用」标记，可停用、可重置
- [ ] 2.5 单测：复用不重复创建；停用后返回明确错误；令牌过期拒绝

## 3. forapi：顶部菜单入口

- [ ] 3.1 `AppHeader` 顶栏新增「视频工作台」，新标签页打开 `https://<子域>/api/auth/linuxdo/start?next=/create&acceptedTerms=true`
- [ ] 3.2 首页 `lp-nav` 同样加一项
- [ ] 3.3 仅在 1.2 开关打开、且非 simple / 后端模式时显示；不在侧栏增加入口

## 4. 视频工作台：集成目录 `video-studio-integration/`

- [ ] 4.1 `upstream.lock`、`build.sh`（fetch → checkout → 还原 → 打补丁 → 构建三个镜像 → tag）、`UPGRADE.md`
- [ ] 4.2 `build.sh` 末尾外部域名扫描，与基线比对，出现新域名即报警（基线包含 `api.liblib.tv`、`app.tapnow.media`、`files.tapnow.media`）
- [ ] 4.3 所有补丁锚点校验，失败时打印文件与锚点

## 5. 视频工作台：登录与计费补丁

- [ ] 5.1 登录页：Linux.do 按钮改为「用 forapi 账号登录」并作为唯一入口；隐藏密码登录、注册、找回密码、注册协议
- [ ] 5.2 后端：非管理员拒绝密码注册与密码登录
- [ ] 5.3 新表 `user_bridge_credentials`（用户 ID、加密密钥、更新时间）
- [ ] 5.4 登录回调：拉完 userinfo 后调 `/oauth/studio-key` 并加密保存
- [ ] 5.5 `resolveProviderConfig` 增加 `userID`；桥接渠道注入用户密钥；取不到即失败不回落；更新全部调用方
- [ ] 5.6 上游 401 时作废会话并提示从 forapi 重新进入
- [ ] 5.7 单测：桥接渠道注入用户密钥；非桥接渠道不变；无密钥失败；用户之间不串用

## 6. 视频工作台：品牌与去第三方

配置（`provision/`，幂等）：

- [ ] 6.1 外观：品牌名、`brandSlug`、浅色 / 深色 logo、SEO、页脚版权、登录页标题与描述、登录页背景（本地素材或空）
- [ ] 6.2 皮肤：新建 forapi 皮肤（浅色 / 深色两套令牌，取主站 `tokens-forapi.css`），设为当前并禁止用户切换
- [ ] 6.3 功能开关：关闭 `credits`、`pluginCenter`、`systemPluginsVisibleToUsers`、`inspirationSourcesVisible`；TapNow / LibTV 画布导入保留（LibTV 是否启用取决于是否有 Token，见待确认 E）
- [ ] 6.4 系统渠道：一条桥接渠道 `https://forapi.cloud`，openai / newapi 协议插件，同步模型目录
- [ ] 6.5 OAuth：Linux.do 登录配置指向 forapi，字段映射 `sub` / `username` / `name` / `email` / `picture`

补丁：

- [ ] 6.6 字体改为 Inter（与主站一致）
- [ ] 6.7 默认外观中的 B 站外链视频与封面改为空
- [ ] 6.8 `/welcome` 重定向到 `/create`，不渲染贡献者、故事、赞助
- [ ] 6.9 渠道默认 UA 改为 `ForAPI-Studio/1.0`
- [ ] 6.10 隐藏管理后台 `settings/system-update` 菜单（在线更新会访问 GitHub）
- [ ] 6.11 3D 预演默认模型打进前端产物，不再请求 jsDelivr
- [ ] 6.12 前端与内置技能中硬编码的「影策」改为读取品牌名或「视频工作台」
- [ ] 6.13 镜像内保留上游 `LICENSE` 与 `NOTICE` 文件

## 6.5 模型目录同步（D8）

- [ ] 6.5.1 forapi `GET /oauth/studio-models`：客户端鉴权，返回全部 active 分组的模型并集与 capability，补测试
- [ ] 6.5.2 工作台 overlay：同步任务（启动 + 每 10 分钟），按 capability 选协议，走 SaveAdminChannelModel
- [ ] 6.5.3 只管理同步出来的模型；上游消失→停用；保留管理员改动；补测试
- [ ] 6.5.4 部署后验证：forapi 任一 active 分组新增模型后，10 分钟内工作台前台出现对应模型

## 7. 部署

- [ ] 7.1 `/opt/video-studio/`：compose（上游 server 版去掉 host-updater）+ `.env`（随机生成数据库、Redis、Agent 令牌）+ 内存上限
- [ ] 7.2 DNS：子域解析到本机（需要你在 Cloudflare 操作）
- [ ] 7.3 证书 + nginx server 块（`nginx -t` 通过才重载，失败自动还原）
- [ ] 7.4 首个管理员：一次性密码创建，执行 `provision`
- [ ] 7.5 forapi 后台填入 Client ID / Secret / 回调地址，打开开关

## 8. 验证

- [ ] 8.1 forapi 已登录 → 顶部菜单「视频工作台」→ 新标签页直接进入创作页，无需输入密码；用户名、头像正确
- [ ] 8.2 首次使用出现一次授权确认，之后不再出现
- [ ] 8.3 界面检查：logo、品牌名、配色、字体与主站一致；页面上没有影策品牌、作者、赞助、GitHub 链接
- [ ] 8.4 网络检查：浏览器开发者工具里，正常使用过程中没有请求第三方域名（模型请求只发往 forapi）；用户主动做的 TapNow / LibTV 导入除外
- [ ] 8.4a 用一个 TapNow 分享链接导入，节点、提示词、媒体正确出现在画布上
- [ ] 8.5 生成一张图、一段视频：forapi 使用记录里出现在该用户的专用密钥下，费用正确
- [ ] 8.6 两个账号同时使用，密钥与费用互不串用
- [ ] 8.7 停用专用密钥后生成失败且提示明确；从顶部菜单重新进入后恢复
- [ ] 8.8 模拟一次上游升级：改 `upstream.lock` 到上游较新的提交跑 `build.sh`，确认流程可行
- [ ] 8.9 new-api、sub2api 现有服务不受影响

## 待确认

- A. 子域名：`video.forapi.cloud` 还是其他？需要你在 Cloudflare 加解析
- B. ~~专用密钥的分组与额度~~ 已定：自动分组模式（用户有权限的全部分组，订阅优先、倍率低优先），额度由管理员配置（0 = 不限）
- C. 品牌名：「ForAPI 视频工作台」还是其他写法？
- D. 3D 预演（世界引擎内测）、时间线字幕转写（需要另装 whisper.cpp）要不要保留？不保留可以用功能开关关掉，部署更轻
- E. LibTV 导入需要一个站点级 LibTV Token：你有的话部署时填上即可启用；没有的话先保留入口，导入时会提示「LibTV 尚未启用」
