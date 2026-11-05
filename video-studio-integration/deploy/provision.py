#!/usr/bin/env python3
"""视频工作台初始化：品牌外观、forapi 配色、功能开关、ForAPI 登录、ForAPI 系统渠道。

可重复执行：每一步都是「读取当前值 → 覆盖需要的字段 → 写回」。
在 compose 内部网络里运行（直连 backend:8080，不经过 Cloudflare）：

  docker run --rm --network video-studio_internal --env-file /opt/video-studio/.env \
    -v /opt/video-studio/integration/deploy:/p:ro python:3.12-alpine python /p/provision.py

需要的环境变量（都在 /opt/video-studio/.env）：
  STUDIO_ADMIN_USERNAME / STUDIO_ADMIN_PASSWORD   工作台本地管理员（首次运行时创建）
  FORAPI_CLIENT_ID / FORAPI_CLIENT_SECRET         与 forapi 后台「视频工作台接入」一致
"""
import base64
import http.cookiejar
import json
import os
import sys
import urllib.error
import urllib.request
import uuid

BASE = os.environ.get("STUDIO_BACKEND", "http://backend:8080/api")
PUBLIC_URL = "https://video.forapi.cloud"
FORAPI_URL = "https://forapi.cloud"
BRAND = "ForAPI 视频工作台"

jar = http.cookiejar.CookieJar()
opener = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar))


def call(method, path, body=None, raw=None, content_type="application/json"):
    data = raw if raw is not None else (json.dumps(body).encode() if body is not None else None)
    req = urllib.request.Request(BASE + path, data=data, method=method)
    if data is not None:
        req.add_header("Content-Type", content_type)
    try:
        with opener.open(req, timeout=60) as resp:
            payload = json.loads(resp.read() or b"{}")
    except urllib.error.HTTPError as err:
        payload = json.loads(err.read() or b"{}")
        raise SystemExit(f"{method} {path} 失败：HTTP {err.code} {payload.get('msg')}")
    return payload.get("data")


def env(name):
    value = os.environ.get(name, "").strip()
    if not value:
        raise SystemExit(f"缺少环境变量 {name}")
    return value


def step(title):
    print(f"==> {title}", flush=True)


def login_admin():
    username, password = env("STUDIO_ADMIN_USERNAME"), env("STUDIO_ADMIN_PASSWORD")
    try:
        call("POST", "/auth/login", {"username": username, "password": password})
        step(f"管理员 {username} 已登录")
        return
    except SystemExit as err:
        if "HTTP 401" not in str(err):
            raise
    # 还没有任何账号时，第一个注册的账号就是管理员（补丁只放行这一次本地注册）
    call("POST", "/auth/register", {"username": username, "password": password, "displayName": "管理员", "acceptedTerms": True})
    step(f"已创建管理员 {username}")


def forapi_logo():
    """直接取主站正在使用的 logo，保证两边一致。"""
    req = urllib.request.Request(FORAPI_URL + "/api/v1/settings/public", headers={"User-Agent": "ForAPI-Studio-Provision/1.0", "Accept": "application/json"})
    with urllib.request.urlopen(req, timeout=30) as resp:
        logo = json.loads(resp.read())["data"]["site_logo"]
    header, encoded = logo.split(",", 1)
    mime = header.split(";")[0].removeprefix("data:")
    return base64.b64decode(encoded), mime


def upload_logo(slot, content, mime):
    boundary = uuid.uuid4().hex
    ext = mime.split("/")[-1].replace("svg+xml", "svg")
    body = (
        f"--{boundary}\r\nContent-Disposition: form-data; name=\"file\"; filename=\"forapi-logo.{ext}\"\r\n"
        f"Content-Type: {mime}\r\n\r\n"
    ).encode() + content + f"\r\n--{boundary}--\r\n".encode()
    resource = call("POST", f"/admin/settings/appearance/assets/{slot}", raw=body, content_type=f"multipart/form-data; boundary={boundary}")
    return resource["resource"]["id"]


