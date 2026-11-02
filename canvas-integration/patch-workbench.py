"""构建前补丁：把画布改成主站「工作台」里的一页。

1. 功能切换交给主站侧栏
   - 去掉画布自己的顶部导航（overlay/src/layouts/user-layout.tsx）
   - 去掉画布首页，"/" 直接进画布列表
   - 主站点侧栏 -> postMessage 让画布内部换页；画布内部换页 -> 报回主站，侧栏高亮跟上
2. 画布项目页顶栏只留「节点插件、配置、快捷键」
   - overlay/src/components/layout/user-status-actions.tsx
   - 下拉菜单里去掉「首页」「文档」两项
   - 明暗、语言跟主站走（overlay 里的 theme store 和 i18n）
3. 配置弹窗
   - 去掉「本地代理」页，并强制关闭代理
   - 「渠道」页换成「密钥」：只能从用户已有的密钥里选（overlay 里的 config-forapi-keys）
   - 去掉配置文件导入 / 导出：导出会把 Key 明文写进文件，导入能绕过密钥选择注入任意渠道
   - 去掉从 URL 参数注入渠道（overlay/src/components/layout/client-root-init.tsx）

overlay/ 下的文件整体覆盖到上游对应位置；这里只做必须改在原文件里的小改动。
每个替换都校验锚点存在——上游改了结构时直接失败，而不是悄悄漏改。
"""
import pathlib
import re
import shutil

here = pathlib.Path(__file__).resolve().parent
web = here.parent / "web"


def edit(rel: str, pairs):
    p = web / rel
    s = p.read_text()
    for old, new in pairs:
        if isinstance(old, re.Pattern):
            s, n = old.subn(new, s, count=1)
            if not n:
                raise SystemExit(f"{rel}: 锚点未匹配 {old.pattern[:60]!r}")
        else:
            if old not in s:
                raise SystemExit(f"{rel}: 锚点未匹配 {old[:60]!r}")
            s = s.replace(old, new, 1)
    p.write_text(s)
    print(f"{rel}: ok")


# ---- overlay：整文件覆盖 ----
overlay = here / "overlay"
for src in sorted(overlay.rglob("*")):
    # 跳过 macOS 打包带进来的 ._* 元数据文件
    if src.is_file() and not src.name.startswith("._"):
        dst = web / src.relative_to(overlay)
        dst.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(src, dst)
        print(f"overlay: {dst.relative_to(web)}")

# ---- 1. 去掉画布首页 ----
edit(
    "src/router.tsx",
    [
        ('import { createBrowserRouter, Outlet } from "react-router-dom";', 'import { createBrowserRouter, Navigate, Outlet } from "react-router-dom";'),
        ('import HomePage from "@/pages/home";\n', ""),
        ('{ path: "/", element: <HomePage /> },', '{ path: "/", element: <Navigate to="/canvas" replace /> },'),
    ],
)

# ---- 1. 切换功能不卸载页面（任务在后台跑完）----
# 路由加显式 id：KeepAliveOutlet 按 id 缓存，快捷键守卫按 id 判断自己是不是当前页
ROUTE_IDS = [
    ('{ path: "/", element: <Navigate to="/canvas" replace /> },', '{ id: "root-redirect", path: "/", element: <Navigate to="/canvas" replace /> },'),
    ('{ path: "/image", element: <ImagePage /> },', '{ id: "image", path: "/image", element: <ImagePage /> },'),
    ('{ path: "/video", element: <VideoPage /> },', '{ id: "video", path: "/video", element: <VideoPage /> },'),
    ('{ path: "/assets", element: <AssetsPage /> },', '{ id: "assets", path: "/assets", element: <AssetsPage /> },'),
    ('{ path: "/prompts", element: <PromptsPage /> },', '{ id: "prompts", path: "/prompts", element: <PromptsPage /> },'),
    ('{ path: "/canvas", element: <CanvasPage /> },', '{ id: "canvas-list", path: "/canvas", element: <CanvasPage /> },'),
    ('{ path: "/canvas/:id", element: <CanvasProjectPage /> },', '{ id: "canvas-project", path: "/canvas/:id", element: <CanvasProjectPage /> },'),
    ('{ path: "/config", element: <ConfigPage /> },', '{ id: "config", path: "/config", element: <ConfigPage /> },'),
]
edit(
    "src/router.tsx",
    [
        ('import { createBrowserRouter, Navigate, Outlet } from "react-router-dom";', 'import { createBrowserRouter, Navigate } from "react-router-dom";'),
        ('import UserLayout from "@/layouts/user-layout";', 'import { KeepAliveOutlet } from "@/components/layout/keep-alive-outlet";\nimport UserLayout from "@/layouts/user-layout";'),
        ("<Outlet />", "<KeepAliveOutlet />"),
        *ROUTE_IDS,
    ],
)

