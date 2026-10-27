<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="hasHomeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Compact Home Page -->
  <div
    v-else-if="compactHomeEnabled"
    data-testid="compact-home"
    class="flex min-h-screen flex-col bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white"
  >
    <header class="border-b border-gray-200 px-4 py-4 sm:px-6 dark:border-dark-800">
      <nav class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 sm:gap-4">
        <div class="flex min-w-0 flex-1 items-center gap-3">
          <img
            :src="siteLogo || '/logo.svg'"
            alt="Logo"
            class="h-9 w-9 shrink-0 rounded-lg object-contain"
          />
          <span class="min-w-0 truncate text-base font-semibold">{{ siteName }}</span>
        </div>
        <div class="flex max-w-full shrink-0 flex-wrap items-center justify-end gap-2">
          <LocaleSwitcher />
          <router-link
            v-if="docsLinkLocal"
            to="/docs"
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </router-link>
          <a
            v-else-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="flex h-10 shrink-0 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="md" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>
          <button
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex min-h-10 shrink-0 items-center justify-center rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-200"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="flex min-w-0 flex-1 items-center justify-center px-4 py-16 sm:px-6">
      <div class="min-w-0 max-w-2xl text-center">
        <img
          :src="siteLogo || '/logo.svg'"
          alt="Logo"
          class="mx-auto mb-6 h-20 w-20 rounded-2xl object-contain"
        />
        <h1 class="[overflow-wrap:anywhere] text-3xl font-bold md:text-4xl">{{ siteName }}</h1>
        <p class="mt-4 whitespace-pre-wrap [overflow-wrap:anywhere] text-base text-gray-600 dark:text-dark-300">{{ siteSubtitle }}</p>
        <router-link
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="mt-8 inline-flex min-h-10 items-center justify-center rounded-lg bg-primary-600 px-5 py-2.5 text-sm font-medium text-white hover:bg-primary-700"
        >
          {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
        </router-link>
      </div>
    </main>

    <footer class="min-w-0 border-t border-gray-200 px-4 py-5 text-center text-sm text-gray-500 [overflow-wrap:anywhere] sm:px-6 dark:border-dark-800 dark:text-dark-400">
      &copy; {{ currentYear }} {{ siteName }}
    </footer>
  </div>

  <!-- Default Home Page — 重新设计：单一轴线版面，配色沿用 forapi token -->
  <div v-else class="lp text-gray-900 dark:text-gray-100">
    <!-- ── Header ─────────────────────────────────────────────── -->
    <header class="lp-nav">
      <div class="lp-wrap lp-nav-inner">
        <router-link to="/home" class="lp-brand">
          <img :src="siteLogo || '/logo.svg'" alt="" class="lp-brand-mark" />
          <span class="lp-brand-name">{{ siteName }}</span>
        </router-link>

        <nav class="lp-nav-links">
          <router-link v-if="docsLinkLocal" to="/docs">{{ t('home.docs') }}</router-link>
          <a v-else-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer">{{ t('home.docs') }}</a>
          <router-link v-if="showModelPlazaEntry" to="/model-plaza">{{ t('nav.modelPlaza') }}</router-link>
        </nav>

        <div class="lp-nav-actions">
          <LocaleSwitcher />
          <button
            class="lp-icon-btn"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="sm" />
            <Icon v-else name="moon" size="sm" />
          </button>
          <router-link :to="authTarget" class="lp-btn lp-btn-primary lp-btn-sm">
            {{ authLabel }}
          </router-link>
        </div>
      </div>
    </header>

    <main>
      <!-- ── Hero ───────────────────────────────────────────────── -->
      <section class="lp-hero">
        <div class="lp-wrap lp-hero-grid">
          <div class="lp-hero-copy">
            <p class="lp-status">
              <span class="lp-status-dot" aria-hidden="true"></span>
              <span class="lp-status-label">{{ t('home.lp.status') }}</span>
              <span class="lp-status-sep" aria-hidden="true">/</span>
              <span class="lp-status-meta">{{ statusMeta }}</span>
            </p>

            <h1 class="lp-h1">
              <span>{{ t('home.lp.headlineOne') }}</span>
              <span>{{ t('home.lp.headlineTwo') }}</span>
            </h1>

            <p class="lp-lede">{{ t('home.lp.lede') }}</p>
            <p class="lp-note">{{ t('home.lp.ledeNote') }}</p>

            <div class="lp-cta-row">
              <router-link :to="authTarget" class="lp-btn lp-btn-primary">
                {{ t('home.lp.ctaPrimary') }}
              </router-link>
              <a :href="'#lp-models'" class="lp-btn lp-btn-quiet">{{ t('home.lp.ctaSecondary') }}</a>
            </div>
          </div>

          <!-- 路由图：纯 SVG，无 JS 循环。全部动效都在这一处 -->
          <figure class="lp-trace" aria-hidden="true">
            <svg viewBox="0 0 460 340" xmlns="http://www.w3.org/2000/svg">
              <defs>
                <marker id="lp-arrow" viewBox="0 0 8 8" refX="6" refY="4" markerWidth="7" markerHeight="7" orient="auto">
                  <path d="M0,0 L8,4 L0,8 z" fill="currentColor" />
                </marker>
              </defs>

              <!-- 左侧：入口 -->
              <g class="lp-trace-node">
                <rect x="8" y="150" width="86" height="40" rx="9" />
                <text x="51" y="175" text-anchor="middle">/v1/chat</text>
              </g>

              <!-- 中间：网关 -->
              <g class="lp-trace-node lp-trace-node-accent">
                <rect x="156" y="132" width="96" height="76" rx="12" />
                <text x="204" y="165" text-anchor="middle">{{ siteName }}</text>
                <text x="204" y="186" text-anchor="middle" class="lp-trace-sub">gateway</text>
              </g>

              <!-- 右侧：上游三条通道 -->
              <g class="lp-trace-node" v-for="ch in channels" :key="ch.key">
                <rect :x="344" :y="ch.y" width="108" height="40" rx="9" :class="{ 'lp-trace-node-dead': ch.dead }" />
                <text :x="398" :y="ch.y + 25" text-anchor="middle">{{ ch.label }}</text>
              </g>

              <!-- 连接线 -->
              <path d="M102 170 C128 170, 132 170, 152 170" class="lp-trace-path" marker-end="url(#lp-arrow)" />

              <path v-for="ch in channels" :key="'p-'+ch.key"
                :d="`M256 170 C300 170, 306 ${ch.y + 20}, 340 ${ch.y + 20}`"
                :class="['lp-trace-path', ch.dead ? 'lp-trace-path-dead' : '']"
                marker-end="url(#lp-arrow)" />

              <!-- 流动的数据包：仅健康的路径上有 -->
              <circle v-for="ch in channels.filter(c => !c.dead)" :key="'d-'+ch.key" r="3" class="lp-trace-packet">
                <animateMotion
                  :dur="ch.dur"
                  repeatCount="indefinite"
                  :path="`M256 170 C300 170, 306 ${ch.y + 20}, 340 ${ch.y + 20}`"
                />
              </circle>

              <!-- 熔断标注 -->
              <g class="lp-trace-flag">
                <circle :cx="300" :cy="channels[0].y + 20" r="9" />
                <text :x="300" :y="channels[0].y + 24" text-anchor="middle">!</text>
              </g>
            </svg>

            <figcaption class="lp-trace-caption">
              <span v-for="(line, i) in traceLines" :key="i" class="lp-trace-line" :class="line.tone">
                <span class="lp-trace-time">{{ line.t }}</span>
                <span>{{ line.text }}</span>
              </span>
            </figcaption>
          </figure>
        </div>

        <!-- 指标条：并列定义，不是卡片 -->
        <div class="lp-wrap">
          <dl class="lp-stats">
            <div v-for="stat in stats" :key="stat.key">
              <dt>{{ stat.label }}</dt>
              <dd>{{ stat.value }}</dd>
            </div>
          </dl>
        </div>
      </section>

      <!-- ── 上游 ───────────────────────────────────────────────── -->
      <section class="lp-section lp-band">
        <div class="lp-wrap">
          <div class="lp-prov-split">
            <header class="lp-sec-head lp-sec-head-flush">
              <h2>{{ t('home.lp.providersTitle') }}</h2>
              <p>{{ t('home.lp.providersSub') }}</p>
            </header>

            <div class="lp-providers" data-testid="home-providers">
              <div v-for="p in providers" :key="p.key" class="lp-provider">
                <span class="lp-provider-mark" :style="{ background: p.tint }">{{ p.mark }}</span>
                <span class="lp-provider-body">
                  <span class="lp-provider-name">{{ p.name }}</span>
                  <span class="lp-provider-share">
                    <span class="lp-provider-bar" :style="{ width: p.share + '%' }"></span>
                  </span>
                </span>
                <span class="lp-provider-count">{{ p.count }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- ── 模型与费率 ─────────────────────────────────────────── -->
      <section class="lp-section" id="lp-models">
        <div class="lp-wrap">
          <header class="lp-sec-head lp-sec-head-row">
            <div>
              <h2>{{ t('home.lp.modelsTitle') }}</h2>
              <p>{{ t('home.lp.modelsSub') }}</p>
            </div>
            <router-link v-if="showModelPlazaEntry" to="/model-plaza" class="lp-link">
              {{ t('home.lp.modelsAll') }}
              <Icon name="chevronRight" size="xs" />
            </router-link>
          </header>

          <table v-if="rateRows.length" class="lp-table">
            <thead>
              <tr>
                <th scope="col">{{ t('home.lp.colModel') }}</th>
                <th scope="col">{{ t('home.lp.colPlatform') }}</th>
                <th scope="col" class="lp-num">{{ t('home.lp.colRate') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in rateRows" :key="row.key">
                <td class="lp-mono">{{ row.model }}</td>
                <td class="lp-muted">{{ row.platform }}</td>
                <td class="lp-num lp-mono">{{ row.rate }}</td>
              </tr>
            </tbody>
          </table>
          <p v-else class="lp-empty">{{ t('home.lp.modelsEmpty') }}</p>
        </div>
      </section>

      <!-- ── 四件能力 ───────────────────────────────────────────── -->
      <section class="lp-section lp-band">
        <div class="lp-wrap">
          <header class="lp-sec-head">
            <h2>{{ t('home.lp.capTitle') }}</h2>
          </header>

          <div class="lp-cap-grid">
            <article v-for="cap in capabilities" :key="cap.key" class="lp-cap">
              <h3>{{ cap.title }}</h3>
              <p>{{ cap.desc }}</p>
            </article>
          </div>
        </div>
      </section>

      <!-- ── 接入 ───────────────────────────────────────────────── -->
      <section class="lp-section">
        <div class="lp-wrap lp-quick">
          <div class="lp-quick-steps">
            <header class="lp-sec-head">
              <h2>{{ t('home.lp.quickTitle') }}</h2>
            </header>
            <ol class="lp-steps">
              <li v-for="(step, i) in steps" :key="step.key">
                <span class="lp-step-num">{{ String(i + 1).padStart(2, '0') }}</span>
                <div>
                  <h3>{{ step.title }}</h3>
                  <p>{{ step.desc }}</p>
                </div>
              </li>
            </ol>
            <router-link v-if="docsLinkLocal" to="/docs" class="lp-link">
              {{ t('home.lp.quickDocs') }}
              <Icon name="chevronRight" size="xs" />
            </router-link>
            <a v-else-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer" class="lp-link">
              {{ t('home.lp.quickDocs') }}
              <Icon name="chevronRight" size="xs" />
            </a>
          </div>

          <div class="lp-code">
            <div class="lp-code-bar">
              <span class="lp-code-file">request.sh</span>
              <button class="lp-code-copy" @click="copyCurl">
                {{ copied ? t('home.lp.quickCopied') : t('home.lp.quickCopy') }}
              </button>
            </div>
<pre class="lp-code-body"><code>{{ curlSnippet }}</code></pre>
            <p class="lp-code-foot">{{ t('home.lp.compatSub') }}</p>
          </div>
        </div>
      </section>

      <!-- ── 收尾 ───────────────────────────────────────────────── -->
      <section class="lp-final">
        <div class="lp-wrap lp-final-inner">
          <div>
            <h2>{{ t('home.lp.ctaTitle') }}</h2>
            <p>{{ t('home.lp.ctaDesc') }}</p>
          </div>
          <router-link :to="authTarget" class="lp-btn lp-btn-primary">{{ t('home.lp.ctaButton') }}</router-link>
        </div>
      </section>
    </main>

    <!-- ── Footer ─────────────────────────────────────────────── -->
    <footer class="lp-footer">
      <div class="lp-wrap lp-footer-inner">
        <div class="lp-footer-brand">
          <img :src="siteLogo || '/logo.svg'" alt="" class="lp-brand-mark lp-brand-mark-sm" />
          <span>{{ siteName }}</span>
          <span class="lp-footer-sub">{{ siteSubtitle }}</span>
        </div>
        <nav class="lp-footer-links">
          <router-link v-if="docsLinkLocal" to="/docs">{{ t('home.docs') }}</router-link>
          <a v-else-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer">{{ t('home.docs') }}</a>
          <router-link v-if="showModelPlazaEntry" to="/model-plaza">{{ t('nav.modelPlaza') }}</router-link>
          <router-link to="/key-usage">{{ t('home.footer.keyUsage') }}</router-link>
        </nav>
        <p class="lp-footer-copy">&copy; {{ currentYear }} {{ siteName }}</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import { isLocalDocsLink, resolveDocsLink } from '@/utils/docsLink'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'
import { platformLabel } from '@/utils/platformColors'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// ── 站点设置 ────────────────────────────────────────────────────
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() =>
  resolveDocsLink(sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')),
)
const docsLinkLocal = computed(() => isLocalDocsLink(docUrl.value))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const compactHomeEnabled = computed(() => appStore.cachedPublicSettings?.compact_home_enabled === true)
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))

const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const isDark = ref(document.documentElement.classList.contains('dark'))

const isAuthenticated = computed(() => authStore.isAuthenticated)
const modelPlazaRequiresAuth = computed(
  () => appStore.cachedPublicSettings?.model_plaza_require_auth === true,
)
const showModelPlazaEntry = computed(
  () => modelPlazaEnabled.value && (isAuthenticated.value || !modelPlazaRequiresAuth.value),
)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')

const authTarget = computed(() => (isAuthenticated.value ? dashboardPath.value : '/login'))
const authLabel = computed(() => (isAuthenticated.value ? t('home.dashboard') : t('home.login')))

// ── 真实数据：模型广场（公开接口，无需登录） ─────────────────────
interface PlazaModel { name: string; platform: string }
interface PlazaGroup { name: string; platform: string; rate_multiplier?: number; models?: PlazaModel[] }

const plaza = ref<{ groups: PlazaGroup[] } | null>(null)

async function loadPlaza() {
  try {
    const res = await fetch('/api/v1/model-plaza', { headers: { Accept: 'application/json' } })
    if (!res.ok) return
    const body = await res.json()
    plaza.value = body?.data ?? null
  } catch {
    // 目录拿不到就退化为占位，不让首页因此报错
  }
}

const modelNames = computed(() => {
  const seen = new Set<string>()
  for (const g of plaza.value?.groups ?? []) {
    for (const m of g.models ?? []) seen.add(m.name)
  }
  return [...seen]
})

// ── 上游：只列出目录里真实存在的平台 ─────────────────────────────
const PLATFORM_MARKS: Record<string, { mark: string; tint: string }> = {
  anthropic: { mark: 'A', tint: 'linear-gradient(135deg,#d97757,#b85a3c)' },
  openai: { mark: 'O', tint: 'linear-gradient(135deg,#10a37f,#0a7c5f)' },
  gemini: { mark: 'G', tint: 'linear-gradient(135deg,#4285f4,#2b6de0)' },
  antigravity: { mark: 'A', tint: 'linear-gradient(135deg,#f43f5e,#be123c)' },
  deepseek: { mark: 'D', tint: 'linear-gradient(135deg,#4d6bfe,#3350d8)' },
  xai: { mark: 'X', tint: 'linear-gradient(135deg,#1f2937,#0b1016)' },
  kimi: { mark: 'K', tint: 'linear-gradient(135deg,#6366f1,#4338ca)' },
  zhipu: { mark: 'Z', tint: 'linear-gradient(135deg,#0ea5e9,#0369a1)' },
  minimax: { mark: 'M', tint: 'linear-gradient(135deg,#f59e0b,#b45309)' },
  qwen: { mark: 'Q', tint: 'linear-gradient(135deg,#8b5cf6,#6d28d9)' },
}

const providers = computed(() => {
  const counts = new Map<string, Set<string>>()
  for (const g of plaza.value?.groups ?? []) {
    const set = counts.get(g.platform) ?? new Set<string>()
    for (const m of g.models ?? []) set.add(m.name)
    counts.set(g.platform, set)
  }
  return [...counts.entries()]
    .filter(([, set]) => set.size > 0)
    .sort((a, b) => b[1].size - a[1].size)
    .map(([key, set]) => ({
      key,
      count: String(set.size),
      size: set.size,
      name: platformLabel(key),
      ...(PLATFORM_MARKS[key] ?? {
        mark: key.slice(0, 1).toUpperCase(),
        tint: 'linear-gradient(135deg,#71717a,#3f3f46)',
      }),
    }))
    .map((p, _i, arr) => {
      const max = Math.max(1, ...arr.map((x) => x.size))
      return { ...p, share: Math.max(8, Math.round((p.size / max) * 100)) }
    })
})


const stats = computed(() => [
  { key: 'models', value: modelNames.value.length ? String(modelNames.value.length) : '—', label: t('home.lp.statModels') },
  { key: 'providers', value: String(providers.value.length), label: t('home.lp.statProviders') },
  { key: 'protocols', value: '3', label: t('home.lp.statProtocols') },
  { key: 'billing', value: t('home.lp.statBillingValue'), label: t('home.lp.statBillingLabel') },
])

const statusMeta = computed(() =>
  t('home.lp.statusMeta', {
    models: modelNames.value.length || '—',
    providers: providers.value.length,
  }),
)

// ── 费率表：按费率升序取前 8 行 ─────────────────────────────────
const rateRows = computed(() => {
  const rows: { key: string; model: string; platform: string; rate: string; sort: number }[] = []
  for (const g of plaza.value?.groups ?? []) {
    const mult = g.rate_multiplier
    if (typeof mult !== 'number') continue
    for (const m of g.models ?? []) {
      rows.push({
        key: `${g.platform}:${m.name}`,
        model: m.name,
        platform: platformLabel(g.platform),
        rate: `×${mult.toFixed(2)}`,
        sort: mult,
      })
    }
  }
  // 同一模型出现在多个分组时保留最便宜的那条
  const best = new Map<string, (typeof rows)[number]>()
  for (const row of rows) {
    const prev = best.get(row.model)
    if (!prev || row.sort < prev.sort) best.set(row.model, row)
  }
  return [...best.values()].sort((a, b) => a.sort - b.sort).slice(0, 8)
})

const capabilities = computed(() => [
  { key: 'route', title: t('home.lp.capRouteTitle'), desc: t('home.lp.capRouteDesc') },
  { key: 'fail', title: t('home.lp.capFailTitle'), desc: t('home.lp.capFailDesc') },
  { key: 'observe', title: t('home.lp.capObserveTitle'), desc: t('home.lp.capObserveDesc') },
  { key: 'bill', title: t('home.lp.capBillTitle'), desc: t('home.lp.capBillDesc') },
])

const steps = computed(() => [
  { key: 'signup', title: t('home.lp.quickStep1'), desc: t('home.lp.quickStep1Desc') },
  { key: 'key', title: t('home.lp.quickStep2'), desc: t('home.lp.quickStep2Desc') },
  { key: 'base', title: t('home.lp.quickStep3'), desc: t('home.lp.quickStep3Desc') },
])

// ── 路由图：一条通道展示熔断，请求改道 ───────────────────────────
const channels = computed(() => [
  { key: 'a', label: 'anthropic', y: 52, dead: true, dur: '0s' },
  { key: 'b', label: 'openai', y: 150, dead: false, dur: '2.8s' },
  { key: 'c', label: 'gemini', y: 248, dead: false, dur: '3.4s' },
])

const traceLines = computed(() => [
  { t: '02ms', text: `POST /v1/chat/completions → openai`, tone: '' },
  { t: '13ms', text: `209 ${t('home.lp.traceLimited')} · ${t('home.lp.traceFailover')} → gemini`, tone: 'lp-trace-warn' },
])

// ── 调用示例：用站点自身域名，示例才不是假的 ─────────────────────
const baseUrl = computed(() => {
  if (typeof window === 'undefined') return 'https://api.example.com'
  return window.location.origin
})

const curlSnippet = computed(
  () => `curl ${baseUrl.value}/v1/chat/completions \\
  -H "Authorization: Bearer $API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "claude-sonnet-5",
    "messages": [{"role": "user", "content": "hi"}]
  }'`,
)

const copied = ref(false)

async function copyCurl() {
  try {
    await navigator.clipboard.writeText(curlSnippet.value)
    copied.value = true
    window.setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch {
    copied.value = false
  }
}

const currentYear = computed(() => new Date().getFullYear())

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
  loadPlaza()
})
</script>

<style scoped>
/* ==========================================================================
   落地页样式
   --------------------------------------------------------------------------
   配色全部取自 theme 层的 forapi token（--p-* 绿 / --n-* 暖中性 / --d-* 暖黑
   / --line 发丝线），这里不引入任何新的色值，只决定「怎么排」。

   版面原则：
     · 单一轴线。内容左对齐，全站同一左边界，不做居中的大标题。
     · 用发丝线和留白分区，不用卡片。整页没有一处 `border-radius + 阴影`
       组合的容器——那种等距圆角方块正是后台模板的味道。
     · 唯一的高饱和面积是收尾区块，其余靠字重和字号建立层级。
     · 动效只有一处（路由图里的数据包），且尊重 prefers-reduced-motion。
   ========================================================================== */

.lp {
  /* 内容左边界与最大行宽：720px 正文 / 1240px 版心 */
  --lp-gutter: 2rem;
  --lp-max: 78rem;

  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: rgb(var(--n-50));
  color: rgb(var(--n-900));
  font-family: var(--font-ui);
  font-feature-settings: 'cv02', 'cv03', 'cv04', 'cv11';
}

.dark .lp {
  background-color: rgb(var(--d-950));
  color: rgb(var(--n-50));
}

.lp-wrap {
  width: 100%;
  max-width: var(--lp-max);
  margin-inline: auto;
  padding-inline: var(--lp-gutter);
}

@media (max-width: 640px) {
  .lp { --lp-gutter: 1.25rem; }
}

/* ── 顶部导航 ───────────────────────────────────────────────────────────── */

.lp-nav {
  position: sticky;
  top: 0;
  z-index: 30;
  height: 3.5rem;
  display: flex;
  align-items: center;
  border-bottom: 1px solid rgb(var(--line-light) / var(--line-a-light));
  background-color: rgb(var(--n-50) / 0.88);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.dark .lp-nav {
  border-bottom-color: rgb(var(--line) / var(--line-a));
  background-color: rgb(var(--d-950) / 0.88);
}

.lp-nav-inner {
  display: flex;
  align-items: center;
  gap: 1rem;
  height: 100%;
}

.lp-brand {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  flex-shrink: 0;
  text-decoration: none;
  color: inherit;
}

.lp-brand-mark {
  width: 1.5rem;
  height: 1.5rem;
  border-radius: 6px;
  object-fit: contain;
}

.lp-brand-mark-sm {
  width: 1.125rem;
  height: 1.125rem;
  border-radius: 4px;
}

.lp-brand-name {
  font-size: 0.9375rem;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.lp-nav-links {
  display: none;
  align-items: center;
  gap: 1.375rem;
  margin-left: 1.5rem;
}

@media (min-width: 900px) {
  .lp-nav-links { display: flex; }
}

.lp-nav-links a {
  font-size: 0.8125rem;
  font-weight: 500;
  color: rgb(var(--n-600));
  text-decoration: none;
  transition: color 120ms ease;
}

.lp-nav-links a:hover { color: rgb(var(--n-900)); }

.dark .lp-nav-links a { color: rgb(var(--n-400)); }
.dark .lp-nav-links a:hover { color: rgb(var(--n-50)); }

.lp-nav-actions {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  margin-left: auto;
}

.lp-icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border-radius: 6px;
  color: rgb(var(--n-500));
  background: transparent;
  border: 0;
  cursor: pointer;
  transition: background-color 120ms ease, color 120ms ease;
}

.lp-icon-btn:hover {
  background-color: rgb(var(--line-light) / 0.07);
  color: rgb(var(--n-900));
}

.dark .lp-icon-btn { color: rgb(var(--n-400)); }
.dark .lp-icon-btn:hover {
  background-color: rgb(255 255 255 / 0.06);
  color: rgb(var(--n-50));
}

/* ── 按钮：两档，只有一档是实心的 ────────────────────────────────────────── */

.lp-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.375rem;
  height: 2.625rem;
  padding-inline: 1.125rem;
  border-radius: var(--r-md);
  font-size: 0.875rem;
  font-weight: 650;
  letter-spacing: -0.005em;
  text-decoration: none;
  white-space: nowrap;
  transition: background-color 140ms ease, border-color 140ms ease, color 140ms ease;
}

.lp-btn-sm {
  height: 2.125rem;
  padding-inline: 0.875rem;
  font-size: 0.8125rem;
}

.lp-btn-primary {
  background-color: rgb(var(--primary));
  color: rgb(var(--primary-ink));
}

.lp-btn-primary:hover {
  background-color: rgb(var(--primary-hover));
}

.lp-btn-quiet {
  border: 1px solid rgb(var(--line-light) / 0.16);
  color: rgb(var(--n-700));
  background: transparent;
}

.lp-btn-quiet:hover {
  border-color: rgb(var(--line-light) / 0.3);
  background-color: rgb(var(--line-light) / 0.04);
  color: rgb(var(--n-900));
}

.dark .lp-btn-quiet {
  border-color: rgb(var(--line) / var(--line-a-hover));
  color: rgb(var(--n-300));
}

.dark .lp-btn-quiet:hover {
  border-color: rgb(var(--line) / 0.34);
  background-color: rgb(255 255 255 / 0.04);
  color: rgb(var(--n-50));
}

.lp-btn:focus-visible,
.lp-link:focus-visible,
.lp-icon-btn:focus-visible,
.lp-code-copy:focus-visible {
  outline: 2px solid rgb(var(--primary));
  outline-offset: 2px;
}

/* ── Hero ──────────────────────────────────────────────────────────────── */

.lp-hero {
  padding-block: 4.5rem 0;
}

@media (min-width: 1024px) {
  .lp-hero { padding-block: 6rem 0; }
}

.lp-hero-grid {
  display: grid;
  gap: 3.5rem;
  align-items: center;
}

@media (min-width: 1024px) {
  .lp-hero-grid {
    grid-template-columns: minmax(0, 1fr) minmax(0, 27rem);
    gap: 4rem;
  }
}

/* 状态行：一个绿点 + 两段文字，用斜线而不是圆点分隔 */
.lp-status {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0 0 1.75rem;
  font-size: 0.75rem;
  color: rgb(var(--n-500));
}

.lp-status-dot {
  width: 0.375rem;
  height: 0.375rem;
  border-radius: 999px;
  background-color: rgb(var(--p-500));
  box-shadow: 0 0 0 3px rgb(var(--p-500) / 0.16);
  flex-shrink: 0;
}

.lp-status-label {
  font-weight: 600;
  color: rgb(var(--n-700));
}

.dark .lp-status { color: rgb(var(--n-500)); }
.dark .lp-status-label { color: rgb(var(--n-300)); }

.lp-status-sep {
  color: rgb(var(--n-300));
}

.dark .lp-status-sep { color: rgb(var(--d-600)); }

.lp-status-meta {
  font-variant-numeric: tabular-nums;
}

/* 主标题：两行，靠字号和字重立起来，不靠颜色 */
.lp-h1 {
  margin: 0;
  font-size: clamp(2.375rem, 5.2vw, 3.75rem);
  font-weight: 800;
  line-height: 1.04;
  letter-spacing: -0.045em;
  text-wrap: balance;
}

.lp-h1 span {
  display: block;
}

/* 第二行降一档灰度，形成「一句话两段重量」的读法 */
.lp-h1 span + span {
  color: rgb(var(--n-500));
}

.dark .lp-h1 span + span {
  color: rgb(var(--n-500));
}

.lp-lede {
  margin: 1.75rem 0 0;
  max-width: 34rem;
  font-size: 1.0625rem;
  line-height: 1.65;
  color: rgb(var(--n-600));
  text-wrap: pretty;
}

.dark .lp-lede { color: rgb(var(--n-400)); }

.lp-note {
  margin: 0.875rem 0 0;
  max-width: 34rem;
  font-size: 0.8125rem;
  line-height: 1.65;
  color: rgb(var(--n-500));
}

.dark .lp-note { color: rgb(var(--n-500)); }

.lp-cta-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.625rem;
  margin-top: 2.25rem;
}

