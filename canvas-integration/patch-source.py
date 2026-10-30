"""构建前：源码补丁 —— 让 SPA 能在子路径下工作。"""
import pathlib
import re

base = pathlib.Path(__file__).resolve().parent.parent / "web"
if not base.exists():
    raise SystemExit(f"找不到 web 目录：{base}")

# 1) router.tsx —— React Router 7 的 createBrowserRouter 接受 basename，
#    从 Vite 的 BASE_URL 取，换挂载路径时不用改源码。
p = base / "src/router.tsx"
s = p.read_text()
if "const basename" not in s:
    s = s.replace(
        "export const router = createBrowserRouter([",
        'const basename = import.meta.env.BASE_URL.replace(/\\/$/, "");\n\n'
        "export const router = createBrowserRouter([",
        1,
    )
    idx = s.rindex("]);")
    s = s[:idx] + "], { basename });" + s[idx + 3:]
    p.write_text(s)
print("router.tsx:", "basename" in p.read_text())

# 2) index.html —— 绝对路径在子路径下会打到站点根
p = base / "index.html"
s = p.read_text()
s = s.replace('href="/logo.svg"', 'href="./logo.svg"')
s = s.replace('src="/config.js"', 'src="./config.js"')
p.write_text(s)
print("index.html:", '"./config.js"' in s)

# 3) plugin-loader —— 唯一一处硬编码的绝对请求
p = base / "src/lib/canvas/plugin-loader.ts"
s = p.read_text()
s = s.replace(
    'fetch("/plugins/index.json")',
    "fetch(`${import.meta.env.BASE_URL}plugins/index.json`)",
)
p.write_text(s)
print("plugin-loader.ts:", "BASE_URL}plugins" in p.read_text())

# 4) 兜底检查：还有没有漏掉的根绝对路径
hits = []
for f in (base / "src").rglob("*"):
    if f.suffix not in {".ts", ".tsx", ".css"}:
        continue
    for i, line in enumerate(f.read_text(errors="ignore").splitlines(), 1):
        if re.search(r'''["'`]/(src|assets|logo|config|plugins)/''', line):
            hits.append(f"{f.relative_to(base)}:{i}: {line.strip()[:100]}")
if hits:
    print(f"警告：仍有 {len(hits)} 处绝对路径引用")
    for h in hits[:20]:
        print("  ", h)
else:
    print("无遗漏的绝对路径引用")
