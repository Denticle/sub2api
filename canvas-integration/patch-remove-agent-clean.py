"""清掉 app-top-nav.tsx 里因移除 Agent 按钮而失效的状态与副作用。

只删确实没人用的东西：useAgentStore 的六个选择器、autoConnectRef、
自动连接的 useEffect，以及因此空转的 import。
"""
import pathlib
import re

p = pathlib.Path(__file__).resolve().parent.parent / "web/src/components/layout/app-top-nav.tsx"
s = p.read_text()

for name in (
    "agentToken",
    "agentEnabled",
    "agentConnected",
    "connectAgent",
    "togglePanel",
    "panelOpen",
):
    s = re.sub(rf"\n\s*const {name} = useAgentStore[^\n]*", "", s)

s = re.sub(r"\n\s*const autoConnectRef = useRef\(false\);", "", s)

# 自动连接的副作用整块
s = re.sub(
    r"\n\s*useEffect\(\(\) => \{\n\s*if \(autoConnectRef[^{]*\{[^}]*\}\n[^}]*?\}, \[[^\]]*\]\);",
    "",
    s,
    flags=re.S,
)

s = s.replace('import { useAgentStore } from "@/stores/use-agent-store";\n', "")
s = s.replace('import { useEffect, useRef, useState } from "react";', 'import { useState } from "react";')
s = s.replace('import { useRef, useState } from "react";', 'import { useState } from "react";')

p.write_text(s)
print("useAgentStore:", "useAgentStore" in s)
print("useEffect:", "useEffect" in s)
print("useRef:", "useRef" in s)
print("autoConnectRef:", "autoConnectRef" in s)
