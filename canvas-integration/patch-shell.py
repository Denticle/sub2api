"""构建前补丁：把画布顶栏降级成「主站里的二级导航」。

割裂感的来源不是配色，是壳：画布自带 logo + 品牌名 + GitHub + 版本号 +
自己的用户菜单，进来像换了个站，也没有回主站的路。

主站现在用同源 iframe 把它装在 /studio 下，侧栏和顶栏都在外面。所以这里：
  - 去掉画布自己的品牌、GitHub 链接、版本号弹窗（那是上游作者的信息）
  - 左侧只留它的工具导航（画布 / 生图 / 视频 / 提示词 / 素材 / 设置），
    作为主站里的二级导航
  - 右侧换成主站登录用户的身份
"""
import pathlib

base = pathlib.Path(__file__).resolve().parent.parent / "web"
p = base / "src/components/layout/app-top-nav.tsx"
s = p.read_text()

HELPER = '''
/** 读主站存下来的登录用户，用于在页头显示身份。 */
function readCurrentUser(): { name: string; email: string } | null {
    try {
        const raw = localStorage.getItem("auth_user");
        if (!raw) return null;
        const u = JSON.parse(raw);
        if (!u) return null;
        const email = String(u.email || "");
        const name = String(u.username || u.email || "");
        return name || email ? { name, email } : null;
    } catch {
        return null;
    }
}
'''

if "readCurrentUser" not in s:
    s = s.replace("export function AppTopNav() {", HELPER + "\nexport function AppTopNav() {", 1)

if "const currentUser = readCurrentUser()" not in s:
    s = s.replace(
        "    const { t } = useTranslation();",
        "    const { t } = useTranslation();\n    const currentUser = readCurrentUser();",
        1,
    )

# 品牌区整块去掉，只留工具导航
import re

old_brand = re.search(
    r'<Link to="/" className="flex h-full shrink-0 items-center gap-2[^"]*">.*?</Link>',
    s,
    re.S,
)
if old_brand:
    s = s.replace(old_brand.group(0), "", 1)
    print("品牌区已移除")
else:
    print("警告：品牌区未匹配")

# 导航左边距修正：没有品牌了，不需要 ml-8
s = s.replace('className="hide-scrollbar ml-8 hidden h-14', 'className="hide-scrollbar hidden h-14')

# 右侧加用户身份
if "currentUser ?" not in s:
    anchor = '''                            <Tooltip title={t(panelOpen ? "topNav.closeAgent" : "topNav.openAgent")}>'''
    chip = '''                            {currentUser ? (
                                <span
                                    className="inline-flex max-w-[200px] items-center gap-2 rounded-md px-1 text-xs text-stone-600 dark:text-stone-300"
                                    title={currentUser.email || currentUser.name}
                                >
                                    <span className="inline-flex size-6 shrink-0 items-center justify-center rounded-full bg-stone-200 text-[11px] font-medium uppercase text-stone-700 dark:bg-stone-700 dark:text-stone-200">
                                        {(currentUser.name || currentUser.email || "?").slice(0, 1)}
                                    </span>
                                    <span className="truncate">{currentUser.name || currentUser.email}</span>
                                </span>
                            ) : null}
''' + anchor
    s = s.replace(anchor, chip, 1)
    print("用户身份已加入")

# 高度：主站外壳已经给了顶栏，这条二级导航压到 44px，别喧宾夺主
s = s.replace(
    '<header className="sticky top-0 z-20 h-14 shrink-0 border-b',
    '<header className="sticky top-0 z-20 h-11 shrink-0 border-b',
)
s = s.replace('h-14 min-w-0 items-center gap-7', 'h-11 min-w-0 items-center gap-6')
s = s.replace('flex h-14 shrink-0 items-center gap-2 text-sm leading-6 transition',
              'flex h-11 shrink-0 items-center gap-2 text-sm leading-6 transition')

p.write_text(s)

# 版本号 / GitHub 从右侧动作区拿掉
ua = base / "src/components/layout/user-status-actions.tsx"
if ua.exists():
    t = ua.read_text()
    t = re.sub(r"\s*<VersionReleaseModal[^/]*/>", "", t)
    t = re.sub(r"\s*<GitHubLink[^/]*/>", "", t)
    ua.write_text(t)
    print("版本号/GitHub 已移除")

print("app-top-nav.tsx 完成")