# forapi 主站配色（frontend/src/theme/tokens-forapi.css）：绿色是动作色，绿底配深绿字
FORAPI_LIGHT = {
    "canvas": "#f8f8f6", "surface": "#ffffff", "surfaceSubtle": "#f0f0ed", "surfaceRaised": "#e6e6e2", "overlay": "#ffffff",
    "text": "#1f1f20", "textMuted": "#52525b", "border": "#e6e6e2",
    "control": "#ffffff", "controlHover": "#f0f0ed", "controlActive": "#e6e6e2", "controlBorder": "#d4d4d0", "controlFocus": "#16a34a",
    "controlDisabledBackground": "#f0f0ed", "controlDisabledForeground": "#a1a19c",
    "switchChecked": "#16a34a", "switchCheckedHover": "#15803d", "switchCheckedHandle": "#ffffff",
    "switchUnchecked": "#d4d4d0", "switchUncheckedHover": "#a1a19c", "switchUncheckedHandle": "#ffffff",
    "primary": "#16a34a", "primaryHover": "#15803d", "primaryActive": "#166534", "primaryForeground": "#ffffff",
    "selected": "#dcfce7", "selectedHover": "#bbf7d0", "selectedActive": "#86efac", "selectedForeground": "#15803d",
    "icon": "#3f3f42", "iconMuted": "#71717a", "iconActive": "#16a34a",
    "success": "#16a34a", "warning": "#ca8a04", "danger": "#dc2626", "dangerHover": "#b91c1c", "dangerActive": "#991b1b", "dangerForeground": "#ffffff",
    "info": "#a78bfa", "workspace": "#f8f8f6", "workspaceGrid": "#e6e6e2",
    "adminBackground": "#f8f8f6", "adminSurface": "#ffffff", "adminSubtle": "#f0f0ed", "adminStrong": "#1f1f20",
    "authBackground": "#f8f8f6", "authPanel": "#ffffff", "authCard": "#ffffff", "authAccent": "#16a34a", "authMuted": "#71717a",
}
FORAPI_DARK = {
    "canvas": "#141413", "surface": "#1a1a19", "surfaceSubtle": "#1f1f20", "surfaceRaised": "#2a2a2c", "overlay": "#1f1f20",
    "text": "#f4f4f3", "textMuted": "#a1a19c", "border": "#2a2a2c",
    "control": "#1f1f20", "controlHover": "#2a2a2c", "controlActive": "#3f3f42", "controlBorder": "#3f3f42", "controlFocus": "#22c55e",
    "controlDisabledBackground": "#1f1f20", "controlDisabledForeground": "#71717a",
    "switchChecked": "#22c55e", "switchCheckedHover": "#4ade80", "switchCheckedHandle": "#062014",
    "switchUnchecked": "#3f3f42", "switchUncheckedHover": "#52525b", "switchUncheckedHandle": "#f4f4f3",
    "primary": "#22c55e", "primaryHover": "#4ade80", "primaryActive": "#16a34a", "primaryForeground": "#062014",
    "selected": "#14301f", "selectedHover": "#183c26", "selectedActive": "#1d4a2e", "selectedForeground": "#4ade80",
    "icon": "#d4d4d0", "iconMuted": "#a1a19c", "iconActive": "#22c55e",
    "success": "#22c55e", "warning": "#eab308", "danger": "#ef4444", "dangerHover": "#f87171", "dangerActive": "#dc2626", "dangerForeground": "#ffffff",
    "info": "#a78bfa", "workspace": "#141413", "workspaceGrid": "#2a2a2c",
    "adminBackground": "#141413", "adminSurface": "#1a1a19", "adminSubtle": "#1f1f20", "adminStrong": "#f4f4f3",
    "authBackground": "#141413", "authPanel": "#1a1a19", "authCard": "#1a1a19", "authAccent": "#22c55e", "authMuted": "#a1a19c",
}


def forapi_skin(themes):
    """以系统默认皮肤为底复制一份，只换颜色；圆角、尺寸等组件参数沿用上游。"""
    base = next((t for t in themes if t.get("locked")), themes[0])
    skin = json.loads(json.dumps(base))
    skin.update({"id": "forapi", "name": "ForAPI", "description": "与 forapi.cloud 主站一致的配色", "locked": False})
    skin["tokens"]["light"].update(FORAPI_LIGHT)
    skin["tokens"]["dark"].update(FORAPI_DARK)
    for mode, palette in (("light", FORAPI_LIGHT), ("dark", FORAPI_DARK)):
        fill = skin["tokens"].get("buttons", {}).get(mode)
        if fill:
            fill.update({"mode": "solid", "start": palette["primary"], "end": palette["primary"], "hoverStart": palette["primaryHover"],
                         "hoverEnd": palette["primaryHover"], "activeStart": palette["primaryActive"], "activeEnd": palette["primaryActive"],
                         "foreground": palette["primaryForeground"]})
    return [t for t in themes if t.get("id") != "forapi"] + [skin]


