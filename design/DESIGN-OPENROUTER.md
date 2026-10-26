# Sub2API × OpenRouter 风格 — 第二版设计

> 第一版（企业海军蓝）见 [DESIGN.md](DESIGN.md)。两版是同一套接入机制的两种皮肤，
> 二选一，不要同时启用。

---

## 一、许可：这次比上次更严格

我查了 OpenRouter 的公开站点（`openrouter.ai` 的线上 CSS，约 535 KB），
**没有找到它公开的前端源代码**（不是开源项目）。这意味着：

| 做法 | 是否可行 |
|---|---|
| 读取它线上的 CSS 变量，理解其设计语言 | 可以 —— 浏览器本来就公开这些内容 |
| 照抄它的 CSS 文件 / 组件代码 | **不行** —— 受著作权保护 |
| 复制它的 logo、图标、插画 | **不行** |
| 沿用它的专有字体（Gordita 是商业字体） | **不行** —— 需单独授权 |
| 借鉴"青黑画布 + 发丝线 + 紫主色"这套视觉方向 | 可以 —— 属于不受保护的设计选择 |

**关键一条**：OpenRouter 用的是 **Gordita**（商业授权字体）做品牌字体，`Plus Jakarta Sans`
做 UI 字体（OFL，可自由使用）。本方案**只用 Plus Jakarta Sans**，不碰 Gordita，
字体文件从 Google Fonts 自行获取，不从它的站点扒。

另外提醒：这是**风格借鉴**，不是复刻。如果你对外宣传时自称"OpenRouter 风格"，
在部分司法辖区可能涉及不正当竞争；建议说"深色技术型控制台风格"。

---

## 二、从 OpenRouter 提取到的四个特征

这些是从它线上 CSS 里读出来的客观事实，不是我的主观描述：

| 特征 | 它的实现 | 证据 |
|---|---|---|
| 画布是青黑，不是炭黑 | `hsl(197.1 53.8% 2.5%)` ≈ `#03080a` | `--background` |
| 卡片比画布浅一档，同色相 | `hsl(197 30% 4.5%)` ≈ `#080d0f` | `--card` |
| 边框是"发丝线" | 墨色 7.8% 透明度，1px | `--border: 197.1 53.8% 2.5% / .078` |
| 主色紫 + 荧光青柠 | `#7624f4` / `#c8ff00` | `--primary: 263.7 90.4% 54.9%`、`--accent: 72.9 100% 50%` |
| 圆角 8px 基准 | `--radius: .5rem`；sm4/md6/lg8/xl12/2xl16 | `--radius-*` |
| 阴影 = 描边 + 柔影 | `0 8px 30px #0000001f, 0 0 1px #0000001a` | `box-shadow` |
| 字体 | Plus Jakarta Sans + Geist Mono | `--font-sans` / `--font-mono` |

第三行那个 `/.078` 是最容易被忽略、又最能体现"高级感"的一处：
**它不是灰色实线，是把前景色（墨）稀释到 8% 画出来的线**。深色模式反过来用白 9%。

---

## 三、一处不能照搬的地方（重要）

OpenRouter 的 `--accent` 是**荧光青柠 `#c8ff00`**。但 sub2api 上游的 `accent-*`
语义完全不同 —— 我统计过，`accent-50…950` 在 sub2api 里**清一色是 slate 灰阶**，
是"辅助中性底色"。

**如果照着 OpenRouter 把 `accent` 映射成青柠，上游所有 `bg-accent-*` 会变成荧光黄**，
页面上会出现大面积高饱和色块。

所以本版的处理是：
- `accent-*` 保持中性灰（不照搬 OpenRouter 的语义）
- 青柠单独开一个 `brand-*` 命名空间，只在**新增**的标记里用

这是"参照"和"照抄"的分界线落地的具体案例。

---

## 四、一个必须处理的对比度陷阱

我算了 WCAG 对比度：

| 组合 | 对比度 | 结论 |
|---|---|---|
| 白字 on `#7624f4` | **6.31:1** | 合格，主按钮可用 |
| `#7624f4` on 白底 | **6.31:1** | 合格，浅色模式链接可用 |
| `#7624f4` on `#03080a` | **3.19:1** | **不合格**，正文不可读 |
| `#a78bfa` (violet-400) on `#03080a` | **7.40:1** | 合格 |
| `#c8ff00` on `#03080a` | **17.03:1** | 极强，但也只适合小面积 |

**结论**：深色模式下的链接、强调文字必须用 `primary-300/400`，不能用 `primary-600`。
这一点已写进 [tailwind.openrouter.js](theme-openrouter/tailwind.openrouter.js) 的注释。

如果直接照抄"主色 = #7624f4"就完事，深色模式的可读性会不达标。

---

## 五、交付物

```
design/
├── DESIGN.md                       # 第一版：企业海军蓝
├── DESIGN-OPENROUTER.md            # 本文件
├── theme/                          # 第一版实现
│   ├── tokens.css
│   ├── tailwind.theme.js
│   └── enterprise.css
├── theme-openrouter/               # 本版实现
│   ├── tokens-openrouter.css       # 配色 / 组件皮肤
│   ├── tailwind.openrouter.js      # 色阶 / 圆角 / 阴影
│   └── layout.css                  # 壳层版式（宽屏）
└── preview/
    ├── index.html                  # 第一版预览
    └── openrouter.html             # 本版预览（浅/深/两版对照）
```

