"""构建前补丁：把无限画布的视觉 token 换成 forapi 的设计令牌。

两块要改：
  1. lib/app-theme.ts  —— antd 的 ThemeConfig（组件库那层：按钮、下拉、表格、选择器）
  2. styles/forapi.css —— 画布自己的 shadcn/tailwind CSS 变量（背景、边框、文字、圆角、字体）

只换 token，不动组件结构和交互——那是上游的事，改了就跟不动上游了。
"""
import pathlib

here = pathlib.Path(__file__).resolve().parent
base = here.parent / "web"

APP_THEME = '''import type { ThemeConfig } from "antd";
import { theme as antdTheme } from "antd";

/**
 * forapi 设计令牌 —— 与主站 frontend/src/theme/tokens-forapi.css 对齐。
 * 浅色主色用 p-600（#16a34a），深色用 p-500（#22c55e）；绿色是动作色，
 * 所以绿底上的文字用 --primary-ink 那种深绿，而不是白色（对比度约 7.9:1）。
 */
const neutral = {
    light: {
        primary: "#16a34a",
        primaryHover: "#15803d",
        primaryText: "#062014",
        elevatedBg: "#ffffff",
        itemHoverBg: "rgba(28, 28, 25, 0.05)",
        itemSelectedBg: "rgba(22, 163, 74, 0.10)",
        itemSelectedHoverBg: "rgba(22, 163, 74, 0.16)",
        itemText: "#1f1f20",
        tableSelectedBg: "rgba(22, 163, 74, 0.06)",
        tableSelectedHoverBg: "rgba(22, 163, 74, 0.10)",
    },
    dark: {
        primary: "#22c55e",
        primaryHover: "#4ade80",
        primaryText: "#062014",
        elevatedBg: "#1a1a19",
        itemHoverBg: "rgba(255, 255, 255, 0.045)",
        itemSelectedBg: "rgba(34, 197, 94, 0.14)",
        itemSelectedHoverBg: "rgba(34, 197, 94, 0.20)",
        itemText: "#f4f4f3",
        tableSelectedBg: "rgba(34, 197, 94, 0.10)",
        tableSelectedHoverBg: "rgba(34, 197, 94, 0.14)",
    },
};

/** forapi 的字体栈：界面用 Inter，数字用 JetBrains Mono。 */
const FONT_FAMILY =
    'Inter, "Inter Variable", -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif';
const FONT_FAMILY_CODE = '"JetBrains Mono", "JetBrains Mono Variable", ui-monospace, SFMono-Regular, Menlo, monospace';

export function getAntThemeConfig(dark: boolean): ThemeConfig {
    const color = dark ? neutral.dark : neutral.light;

    return {
        algorithm: dark ? antdTheme.darkAlgorithm : antdTheme.defaultAlgorithm,
        cssVar: { key: dark ? "forapi-dark" : "forapi-light" },
        token: {
            colorPrimary: color.primary,
            colorInfo: color.primary,
            colorSuccess: color.primary,
            colorLink: color.primary,
            colorLinkHover: color.primaryHover,
            colorLinkActive: color.primary,
            colorTextLightSolid: color.primaryText,
            colorBgElevated: color.elevatedBg,
            controlItemBgHover: color.itemHoverBg,
            controlItemBgActive: color.itemSelectedBg,
            controlItemBgActiveHover: color.itemSelectedHoverBg,
            // forapi：8px 小圆角 / 12px 大圆角，控制线 1px 发丝
            borderRadius: 8,
            borderRadiusLG: 12,
            borderRadiusSM: 6,
            fontFamily: FONT_FAMILY,
            fontFamilyCode: FONT_FAMILY_CODE,
            lineWidth: 1,
        },
        components: {
            Button: {
                primaryShadow: "none",
                defaultShadow: "none",
                fontWeight: 500,
            },
            Dropdown: {
                colorBgElevated: color.elevatedBg,
                colorText: color.itemText,
                controlItemBgHover: color.itemHoverBg,
                controlItemBgActive: color.itemSelectedBg,
                controlItemBgActiveHover: color.itemSelectedHoverBg,
            },
            Menu: {
                popupBg: color.elevatedBg,
                itemActiveBg: color.itemSelectedBg,
                itemHoverBg: color.itemHoverBg,
                itemSelectedBg: color.itemSelectedBg,
                itemSelectedColor: color.itemText,
                darkPopupBg: neutral.dark.elevatedBg,
                darkItemHoverBg: neutral.dark.itemHoverBg,
                darkItemSelectedBg: neutral.dark.itemSelectedBg,
                darkItemSelectedColor: neutral.dark.itemText,
            },
            Select: {
                optionActiveBg: color.itemHoverBg,
                optionSelectedBg: color.itemSelectedBg,
                optionSelectedColor: color.itemText,
            },
            Table: {
                rowSelectedBg: color.tableSelectedBg,
                rowSelectedHoverBg: color.tableSelectedHoverBg,
            },
            Card: {
                borderRadiusLG: 12,
                paddingLG: 20,
            },
            Modal: {
                borderRadiusLG: 12,
            },
        },
    };
}
'''

