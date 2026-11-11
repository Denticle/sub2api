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
# v1.8.x 新增的电商 Agent 调用点：同样要显式传入任务归属用户，
# 否则会悄悄用上渠道公共密钥（上游以后再新增调用点时会再次编译失败）。
patch(
    "backend/internal/app/ecommerce_agent_provider.go",
    "\tresolved, err := s.resolveProviderConfig(providerConfig{ChannelID: channelID, ChannelModelKey: modelKey, Model: modelKey})\n",
    "\tresolved, err := s.resolveProviderConfig(userID, providerConfig{ChannelID: channelID, ChannelModelKey: modelKey, Model: modelKey})\n",
)
patch(
    "backend/internal/app/ecommerce_agent.go",
    "func (s *Service) ecommerceLiveImageConfig(req ecommerceagent.Request, cfg ecommerceProviderConfig) (ecommerceProviderConfig, error) {",
    "func (s *Service) ecommerceLiveImageConfig(userID string, req ecommerceagent.Request, cfg ecommerceProviderConfig) (ecommerceProviderConfig, error) {",
)
patch(
    "backend/internal/app/ecommerce_agent.go",
    "\tprovider, err := s.resolveProviderConfig(providerConfig{ChannelID: channelID, ChannelModelKey: modelKey, Model: modelKey})\n",
    "\tprovider, err := s.resolveProviderConfig(userID, providerConfig{ChannelID: channelID, ChannelModelKey: modelKey, Model: modelKey})\n",
)
patch(
    "backend/internal/app/ecommerce_agent.go",
    "\t\tcfg, err = s.ecommerceLiveImageConfig(req, cfg)\n",
    "\t\tcfg, err = s.ecommerceLiveImageConfig(userID, req, cfg)\n",
    count=2,
)
patch(
    "backend/internal/app/ecommerce_agent_test.go",
    "svc.ecommerceLiveImageConfig(req, ecommerceConfig())",
    "svc.ecommerceLiveImageConfig(\"\", req, ecommerceConfig())",
    count=2,
)
# 后台「拉取模型目录」：桥接渠道用当前管理员自己的 forapi 密钥
# v1.8.4 起两处调用分属不同文件（catalog 拉取 / snapshot 读目录），分别打补丁。
patch(
    "backend/internal/app/channel_models_catalog.go",
    "\tmodels, err := s.FetchChannelModels(ctx, actor, ChannelModelsRequest{BaseURL: channel.BaseURL, APIKey: channel.APIKey, APIFormat: channel.APIFormat, ProxyURL: channel.ProxyURL, Headers: headers})\n",
    "\tchannelKey, err := s.ChannelAPIKeyForUser(channel, actor.ID)\n"
    "\tif err != nil {\n\t\treturn nil, err\n\t}\n"
    "\tmodels, err := s.FetchChannelModels(ctx, actor, ChannelModelsRequest{BaseURL: channel.BaseURL, APIKey: channelKey, APIFormat: channel.APIFormat, ProxyURL: channel.ProxyURL, Headers: headers})\n",
)
patch(
    "backend/internal/app/channel_catalog_snapshot.go",
    "\tmodels, err := s.FetchChannelModels(ctx, actor, ChannelModelsRequest{BaseURL: channel.BaseURL, APIKey: channel.APIKey, APIFormat: channel.APIFormat, ProxyURL: channel.ProxyURL, Headers: headers})\n",
    "\tchannelKey, err := s.ChannelAPIKeyForUser(channel, actor.ID)\n"
    "\tif err != nil {\n\t\treturn nil, err\n\t}\n"
    "\tmodels, err := s.FetchChannelModels(ctx, actor, ChannelModelsRequest{BaseURL: channel.BaseURL, APIKey: channelKey, APIFormat: channel.APIFormat, ProxyURL: channel.ProxyURL, Headers: headers})\n",
)
# 同步代理（画布里的对话/改写等直接请求）：同样换成请求者的密钥
# v1.8.4 起鉴权统一收进 ApplySystemProxyAuth，按协议声明写请求头。
# 最小侵入做法：在鉴权之前把 channel.APIKey 换成请求者自己的密钥，
# 由上游的 ApplySystemProxyAuth 负责按协议写请求头，不再单独改每个 header。
patch(
    "backend/internal/handler/auth.go",
    "\tdefer releaseChannel()\n\tif c.Request.Method == http.MethodPost {\n",
    "\tdefer releaseChannel()\n"
    "\tchannelKey, err := svc.ChannelAPIKeyForUser(channel, user.ID)\n"
    "\tif err != nil {\n\t\tfailService(c, err)\n\t\treturn\n\t}\n"
    "\t// 让上游按协议写鉴权时使用请求者自己的密钥（channel 每次从库里取，可安全改写）。\n"
    "\tchannel.APIKey = channelKey\n"
    "\tif c.Request.Method == http.MethodPost {\n",
)