/* ── 路由图 ─────────────────────────────────────────────────────────────── */

.lp-trace {
  margin: 0;
  padding: 1.5rem 1.25rem 1rem;
  border: 1px solid rgb(var(--line-light) / var(--line-a-light));
  border-radius: var(--r-lg);
  background-color: rgb(var(--n-100) / 0.6);
}

.dark .lp-trace {
  border-color: rgb(var(--line) / var(--line-a));
  background-color: rgb(var(--d-900) / 0.7);
}

.lp-trace svg {
  display: block;
  width: 100%;
  height: auto;
  overflow: visible;
}

.lp-trace-node rect {
  fill: rgb(var(--n-50));
  stroke: rgb(var(--line-light) / 0.16);
  stroke-width: 1;
}

.dark .lp-trace-node rect {
  fill: rgb(var(--d-800));
  stroke: rgb(var(--line) / var(--line-a-hover));
}

.lp-trace-node-accent rect {
  fill: rgb(var(--p-500) / 0.1);
  stroke: rgb(var(--p-500) / 0.38);
}

.dark .lp-trace-node-accent rect {
  fill: rgb(var(--p-500) / 0.14);
  stroke: rgb(var(--p-500) / 0.45);
}

.lp-trace-node text {
  font-size: 10.5px;
  font-family: var(--font-numeric);
  fill: rgb(var(--n-700));
  letter-spacing: -0.01em;
}

