#!/usr/bin/env python3
"""把 forapi 集成补丁打到干净的上游源码上。

每条补丁 = (文件, 原文锚点, 替换内容, 期望出现次数)。锚点找不到或次数不对就整体失败，
不做模糊匹配：上游改了这段代码时宁可构建失败，也不能悄悄漏掉一处（尤其是密钥注入）。

用法：apply.py <上游源码目录>
"""
import sys
from pathlib import Path

PATCHES = []


def patch(path, old, new, count=1):
    PATCHES.append((path, old, new, count))


# ---------------------------------------------------------------- 后端：密钥桥接
# resolveProviderConfig 增加 userID 参数：所有调用点都必须显式传入任务归属用户，
# 上游以后新增调用点时会编译失败，而不是悄悄用上渠道公共密钥。
patch(
    "backend/internal/app/provider.go",
    "func (s *Service) resolveProviderConfig(config providerConfig) (providerConfig, error) {",
    "func (s *Service) resolveProviderConfig(userID string, config providerConfig) (providerConfig, error) {",
)
patch(
    "backend/internal/app/provider.go",
    "\tconfig.APIKey = channel.APIKey\n\tconfig.SecretKey = channel.SecretKey\n",
    "\tconfig.APIKey, err = s.ChannelAPIKeyForUser(channel, userID)\n"
    "\tif err != nil {\n\t\treturn providerConfig{}, err\n\t}\n"
    "\tconfig.SecretKey = channel.SecretKey\n",
)
patch(
    "backend/internal/app/provider.go",
    "\tconfig, err := s.resolveProviderConfig(input.Config)\n",
    "\tconfig, err := s.resolveProviderConfig(userID, input.Config)\n",
)
patch(
    "backend/internal/app/provider_task_cancellation.go",
    "\tconfig, err := s.resolveProviderConfig(input.Config)\n",
    "\tconfig, err := s.resolveProviderConfig(task.UserID, input.Config)\n",
)
patch(
    "backend/internal/app/provider_task_recovery.go",
    "\tconfig, err := s.resolveProviderConfig(input.Config)\n",
    "\tconfig, err := s.resolveProviderConfig(task.UserID, input.Config)\n",
)
patch(
    "backend/internal/app/task_worker.go",
    "\tresolved, resolveErr := s.resolveProviderConfig(input.Config)\n",
    "\tresolved, resolveErr := s.resolveProviderConfig(task.UserID, input.Config)\n",
)
patch(
    "backend/internal/app/task_media_recovery.go",
    "\t\t\t\tconfig, err = s.resolveProviderConfig(config)\n",
    "\t\t\t\tconfig, err = s.resolveProviderConfig(task.UserID, config)\n",
)
patch(
    "backend/internal/app/admin_channel_test.go",
    "svc.resolveProviderConfig(providerConfig{",
    "svc.resolveProviderConfig(\"\", providerConfig{",
)
# 后台「拉取模型目录」：桥接渠道用当前管理员自己的 forapi 密钥
patch(
    "backend/internal/app/channel_models_catalog.go",
    "\tmodels, err := s.FetchChannelModels(ctx, actor, ChannelModelsRequest{BaseURL: channel.BaseURL, APIKey: channel.APIKey, APIFormat: channel.APIFormat, Headers: headers})\n",
    "\tchannelKey, err := s.ChannelAPIKeyForUser(channel, actor.ID)\n"
    "\tif err != nil {\n\t\treturn nil, err\n\t}\n"
    "\tmodels, err := s.FetchChannelModels(ctx, actor, ChannelModelsRequest{BaseURL: channel.BaseURL, APIKey: channelKey, APIFormat: channel.APIFormat, Headers: headers})\n",
    count=2,
)
# 同步代理（画布里的对话/改写等直接请求）：同样换成请求者的密钥
patch(
    "backend/internal/handler/auth.go",
    "\tdefer releaseChannel()\n\tif c.Request.Method == http.MethodPost {\n",
    "\tdefer releaseChannel()\n"
    "\tchannelKey, err := svc.ChannelAPIKeyForUser(channel, user.ID)\n"
    "\tif err != nil {\n\t\tfailService(c, err)\n\t\treturn\n\t}\n"
    "\tif c.Request.Method == http.MethodPost {\n",
)
patch("backend/internal/handler/auth.go", 'upstreamReq.Header.Set("x-goog-api-key", channel.APIKey)', 'upstreamReq.Header.Set("x-goog-api-key", channelKey)')
patch("backend/internal/handler/auth.go", 'upstreamReq.Header.Set("x-api-key", channel.APIKey)', 'upstreamReq.Header.Set("x-api-key", channelKey)')
patch("backend/internal/handler/auth.go", 'upstreamReq.Header.Set("Authorization", "Bearer "+channel.APIKey)', 'upstreamReq.Header.Set("Authorization", "Bearer "+channelKey)')

