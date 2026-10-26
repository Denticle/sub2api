/**
 * Sub2API — Forapi 主题 overlay
 * ---------------------------------------------------------------------------
 * 取色来自线上 forapi.cloud（--seed-* / --admin-* token）：
 *   主色绿 #22c55e（深色）/ #16a34a（浅色），紫 #a78bfa 只作辅助色。
 *
 * 用法：合并进 tailwind.config.js 的 theme.extend，且必须放在 extend 的**最后**，
 * 否则上游的 colors 键会盖掉这里。
 */

export const forapiTheme = {
  colors: {
    /* 品牌绿。500 = #22c55e（forapi 深色主色），600 = #16a34a（浅色主色）。
       注意：主按钮是「绿底 + 深绿字 #062014」，不是白字 ——
       白字在 #22c55e 上只有 2.3:1，深绿字有 7.9:1。 */
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
       不是强调色。若把它映射成品牌色，会有大面积的 bg-accent-* 变成彩色块。
       这里保持中性，品牌绿单独走 primary-*。 */
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

    /* 辅助紫：对应 forapi 的 --seed-accent，只用于高亮，不当主色 */
    brand: {
      DEFAULT: 'rgb(var(--accent) / <alpha-value>)',
      accent: 'rgb(var(--accent) / <alpha-value>)',
      cyan: 'rgb(var(--cyan) / <alpha-value>)',
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
    /* dark-* 走暖黑阶（#141413 → #272726），与 forapi 的画布同族。
       注意 50–500 映射到浅色中性色，600 起才切到暖黑表面。 */
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

  /* forapi 的档位：radius-sm 8 / radius 12 / radius-lg 20。
     上游 lg 默认 8px 恰好是 forapi 的小号档，xl 收到 10、2xl 收到 12。 */
  borderRadius: {
    none: '0',
    sm: '0.25rem',     /* 4px */
    DEFAULT: '0.25rem',
    md: '0.375rem',    /* 6px */
    lg: '0.5rem',      /* 8px —— forapi --seed-radius-sm */
    xl: '0.625rem',    /* 10px */
    '2xl': '0.75rem',  /* 12px —— forapi --seed-radius */
    '3xl': '1.25rem',  /* 20px —— forapi --seed-radius-lg */
    '4xl': '1.5rem',
    full: '9999px',    /* 头像/图标底仍需要正圆 */
  },

  fontFamily: {
    /* forapi 用 Inter + JetBrains Mono */
    sans: [
      'Inter Variable', 'Inter', 'system-ui', '-apple-system',
      'BlinkMacSystemFont', 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB',
      'Microsoft YaHei', 'sans-serif',
    ],
    mono: [
      'JetBrains Mono Variable', 'JetBrains Mono', 'SFMono-Regular',
      'Consolas', 'ui-monospace', 'Menlo', 'monospace',
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
    /* 上游的品牌发光阴影：改成绿色投影，保留语义但不发光 */
    glow: '0 2px 8px rgb(34 197 94 / 0.24)',
    'glow-lg': '0 4px 16px rgb(34 197 94 / 0.30)',
    card: 'var(--shadow-card)',
    'card-hover': 'var(--shadow-card-hover)',
    glass: 'var(--shadow-card)',
    'glass-sm': 'var(--shadow-card)',
    'inner-glow': 'none',
  },

  backgroundImage: {
    /* 关掉上游的渐变，forapi 风格靠平面 + 发丝线，不靠渐变。
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