# 画布项目页的全局快捷键挂在 window 上，页面隐藏了也还在监听。
# 只在画布项目页是当前页时响应——否则在生图页按 Delete / ⌘Z 会改后台的画布。
GUARD = 'if (!isRouteActive("canvas-project")) return;'
IMPORT = 'import { isRouteActive } from "@/components/layout/keep-alive-outlet";\n'
edit(
    "src/pages/canvas/project.tsx",
    [
        ('import { useNavigate, useParams, useSearchParams } from "react-router-dom";\n', 'import { useNavigate, useParams, useSearchParams } from "react-router-dom";\n' + IMPORT),
        (
            "        const handleKeyDown = (event: KeyboardEvent) => {\n            const target = event.target instanceof Element ? event.target : null;\n            if (event.target instanceof HTMLInputElement",
            "        const handleKeyDown = (event: KeyboardEvent) => {\n            " + GUARD + "\n            const target = event.target instanceof Element ? event.target : null;\n            if (event.target instanceof HTMLInputElement",
        ),
        (
            '        const exit = (event: KeyboardEvent) => {\n            if (event.key !== "Escape") return;',
            '        const exit = (event: KeyboardEvent) => {\n            ' + GUARD + '\n            if (event.key !== "Escape") return;',
        ),
    ],
)
edit(
    "src/components/canvas/infinite-canvas.tsx",
    [
        ('import React, { useEffect, useRef, useState } from "react";\n', 'import React, { useEffect, useRef, useState } from "react";\n' + IMPORT),
        (
            '        const handleKeyDown = (event: KeyboardEvent) => {\n            if (event.key === "Control") setIsControlPressed(true);',
            '        const handleKeyDown = (event: KeyboardEvent) => {\n            ' + GUARD + '\n            if (event.key === "Control") setIsControlPressed(true);',
        ),
        (
            "        const handleKeyUp = (event: KeyboardEvent) => {\n            if (event.code === \"Space\") {",
            "        const handleKeyUp = (event: KeyboardEvent) => {\n            " + GUARD + "\n            if (event.code === \"Space\") {",
        ),
    ],
)

# ---- 1. 与主站互通路由 ----
EMBED = '''
// 嵌在主站「工作台」里时：
//   主站 -> 画布：{ type: "studio:navigate", path }，侧栏切功能时让画布在内部换页（不整页重载）
//   画布 -> 主站：{ type: "studio:route", path }，画布里自己换页时报回去，主站侧栏高亮跟上
// 只认同源、且来自父窗口的消息。
if (window.parent !== window) {
    window.addEventListener("message", (event) => {
        if (event.origin !== window.location.origin || event.source !== window.parent) return;
        const data = event.data;
        if (data && data.type === "studio:navigate" && typeof data.path === "string" && data.path.startsWith("/")) {
            void router.navigate(data.path);
        }
    });
    // router.state.location.pathname 带着 basename（/studio/canvas），主站只认 /canvas 这一层，报之前去掉
    const base = import.meta.env.BASE_URL.replace(/\\/$/, "");
    const report = (pathname: string) => {
        const path = base && pathname.startsWith(base) ? pathname.slice(base.length) || "/" : pathname;
        window.parent.postMessage({ type: "studio:route", path }, window.location.origin);
    };
    report(router.state.location.pathname);
    router.subscribe((state) => report(state.location.pathname));
}
'''
edit(
    "src/main.tsx",
    [("initAnalytics();", "initAnalytics();\n" + EMBED)],
)