# 未上传登录页视频时，后端默认返回 B 站 CDN 上的演示视频：改为不使用视频
patch(
    "backend/internal/app/appearance.go",
    '\tdefaultAppearanceVideoURL  = "https://boss-shjd.biliapi.net/updream/aniforge/video/video_bbcb00bd-650d-4249-9346-5cd21fd2484c_m1hc-u0-1pu13x-3v1s.mp4"\n'
    '\tdefaultAppearancePosterURL = "https://i0.hdslb.com/bfs/aitool/aniforge/image/02933f26-5f1b-49ff-a811-b7f95ee5e5b8_m1hc-u0-sau.jpg"\n',
    '\tdefaultAppearanceVideoURL  = ""\n\tdefaultAppearancePosterURL = ""\n',
)

# ---------------------------------------------------------------- 后端：登录
# OAuth 回调：确认账号可用后立刻换取 forapi 专用密钥，失败则整次登录失败
patch(
    "backend/internal/auth/linuxdo.go",
    "\tloginLogUser(event, user)\n\tif user.Status != model.UserStatusActive {\n\t\treturn nil, kernel.Forbidden(\"该账号已被禁用\")\n\t}\n",
    "\tloginLogUser(event, user)\n\tif user.Status != model.UserStatusActive {\n\t\treturn nil, kernel.Forbidden(\"该账号已被禁用\")\n\t}\n"
    "\tif err := s.storeForAPIBridgeKey(setting, accessToken, user.ID); err != nil {\n\t\treturn nil, err\n\t}\n",
)
# 密码登录只留给本地管理员；普通用户一律走 forapi
patch(
    "backend/internal/auth/auth.go",
    "\tif !verifyPassword(req.Password, user.PasswordHash) {\n\t\treturn nil, kernel.Unauthorized(\"用户名、邮箱或密码不正确\")\n\t}\n",
    "\tif !verifyPassword(req.Password, user.PasswordHash) {\n\t\treturn nil, kernel.Unauthorized(\"用户名、邮箱或密码不正确\")\n\t}\n"
    "\tif user.Role != model.UserRoleAdmin {\n\t\treturn nil, kernel.Forbidden(\"请使用 ForAPI 账号登录\")\n\t}\n",
)
# 本地注册只用于创建第一个管理员；之后新用户只能经 forapi 首次登录创建
patch(
    "backend/internal/auth/auth.go",
    "\tvar verifiedCode *model.EmailVerificationCode\n\tvar verification *model.AuthVerification\n",
    "\tif count > 0 {\n\t\treturn nil, kernel.Forbidden(\"请使用 ForAPI 账号登录，首次登录会自动创建账号\")\n\t}\n"
    "\tvar verifiedCode *model.EmailVerificationCode\n\tvar verification *model.AuthVerification\n",
)