.dark .lp-trace-node text { fill: rgb(var(--n-300)); }

.lp-trace-sub {
  font-size: 9px !important;
  fill: rgb(var(--n-500)) !important;
  letter-spacing: 0.02em;
}

.dark .lp-trace-sub { fill: rgb(var(--n-500)) !important; }

.lp-trace-node-dead {
  stroke-dasharray: 3 3;
  opacity: 0.55;
}

.lp-trace-path {
  fill: none;
  stroke: rgb(var(--line-light) / 0.22);
  stroke-width: 1.25;
  color: rgb(var(--line-light) / 0.22);
}

.dark .lp-trace-path {
  stroke: rgb(var(--line) / 0.3);
  color: rgb(var(--line) / 0.3);
}

.lp-trace-path-dead {
  stroke-dasharray: 3 4;
  opacity: 0.5;
}

.lp-trace-packet {
  fill: rgb(var(--p-500));
  filter: drop-shadow(0 0 4px rgb(var(--p-500) / 0.7));
}

.lp-trace-flag circle {
  fill: rgb(var(--danger));
}

.lp-trace-flag text {
  font-size: 11px;
  font-weight: 800;
  fill: #fff;
}

.lp-trace-caption {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
  margin-top: 0.875rem;
  padding-top: 0.875rem;
  border-top: 1px solid rgb(var(--line-light) / var(--line-a-light));
}