p = base / "src/lib/app-theme.ts"
p.write_text(APP_THEME)
print("app-theme.ts 已替换")

CSS = '''/* forapi 主题覆盖
   --------------------------------------------------------------------------
   无限画布自己用 shadcn + tailwind 的 CSS 变量（oklch）。这里把底层的
   background / foreground / border / primary / radius / 字体整体换成 forapi
   的令牌，让画布的画布区、面板、按钮和主站是同一套颜色。

   只覆盖变量，不碰组件结构。深色值按 tokens-forapi.css 的 d-* 阶梯。
*/

:root {
    --radius: 8px;

    /* 面：n-50 页面底 / 白卡片 */
    --background: #f8f8f6;
    --foreground: #1f1f20;
    --card: #ffffff;
    --card-foreground: #1f1f20;
    --popover: #ffffff;
    --popover-foreground: #1f1f20;

    /* 动作色：p-600 */
    --primary: #16a34a;
    --primary-foreground: #062014;

    --secondary: #f0f0ed;
    --secondary-foreground: #27272a;
    --muted: #f0f0ed;
    --muted-foreground: #71717a;
    --accent: #f0f0ed;
    --accent-foreground: #27272a;
    --destructive: #ef4444;

    /* 发丝线：n-200 实色，比 oklch 半透明更接近主站的 1px 描边观感 */
    --border: #e6e6e2;
    --input: #e6e6e2;
    --ring: #16a34a;

    --sidebar: #f8f8f6;
    --sidebar-foreground: #1f1f20;
    --sidebar-primary: #16a34a;
    --sidebar-primary-foreground: #062014;
    --sidebar-accent: #f0f0ed;
    --sidebar-accent-foreground: #27272a;
    --sidebar-border: #e6e6e2;
    --sidebar-ring: #16a34a;

    --font-sans: Inter, "Inter Variable", -apple-system, BlinkMacSystemFont,
        "Segoe UI", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
    --font-mono: "JetBrains Mono", "JetBrains Mono Variable", ui-monospace,
        SFMono-Regular, Menlo, monospace;
}

.dark {
    --background: #141413;
    --foreground: #f4f4f3;
    --card: #1f1f1e;
    --card-foreground: #f4f4f3;
    --popover: #1a1a19;
    --popover-foreground: #f4f4f3;

    --primary: #22c55e;
    --primary-foreground: #062014;

    --secondary: #272726;
    --secondary-foreground: #e4e4e3;
    --muted: #272726;
    --muted-foreground: #a1a19c;
    --accent: #272726;
    --accent-foreground: #f4f4f3;
    --destructive: #ef4444;

    --border: #3a3a38;
    --input: #3a3a38;
    --ring: #22c55e;

    --sidebar: #1a1a19;
    --sidebar-foreground: #f4f4f3;
    --sidebar-primary: #22c55e;
    --sidebar-primary-foreground: #062014;
    --sidebar-accent: #272726;
    --sidebar-accent-foreground: #f4f4f3;
    --sidebar-border: #3a3a38;
    --sidebar-ring: #22c55e;
}
'''

css_path = base / "src/styles/forapi.css"
css_path.write_text(CSS)
print("forapi.css 已写入")

# 在 main.tsx 里最后引入，保证覆盖 globals.css 的变量
p = base / "src/main.tsx"
s = p.read_text()
if "forapi.css" not in s:
    s = s.replace(
        'import "./styles/globals.css";',
        'import "./styles/globals.css";\nimport "./styles/forapi.css";',
        1,
    )
    p.write_text(s)
print("main.tsx:", "forapi.css" in p.read_text())

# 字体栈：它原本在 main.tsx 里硬写 SF Pro，换成 forapi 的 Inter
p = base / "src/main.tsx"
s = p.read_text()
s = s.replace(
    """document.body.style.fontFamily = '"SF Pro Display","SF Pro Text","PingFang SC","Microsoft YaHei","Helvetica Neue",sans-serif';""",
    """document.body.style.fontFamily = 'Inter,"Inter Variable",-apple-system,BlinkMacSystemFont,"PingFang SC","Microsoft YaHei","Helvetica Neue",sans-serif';""",
)
p.write_text(s)
print("font family:", "Inter Variable" in p.read_text())
