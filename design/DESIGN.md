# Sub2API 企业风格改造方案

> 目标：把 sub2api 前端做成 new-api 那种企业级控制台观感，同时
> (a) 不触碰 new-api 的授权限制，(b) 不阻断后续合并上游更新。

---

## 一、许可边界：你真正要避开的是什么

先把两边的授权查清楚，因为这里决定了哪些做法是合法的。

| 项目 | 许可证 | 关键约束 |
|---|---|---|
| **sub2api**（你要改的） | **LGPL-3.0 or later** | 修改后分发须以同许可开放对应源码；须保留版权声明 |
| **new-api**（你参考的） | **AGPLv3 + §7 附加条款** | 修改后分发：①保留署名 "Frontend design and development by New API contributors." ②界面上保留指向原仓库的可见链接 ③整个网络服务须开源 |

### 法律上受保护的 vs 不受保护的

**不受保护（可以自由借鉴）**
- "企业风格"本身：配色倾向、信息密度、圆角尺度、留白节奏、层级手法
- 设计模式：侧栏 + 顶栏 + 统计卡 + 数据表、表头小写字距、等宽数字对齐
- 通用 UI 语言：这是行业共有的功能性与美学选择，不构成著作权客体

**受保护（不能复制）**
- new-api 的源代码，包括 `web/src/styles/theme.css`、`index.css`、组件 `.tsx`
- 它的具体 token 命名体系与取值组合（照抄即构成复制）
- 自定义图标集、插画、截图、logo
- 字体文件（Public Sans / Lora 是 OFL，可单独获取，但不建议沿用同一套——见下）
- 商标与项目名

### 结论

**"做成 new-api 那种企业风格"是合法的；"把 new-api 的前端拿过来"不合法。**

所以本方案的做法是：**只取设计意图（企业控制台的共性语言），从零实现 token 与组件皮肤**。交付的 CSS 没有一行来自 new-api，命名、取值、结构都是独立设计的。

### 三条必须守住的线

1. **不要 copy 任何 `.css` / `.tsx` / 图标资源**。风格可以像，代码不能像。
2. **不要用同一套字体组合**。new-api 用 Public Sans + Lora（且有衬线体做正文字轴）。沿用同一组合会让人直接联想到它；本方案改用 Inter + JetBrains Mono。
3. **不要复刻它的 token 命名**（`--sidebar-primary` / `--table-disabled-border` 那一套）。命名是识别度很高的部分，本方案用 `--brand-* / --n-* / --d-*`。

> 商业用途（对外售卖、SaaS 交付）请让法务过一遍 LGPL 的分发义务，
> 尤其是"前端皮肤是否构成衍生作品"。本文件不构成法律意见。

---

## 二、差异化锚点：怎么做到"像但不等于"

new-api 的实际视觉语言（从它的源码确认）：

- 主色：亮天蓝 `oklch(0.692 0.141 243.7)` ≈ `#4ba3e3`，饱和度中高、亮度高
- 中性色：**完全去饱和**（`oklch(x 0 0)`），纯灰
- 深色画布：炭黑 `oklch(0.235 0 0)`
- 圆角基准：`--radius: 1rem`，`rounded-2xl` 相当于 1.8rem，偏圆润
- 字体：Public Sans（UI）+ Lora（可选衬线正文）
- 风格：shadcn 流派，柔和、留白多

本方案的四个反向选择，每一处都同时是"企业感"的来源和"不像 new-api"的保证：

| 维度 | new-api | 本方案 | 理由 |
|---|---|---|---|
| 主色 | 亮天蓝，高亮度 | **深海军蓝** `#164cba`，低亮度高饱和 | 深色更耐看，长时间盯表格不刺眼；且刻意避开同一色相区间 |
| 中性色 | 纯去饱和灰 | **冷调 Slate**（带 ~215° 色相） | 纯灰在大面积下显"廉价塑料"；冷灰是 IBM Carbon / 阿里云控制台的共同选择 |
| 深色画布 | 炭黑，无彩 | **深海军蓝** `#0b111c` | 深色模式仍保留品牌气质，而不是变灰 |
| 圆角 | 1rem 基准（大） | **6–8px 基准**（小） | 圆角越小越"工具"，越大越"消费级"。这是观感变化最大的一项 |
| 阴影 | shadcn 默认柔和阴影 | **几乎不用阴影，1px 边框分层** | 密集表格下阴影会脏；弹出层才用重阴影 |
| 字体 | Public Sans + Lora | **Inter + JetBrains Mono** | 数字等宽对齐是数据后台的刚需 |

---

## 三、交付物

```
design/
├── DESIGN.md                    # 本文件：方案、边界、接入步骤
├── theme/
│   ├── tokens.css               # 变量 + 组件皮肤（上游 .vue 零改动）
│   ├── tailwind.theme.js        # Tailwind scale 重映射（圆角/色板/阴影）
│   └── enterprise.css           # 入口（目前只 @import tokens.css）
└── preview/
    └── index.html               # 可交互预览（浅色/深色/新旧对照）
```

