"""构建前补丁：全局移除 Agent（Codex 接入）功能。

它在画布里有三个入口，全部摘掉：
  1. layouts/user-layout.tsx —— 右侧撑开的 AgentPanel
  2. components/canvas/canvas-top-bar.tsx —— 顶栏的「Codex 未连接」状态 + Agent 按钮
  3. pages/canvas/project.tsx —— 进画布时自动展开 Agent 面板的副作用

只摘渲染入口和副作用，不删组件文件本身：
  - 删文件会牵连 stores/use-agent-store、agent-url-bootstrap、canvas-agent-ops
    一整条依赖链，上游一更新就冲突
  - 留着文件不影响任何东西——没被引用的模块不会进包（Vite 会摇掉）
这样上游同步时这个补丁的冲突面最小。
"""
import pathlib
import re

base = pathlib.Path(__file__).resolve().parent.parent / "web"


def patch(rel: str, fn):
    p = base / rel
    s = p.read_text()
    out = fn(s)
    p.write_text(out)
    return out != s


# 1) 布局层：整个 AgentPanel 组件从渲染树上摘掉
def fix_layout(s: str) -> str:
    if "AgentPanel" not in s:
        return s
    s = s.replace('import { AgentPanel } from "@/components/agent/agent-panel";\n', "")
    s = re.sub(r"\n\s*<AgentPanel />", "", s)
    return s


print("user-layout.tsx:", patch("src/layouts/user-layout.tsx", fix_layout))


# 2) 画布顶栏：去掉状态文字和 Agent 按钮
def fix_top_bar(s: str) -> str:
    if "CompactAgentStatus" not in s and "<Bot" not in s:
        return s

    # 顶栏标题右侧的状态
    s = re.sub(r"\s*<CompactAgentStatus[^>]*/>", "", s)

    # 右侧的 Agent 按钮（从 <Button ... icon={<Bot .../>} ... >Agent</Button> 整段）
    s = re.sub(
        r"\s*<Button\s*\n\s*type=\"text\"\s*\n\s*className=\"!h-10 !rounded-xl !px-3 !font-medium\"\s*\n"
        r"\s*style=\{\{ background: agentOpen[^}]*\}\}\s*\n"
        r"\s*icon=\{<Bot[^/]*/>\}\s*\n"
        r"\s*onClick=\{onToggleAgent\}\s*\n\s*>\s*Agent\s*</Button>",
        "",
        s,
    )

    # 冷掉的 CompactAgentStatus 函数定义
    s = re.sub(r"\nfunction CompactAgentStatus\([^)]*\)[^{]*\{.*?\n\}\n", "\n", s, flags=re.S)

    # 不再用到的 import
    s = s.replace('import { Bot, ', "import { ").replace(", Bot, ", ", ").replace(", Bot }", " }")

    return s


print("canvas-top-bar.tsx:", patch("src/components/canvas/canvas-top-bar.tsx", fix_top_bar))


# 3) project.tsx：去掉自动展开面板的副作用
def fix_project(s: str) -> str:
    if "openAgentPanel" not in s:
        return s
    # 进画布时自动 openPanel 的那一行
    s = re.sub(
        r"\n\s*if \(!searchParams\.has\(\"agentUrl\"\)[^\n]*openAgentPanel\(\);\n",
        "\n",
        s,
    )
    return s


print("project.tsx:", patch("src/pages/canvas/project.tsx", fix_project))
