"""清理：patch-shell.py 重复执行后残留的「返回主站」块。

主站侧栏就在 iframe 外面，返回路径已经有了，画布内部不需要再放一个。
"""
import pathlib
import re

base = pathlib.Path(__file__).resolve().parent.parent / "web"
p = base / "src/components/layout/app-top-nav.tsx"
s = p.read_text()

# 去掉 href="/" 的那个返回链接 + 紧跟的分隔线
s = re.sub(
    r'\s*<a\s+href="/"\s+className="flex h-full shrink-0 items-center gap-1\.5.*?</a>\s*<span className="mx-3 h-5 w-px[^"]*" />',
    "",
    s,
    flags=re.S,
)

# 没有返回链接了，移除按钮的 ml-3
s = s.replace('className="ml-3 inline-flex size-8 shrink-0', 'className="inline-flex size-8 shrink-0')

p.write_text(s)
print("返回链接残留:", 'href="/"' in s and "返回主站" in s)
print("空行清理完成")
