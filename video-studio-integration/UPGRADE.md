# ForAPI 视频工作台：集成层说明与升级流程

上游：[ddcat-ai/open-ai-canvas](https://github.com/ddcat-ai/open-ai-canvas)（影策）。本目录不 fork 上游，
每次构建都从 `upstream.lock` 锁定的 commit 导出干净源码，再叠加本目录的改动。

## 目录

| 路径 | 作用 |
|---|---|
| `upstream.lock` | 上游 commit、版本号；overlay 覆盖的上游文件的 sha256 |
| `patches/apply.py` | 锚点补丁：原文必须精确出现指定次数，否则整体失败、源码不动 |
| `overlay/` | 整文件放入：新增的桥接代码与测试、重写的登录/注册页 |
| `domain-allowlist.txt` | 前端允许出现的外部域名，新域名会让构建失败 |
| `build.sh` | 导出 → 校验 → 打补丁 → 复制 overlay → 扫描域名 → 构建镜像 |
| `deploy/` | compose、nginx 站点配置、初始化脚本 |

部署机：`/opt/video-studio/`（`integration/` 是本目录的副本，`.env` 是密钥，权限 600，不入库）。

## 改了什么

- **登录**：普通用户只能用 ForAPI 账号登录（上游通用 OAuth「Linux.do 登录」指向 forapi.cloud）；
  首次登录自动建号；本地注册只放行第一个账号（管理员）；密码登录只对管理员开放（`/login?admin=1`）。
- **计费**：上游积分关闭。登录回调时用同一个 access token 调 forapi `/oauth/studio-key`，取得该用户的
  「视频工作台」专用密钥，加密存在 `system_settings`（key `forapi_bridge_key:<用户ID>`）。
  渠道 API Key 填 `forapi-user-key` 的系统渠道，生成时换成发起用户自己的密钥；找不到就报错，不回退。
  `resolveProviderConfig` 增加了 `userID` 参数，上游新增调用点会编译失败，而不是漏掉密钥替换。
- **品牌**：默认文案、AI 助手提示词、身份标识改为 ForAPI；登录页默认 B 站视频去掉；
  logo、配色、SEO 由 `deploy/provision.py` 写入后台设置（logo 直接取 forapi.cloud 正在用的）。
- **第三方**：欢迎页（上游介绍、GitHub）关闭；模型目录服务 `CANVAS_PROVIDER_CATALOG_DISABLED=1`；
  默认 UA 改为 `ForAPI-Studio/1.0`；「系统更新」入口隐藏（它会从 GitHub 原地升级，绕过本集成层）。
  TapNow、LibTV、Eagle 素材导入属于创作功能，保留。

## 升级上游

```bash
ssh newapi-production
cd /opt/video-studio/integration
git -C ../build/upstream.git fetch origin && git -C ../build/upstream.git log --oneline -5 main
```

1. 把 `upstream.lock` 的 `UPSTREAM_COMMIT` / `UPSTREAM_VERSION` 改成新版本（在本仓库改，再同步到部署机）。
2. `./build.sh --check-only`，按失败类型处理：
   - **补丁锚点不匹配**：读上游这段的新代码，更新 `apply.py` 里对应的原文和替换内容。
     密钥相关的补丁（`resolveProviderConfig`、`channel.APIKey`、`proxySystemRequestPath`）
     必须逐处确认，再 `grep -rn "channel.APIKey\|resolveProviderConfig(" backend/internal` 检查有没有新的取密钥路径。
   - **overlay 上游已修改**：`git -C ../build/upstream.git diff <旧commit> <新commit> -- <文件>`，
     把需要的上游改动合进 overlay 版本，再更新 `upstream.lock` 里的 sha256。
   - **新外部域名**：确认用途。追踪、CDN、推广要用补丁去掉；正常功能再加进 `domain-allowlist.txt`。
3. 跑桥接测试（需要 cgo）：
   ```bash
   docker run --rm -v /opt/video-studio/build/src/backend:/src -w /src -e CGO_ENABLED=1 golang:1.25 \
     go test ./internal/auth/ ./internal/app/ -run 'ForAPI|Bridge|OrdinaryChannel|PasswordLoginOnlyForAdmins|StudioKeyURL'
   ```
   上游自带测试里有一批会因本集成的行为改变而失败（普通用户密码登录/注册被拒、OAuth 回调多了取密钥一步），
   这是预期的；判断标准是「打补丁前后的失败列表之差」只包含这类用例。

   前端同理，用「干净上游」做基线再对比失败集合：
   ```bash
   rm -rf /tmp/webbase && mkdir -p /tmp/webbase
   git -C /opt/video-studio/build/upstream.git archive "$UPSTREAM_COMMIT" | tar -x -C /tmp/webbase
   for t in base patched; do
     S=/tmp/webbase; [ "$t" = patched ] && S=/opt/video-studio/build/src
     docker run --rm -v "$S":/app -w /app/web oven/bun:1.3.13 sh -c \
       "bun install --frozen-lockfile --cache-dir=/root/.bun/install/cache >/dev/null 2>&1; bun test" \
       > "/tmp/web-$t.txt" 2>&1
     grep -E '^\(fail\)' "/tmp/web-$t.txt" | sed -E 's/ \[[0-9.]+ms\]$//' | sort -u > "/tmp/web-$t-set.txt"
   done
   comm -13 /tmp/web-base-set.txt /tmp/web-patched-set.txt   # 新增失败，应为空
   ```

   注意 `test/channel-model-templates.interaction.test.ts` 里的「完整模型编辑器」用例是**上游自带的抖动测试**：
   全量并发跑时偶尔会因为 5 秒超时而失败，单独跑一定通过，干净上游同样会抖。
   遇到这组失败要重跑确认，不要当成回归。
4. `./build.sh` 构建，记下输出的镜像标签；`/opt/video-studio/.env` 里 `VIDEO_STUDIO_TAG` 改成新标签，
   `cd /opt/video-studio && docker compose up -d`（migrate 容器会先跑数据库迁移）。
5. 回滚：`.env` 改回上一个标签再 `docker compose up -d`。数据库迁移不可逆，升级前先
   `docker exec video-studio-postgres pg_dump -U video_studio video_studio | gzip > /root/video-studio-$(date +%F).sql.gz`。

## 初始化 / 重新应用配置

```bash
docker run --rm --network video-studio_internal --env-file /opt/video-studio/.env \
  -v /opt/video-studio/integration/deploy:/p:ro python:3.12-alpine python /p/provision.py
```

可重复执行。之后在工作台后台给「ForAPI」系统渠道添加模型（模型 Key 用 forapi 的模型名）。
后台「拉取模型」使用当前管理员自己的 ForAPI 密钥，本地管理员没有，所以需手动添加，
或让一个 ForAPI 账号登录后在用户管理里设为管理员再拉取。

## forapi 侧

后台 → 视频工作台接入：Client ID / Secret 与部署机 `.env` 的 `FORAPI_CLIENT_ID` / `FORAPI_CLIENT_SECRET` 一致，
回调地址 `https://video.forapi.cloud/api/auth/linuxdo/callback`，工作台地址 `https://video.forapi.cloud`。
