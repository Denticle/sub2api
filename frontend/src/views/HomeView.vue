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
          <a
            v-if="docUrl"
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

  <!-- Default Home Page — OpenRouter 式版面：通栏导航 + 居中主视觉 + 分区内容 -->
  <div v-else class="flex min-h-screen flex-col bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-gray-100">
    <!-- Header -->
    <header class="sticky top-0 z-30 border-b border-gray-200/70 bg-gray-50/85 backdrop-blur-xl dark:border-dark-800/80 dark:bg-dark-950/85">
      <nav class="flex h-14 w-full items-center gap-3 px-6 lg:px-10">
        <router-link to="/home" class="flex shrink-0 items-center gap-2">
          <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-7 w-7 rounded-lg object-contain" />
          <span class="text-[15px] font-bold tracking-tight">{{ siteName }}</span>
        </router-link>

        <div class="ml-auto flex items-center gap-1">
          <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer" class="hidden rounded-lg px-2.5 py-1.5 text-[13px] font-medium text-gray-600 transition-colors hover:bg-gray-200/60 hover:text-gray-900 sm:block dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white">{{ t('home.docs') }}</a>
          <router-link v-if="showModelPlazaEntry" to="/model-plaza" class="hidden rounded-lg px-2.5 py-1.5 text-[13px] font-medium text-gray-600 transition-colors hover:bg-gray-200/60 hover:text-gray-900 sm:block dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white">{{ t('nav.modelPlaza') }}</router-link>
          <LocaleSwitcher />
          <button class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-200/60 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white" :title="isDark ? t('home.switchToLight') : t('home.switchToDark')" @click="toggleTheme">
            <Icon v-if="isDark" name="sun" size="sm" />
            <Icon v-else name="moon" size="sm" />
          </button>
          <router-link v-if="isAuthenticated" :to="dashboardPath" class="btn btn-primary ml-1">{{ t('home.dashboard') }}</router-link>
          <router-link v-else to="/login" class="btn btn-primary ml-1">{{ t('home.login') }}</router-link>
        </div>
      </nav>
    </header>

    <main class="flex-1">
      <!-- Hero -->
      <section class="relative overflow-hidden border-b border-gray-200/70 dark:border-dark-800/80">
        <div class="pointer-events-none absolute inset-x-0 -top-48 h-96 bg-[radial-gradient(ellipse_at_center,rgba(34,197,94,0.12),transparent_65%)]"></div>

        <div class="relative w-full px-6 pb-14 pt-16 sm:pt-20 lg:px-10">
          <div class="grid items-center gap-12 lg:grid-cols-[minmax(0,1fr)_minmax(0,520px)]">
            <!-- 左：文案 -->
            <div class="text-center lg:text-left">
              <span class="inline-flex items-center gap-1.5 rounded-full border border-gray-200 bg-white px-2.5 py-1 text-[11.5px] font-semibold text-gray-600 dark:border-dark-800 dark:bg-dark-900 dark:text-dark-400">
                <span class="h-1.5 w-1.5 rounded-full bg-primary-500"></span>
                {{ t('home.hero.eyebrow') }}
              </span>

              <h1 class="mt-5 text-[clamp(2rem,4.4vw,3.25rem)] font-extrabold leading-[1.08] tracking-[-0.04em]">
                {{ siteName }}
              </h1>
              <p class="mt-4 max-w-xl text-[clamp(0.95rem,1.4vw,1.0625rem)] leading-relaxed text-gray-600 lg:mx-0 dark:text-dark-400">
                {{ siteSubtitle }}
              </p>
              <p class="mt-2 max-w-xl text-[13px] leading-relaxed text-gray-500 dark:text-dark-500">
                {{ t('home.hero.description') }}
              </p>

              <div class="mt-7 flex flex-wrap items-center justify-center gap-3 lg:justify-start">
                <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="btn btn-primary h-10 px-5 text-[13.5px]">
                  {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
                  <Icon name="arrowRight" size="sm" class="ml-1" />
                </router-link>
                <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer" class="btn btn-secondary h-10 px-5 text-[13.5px]">
                  <Icon name="book" size="sm" class="mr-1" />
                  {{ t('home.viewDocs') }}
                </a>
                <router-link v-else-if="showModelPlazaEntry" to="/model-plaza" class="btn btn-secondary h-10 px-5 text-[13.5px]">
                  {{ t('nav.modelPlaza') }}
                </router-link>
              </div>

              <div class="mt-7 flex flex-wrap items-center justify-center gap-x-6 gap-y-2 text-[12.5px] text-gray-500 lg:justify-start dark:text-dark-500">
                <span v-for="tag in tags" :key="tag.key" class="inline-flex items-center gap-1.5">
                  <Icon :name="tag.icon" size="xs" class="text-primary-600 dark:text-primary-400" />
                  {{ tag.label }}
                </span>
              </div>
            </div>

            <!-- 右：调用示例 -->
            <div class="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm dark:border-dark-800 dark:bg-dark-900">
              <div class="flex items-center gap-1.5 border-b border-gray-200/80 px-3.5 py-2.5 dark:border-dark-800">
                <span class="h-2.5 w-2.5 rounded-full bg-[#ff5f57]"></span>
                <span class="h-2.5 w-2.5 rounded-full bg-[#febc2e]"></span>
                <span class="h-2.5 w-2.5 rounded-full bg-[#28c840]"></span>
                <span class="ml-2 text-[11.5px] font-medium text-gray-400 dark:text-dark-600">request.sh</span>
                <span class="ml-auto rounded bg-gray-100 px-1.5 py-0.5 font-mono text-[10.5px] text-gray-500 dark:bg-dark-800 dark:text-dark-500">cURL</span>
              </div>
              <pre class="overflow-x-auto px-4 py-4 font-mono text-[12px] leading-[1.75] text-gray-600 dark:text-dark-400"><code><span class="text-gray-400 dark:text-dark-600"># 只需换掉 base_url，无需改调用代码</span>
<span class="text-primary-700 dark:text-primary-400">curl</span> {{ baseUrl }}/v1/chat/completions \
  -H <span class="text-amber-600 dark:text-amber-400">"Authorization: Bearer $API_KEY"</span> \
  -H <span class="text-amber-600 dark:text-amber-400">"Content-Type: application/json"</span> \
  -d <span class="text-amber-600 dark:text-amber-400">'{
    "model": "claude-sonnet-4-5",
    "messages": [{"role": "user", "content": "hi"}]
  }'</span></code></pre>
              <div class="border-t border-gray-200/80 px-4 py-2.5 text-[11.5px] text-gray-500 dark:border-dark-800 dark:text-dark-500">
                {{ t('home.hero.compatible') }}
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 指标条 -->
      <section class="border-b border-gray-200/70 dark:border-dark-800/80">
        <div class="grid w-full grid-cols-2 divide-x divide-gray-200/70 sm:grid-cols-4 dark:divide-dark-800/80">
          <div v-for="stat in stats" :key="stat.key" class="px-6 py-6 lg:px-10">
            <div class="font-mono text-[22px] font-semibold tracking-tight">{{ stat.value }}</div>
            <div class="mt-1 text-[12.5px] text-gray-500 dark:text-dark-500">{{ stat.label }}</div>
          </div>
        </div>
      </section>

      <!-- 支持的模型 -->
      <section data-testid="home-providers" class="w-full px-6 pb-16 pt-14 lg:px-10">
        <div class="mb-5 flex items-end justify-between gap-4">
          <div>
            <h2 class="text-[17px] font-bold tracking-tight">{{ t('home.providers.title') }}</h2>
            <p class="mt-1 text-[13px] text-gray-500 dark:text-dark-500">{{ t('home.providers.description') }}</p>
          </div>
          <router-link v-if="showModelPlazaEntry" to="/model-plaza" class="inline-flex items-center gap-1 text-[13px] font-semibold text-primary-700 hover:underline dark:text-primary-400">
            {{ t('home.providers.viewAll') }}
            <Icon name="chevronRight" size="xs" />
          </router-link>
        </div>

        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          <div
            v-for="provider in providers"
            :key="provider.key"
            class="group flex items-start gap-3 rounded-xl border border-gray-200 bg-white p-4 transition-colors hover:border-gray-300 dark:border-dark-800 dark:bg-dark-900 dark:hover:border-dark-700"
          >
            <span class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-[13px] font-bold text-white" :style="{ background: provider.tint }">{{ provider.mark }}</span>
            <div class="min-w-0">
              <div class="flex items-center gap-2">
                <span class="truncate text-[13.5px] font-semibold">{{ provider.name }}</span>
                <span class="shrink-0 rounded-full bg-gray-100 px-1.5 py-0.5 text-[10px] font-bold text-gray-500 dark:bg-dark-800 dark:text-dark-500">{{ provider.tag }}</span>
              </div>
              <p class="mt-1 text-[12.5px] leading-relaxed text-gray-500 dark:text-dark-500">{{ provider.desc }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 三步接入 -->
      <section class="border-y border-gray-200/70 bg-white/60 dark:border-dark-800/80 dark:bg-dark-900/40">
        <div class="w-full px-6 py-16 lg:px-10">
          <div class="max-w-xl">
            <h2 class="text-[17px] font-bold tracking-tight">{{ t('home.steps.title') }}</h2>
            <p class="mt-1.5 text-[13px] text-gray-500 dark:text-dark-500">{{ t('home.steps.subtitle') }}</p>
          </div>

          <div class="mt-9 grid gap-8 md:grid-cols-3">
            <div v-for="(step, i) in steps" :key="step.key">
              <div class="flex items-center gap-2.5">
                <span class="flex h-6 w-6 items-center justify-center rounded-full bg-primary-500/12 text-[11.5px] font-bold text-primary-700 dark:text-primary-400">{{ i + 1 }}</span>
                <h3 class="text-[14px] font-semibold">{{ step.title }}</h3>
              </div>
              <p class="mt-2 pl-[34px] text-[13px] leading-relaxed text-gray-600 dark:text-dark-400">{{ step.desc }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 能力 -->
      <section class="w-full px-6 py-16 lg:px-10">
        <div class="max-w-xl">
          <h2 class="text-[17px] font-bold tracking-tight">{{ t('home.features.title') }}</h2>
          <p class="mt-1.5 text-[13px] text-gray-500 dark:text-dark-500">{{ t('home.features.subtitle') }}</p>
        </div>

        <div class="mt-8 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <div v-for="feature in features" :key="feature.key" class="rounded-xl border border-gray-200 bg-white p-5 dark:border-dark-800 dark:bg-dark-900">
            <span class="flex h-8 w-8 items-center justify-center rounded-lg bg-gray-100 text-gray-600 dark:bg-dark-800 dark:text-dark-400">
              <Icon :name="feature.icon" size="sm" />
            </span>
            <h3 class="mt-3 text-[14px] font-semibold">{{ feature.title }}</h3>
            <p class="mt-1.5 text-[13px] leading-relaxed text-gray-600 dark:text-dark-400">{{ feature.desc }}</p>
          </div>
        </div>
      </section>

      <!-- 对比 -->
      <section class="border-y border-gray-200/70 bg-white/60 dark:border-dark-800/80 dark:bg-dark-900/40">
        <div class="w-full px-6 py-16 lg:px-10">
          <h2 class="text-[17px] font-bold tracking-tight">{{ t('home.comparison.title') }}</h2>

          <div class="mt-6 overflow-x-auto rounded-xl border border-gray-200 dark:border-dark-800">
            <table class="w-full min-w-[640px] border-collapse bg-white text-[13px] dark:bg-dark-900">
              <thead>
                <tr class="border-b border-gray-200 dark:border-dark-800">
                  <th class="px-4 py-3 text-left text-[12px] font-bold text-gray-500 dark:text-dark-500">{{ t('home.comparison.headers.feature') }}</th>
                  <th class="px-4 py-3 text-left text-[12px] font-bold text-gray-500 dark:text-dark-500">{{ t('home.comparison.headers.official') }}</th>
                  <th class="px-4 py-3 text-left text-[12px] font-bold text-primary-700 dark:text-primary-400">{{ siteName }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in comparisonRows" :key="row.key" class="border-b border-gray-100 last:border-0 dark:border-dark-800/60">
                  <td class="px-4 py-3 font-medium text-gray-700 dark:text-dark-300">{{ row.feature }}</td>
                  <td class="px-4 py-3 text-gray-500 dark:text-dark-500">{{ row.official }}</td>
                  <td class="px-4 py-3 font-medium text-gray-900 dark:text-gray-100">{{ row.us }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <!-- 收尾 CTA -->
      <section class="w-full px-6 py-20 lg:px-10">
        <div class="rounded-2xl border border-gray-200 bg-white px-6 py-12 text-center dark:border-dark-800 dark:bg-dark-900">
          <h2 class="text-[22px] font-extrabold tracking-tight">{{ t('home.cta.title') }}</h2>
          <p class="mx-auto mt-2 max-w-lg text-[13px] leading-relaxed text-gray-600 dark:text-dark-400">{{ t('home.cta.description') }}</p>
          <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="btn btn-primary mt-6 inline-flex h-10 px-6 text-[13.5px]">
            {{ isAuthenticated ? t('home.goToDashboard') : t('home.cta.button') }}
          </router-link>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <footer class="border-t border-gray-200/70 dark:border-dark-800/80">
      <div class="w-full px-6 py-10 lg:px-10">
        <div class="flex flex-col gap-8 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <div class="flex items-center gap-2">
              <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-6 w-6 rounded-md object-contain" />
              <span class="text-[13.5px] font-semibold">{{ siteName }}</span>
            </div>
            <p class="mt-2 max-w-xs text-[12.5px] leading-relaxed text-gray-500 dark:text-dark-500">{{ siteSubtitle }}</p>
          </div>

          <div class="flex flex-wrap gap-x-10 gap-y-4 text-[13px]">
            <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer" class="text-gray-600 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white">{{ t('home.docs') }}</a>
            <router-link v-if="showModelPlazaEntry" to="/model-plaza" class="text-gray-600 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white">{{ t('nav.modelPlaza') }}</router-link>
            <router-link to="/key-usage" class="text-gray-600 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white">{{ t('home.footer.keyUsage') }}</router-link>
            <a :href="githubUrl" target="_blank" rel="noopener noreferrer" class="text-gray-600 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white">GitHub</a>
          </div>
        </div>

        <p class="mt-8 text-[12px] text-gray-500 dark:text-dark-500">&copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}</p>
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
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const compactHomeEnabled = computed(() => appStore.cachedPublicSettings?.compact_home_enabled === true)
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// GitHub URL
const githubUrl = 'https://github.com/Wei-Shaw/sub2api'

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const modelPlazaRequiresAuth = computed(
  () => appStore.cachedPublicSettings?.model_plaza_require_auth === true,
)
const showModelPlazaEntry = computed(
  () => modelPlazaEnabled.value && (isAuthenticated.value || !modelPlazaRequiresAuth.value),
)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
type IconName = InstanceType<typeof Icon>['$props']['name']

// 指标条
const stats = computed(() => [
  { key: 'providers', value: String(providers.value.length), label: t('home.stats.providers') },
  { key: 'protocol', value: '4', label: t('home.stats.protocols') },
  { key: 'billing', value: '24h', label: t('home.stats.billing') },
  { key: 'uptime', value: '99.9%', label: t('home.stats.uptime') },
])

// 顶部标签
const tags = computed(() => [
  { key: 'sub2api', icon: 'swap' as IconName, label: t('home.tags.subscriptionToApi') },
  { key: 'sticky', icon: 'shield' as IconName, label: t('home.tags.stickySession') },
  { key: 'billing', icon: 'chart' as IconName, label: t('home.tags.realtimeBilling') },
])

// 调用的 base_url：展示站点自身的地址，示例才不是假的
const baseUrl = computed(() => {
  if (typeof window === 'undefined') return 'https://api.example.com'
  return window.location.origin
})

// 首页内容块：不写死在模板里，便于按语言维护文案
const providers = computed<{
  key: string
  mark: string
  tint: string
  name: string
  desc: string
  tag: string
}[]>(() => [
  {
    key: 'claude',
    mark: 'C',
    tint: 'linear-gradient(135deg,#d97757,#b85a3c)',
    name: t('home.providers.claude'),
    desc: t('home.providers.claudeDesc'),
    tag: t('home.providers.tagChat'),
  },
  {
    key: 'gpt',
    mark: 'G',
    tint: 'linear-gradient(135deg,#10a37f,#0a7c5f)',
    name: 'GPT',
    desc: t('home.providers.gptDesc'),
    tag: t('home.providers.tagChat'),
  },
  {
    key: 'gemini',
    mark: 'G',
    tint: 'linear-gradient(135deg,#4285f4,#2b6de0)',
    name: t('home.providers.gemini'),
    desc: t('home.providers.geminiDesc'),
    tag: t('home.providers.tagMultimodal'),
  },
  {
    key: 'antigravity',
    mark: 'A',
    tint: 'linear-gradient(135deg,#f43f5e,#be123c)',
    name: t('home.providers.antigravity'),
    desc: t('home.providers.antigravityDesc'),
    tag: t('home.providers.tagAgent'),
  },
  {
    key: 'deepseek',
    mark: 'D',
    tint: 'linear-gradient(135deg,#4d6bfe,#3350d8)',
    name: 'DeepSeek',
    desc: t('home.providers.deepseekDesc'),
    tag: t('home.providers.tagChat'),
  },
  {
    key: 'grok',
    mark: 'X',
    tint: 'linear-gradient(135deg,#1f2937,#0b1016)',
    name: 'Grok',
    desc: t('home.providers.grokDesc'),
    tag: t('home.providers.tagChat'),
  },
  {
    key: 'kimi',
    mark: 'K',
    tint: 'linear-gradient(135deg,#6366f1,#4338ca)',
    name: 'Kimi',
    desc: t('home.providers.kimiDesc'),
    tag: t('home.providers.tagLongContext'),
  },
  {
    key: 'glm',
    mark: 'Z',
    tint: 'linear-gradient(135deg,#0ea5e9,#0369a1)',
    name: 'Zhipu GLM',
    desc: t('home.providers.glmDesc'),
    tag: t('home.providers.tagChat'),
  },
])

const steps = computed(() => [
  { key: 'signup', title: t('home.steps.signup'), desc: t('home.steps.signupDesc') },
  { key: 'credits', title: t('home.steps.credits'), desc: t('home.steps.creditsDesc') },
  { key: 'key', title: t('home.steps.key'), desc: t('home.steps.keyDesc') },
])

const features = computed<{ key: string; icon: IconName; title: string; desc: string }[]>(() => [
  {
    key: 'gateway',
    icon: 'server',
    title: t('home.features.unifiedGateway'),
    desc: t('home.features.unifiedGatewayDesc'),
  },
  {
    key: 'pool',
    icon: 'users',
    title: t('home.features.multiAccount'),
    desc: t('home.features.multiAccountDesc'),
  },
  {
    key: 'billing',
    icon: 'chart',
    title: t('home.features.balanceQuota'),
    desc: t('home.features.balanceQuotaDesc'),
  },
])

// 对比表
interface ComparisonRow { key: string; feature: string; official: string; us: string }

// 逐字段取文案，不依赖 tm() 的对象返回形态（不同 i18n 配置下不稳定）
const comparisonRows = computed<ComparisonRow[]>(() => {
  const keys = ['pricing', 'models', 'management', 'stability', 'control'] as const
  return keys.map((key) => ({
    key,
    feature: t(`home.comparison.items.${key}.feature`),
    official: t(`home.comparison.items.${key}.official`),
    us: t(`home.comparison.items.${key}.us`),
  }))
})

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// Initialize theme
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

  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
/* 首页不需要额外样式：版式全部由 theme 层与 Tailwind 工具类承担。
   上一版的终端动画（.terminal-* / .code-line 等）随版面改为一屏式落地页后
   已无引用，删掉避免死代码。 */
</style>