接入步骤与第一版相同（改 `main.ts` 一行 + 合并 `tailwind.config.js`），
另外**需要额外接入 `layout.css` 并改动两个壳层文件**（见下一节）。

> ⚠️ **不要同时挂载两套 tokens.css**。它们都定义了 `.btn-primary` / `.card` 等同名类，
> 后加载的会覆盖前者，但变量会混在一起，得到难以排查的混合色。
> 如果要做"可切换皮肤"，正确的做法是给 `<html>` 加 `data-theme` 属性并把两套规则
> 写成 `[data-theme="navy"] .btn-primary` 这种作用域形式。

---

## 五之二、宽屏版式（这一节是 v3 修正的核心）

第一版的预览把整个应用塞进了一个 `max-width: 1320px` 的圆角容器里，
在 1440 屏上只渲染出 1272px，右侧露出灰边。**这不是配色问题，是布局问题** ——
观感"还是 sub2api"的根源在这里。

### OpenRouter 线上实测的版式参数

| 项 | 它的实现 | 换算 |
|---|---|---|
| 顶栏 | `w-full`，仅 `px-4 lg:px-6` | 通栏，无 max-width，16/24px 内边距 |
| 侧栏列 | `w-60` | **240px** |
| 内容区上限 | `max-w-7xl` + `mx-auto` | **1280px 居中** |
| 内容内边距 | `px-6 md:px-8` | 24/32px |
| 顶栏高度 | `--navbar-height: 3.5rem` | 56px |
| 底部留白 | `--dashboard-page-pad-b` | 2.5rem / 2rem |

### 与 sub2api 默认值的差异

| 项 | sub2api 默认 | 改成 |
|---|---|---|
| 侧栏宽度 | `w-64` = **256px** | 240px |
| 主区内边距 | `p-4 md:p-6 lg:p-8` | 24/32px + 1280 居中容器 |
| 内容宽度 | 全通栏，无上限 | max-width 1280 居中 |

### 一个反直觉的结论

**"宽屏"不等于"更宽"。** OpenRouter 给内容加了 1280px 上限并居中，
屏比 1280 宽时两侧留白对称。如果照 sub2api 那样全通栏，2K/4K 屏上表格行会被
拉到 2000px 以上，扫读极困难 —— 那才是真正"不像 OpenRouter"的地方。

实测结果（`layout.css` + 预览）：
- 1440 视口：内容区 1200px，正好用满主区（1440 − 240 侧栏）
- 1600 视口：内容区 1280px 居中，左右各留 40px
- 超宽屏（≥1600）：上限放宽到 1408px，避免两侧留白过大

### 壳层接入（**这次需要改 .vue**）

配色层是零 `.vue` 改动，但版式绕不开。需要动 `frontend/src/components/layout/` 里两个文件：

```diff
  <!-- AppLayout.vue -->
  <div
    class="relative min-h-screen transition-all duration-300"
-   :class="[sidebarCollapsed ? 'lg:ml-[72px]' : 'lg:ml-64']"
+   :class="[sidebarCollapsed ? 'lg:ml-[72px]' : 'lg:ml-60']"
  >
    <AppHeader />
    <main class="p-4 md:p-6 lg:p-8">
-     <slot />
+     <div class="app-container">
+       <slot />
+     </div>
    </main>
```

```diff
  <!-- AppSidebar.vue:5 -->
- sidebarCollapsed ? 'w-[72px]' : 'w-64',
+ sidebarCollapsed ? 'w-[72px]' : 'w-60',
```

**两处必须同时改**：`w-64`→`w-60` 是侧栏宽度（256→240），
`lg:ml-64`→`lg:ml-60` 是主区左偏移。只改一处会让主区和侧栏错位 16px。

这是本方案第一次产生**需要维护的 merge 冲突面**。好消息是：
`AppLayout.vue` 只有 52 行、`AppSidebar.vue` 的这行在文件第 5 行，
两者都是"改一个字符串"级别的冲突，解起来比改样式表简单得多。

---

## 六、两版怎么选

| | 企业海军蓝 | OpenRouter 风 |
|---|---|---|
| 气质 | 稳重、克制、偏传统 SaaS | 技术感、有锋芒、偏开发者工具 |
| 深色模式 | 深海军蓝，仍像"办公软件" | 青黑，像"终端/IDE" |
| 品牌色 | 蓝，安全但没记忆点 | 紫 + 青柠，辨识度高 |
| 适合 | 面向企业采购、需要"可信"的客户 | 面向开发者、需要"酷"的产品 |
| 风险 | 容易被认为平庸 | 紫色在部分行业（金融）不合规；青柠需严格控制用量 |

**我的建议**：sub2api 的用户是开发者（API 网关），**OpenRouter 这版更匹配受众**。
但要守住两条 —— 青柠总用量不超过页面面积的 1%；深色模式强调文字一律用 `violet-400`。

---

## 七、仍未验证的部分

与第一版相同：
- 没在真实仓库里跑过构建（`/Users/denticle/Documents/sub2api` 目前为空目录）
- 用 Tailwind 3.4.17 单独验证过 `backgroundImage` 关闭与 `rgb(var()/ <alpha-value>)` 的
  编译产物，但**没有对着 sub2api 的实际构建跑一遍**
- 壳层文件（`AppSidebar.vue` 等 5 个）仍未改，预览里的侧栏是重构目标而非成品