.dark .lp-trace-caption {
  border-top-color: rgb(var(--line) / var(--line-a));
}

.lp-trace-line {
  display: flex;
  gap: 0.625rem;
  font-family: var(--font-numeric);
  font-size: 10.5px;
  line-height: 1.5;
  color: rgb(var(--n-500));
}

.lp-trace-time {
  flex-shrink: 0;
  width: 2.75rem;
  color: rgb(var(--n-400));
  font-variant-numeric: tabular-nums;
}

.lp-trace-warn {
  color: rgb(var(--warning));
}

.dark .lp-trace-warn {
  color: rgb(var(--warning));
}

/* ── 指标条：并列定义，靠一条上边线收口 ───────────────────────────────── */

.lp-stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0;
  margin: 4.5rem 0 0;
  border-top: 1px solid rgb(var(--line-light) / var(--line-a-light));
}

@media (min-width: 768px) {
  .lp-stats { grid-template-columns: repeat(4, minmax(0, 1fr)); }
}

.dark .lp-stats {
  border-top-color: rgb(var(--line) / var(--line-a));
}

.lp-stats > div {
  padding: 1.5rem 1.25rem 1.75rem 0;
}

.lp-stats > div + div {
  padding-left: 1.25rem;
  border-left: 1px solid rgb(var(--line-light) / var(--line-a-light));
}