预览已在浏览器中验证：浅色、深色两种模式渲染正常，1440px 下三栏布局完整。

---

## 四、接入步骤

### 1. 复制三个样式文件

```bash
mkdir -p frontend/src/theme
cp design/theme/tokens.css    frontend/src/theme/
cp design/theme/enterprise.css frontend/src/theme/
cp design/theme/tailwind.theme.js frontend/src/theme/
```

### 2. 挂载 CSS（唯一需要改的上游文件，1 行）

`frontend/src/main.ts`：

```diff
  import './style.css'
+ import './theme/enterprise.css'
```

### 3. 合并 Tailwind 主题

`frontend/tailwind.config.js` —— 把 `theme.extend` 里**已有的** `colors` / `borderRadius` /
`fontFamily` / `boxShadow` / `backgroundImage` 五个键替换为 `enterpriseTheme`
（其余 `animation` / `keyframes` 等保持不动）：

```diff
+ import { enterpriseTheme } from './src/theme/tailwind.theme.js'
+
  export default {
    content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
    darkMode: 'class',
    theme: {
-     extend: {
-       colors: { primary: {...}, accent: {...}, dark: {...} },
-       fontFamily: {...},
-       boxShadow: {...},
-       backgroundImage: {...},
-       borderRadius: { '4xl': '2rem' },
-       animation: {...},
-       keyframes: {...},
-     }
+     extend: {
+       ...enterpriseTheme,
+       animation: {...},
+       keyframes: {...},
+     }
    },
  }
```

> `tailwind.config.js` 上游最后一次改动是 2026-01，属于冷点文件，适合承载改动。

### 4. 验证

```bash
pnpm --dir frontend run build
```

```bash
pnpm --dir frontend run test:run
```

---

## 五、已知的坑（都已在实现里处理）

1. **`rgb(var(--x) / <alpha-value>)` 必须带**。上游有 2000+ 处 `bg-primary-500/20`
   这类透明度修饰符，写成裸 `rgb(var(--x))` 会让这些全部失效。
2. **`gray-*` 与 `dark-*` 不能合并成一个色阶**。上游用 `gray-*` 当浅色模式中性色、
   `dark-*` 当深色模式表面色（`dark:bg-dark-800`）。合并后深色模式会失去层级。
3. **`gray-*` 必须定义完整 50→950**。`tailwind.theme.js` 把上游全部档位都映射到
   `--n-*`，缺一档就有一条工具类变成无效声明。
4. **`.stat-card` 用 `@apply card` 内联了样式**，只覆盖 `.card` 管不到它，
   所以 `tokens.css` 里两者并列声明。
5. **`rounded-full` 不能改**。542 处头像/图标底依赖它保持圆形；徽标改方角
   直接在 `.badge` 上加 `border-radius`（`.badge` 的 `@apply rounded-full` 是内联进
   规则体的，本文件后加载可覆盖）。
6. **`backgroundImage: { 'gradient-to-br': 'none' }` 需要验证**。已用 Tailwind 3.4.17
   实测，编译结果是 `background-image: none`，上游 76 处渐变会退化为无色。
   **但 `from-primary-500` / `to-primary-600` 仍会输出 `--tw-gradient-*` 变量**
   （无害，因为没有 `background-image` 消费它们）。如果发现某个按钮变成透明，
   把那个页面的 `bg-gradient-*` 换成实心色即可。
7. **`designSystem.structure.spec.ts` 会挂**。该测试对源码做字符串断言
   （要求存在 `btn btn-secondary`、`rounded-3xl`、`ring-1 ring-gray-900/5`）。
   改配色不影响，但如果你后续调整了 channel-monitor-v2 的 markup 就会失败。

---

## 六、还没做的（下一步）

- **Phase 3 壳层**：`AppSidebar.vue`(1105 行) / `AppHeader.vue`(407) /
  `AppLayout.vue`(52) / `AuthLayout.vue`(88)。这 5 个文件决定 80+ 页面的骨架，
  是拿到最后 20% 观感的关键，但也是上游改动最频繁的文件（Sidebar 改到 2026-09）。
  预览里的侧栏设计可以直接作为重构目标。
- **图表配色**：`src/components/charts/` 有 46 处硬编码 hex（`#3b82f6`、`#14b8a6` 等），
  不在 token 体系内，需要单独建一个 `chartPalette.ts`。
- **登录页 / 公开页**：`AuthLayout.vue` 是品牌露出最集中的地方，值得单独设计。
- **要不要引入组件库**：Element Plus / Arco / Semi 能快速拿到企业感，但会把
  diff 面从 3 个文件扩到全部页面，上游同步成本翻倍。**不建议**。