# ---------------------------------------------------------------- 前端：品牌与第三方
patch("web/src/stores/use-appearance-store.ts", 'brandName: "影策",', 'brandName: "ForAPI 视频工作台",')
patch("web/src/stores/use-appearance-store.ts", 'seoTitle: "影策",', 'seoTitle: "ForAPI 视频工作台",')
patch(
    "web/src/stores/use-appearance-store.ts",
    'seoDescription: "影策，面向 AI 影视与短剧创作的工作台。",',
    'seoDescription: "ForAPI 视频工作台，面向 AI 影视与短剧创作。",',
)
patch(
    "web/src/stores/use-appearance-store.ts",
    "footerCopyright: `© ${new Date().getFullYear()} 影策. All rights reserved.`,",
    "footerCopyright: `© ${new Date().getFullYear()} ForAPI. All rights reserved.`,",
)
# 登录页默认背景视频来自 B 站 CDN：清空，背景改用纯色
patch(
    "web/src/stores/use-appearance-store.ts",
    'authVideoUrl: "https://boss-shjd.biliapi.net/updream/aniforge/video/video_bbcb00bd-650d-4249-9346-5cd21fd2484c_m1hc-u0-1pu13x-3v1s.mp4",',
    'authVideoUrl: "",',
)
patch(
    "web/src/stores/use-appearance-store.ts",
    'authVideoPosterUrl: "https://i0.hdslb.com/bfs/aitool/aniforge/image/02933f26-5f1b-49ff-a811-b7f95ee5e5b8_m1hc-u0-sau.jpg",',
    'authVideoPosterUrl: "",',
)
patch(
    "web/src/components/channel-headers-editor.tsx",
    'const DEFAULT_USER_AGENT = "Yingce/1.0 (+https://github.com/ddcat-ai/open-ai-canvas)";',
    'const DEFAULT_USER_AGENT = "ForAPI-Studio/1.0";',
)
patch("web/src/lib/sidebar-checkin.ts", 'brandName.trim() || "影策"', 'brandName.trim() || "ForAPI 视频工作台"')
patch("web/src/lib/creation/creative-agent-tools.ts", "你是影策的通用创作助手", "你是 ForAPI 视频工作台的通用创作助手")
patch("web/src/lib/timeline/ai-command-schema.ts", "你是影策视频剪辑工作台的 AI 剪辑助手", "你是 ForAPI 视频工作台的 AI 剪辑助手")
# 身份标识：底层仍是通用 OAuth 登录（沿用 linuxdo 命名），界面上显示为 ForAPI
patch(
    "web/src/components/layout/identity-provider-badge.tsx",
    'title={user.identityUsername ? `Linux.do · @${user.identityUsername}` : "Linux.do 用户"}',
    'title={user.identityUsername ? `ForAPI · @${user.identityUsername}` : "ForAPI 用户"}',
)
patch("web/src/components/layout/identity-provider-badge.tsx", '{compact ? "L" : "Linux.do"}', '{compact ? "F" : "ForAPI"}')
patch("web/src/pages/settings/profile-pane.tsx", "Linux.do", "ForAPI", count=4)
# 登录卡片去掉「注册」页签：账号首次用 ForAPI 登录时自动创建
patch(
    "web/src/pages/auth/auth-scene.tsx",
    '    { key: "login", label: "登录" },\n    { key: "register", label: "注册" },\n',
    '    { key: "login", label: "登录" },\n',
)
# 没有配置登录页视频时，不显示「播放品牌影片」按钮
patch(
    "web/src/pages/auth/auth-scene.tsx",
    '                        <button type="button" className="auth-scene-hero-control inline-flex',
    '                        {appearance.authVideoUrl ? <button type="button" className="auth-scene-hero-control inline-flex',
)
patch(
    "web/src/pages/auth/auth-scene.tsx",
    '                            {videoPlaying ? "创作正在发生" : "播放品牌影片"}\n                        </button>\n',
    '                            {videoPlaying ? "创作正在发生" : "播放品牌影片"}\n                        </button> : null}\n',
)
# 系统更新会去 GitHub 拉上游版本并原地升级，与本仓库的补丁构建冲突：隐藏入口
patch(
    "web/src/pages/admin/components/admin-shell.tsx",
    '            { path: "/admin/settings/system-update", label: "系统更新", shortLabel: "更新", icon: <RefreshCw className="size-4" /> },\n',
    "",
)


# ---------------------------------------------------------------- 后端：ForAPI 模型目录同步
# 启动后台同步：桥接渠道的模型由 forapi「工作台」分组决定，管理员不再手工添加。
patch(
    "backend/cmd/server/main.go",
    "\tsvc.StartWorker()\n",
    "\tsvc.StartWorker()\n\tgo svc.RunForAPIModelSync(ctx)\n",
)


def main():
    if len(sys.argv) != 2:
        sys.exit("用法：apply.py <上游源码目录>")
    root = Path(sys.argv[1])
    # 先全部检查，再统一写入：任何一条失败都不留下半打补丁的源码
    contents = {}
    errors = []
    for path, old, new, count in PATCHES:
        file = root / path
        if path not in contents:
            if not file.is_file():
                errors.append(f"{path}: 文件不存在")
                contents[path] = None
                continue
            contents[path] = file.read_text(encoding="utf-8")
        text = contents[path]
        if text is None:
            continue
        found = text.count(old)
        if found != count:
            errors.append(f"{path}: 锚点出现 {found} 次（期望 {count}）：{old.strip()[:80]!r}")
            continue
        contents[path] = text.replace(old, new)
    if errors:
        print("补丁失败，源码未修改：", file=sys.stderr)
        for line in errors:
            print("  - " + line, file=sys.stderr)
        sys.exit(1)
    for path, text in contents.items():
        (root / path).write_text(text, encoding="utf-8")
    print(f"已应用 {len(PATCHES)} 处补丁，涉及 {len(contents)} 个文件")


if __name__ == "__main__":
    main()