.dark .lp-stats > div + div {
  border-left-color: rgb(var(--line) / var(--line-a));
}

@media (max-width: 767px) {
  .lp-stats > div:nth-child(odd) {
    padding-left: 0;
    border-left: 0;
  }
  .lp-stats > div:nth-child(even) {
    padding-left: 1.25rem;
    border-left: 1px solid rgb(var(--line-light) / var(--line-a-light));
  }
  .dark .lp-stats > div:nth-child(even) {
    border-left-color: rgb(var(--line) / var(--line-a));
  }
}

.lp-stats dt {
  font-size: 0.75rem;
  font-weight: 500;
  color: rgb(var(--n-500));
}

.lp-stats dd {
  margin: 0.5rem 0 0;
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.035em;
  font-variant-numeric: tabular-nums;
  line-height: 1;
}

/* ── 分区 ──────────────────────────────────────────────────────────────── */

.lp-section {
  padding-block: 5rem;
  border-top: 1px solid rgb(var(--line-light) / var(--line-a-light));
}

.dark .lp-section {
  border-top-color: rgb(var(--line) / var(--line-a));
}

.lp-band {
  background-color: rgb(var(--n-100) / 0.72);
}

.dark .lp-band {
  background-color: rgb(var(--d-900) / 0.55);
}

