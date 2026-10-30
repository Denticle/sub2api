"""构建前补丁：把首页两处写死的品牌色换成 forapi 的令牌。

上游在 home/index.tsx 里硬编码了 `#FF9800`（橙色下划线）和 `#87CEFA`（天蓝高亮），
这两个不走 antd 主题也不走 CSS 变量，只能直接替换。用 CSS 变量而不是写死颜色，
这样深浅色模式各自取到对应的绿，不用维护两份。
"""
import pathlib

base = pathlib.Path(__file__).resolve().parent.parent / "web"
p = base / "src/pages/home/index.tsx"
s = p.read_text()

before = s
# 下划线：浅色 p-600 / 深色 p-500 —— 用 color-mix 不行，直接读变量
s = s.replace(
    'color="#FF9800"',
    'color="var(--primary)"',
)
# 高亮底色：橙色换掉，用同色系更淡的一档，避免和主按钮抢
s = s.replace(
    'color="#87CEFA"',
    'color="var(--primary)"',
)
p.write_text(s)

print("home/index.tsx 已替换:", s != before)
print("残留硬编码色:", "FF9800" in s or "87CEFA" in s)
