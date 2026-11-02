"""构建后：产物补丁 —— 延迟应用脚本，注入 bootstrap。

为什么必须在构建后做：Vite 会重建 <script> 标签并丢掉自定义属性，
写在源码上的 data-* 标记到不了产物。
"""
import hashlib
import pathlib
import re

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

m = re.search(r'data-studio-app-src="([^"]+)/assets/index-[^"]+\.js"', s)
if not m:
    raise SystemExit("index.html 里找不到应用脚本，无法推出 base")
base = m.group(1)  # 例如 /studio

# bootstrap 放进 assets/ 并按内容加哈希：
#   - 文件名随内容变，nginx 和 Cloudflare 可以把它当不可变资源长缓存，
#     改了 bootstrap 也不会有人拿到旧的（之前不带哈希时 Cloudflare 缓存了 4 小时）
#   - 用绝对路径：相对路径 ./bootstrap.js 在 /studio/canvas/xxx 这种深层地址刷新时
#     会解析成 /studio/canvas/bootstrap.js，拿回来的是 index.html
data = (here / "bootstrap.js").read_bytes()
name = f"bootstrap-{hashlib.sha256(data).hexdigest()[:10]}.js"
(dist / "assets" / name).write_bytes(data)
(dist / "bootstrap.js").unlink(missing_ok=True)

s = re.sub(r'\n\s*<script src="[^"]*bootstrap[^"]*\.js"></script>', "", s)
s = s.replace('<script src="./config.js"></script>', f'<script src="{base}/config.js"></script>')
s = s.replace(f'<script src="{base}/config.js"></script>', f'<script src="{base}/config.js"></script>\n        <script src="{base}/assets/{name}"></script>', 1)
s = s.replace('href="./logo.svg"', f'href="{base}/logo.svg"')

p.write_text(s)

final = p.read_text()
print("应用脚本已延迟:", "data-studio-app-src" in final)
print("bootstrap 已注入:", f"{base}/assets/{name}" in final)
print("bootstrap.js 已就位:", (dist / "assets" / name).exists())