# ---- 2. 画布项目页顶栏 ----
edit(
    "src/components/canvas/canvas-top-bar.tsx",
    [
        (re.compile(r'\n\s*\{ key: "home",[^\n]*onClick: onHome \},'), ""),
        (re.compile(r'\n\s*\{ key: "docs",[^\n]*\},'), ""),
        # Agent 按钮已由 patch-remove-agent 去掉，它前面那条分隔线跟着去
        (re.compile(r'\n\s*<span className="h-6 w-px" style=\{\{ background: theme\.toolbar\.border \}\} />'), ""),
    ],
)

# 明暗跟主站走：首屏防闪烁脚本读主站的 key
edit(
    "index.html",
    [
        (
            'var s = JSON.parse(localStorage.getItem("infinite-canvas:theme_store") || "{}");\n                var t = s.state && s.state.theme === "light" ? "light" : "dark";',
            'var h = localStorage.getItem("theme");\n                var t = h === "dark" || h === "light" ? h : (window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light");',
        )
    ],
)

# ---- 3. 配置弹窗 ----
edit(
    "src/components/layout/app-config-modal.tsx",
    [
        ('import { ChannelEditorDrawer } from "@/components/layout/channel-editor-drawer";\n', ""),
        (
            'import { ConfigLocalProxy } from "@/components/layout/config-local-proxy";\n',
            'import { ConfigForapiKeys, ConfigForapiKeysTabLabel } from "@/components/layout/config-forapi-keys";\n',
        ),
        # 导入 / 导出整行工具条
        (
            re.compile(r'\n\s*<div className="flex flex-wrap items-center justify-between gap-3 border-b border-stone-200 pb-3 dark:border-stone-800">.*?\n            </div>\n(?=\s*<Tabs)', re.S),
            "\n",
        ),
        # 渠道页 -> 密钥页
        (
            re.compile(r'label: t\("config\.tabs\.channels"\),\n(\s*)children: \(\n.*?\n\s*\),\n(\s*)\},\n(\s*)\{\n\s*key: "local-proxy",\n\s*label: t\("config\.tabs\.localProxy"\),\n\s*children: <ConfigLocalProxy />,\n\s*\},', re.S),
            r'label: <ConfigForapiKeysTabLabel />,\n\1children: <ConfigForapiKeys />,\n\2},',
        ),
        (re.compile(r'\n\s*<ChannelEditorDrawer [^\n]*/>'), ""),
        # 「完成」判断渠道是否就绪：只有一条本站渠道，有 Key 有模型就算好
        ("channel.baseUrl.trim() && channel.apiKey.trim() && channel.models.length", "channel.apiKey.trim() && channel.models.length"),
    ],
)

# 弹窗副标题原来写着「渠道聚合」，渠道已经没了
edit(
    "src/i18n/locales/zh-CN.ts",
    [(re.compile(r'modalDescription: "[^"]*"'), 'modalDescription: "密钥、默认模型、同步与本地存储"')],
)
edit(
    "src/i18n/locales/en-US.ts",
    [(re.compile(r'modalDescription: "[^"]*"'), 'modalDescription: "API key, default models, sync and local storage"')],
)

# 本站生图开了 JSON keepalive（防 Cloudflare 120 秒断开）：心跳一旦发出，状态码已经是 200，
# 之后上游再失败只能在 200 里返回 { error: { message } }。上游代码只认 HTTP 错误码和 code/msg，
# 碰到这种 200 会报一句「没有返回图片」，把真正的原因吞掉——这里先读 error 字段。
edit(
    "src/services/api/image.ts",
    [
        (
            '    if (typeof payload.code === "number" && payload.code !== 0) {\n        throw new Error(payload.msg || apiText("requestFailed"));\n    }',
            '    if (typeof payload.code === "number" && payload.code !== 0) {\n        throw new Error(payload.msg || apiText("requestFailed"));\n    }\n'
            '    const bodyError = (payload as Record<string, unknown>).error;\n'
            '    if (bodyError && !payload.data) {\n'
            '        const message = typeof bodyError === "string" ? bodyError : (bodyError as { message?: unknown }).message;\n'
            '        throw new Error(typeof message === "string" && message ? message : apiText("requestFailed"));\n'
            '    }',
        ),
    ],
)

# 历史配置里开着的本地代理一律关掉——请求只走本站
edit(
    "src/stores/use-config-store.ts",
    [("proxyEnabled: Boolean(config.proxyEnabled),", "proxyEnabled: false,")],
)

print("patch-workbench 完成")