def provision_appearance():
    current = call("GET", "/admin/settings/appearance")["setting"]
    content, mime = forapi_logo()
    logo_id = upload_logo("logo", content, mime)
    dark_logo_id = upload_logo("logo-dark", content, mime)
    setting = {k: v for k, v in current.items() if k not in ("public", "configured", "updatedBy", "createdAt", "updatedAt")}
    setting.update({
        "brandName": BRAND, "brandSlug": "forapi-video",
        "authHeroTitle": "用 ForAPI 账号，开始视频创作",
        "authHeroDescription": "画布、分镜、图片与视频生成在一个工作台完成，费用从你的 ForAPI 余额扣除。",
        "logoResourceId": logo_id, "darkLogoResourceId": dark_logo_id, "logoFrameEnabled": False,
        # 不使用默认的第三方登录背景视频
        "authVideoResourceId": "", "authVideoPosterResourceId": "", "authVideoAutoplay": False,
        "seoTitle": BRAND, "seoDescription": "ForAPI 视频工作台：AI 影视、短剧与营销视频创作。", "seoKeywords": "ForAPI,AI 视频,短剧,分镜,画布",
        "footerCopyright": "© ForAPI. All rights reserved.",
        "skinThemes": forapi_skin(current.get("skinThemes") or []), "skinId": "forapi",
    })
    canvas = dict(setting.get("canvas") or {})
    canvas["agentName"] = "视频助手"
    setting["canvas"] = canvas
    call("PATCH", "/admin/settings/appearance", setting)
    step("品牌外观与配色已更新")


def provision_features():
    current = call("GET", "/admin/settings/features")["features"]
    features = {k: v for k, v in current.items() if k not in ("configured", "updatedBy", "updatedAt")}
    features.update({
        "welcomeEnabled": False,          # 欢迎页含上游项目介绍和 GitHub 链接
        "creditsEnabled": False,          # 计费在 forapi 侧按用户自己的密钥完成
        "customChannelsEnabled": False,   # 只能用 ForAPI 渠道
        "pluginCenterEnabled": False,
        "systemPluginsVisibleToUsers": False,
        "inspirationSourcesVisible": False,
        "timelineTranscriptionEnabled": False,
        "shortDramaEnabled": True, "taskCenterEnabled": True,
        "imageWorkbenchEnabled": True, "videoWorkbenchEnabled": True,
    })
    call("PATCH", "/admin/settings/features", features)
    step("功能开关已更新")


def provision_login():
    current = call("GET", "/admin/settings/registration")
    call("PATCH", "/admin/settings/registration", {"enabled": True, "agreementTitle": current.get("agreementTitle"), "agreementContent": current.get("agreementContent")})
    call("PATCH", "/admin/settings/linuxdo", {
        "enabled": True,
        "clientId": env("FORAPI_CLIENT_ID"), "clientSecret": env("FORAPI_CLIENT_SECRET"),
        "authorizationUrl": FORAPI_URL + "/oauth/authorize",
        "tokenUrl": FORAPI_URL + "/oauth/token",
        "userInfoUrl": FORAPI_URL + "/oauth/userinfo",
        "redirectUrl": PUBLIC_URL + "/api/auth/linuxdo/callback",
        "scopes": ["profile"],
        "clientAuthMethod": "client_secret_basic",
        "subjectField": "sub", "usernameField": "username", "displayNameField": "name",
        "emailField": "email", "avatarField": "picture",
    })
    step("ForAPI 登录已启用（注册仅限 ForAPI 首次登录）")


def provision_channel():
    channels = call("GET", "/admin/channels") or []
    if isinstance(channels, dict):
        channels = channels.get("channels") or channels.get("items") or []
    existing = next((c for c in channels if c.get("name") == "ForAPI"), None)
    body = {
        "name": "ForAPI", "baseUrl": FORAPI_URL + "/v1",
        # 占位值：生成时换成发起用户自己的 forapi 密钥（见 overlay forapi_bridge.go）
        "apiKey": "forapi-user-key",
        "enabled": True,
    }
    if existing:
        call("PATCH", f"/admin/channels/{existing['id']}", body)
        step("ForAPI 系统渠道已更新")
    else:
        call("POST", "/admin/channels", body)
        step("已创建 ForAPI 系统渠道（模型需在后台添加，见 UPGRADE.md）")


def main():
    login_admin()
    provision_appearance()
    provision_features()
    provision_login()
    provision_channel()
    print("完成")


if __name__ == "__main__":
    try:
        main()
    except SystemExit as err:
        print(err, file=sys.stderr)
        sys.exit(1)
