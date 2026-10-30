"""构建后：产物补丁 —— 延迟应用脚本，注入 bootstrap。

为什么必须在构建后做：Vite 会重建 <script> 标签并丢掉自定义属性，
写在源码上的 data-* 标记到不了产物。
"""
import pathlib
import re
import shutil

here = pathlib.Path(__file__).resolve().parent
dist = here.parent / "web" / "dist"
if not dist.exists():
    raise SystemExit(f"找不到构建产物：{dist}")

p = dist / "index.html"
s = p.read_text()

# 关键：不能只摘掉 src。<script type="module" src> 在 <head> 里是异步执行的，
# bootstrap 在 <body> 末尾，等它跑起来时应用可能已经开始执行，登录守卫就形同虚设。
# 所以整体换成不会执行的数据标签，由 bootstrap 决定何时注入。
s = re.sub(
    r'<script type="module"[^>]*src="([^"]+index-[^"]+\.js)"[^>]*></script>',
    r'<script type="application/json" data-studio-app-src="\1"></script>',
    s,
)

if "./bootstrap.js" not in s:
    s = s.replace(
        '<script src="./config.js"></script>',
        '<script src="./config.js"></script>\n        <script src="./bootstrap.js"></script>',
    )

p.write_text(s)

shutil.copyfile(here / "bootstrap.js", dist / "bootstrap.js")

final = p.read_text()
print("应用脚本已延迟:", "data-studio-app-src" in final)
print("bootstrap 已注入:", "./bootstrap.js" in final)
print("bootstrap.js 已就位:", (dist / "bootstrap.js").exists())