.lp-sec-head {
  margin-bottom: 2.5rem;
  max-width: 40rem;
}

.lp-sec-head-row {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  justify-content: space-between;
  gap: 1rem;
  max-width: none;
}

.lp-sec-head h2 {
  margin: 0;
  font-size: 1.375rem;
  font-weight: 750;
  letter-spacing: -0.03em;
  line-height: 1.2;
}

.lp-sec-head p {
  margin: 0.625rem 0 0;
  font-size: 0.875rem;
  line-height: 1.65;
  color: rgb(var(--n-600));
  text-wrap: pretty;
}

.dark .lp-sec-head p { color: rgb(var(--n-400)); }

.lp-link {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.8125rem;
  font-weight: 600;
  color: rgb(var(--p-700));
  text-decoration: none;
  white-space: nowrap;
}

.lp-link:hover { text-decoration: underline; }

.dark .lp-link { color: rgb(var(--p-400)); }

/* ── 上游：一行流式排列，不是卡片 ──────────────────────────────────────── */

.lp-providers {
  display: block;
}

.lp-provider-mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  border-radius: 999px;
  font-size: 0.6875rem;
  font-weight: 800;
  color: #fff;
  flex-shrink: 0;
}

.lp-provider-name {
  font-size: 0.8125rem;
  font-weight: 600;
  letter-spacing: -0.005em;
}

