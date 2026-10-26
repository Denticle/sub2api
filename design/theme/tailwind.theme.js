/**
 * Sub2API Enterprise — Tailwind theme overlay
 * ---------------------------------------------------------------------------
 * 放在 frontend/tailwind.config.js 里 theme.extend 的位置，或抽成独立文件：
 *
 *   import { enterpriseTheme } from './theme/tailwind.theme.js'
 *   export default {
 *     content: [...],
 *     darkMode: 'class',
 *     theme: { extend: { ...enterpriseTheme } },
 *   }
 *
 * 为什么要改这里而不是全在 CSS 里：
 *   1) 圆角 —— 上游大量 `rounded-lg` / `rounded-2xl` 是 utility 类，
 *      utilities 层优先级高于 components 层，CSS 覆盖不掉，只能重映射 scale。
 *   2) 颜色 —— `bg-primary-500/20` 这类带透明度修饰符的用法有 2000+ 处，
 *      必须用 rgb(var(--x) / <alpha-value>) 才能让 `/20` 继续生效。
 */

export const enterpriseTheme = {
  colors: {
    /* 品牌色：深海军蓝。语义保持不变（50→950），只换取值。 */
    primary: {
      50: 'rgb(var(--brand-50) / <alpha-value>)',
      100: 'rgb(var(--brand-100) / <alpha-value>)',
      200: 'rgb(var(--brand-200) / <alpha-value>)',
      300: 'rgb(var(--brand-300) / <alpha-value>)',
      400: 'rgb(var(--brand-400) / <alpha-value>)',
      500: 'rgb(var(--brand-500) / <alpha-value>)',
      600: 'rgb(var(--brand-600) / <alpha-value>)',
      700: 'rgb(var(--brand-700) / <alpha-value>)',
      800: 'rgb(var(--brand-800) / <alpha-value>)',
      900: 'rgb(var(--brand-900) / <alpha-value>)',
      950: 'rgb(var(--brand-950) / <alpha-value>)',
    },

    /* 上游 accent-* 清一色是 slate 灰阶，仅作辅助底色。
       保持"辅助中性"的语义，映射到同一套冷灰。 */
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

    /* 注意：上游把 gray-* 当"浅色模式中性色"、dark-* 当"深色模式表面色"，
       两者用途不同，必须分开映射，不能合并成一个色阶。
       合并会让深色模式下的 body 与 card 变成同一个暗度，层级消失。 */
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

  borderRadius: {
    /* 企业控制台：6px 基准。上游默认 lg=8/xl=12/2xl=16 偏圆润。 */
    none: '0',
    sm: '0.1875rem',   /* 3px */
    DEFAULT: '0.25rem',/* 4px  —— 上游 rounded 446 次 */
    md: '0.375rem',    /* 6px  —— 190 次 */
    lg: '0.5rem',      /* 8px  —— 835 次，最主要的一档 */
    xl: '0.625rem',    /* 10px —— 268 次 */
    '2xl': '0.75rem',  /* 12px —— 84 次（卡片） */
    '3xl': '1rem',     /* 16px —— 36 次 */
    '4xl': '1.25rem',  /* 20px（上游自定义为 2rem，这里收紧） */
    full: '9999px',    /* 徽标/头像/图标底，保留 */
  },

  fontFamily: {
    sans: [
      'Inter', 'Geist', 'system-ui', '-apple-system', 'BlinkMacSystemFont',
      'Segoe UI', 'Roboto', 'PingFang SC', 'Hiragino Sans GB',
      'Microsoft YaHei', 'sans-serif',
    ],
    mono: [
      'JetBrains Mono', 'ui-monospace', 'SFMono-Regular', 'Menlo', 'Monaco',
      'Consolas', 'monospace',
    ],
  },

  boxShadow: {
    /* 保留上游的键名，只换取值，避免 152 处 shadow-sm 等引用失效。 */
    sm: '0 1px 2px 0 rgb(15 23 42 / 0.04)',
    DEFAULT: '0 1px 2px 0 rgb(15 23 42 / 0.05)',
    md: '0 2px 6px -1px rgb(15 23 42 / 0.08)',
    lg: '0 6px 16px -4px rgb(15 23 42 / 0.10)',
    xl: '0 12px 28px -8px rgb(15 23 42 / 0.14)',
    '2xl': '0 20px 48px -12px rgb(15 23 42 / 0.20)',
    /* 上游自有的品牌色阴影，取消发光效果 */
    glow: '0 1px 2px 0 rgb(15 23 42 / 0.08)',
    'glow-lg': '0 2px 6px 0 rgb(15 23 42 / 0.12)',
    card: 'var(--shadow-card)',
    'card-hover': 'var(--shadow-card-hover)',
    glass: 'var(--shadow-card)',
    'glass-sm': 'var(--shadow-card)',
    'inner-glow': 'none',
  },

  backgroundImage: {
    /* 企业风格不使用渐变。上游 49 处 bg-gradient-to-br 需要在这里
       退化成单色，否则会出现彩条。 */
    'gradient-radial': 'none',
    'gradient-primary': 'none',
    'gradient-dark': 'none',
    'gradient-glass': 'none',
    'mesh-gradient': 'none',
  },

  transitionDuration: {
    DEFAULT: '150ms',
  },
}
