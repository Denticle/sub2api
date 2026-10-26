/**
 * Sub2API — OpenRouter-inspired Tailwind theme overlay
 * ---------------------------------------------------------------------------
 * 用法与 theme/tailwind.theme.js 相同：合并进 tailwind.config.js 的 theme.extend。
 * 两套方案（海军蓝企业版 / OpenRouter 版）二选一，不要同时合并——
 * 它们的 colors / borderRadius 会互相覆盖。
 */

export const openrouterTheme = {
  colors: {
    /* 品牌紫。600 = #7624f4，白字对比度 6.31:1。
       深色模式下的链接文字请用 300/400（#a78bfa 在 #03080a 上 7.4:1），
       不要用 600 —— 紫在青黑底上只有 3.19:1，正文不可读。 */
    primary: {
      50: 'rgb(var(--p-50) / <alpha-value>)',
      100: 'rgb(var(--p-100) / <alpha-value>)',
      200: 'rgb(var(--p-200) / <alpha-value>)',
      300: 'rgb(var(--p-300) / <alpha-value>)',
      400: 'rgb(var(--p-400) / <alpha-value>)',
      500: 'rgb(var(--p-500) / <alpha-value>)',
      600: 'rgb(var(--p-600) / <alpha-value>)',
      700: 'rgb(var(--p-700) / <alpha-value>)',
      800: 'rgb(var(--p-800) / <alpha-value>)',
      900: 'rgb(var(--p-900) / <alpha-value>)',
      950: 'rgb(var(--p-950) / <alpha-value>)',
    },

    /* 关键：上游的 accent-* 是「中性辅助底色」（清一色 slate 灰阶），
       不是强调色。若把它映射成青柠，会有大面积的 bg-accent-* 变成荧光黄。
       这里保持中性，青柠单独走 brand-*。 */
    accent: {
      50: 'rgb(var(--n-50) / <alpha-value>)',
      100: 'rgb(var(--n-100) / <alpha-value>)',
      200: 'rgb(var(--n-200) / <alpha-value>)',
      300: 'rgb(var(--n-300) / <alpha-value>)',
      400: 'rgb(var(--n-400) / <alpha-value>)',
      500: 'rgb(var(--n-500) / <alpha-value>)',
      600: 'rgb(var(--n-600) / <alpha-value>)',
      700: 'rgb(var(--n-700) / <alpha-value>)',
      800: 'rgb(var(--n-800) / <alpha-value>)',
      900: 'rgb(var(--n-900) / <alpha-value>)',
      950: 'rgb(var(--n-950) / <alpha-value>)',
    },

    /* 青柠：只供你后续新增的标记使用，上游没有 brand-* 类，不会误伤。 */
    brand: {
      DEFAULT: 'rgb(var(--lime) / <alpha-value>)',
      lime: 'rgb(var(--lime) / <alpha-value>)',
    },

    gray: {
      50: 'rgb(var(--n-50) / <alpha-value>)',
      100: 'rgb(var(--n-100) / <alpha-value>)',
      200: 'rgb(var(--n-200) / <alpha-value>)',
      300: 'rgb(var(--n-300) / <alpha-value>)',
      400: 'rgb(var(--n-400) / <alpha-value>)',
      500: 'rgb(var(--n-500) / <alpha-value>)',
      600: 'rgb(var(--n-600) / <alpha-value>)',
      700: 'rgb(var(--n-700) / <alpha-value>)',
      800: 'rgb(var(--n-800) / <alpha-value>)',
      900: 'rgb(var(--n-900) / <alpha-value>)',
      950: 'rgb(var(--n-950) / <alpha-value>)',
    },
    /* dark-* 走青黑阶，不是 zinc —— 这是"青黑画布"能否成立的关键。 */
    dark: {
      50: 'rgb(var(--n-50) / <alpha-value>)',
      100: 'rgb(var(--n-100) / <alpha-value>)',
      200: 'rgb(var(--n-200) / <alpha-value>)',
      300: 'rgb(var(--n-300) / <alpha-value>)',
      400: 'rgb(var(--n-400) / <alpha-value>)',
      500: 'rgb(var(--n-500) / <alpha-value>)',
      600: 'rgb(var(--d-600) / <alpha-value>)',
      700: 'rgb(var(--d-700) / <alpha-value>)',
      800: 'rgb(var(--d-800) / <alpha-value>)',
      900: 'rgb(var(--d-900) / <alpha-value>)',
      950: 'rgb(var(--d-950) / <alpha-value>)',
    },
  },

  /* OpenRouter 的档位：sm 4 / md 6 / lg 8 / xl 12 / 2xl 16。
     上游默认 lg=8 恰好一致，主要需要收的是 xl(12→10) 与 2xl(16→12)。 */
  borderRadius: {
    none: '0',
    sm: '0.25rem',     /* 4px */
    DEFAULT: '0.25rem',
    md: '0.375rem',    /* 6px */
    lg: '0.5rem',      /* 8px */
    xl: '0.625rem',    /* 10px */
    '2xl': '0.75rem',  /* 12px */
    '3xl': '1rem',     /* 16px */
    '4xl': '1.25rem',
    full: '9999px',    /* 头像/图标底仍需要正圆 */
  },

  fontFamily: {
    sans: [
      'Plus Jakarta Sans', 'Inter', 'system-ui', '-apple-system',
      'BlinkMacSystemFont', 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB',
      'Microsoft YaHei', 'sans-serif',
    ],
    mono: [
      'Geist Mono', 'JetBrains Mono', 'ui-monospace', 'SFMono-Regular',
      'Menlo', 'Consolas', 'monospace',
    ],
  },

  boxShadow: {
    /* OpenRouter 的签名阴影：描边 + 柔影两层叠加 */
    sm: '0 1px 2px rgb(0 0 0 / 0.06)',
    DEFAULT: '0 1px 3px rgb(0 0 0 / 0.10), 0 1px 2px -1px rgb(0 0 0 / 0.10)',
    md: '0 2px 8px rgb(0 0 0 / 0.08), 0 0 1px rgb(0 0 0 / 0.10)',
    lg: '0 4px 20px rgb(0 0 0 / 0.08), 0 0 1px rgb(0 0 0 / 0.10)',
    xl: '0 8px 30px rgb(0 0 0 / 0.12), 0 0 1px rgb(0 0 0 / 0.10)',
    '2xl': '0 16px 48px rgb(0 0 0 / 0.16), 0 0 1px rgb(0 0 0 / 0.12)',
    /* 上游的品牌发光阴影：改成紫色投影，保留"主色按钮浮起"的语义但不发光 */
    glow: '0 2px 8px rgb(118 36 244 / 0.24)',
    'glow-lg': '0 4px 16px rgb(118 36 244 / 0.30)',
    card: 'var(--shadow-card)',
    'card-hover': 'var(--shadow-card-hover)',
    glass: 'var(--shadow-card)',
    'glass-sm': 'var(--shadow-card)',
    'inner-glow': 'none',
  },

  backgroundImage: {
    /* 关掉上游 76 处渐变，OpenRouter 风格靠平面 + 发丝线，不靠渐变。
       注意：bg-gradient-to-* 的 8 个方向来自 Tailwind 的**默认主题**，
       不是上游 tailwind.config.js 里写的（上游只自定义了下面 5 个）。
       必须逐个列全，否则漏掉的键会回落到默认渐变。 */
    'gradient-to-t': 'none',
    'gradient-to-tr': 'none',
    'gradient-to-r': 'none',
    'gradient-to-br': 'none',
    'gradient-to-b': 'none',
    'gradient-to-bl': 'none',
    'gradient-to-l': 'none',
    'gradient-to-tl': 'none',
    'gradient-radial': 'none',
    'gradient-primary': 'none',
    'gradient-dark': 'none',
    'gradient-glass': 'none',
    'mesh-gradient': 'none',
  },
}