.lp-provider-count {
  font-family: var(--font-numeric);
  font-size: 0.6875rem;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--n-500));
}

/* 上游：左侧说明 + 右侧分布条，两栏吃掉整行宽度 */
.lp-prov-split {
  display: grid;
  gap: 2.5rem;
  align-items: start;
}

@media (min-width: 900px) {
  .lp-prov-split {
    grid-template-columns: minmax(0, 22rem) minmax(0, 1fr);
    gap: 4rem;
  }
}

.lp-sec-head-flush { margin-bottom: 0; }

.lp-provider {
  display: grid;
  grid-template-columns: 1.5rem minmax(0, 1fr) auto;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 0;
  border-bottom: 1px solid rgb(var(--line-light) / var(--line-a-light));
  border-radius: 0;
  background-color: transparent;
}

.dark .lp-provider {
  border-bottom-color: rgb(var(--line) / var(--line-a));
  background-color: transparent;
}

.lp-provider:last-child { border-bottom: 0; }

.lp-provider-body {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
  min-width: 0;
}

.lp-provider-share {
  display: block;
  height: 2px;
  width: 100%;
  background-color: rgb(var(--line-light) / 0.1);
  border-radius: 999px;
  overflow: hidden;
}

.dark .lp-provider-share {
  background-color: rgb(255 255 255 / 0.08);
}

.lp-provider-bar {
  display: block;
  height: 100%;
  background-color: rgb(var(--p-500) / 0.55);
  border-radius: 999px;
}

/* ── 费率表 ────────────────────────────────────────────────────────────── */

.lp-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8125rem;
}

.lp-table thead th {
  padding: 0 0 0.625rem;
  text-align: left;
  font-size: 0.6875rem;
  font-weight: 600;
  color: rgb(var(--n-500));
  border-bottom: 1px solid rgb(var(--line-light) / var(--line-a-light));
}

.dark .lp-table thead th {
  border-bottom-color: rgb(var(--line) / var(--line-a));
}

.lp-table tbody td {
  padding: 0.8125rem 0;
  border-bottom: 1px solid rgb(var(--line-light) / var(--line-a-light));
  vertical-align: baseline;
}

.dark .lp-table tbody td {
  border-bottom-color: rgb(var(--line) / var(--line-a));
}

.lp-table tbody tr:last-child td { border-bottom: 0; }

.lp-mono {
  font-family: var(--font-numeric);
  letter-spacing: -0.01em;
}

.lp-muted { color: rgb(var(--n-500)); }

.lp-num {
  text-align: right;
  font-variant-numeric: tabular-nums;
}

.lp-empty {
  margin: 0;
  font-size: 0.8125rem;
  color: rgb(var(--n-500));
}

/* ── 能力：2×2，只用发丝线分隔，没有容器 ──────────────────────────────── */

.lp-cap-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 0;
  border-top: 1px solid rgb(var(--line-light) / var(--line-a-light));
}

.dark .lp-cap-grid {
  border-top-color: rgb(var(--line) / var(--line-a));
}