# 未上传登录页视频时，后端默认返回 B 站 CDN 上的演示视频：改为不使用视频
patch(
    "backend/internal/app/appearance.go",
    '\tdefaultAppearanceVideoURL  = "https://boss-shjd.biliapi.net/updream/aniforge/video/video_bbcb00bd-650d-4249-9346-5cd21fd2484c_m1hc-u0-1pu13x-3v1s.mp4"\n'
    '\tdefaultAppearancePosterURL = "https://i0.hdslb.com/bfs/aitool/aniforge/image/02933f26-5f1b-49ff-a811-b7f95ee5e5b8_m1hc-u0-sau.jpg"\n',
    '\tdefaultAppearanceVideoURL  = ""\n\tdefaultAppearancePosterURL = ""\n',
)

# ---------------------------------------------------------------- 后端：登录
# OAuth 回调：确认账号可用后立刻换取 forapi 专用密钥，失败则整次登录失败
# v1.8.4 起上游把登录拆成 CompleteLinuxDOLogin（换 token）→ finishOAuthLogin（建会话），
# 且 finishOAuthLogin 由 linuxdo 与 oidc 共用：把 forapi 的 accessToken 作为参数传进来，
# 为空表示这条路径不注入密钥（oidc 不经过 Linux.do 换取 forapi 密钥）。
patch(
    "backend/internal/auth/linuxdo.go",
    "func (s *Service) finishOAuthLogin(state *model.OAuthState, event *model.LoginLog, profile oauthProfile, allowPrivateSignup bool) (*OAuthCallbackResult, error) {",
    "func (s *Service) finishOAuthLogin(state *model.OAuthState, event *model.LoginLog, profile oauthProfile, allowPrivateSignup bool, forAPIAccessToken string) (*OAuthCallbackResult, error) {",
)
patch(
    "backend/internal/auth/linuxdo.go",
    "\tloginLogUser(event, user)\n\tif user.Status != model.UserStatusActive {\n\t\treturn nil, kernel.Forbidden(\"该账号已被禁用\")\n\t}\n",
    "\tloginLogUser(event, user)\n\tif user.Status != model.UserStatusActive {\n\t\treturn nil, kernel.Forbidden(\"该账号已被禁用\")\n\t}\n"
    "\tif forAPIAccessToken != \"\" {\n"
    "\t\t// finishOAuthLogin 拿不到上游的 setting 变量，这里重新读一次 forapi 的 OAuth 配置。\n"
    "\t\t_, forAPISetting, settingErr := s.readLinuxDOSetting()\n"
    "\t\tif settingErr != nil {\n"
    "\t\t\treturn nil, settingErr\n"
    "\t\t}\n"
    "\t\tif err := s.storeForAPIBridgeKey(forAPISetting, forAPIAccessToken, user.ID); err != nil {\n"
    "\t\t\treturn nil, err\n\t\t}\n"
    "\t}\n",
)
# linuxdo 路径把 accessToken 传下去
patch(
    "backend/internal/auth/linuxdo.go",
    "\treturn s.finishOAuthLogin(state, event, oauthProfile{\n"
    "\t\tprovider: \"linuxdo\", subject: subject, username: providerUsername, name: displayName,\n"
    "\t\temail: profileString(profile, setting.EmailField), avatarURL: profileString(profile, setting.AvatarField),\n"
    "\t}, false)\n",
    "\treturn s.finishOAuthLogin(state, event, oauthProfile{\n"
    "\t\tprovider: \"linuxdo\", subject: subject, username: providerUsername, name: displayName,\n"
    "\t\temail: profileString(profile, setting.EmailField), avatarURL: profileString(profile, setting.AvatarField),\n"
    "\t}, false, accessToken)\n",
)
# oidc 路径没有 Linux.do 的 accessToken，传空串跳过注入
patch(
    "backend/internal/auth/oidc.go",
    "\t}, cfg.allowedGroup != \"\")\n",
    "\t}, cfg.allowedGroup != \"\", \"\")\n",
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

# ---------------------------------------------------------------- 后端：适配上游测试
# 上游测试用的假 OAuth 服务器只实现了 /token 与 /profile，没有 /studio-key；
# 密钥注入会向同站的 /studio-key 取密钥，取不到就让整次登录失败，所以假服务器要补上。
patch(
    "backend/internal/auth/auth_registration_test.go",
    "\t\"net/url\"\n\t\"strconv\"\n\t\"testing\"\n",
    "\t\"net/url\"\n\t\"strconv\"\n\t\"strings\"\n\t\"testing\"\n",
)
patch(
    "backend/internal/auth/auto_registration_test.go",
    "\t\"errors\"\n\t\"testing\"\n",
    "\t\"errors\"\n\t\"strings\"\n\t\"testing\"\n",
)
patch(
    "backend/internal/auth/gmail_registration_test.go",
    "\t\"strings\"\n\t\"testing\"\n\t\"time\"\n",
    "\t\"testing\"\n\t\"time\"\n",
)
patch(
    "backend/internal/auth/auth_registration_test.go",
    "\t\tcase \"/profile\":\n\t\t\tpayload = map[string]string{\"id\": \"42\", \"username\": \"linux-user\", \"name\": \"Linux User\"}\n",
    "\t\tcase \"/profile\":\n\t\t\tpayload = map[string]string{\"id\": \"42\", \"username\": \"linux-user\", \"name\": \"Linux User\"}\n"
    "\t\tcase \"/studio-key\":\n\t\t\tpayload = map[string]string{\"api_key\": \"test-studio-key\"}\n",
)
patch(
    "backend/internal/auth/login_status_test.go",
    "\t\t\tbody := `{\"id\":\"42\",\"username\":\"member\"}`\n"
    "\t\t\tif req.URL.Path == \"/token\" {\n\t\t\t\tbody = `{\"access_token\":\"test-token\"}`\n\t\t\t}\n",
    "\t\t\tbody := `{\"id\":\"42\",\"username\":\"member\"}`\n"
    "\t\t\tif req.URL.Path == \"/token\" {\n\t\t\t\tbody = `{\"access_token\":\"test-token\"}`\n\t\t\t}\n"
    "\t\t\tif req.URL.Path == \"/studio-key\" {\n\t\t\t\tbody = `{\"api_key\":\"test-studio-key\"}`\n\t\t\t}\n",
)
# 普通用户的本地注册与密码登录已被关闭（只保留首个管理员的本地注册），
# 相应测试改用管理员账号，走仍然保留的那条路径。
patch(
    "backend/internal/auth/auth_registration_test.go",
    "\tuser := model.User{ID: \"existing-user\", Username: \"existing-user\", PasswordHash: hash, Role: model.UserRoleUser, Status: model.UserStatusActive}\n",
    "\tuser := model.User{ID: \"existing-user\", Username: \"existing-user\", PasswordHash: hash, Role: model.UserRoleAdmin, Status: model.UserStatusActive}\n",
)
patch(
    "backend/internal/auth/login_status_test.go",
    "\tif err := db.Create(&model.User{ID: \"member\", Username: \"member\", DisplayName: \"Before\", PasswordHash: hash, Role: model.UserRoleUser, Status: model.UserStatusActive}).Error; err != nil {\n",
    "\tif err := db.Create(&model.User{ID: \"member\", Username: \"member\", DisplayName: \"Before\", PasswordHash: hash, Role: model.UserRoleAdmin, Status: model.UserStatusActive}).Error; err != nil {\n",
)
patch(
    "backend/internal/auth/login_status_test.go",
    "\tif err := svc.repo.CreateOAuthUser(&model.User{ID: \"member\", Username: \"member\", DisplayName: \"Before\", Role: model.UserRoleUser, Status: model.UserStatusActive}, &model.UserIdentity{ID: \"identity\", UserID: \"member\", Provider: \"linuxdo\", Subject: \"42\"}); err != nil {\n",
    "\tif err := svc.repo.CreateOAuthUser(&model.User{ID: \"member\", Username: \"member\", DisplayName: \"Before\", Role: model.UserRoleAdmin, Status: model.UserStatusActive}, &model.UserIdentity{ID: \"identity\", UserID: \"member\", Provider: \"linuxdo\", Subject: \"42\"}); err != nil {\n",
)

# 上游这几个测试验证的是「普通用户也能本地注册」；该能力已关闭
# （首个管理员之外的账号一律走 forapi 首次登录），所以改为断言新行为：
# 注册被拒、且不会留下用户或消耗验证码。
patch(
    "backend/internal/auth/auth_registration_test.go",
    "\tresult, err := svc.Register(RegisterRequest{Username: \"member\", Email: \"member@example.com\", EmailCode: code, Password: \"password\", AcceptedTerms: true})\n"
    "\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n"
    "\tif result.User.Role != model.UserRoleUser || result.Session == \"\" {\n\t\tt.Fatalf(\"Register() result = %#v\", result)\n\t}\n",
    "\tif _, err := svc.Register(RegisterRequest{Username: \"member\", Email: \"member@example.com\", EmailCode: code, Password: \"password\", AcceptedTerms: true}); err == nil || !strings.Contains(err.Error(), \"请使用 ForAPI 账号登录\") {\n"
    "\t\tt.Fatalf(\"Register() must be closed once an administrator exists: %v\", err)\n\t}\n"
    "\tvar count int64\n"
    "\tif err := db.Model(&model.User{}).Count(&count).Error; err != nil || count != 1 {\n\t\tt.Fatalf(\"rejected registration created a user: count=%d err=%v\", count, err)\n\t}\n",
)
patch(
    "backend/internal/auth/auto_registration_test.go",
    "\tresult, err := svc.Register(RegisterRequest{Username: \"member\", Password: \"password\", AcceptedTerms: true})\n"
    "\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n"
    "\tif result.User.Email != \"\" || result.User.EmailVerifiedAt != nil {\n\t\tt.Fatalf(\"user = %#v\", result.User)\n\t}\n"
    "\t// 填写的邮箱按未验证保存，后续在账号设置里验证绑定。\n"
    "\tresult, err = svc.Register(RegisterRequest{Email: \"Writer.One@example.com\", Password: \"password\", AcceptedTerms: true})\n"
    "\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n"
    "\tif result.User.Username != \"writer-one\" || result.User.Email != \"writer.one@example.com\" || result.User.EmailVerifiedAt != nil {\n\t\tt.Fatalf(\"email user = %#v\", result.User)\n\t}\n",
    "\tif _, err := svc.Register(RegisterRequest{Username: \"member\", Password: \"password\", AcceptedTerms: true}); err == nil || !strings.Contains(err.Error(), \"请使用 ForAPI 账号登录\") {\n"
    "\t\tt.Fatalf(\"Register() must be closed once an administrator exists: %v\", err)\n\t}\n"
    "\tvar count int64\n"
    "\tif err := db.Model(&model.User{}).Count(&count).Error; err != nil || count != 1 {\n\t\tt.Fatalf(\"rejected registration created a user: count=%d err=%v\", count, err)\n\t}\n",
)
patch(
    "backend/internal/auth/auto_registration_test.go",
    "\tif _, err := svc.Register(RegisterRequest{Username: \"member\", Password: \"password\", AcceptedTerms: true}); err == nil {\n\t\tt.Fatal(\"Register() without email succeeded\")\n\t}\n"
    "\tif _, err := svc.Register(RegisterRequest{Email: \"member@example.com\", Password: \"password\", AcceptedTerms: true}); err != nil {\n\t\tt.Fatal(err)\n\t}\n",
    "\tif _, err := svc.Register(RegisterRequest{Username: \"member\", Password: \"password\", AcceptedTerms: true}); err == nil {\n\t\tt.Fatal(\"Register() without email succeeded\")\n\t}\n"
    "\tif _, err := svc.Register(RegisterRequest{Email: \"member@example.com\", Password: \"password\", AcceptedTerms: true}); err == nil || !strings.Contains(err.Error(), \"请使用 ForAPI 账号登录\") {\n"
    "\t\tt.Fatalf(\"Register() must be closed once an administrator exists: %v\", err)\n\t}\n",
)
patch(
    "backend/internal/auth/email_domain_test.go",
    "\tif _, err := svc.Register(RegisterRequest{Username: \"member\", Email: \"member@example.com\", Password: \"strong-password\", EmailCode: deliveredCode, AcceptedTerms: true}); err == nil || !strings.Contains(err.Error(), \"白名单\") {\n"
    "\t\tt.Fatalf(\"registration must recheck current domain policy: %v\", err)\n\t}\n",
    "\tif _, err := svc.Register(RegisterRequest{Username: \"member\", Email: \"member@example.com\", Password: \"strong-password\", EmailCode: deliveredCode, AcceptedTerms: true}); err == nil {\n"
    "\t\tt.Fatal(\"registration must be closed once an administrator exists\")\n\t}\n",
)
patch(
    "backend/internal/auth/gmail_registration_test.go",
    "\t_, err = svc.Register(RegisterRequest{Username: \"new-member\", Email: email, EmailCode: code, Password: \"password\", AcceptedTerms: true})\n"
    "\tif err == nil || !strings.Contains(err.Error(), \"邮箱已被注册\") {\n\t\tt.Fatalf(\"registration accepted a second account for the same Gmail mailbox or failed unexpectedly: %v\", err)\n\t}\n",
    "\t_, err = svc.Register(RegisterRequest{Username: \"new-member\", Email: email, EmailCode: code, Password: \"password\", AcceptedTerms: true})\n"
    "\tif err == nil {\n\t\tt.Fatal(\"registration must be closed once an administrator exists\")\n\t}\n",
)

# 上游 resolveProviderConfig 增加了 userID 参数（用于按用户取桥接渠道密钥），
# 这些测试用的都是普通渠道（密钥就在渠道上），传空 userID 即可。
patch(
    "backend/internal/app/media_transport_test.go",
    "\tresolved, err := s.resolveProviderConfig(config)\n",
    "\tresolved, err := s.resolveProviderConfig(\"\", config)\n",
)
patch(
    "backend/internal/app/text_thinking_test.go",
    "\tresolved, err := svc.resolveProviderConfig(taskInput.Config)\n",
    "\tresolved, err := svc.resolveProviderConfig(\"\", taskInput.Config)\n",
)
patch(
    "backend/internal/app/text_thinking_test.go",
    "\tresolved, err = svc.resolveProviderConfig(spoofed)\n",
    "\tresolved, err = svc.resolveProviderConfig(\"\", spoofed)\n",
)
patch(
    "backend/internal/app/text_thinking_test.go",
    "\tif _, err := svc.resolveProviderConfig(taskInput.Config); err == nil {\n",
    "\tif _, err := svc.resolveProviderConfig(\"\", taskInput.Config); err == nil {\n",
)
# 密码登录只保留给管理员（普通用户走 forapi 首次登录），
# 这两个测试要验证登录/登出的日志与事务，所以先把用户提为管理员。
patch(
    "backend/internal/handler/login_logs_test.go",
    "\trouter, db, _ := profileHandlerFixture(t)\n"
    "\trec := profileRequest(router, \"POST\", \"/api/auth/login\", `{\"username\":\"profile-user\",\"password\":\"old-password\"}`, \"\")\n",
    "\trouter, db, _ := profileHandlerFixture(t)\n"
    "\tif err := db.Model(&model.User{}).Where(\"id = ?\", \"profile-user\").Update(\"role\", model.UserRoleAdmin).Error; err != nil {\n\t\tt.Fatal(err)\n\t}\n"
    "\trec := profileRequest(router, \"POST\", \"/api/auth/login\", `{\"username\":\"profile-user\",\"password\":\"old-password\"}`, \"\")\n",
)
patch(
    "backend/internal/handler/login_logs_test.go",
    "\trouter, db, cookie := profileHandlerFixture(t)\n"
    "\tif err := db.Exec(`CREATE TRIGGER reject_success_login_log",
    "\trouter, db, cookie := profileHandlerFixture(t)\n"
    "\tif err := db.Model(&model.User{}).Where(\"id = ?\", \"profile-user\").Update(\"role\", model.UserRoleAdmin).Error; err != nil {\n\t\tt.Fatal(err)\n\t}\n"
    "\tif err := db.Exec(`CREATE TRIGGER reject_success_login_log",
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
# 身份标识：底层仍是通用 OAuth 登录（沿用 linuxdo 命名），界面上显示为 ForAPI
# v1.8.4 起文案收敛到 i18n，改 zh-CN 源语言文件即可覆盖全站。
patch(
    "web/src/components/layout/identity-provider-badge.tsx",
    'title={user.identityUsername ? `Linux.do · @${user.identityUsername}` : t("account.layout.identity.linuxdoUser")}',
    'title={user.identityUsername ? `ForAPI · @${user.identityUsername}` : t("account.layout.identity.linuxdoUser")}',
)
patch("web/src/components/layout/identity-provider-badge.tsx", '{compact ? "L" : "Linux.do"}', '{compact ? "F" : "ForAPI"}')
patch("web/src/i18n/messages/areas/account/zh-CN.ts", 'linuxdoUser: "Linux.do 用户",', 'linuxdoUser: "ForAPI 用户",')
patch("web/src/i18n/messages/zh-CN.ts", 'linuxdo: "使用 Linux.do 登录",', 'linuxdo: "使用 ForAPI 登录",')
patch("web/src/i18n/messages/zh-CN.ts", 'linuxdo: "使用 Linux.do 注册 / 登录",', 'linuxdo: "使用 ForAPI 注册 / 登录",')
# 系统更新会去 GitHub 拉上游版本并原地升级，与本仓库的补丁构建冲突：隐藏入口
patch(
    "web/src/pages/admin/components/admin-shell.tsx",
    '            { path: "/admin/settings/system-update", keywords: "系统 更新 升级 版本 部署", label: "系统更新", shortLabel: "更新", icon: <RefreshCw className="size-4" /> },\n',
    "",
)
# 上游这两个测试断言的还是改动前的状态：
# 品牌名默认值已改成 ForAPI 视频工作台，「系统更新」入口也已按上文隐藏。
patch(
    "web/test/workspace-sidebar-checkin.test.ts",
    '        expect(sidebarCheckinTitle("")).toBe("影策加油站");\n',
    '        expect(sidebarCheckinTitle("")).toBe("ForAPI 视频工作台加油站");\n',
)
patch(
    "web/test/admin-ui-regressions.test.ts",
    '    for (const path of ["/admin/logs", "/admin/login-logs", "/admin/users", "/admin/user-groups", "/admin/agent-lessons", "/admin/settings/system-update"]) {\n',
    '    for (const path of ["/admin/logs", "/admin/login-logs", "/admin/users", "/admin/user-groups", "/admin/agent-lessons"]) {\n',
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