@media (min-width: 768px) {
  .lp-cap-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

.lp-cap {
  padding: 2rem 2rem 2.25rem 0;
  border-bottom: 1px solid rgb(var(--line-light) / var(--line-a-light));
}

.dark .lp-cap {
  border-bottom-color: rgb(var(--line) / var(--line-a));
}

@media (min-width: 768px) {
  .lp-cap { padding-right: 2.5rem; }
  .lp-cap:nth-child(even) {
    padding-right: 0;
    padding-left: 2.5rem;
    border-left: 1px solid rgb(var(--line-light) / var(--line-a-light));
  }
  .dark .lp-cap:nth-child(even) {
    border-left-color: rgb(var(--line) / var(--line-a));
  }
  .lp-cap:nth-last-child(-n + 2) { border-bottom: 0; }
}

.lp-cap h3 {
  margin: 0;
  font-size: 0.9375rem;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.lp-cap p {
  margin: 0.625rem 0 0;
  font-size: 0.8125rem;
  line-height: 1.7;
  color: rgb(var(--n-600));
  text-wrap: pretty;
}

.dark .lp-cap p { color: rgb(var(--n-400)); }

/* ── 接入：三步是真序列，所以用编号 ──────────────────────────────────── */

.lp-quick {
  display: grid;
  gap: 3rem;
  align-items: start;
}

@media (min-width: 1024px) {
  .lp-quick { grid-template-columns: minmax(0, 1fr) minmax(0, 29rem); gap: 4rem; }
}

.lp-steps {
  list-style: none;
  margin: 0 0 1.75rem;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.lp-steps li {
  display: grid;
  grid-template-columns: 2rem minmax(0, 1fr);
  gap: 0.875rem;
  align-items: baseline;
}

.lp-step-num {
  font-family: var(--font-numeric);
  font-size: 0.6875rem;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: rgb(var(--p-600));
  padding-top: 0.125rem;
}

.dark .lp-step-num { color: rgb(var(--p-400)); }

.lp-steps h3 {
  margin: 0;
  font-size: 0.875rem;
  font-weight: 650;
  letter-spacing: -0.015em;
}

.lp-steps p {
  margin: 0.25rem 0 0;
  font-size: 0.8125rem;
  line-height: 1.65;
  color: rgb(var(--n-600));
}

.dark .lp-steps p { color: rgb(var(--n-400)); }

.lp-code {
  border: 1px solid rgb(var(--line-light) / var(--line-a-light));
  border-radius: var(--r-lg);
  overflow: hidden;
  background-color: rgb(var(--n-100) / 0.7);
}

.dark .lp-code {
  border-color: rgb(var(--line) / var(--line-a));
  background-color: rgb(var(--d-900));
}

.lp-code-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 0.75rem 0.5rem 0.875rem;
  border-bottom: 1px solid rgb(var(--line-light) / var(--line-a-light));
}

.dark .lp-code-bar {
  border-bottom-color: rgb(var(--line) / var(--line-a));
}

.lp-code-file {
  font-family: var(--font-numeric);
  font-size: 0.6875rem;
  color: rgb(var(--n-500));
}

.lp-code-copy {
  border: 0;
  background: transparent;
  padding: 0.25rem 0.5rem;
  border-radius: 5px;
  font-size: 0.6875rem;
  font-weight: 600;
  color: rgb(var(--n-600));
  cursor: pointer;
  transition: background-color 120ms ease, color 120ms ease;
}

.lp-code-copy:hover {
  background-color: rgb(var(--line-light) / 0.07);
  color: rgb(var(--n-900));
}

.dark .lp-code-copy { color: rgb(var(--n-400)); }
.dark .lp-code-copy:hover {
  background-color: rgb(255 255 255 / 0.06);
  color: rgb(var(--n-50));
}

.lp-code-body {
  margin: 0;
  padding: 1rem 1.125rem;
  overflow-x: auto;
  font-family: var(--font-numeric);
  font-size: 0.75rem;
  line-height: 1.8;
  color: rgb(var(--n-700));
  tab-size: 2;
}

.dark .lp-code-body { color: rgb(var(--n-300)); }

.lp-code-foot {
  margin: 0;
  padding: 0.625rem 1.125rem 0.75rem;
  border-top: 1px solid rgb(var(--line-light) / var(--line-a-light));
  font-size: 0.6875rem;
  color: rgb(var(--n-500));
}

.dark .lp-code-foot {
  border-top-color: rgb(var(--line) / var(--line-a));
}

/* ── 收尾：全页唯一的高饱和面 ─────────────────────────────────────────── */

.lp-final {
  border-top: 1px solid rgb(var(--line-light) / var(--line-a-light));
  background-color: rgb(var(--p-950));
  color: rgb(var(--n-50));
}

.dark .lp-final {
  border-top-color: rgb(var(--line) / var(--line-a));
  background-color: rgb(var(--p-950));
}

.lp-final-inner {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 2rem;
  padding-block: 3.5rem;
}

.lp-final h2 {
  margin: 0;
  font-size: clamp(1.5rem, 2.6vw, 2rem);
  font-weight: 780;
  letter-spacing: -0.035em;
  line-height: 1.15;
}

.lp-final p {
  margin: 0.75rem 0 0;
  max-width: 34rem;
  font-size: 0.875rem;
  line-height: 1.65;
  color: rgb(var(--p-200));
}

/* 深绿底上的按钮反差：亮绿底 + 深墨绿字，对比度约 8:1 */
.lp-final .lp-btn-primary {
  background-color: rgb(var(--p-300));
  color: rgb(var(--p-950));
}

.lp-final .lp-btn-primary:hover {
  background-color: rgb(var(--p-200));
}

/* ── 页脚 ──────────────────────────────────────────────────────────────── */

.lp-footer {
  border-top: 1px solid rgb(var(--line-light) / var(--line-a-light));
  background-color: rgb(var(--n-50));
}

.dark .lp-footer {
  border-top-color: rgb(var(--line) / var(--line-a));
  background-color: rgb(var(--d-950));
}

.lp-footer-inner {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 1rem 1.5rem;
  padding-block: 1.75rem;
}

.lp-footer-brand {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.8125rem;
  font-weight: 650;
}

.lp-footer-sub {
  font-weight: 400;
  color: rgb(var(--n-500));
}

.lp-footer-links {
  display: flex;
  flex-wrap: wrap;
  gap: 1.25rem;
  margin-left: auto;
}

.lp-footer-links a {
  font-size: 0.8125rem;
  color: rgb(var(--n-600));
  text-decoration: none;
}

.lp-footer-links a:hover { color: rgb(var(--n-900)); }

.dark .lp-footer-links a { color: rgb(var(--n-400)); }
.dark .lp-footer-links a:hover { color: rgb(var(--n-50)); }

.lp-footer-copy {
  margin: 0;
  font-size: 0.75rem;
  color: rgb(var(--n-400));
  flex-basis: 100%;
}

/* ── 动效：只有一处，且可关掉 ─────────────────────────────────────────── */

@media (prefers-reduced-motion: reduce) {
  .lp-trace svg animateMotion {
    display: none;
  }
  .lp * {
    transition-duration: 0.01ms !important;
  }
}
</style>
