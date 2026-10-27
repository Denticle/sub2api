<template>
  <div class="docs-page">
    <header class="docs-topbar">
      <router-link to="/home" class="docs-brand">
        <img :src="siteLogo || '/logo.svg'" alt="" class="docs-brand-mark" />
        <span class="docs-brand-name">{{ siteName }}</span>
      </router-link>

      <div class="docs-topbar-actions">
        <LocaleSwitcher />
        <button
          type="button"
          class="docs-icon-btn"
          :aria-label="isDark ? t('nav.lightMode') : t('nav.darkMode')"
          :title="isDark ? t('nav.lightMode') : t('nav.darkMode')"
          @click="toggleTheme"
        >
          <Icon :name="isDark ? 'sun' : 'moon'" size="sm" />
        </button>
        <router-link :to="authTarget" class="docs-login">{{ authLabel }}</router-link>
      </div>
    </header>

    <div class="docs-layout">
      <aside class="docs-sidebar">
        <nav class="docs-nav" :aria-label="t('docs.nav.menuLabel')">
          <p class="docs-nav-section">{{ t('docs.nav.overview') }}</p>
          <a href="#intro" :class="{ active: activeSection === 'intro' }">{{ t('docs.nav.intro') }}</a>
          <a href="#quickstart" :class="{ active: activeSection === 'quickstart' }">{{ t('docs.nav.quickstart') }}</a>
          <p class="docs-nav-section">{{ t('docs.nav.apiReference') }}</p>
          <a href="#endpoints" :class="{ active: activeSection === 'endpoints' }">{{ t('docs.nav.endpoints') }}</a>
          <a href="#auth" :class="{ active: activeSection === 'auth' }">{{ t('docs.nav.auth') }}</a>
          <p class="docs-nav-section">{{ t('docs.nav.integration') }}</p>
          <a href="#codex-cli" :class="{ active: activeSection === 'codex-cli' }">{{ t('docs.nav.codexCli') }}</a>
          <a href="#codex-ws" :class="{ active: activeSection === 'codex-ws' }">{{ t('docs.nav.codexWs') }}</a>
          <a href="#claude-code" :class="{ active: activeSection === 'claude-code' }">{{ t('docs.nav.claudeCode') }}</a>
          <a href="#opencode" :class="{ active: activeSection === 'opencode' }">{{ t('docs.nav.opencode') }}</a>
          <p class="docs-nav-section">{{ t('docs.nav.tools') }}</p>
          <a href="#ccswitch" :class="{ active: activeSection === 'ccswitch' }">{{ t('docs.nav.ccswitch') }}</a>
          <p class="docs-nav-section">{{ t('docs.nav.other') }}</p>
          <a href="#usage" :class="{ active: activeSection === 'usage' }">{{ t('docs.nav.usage') }}</a>
          <a href="#faq" :class="{ active: activeSection === 'faq' }">{{ t('docs.nav.faq') }}</a>
        </nav>
      </aside>

      <main class="docs-main">
        <header class="docs-hero">
          <p class="docs-badge">
            <span class="docs-badge-dot" aria-hidden="true"></span>
            {{ t('docs.badge') }}
          </p>
          <h1 class="docs-h1">
            <span>{{ t('docs.heroTitleLead') }}</span>
            <span class="docs-h1-accent">{{ t('docs.heroTitleAccent', { site: siteName }) }}</span>
          </h1>
          <p class="docs-hero-desc">{{ t('docs.heroDesc') }}</p>
          <div class="docs-hero-meta">
            <span>{{ t('docs.metaOpenAI') }}</span>
            <span>{{ t('docs.metaAnthropic') }}</span>
            <span>{{ t('docs.metaBilling') }}</span>
          </div>
        </header>

        <!-- 01 简介 -->
        <section id="intro" class="docs-section">
          <div class="docs-section-head">
            <span class="docs-section-num">01</span>
            <h2 class="docs-h2">{{ t('docs.intro.title') }}</h2>
          </div>
          <p>{{ t('docs.intro.lead', { site: siteName }) }}</p>
          <div class="docs-feature-grid">
            <div v-for="f in introFeatures" :key="f.title" class="docs-feature">
              <h3 class="docs-feature-title">{{ f.title }}</h3>
              <p class="docs-feature-desc">{{ f.desc }}</p>
            </div>
          </div>
          <p>
            {{ t('docs.intro.siteLabel') }}
            <a :href="siteRoot" target="_blank" rel="noopener noreferrer" class="docs-link">{{ siteRoot }}</a>
          </p>
        </section>

        <!-- 02 快速开始 -->
        <section id="quickstart" class="docs-section">
          <div class="docs-section-head">
            <span class="docs-section-num">02</span>
            <h2 class="docs-h2">{{ t('docs.quickstart.title') }}</h2>
          </div>

          <h3 class="docs-h3">2.1 {{ t('docs.quickstart.registerTitle') }}</h3>
          <p>{{ t('docs.quickstart.registerDesc', { site: siteRoot }) }}</p>

          <h3 class="docs-h3">2.2 {{ t('docs.quickstart.keyTitle') }}</h3>
          <ol class="docs-steps">
            <li>
              <span class="docs-step-num">1</span>
              <i18n-t keypath="docs.quickstart.keyStep1" tag="span" scope="global">
                <template #strong><strong>{{ t('docs.quickstart.keyStep1Strong') }}</strong></template>
              </i18n-t>
            </li>
            <li>
              <span class="docs-step-num">2</span>
              <i18n-t keypath="docs.quickstart.keyStep2" tag="span" scope="global">
                <template #strong><strong>{{ t('docs.quickstart.keyStep2Strong') }}</strong></template>
              </i18n-t>
            </li>
            <li>
              <span class="docs-step-num">3</span>
              <i18n-t keypath="docs.quickstart.keyStep3" tag="span" scope="global">
                <template #code><code class="docs-code-inline">sk-</code></template>
              </i18n-t>
            </li>
          </ol>

          <h3 class="docs-h3">2.3 {{ t('docs.quickstart.envTitle') }}</h3>
          <p>{{ t('docs.quickstart.envDesc') }}</p>
          <CodeBlock lang="bash" :code="envSnippet" />
        </section>

        <!-- 03 API 端点 -->
        <section id="endpoints" class="docs-section">
          <div class="docs-section-head">
            <span class="docs-section-num">03</span>
            <h2 class="docs-h2">{{ t('docs.endpoints.title') }}</h2>
          </div>
          <p>{{ t('docs.endpoints.desc', { site: siteName }) }}</p>

          <div class="docs-endpoint-card">
            <span class="docs-endpoint-badge">{{ t('docs.endpoints.badge') }}</span>
            <p class="docs-endpoint-label">{{ t('docs.endpoints.label') }}</p>
            <code class="docs-endpoint-url">{{ endpoint }}</code>
            <p class="docs-endpoint-desc">{{ t('docs.endpoints.cardDesc') }}</p>
          </div>

          <div class="docs-notice">
            <p class="docs-notice-title">{{ t('docs.endpoints.noticeTitle') }}</p>
            <i18n-t keypath="docs.endpoints.noticeText" tag="p" scope="global" class="docs-notice-text">
              <template #code><code class="docs-code-inline">{{ endpoint }}</code></template>
              <template #code2><code class="docs-code-inline">/v1</code></template>
            </i18n-t>
          </div>
        </section>

        <!-- 04 认证方式 -->
        <section id="auth" class="docs-section">
          <div class="docs-section-head">
            <span class="docs-section-num">04</span>
            <h2 class="docs-h2">{{ t('docs.auth.title') }}</h2>
          </div>
          <p>{{ t('docs.auth.desc') }}</p>

          <h3 class="docs-h3">4.1 {{ t('docs.auth.formatTitle') }}</h3>
          <i18n-t keypath="docs.auth.formatDesc" tag="p" scope="global">
            <template #code><code class="docs-code-inline">sk-</code></template>
          </i18n-t>
          <CodeBlock lang="" code="sk-************************" />

          <h3 class="docs-h3">4.2 {{ t('docs.auth.headerTitle') }}</h3>
          <CodeBlock lang="" code="Authorization: Bearer YOUR_API_KEY" />

          <h3 class="docs-h3">4.3 {{ t('docs.auth.curlTitle') }}</h3>
          <CodeBlock lang="bash" :code="curlSnippet" />
        </section>

        <!-- 05 集成方式 -->
        <section id="codex-cli" class="docs-section">
          <div class="docs-section-head">
            <span class="docs-section-num">05</span>
            <h2 class="docs-h2">{{ t('docs.integration.title') }}</h2>
          </div>

          <h3 class="docs-h3"><span class="docs-h3-num">5.1</span> {{ t('docs.integration.codexTitle') }}</h3>
          <p>{{ t('docs.integration.codexDesc') }}</p>
          <ol class="docs-steps">
            <li>
              <span class="docs-step-num">1</span>
              <i18n-t keypath="docs.integration.codexStep1" tag="span" scope="global">
                <template #code><code class="docs-code-inline">~/.codex/config.toml</code></template>
              </i18n-t>
            </li>
            <li>
              <span class="docs-step-num">2</span>
              <i18n-t keypath="docs.integration.codexStep2" tag="span" scope="global">
                <template #code><code class="docs-code-inline">base_url</code></template>
              </i18n-t>
            </li>
            <li>
              <span class="docs-step-num">3</span>
              <i18n-t keypath="docs.integration.codexStep3" tag="span" scope="global">
                <template #code><code class="docs-code-inline">~/.codex/auth.json</code></template>
              </i18n-t>
            </li>
          </ol>
          <p class="docs-strong-line">{{ t('docs.integration.codexExample') }}</p>
          <CodeBlock lang="toml" :code="codexConfig" />
          <CodeBlock lang="json" :code="codexAuth" />
        </section>

        <section id="codex-ws" class="docs-section">
          <h3 class="docs-h3"><span class="docs-h3-num">5.2</span> {{ t('docs.integration.codexWsTitle') }}</h3>
          <p>{{ t('docs.integration.codexWsDesc') }}</p>
          <CodeBlock lang="toml" :code="codexWsConfig" />
          <div class="docs-notice">
            <p class="docs-notice-title">{{ t('docs.integration.codexWsNoticeTitle') }}</p>
            <i18n-t keypath="docs.integration.codexWsNoticeText" tag="p" scope="global" class="docs-notice-text">
              <template #code><code class="docs-code-inline">base_url</code></template>
              <template #code2><code class="docs-code-inline">wss://</code></template>
            </i18n-t>
          </div>
        </section>

        <section id="claude-code" class="docs-section">
          <h3 class="docs-h3"><span class="docs-h3-num">5.3</span> {{ t('docs.integration.claudeTitle') }}</h3>
          <p>{{ t('docs.integration.claudeDesc') }}</p>

          <h4 class="docs-h4">{{ t('docs.integration.claudeEnvTitle') }}</h4>
          <CodeBlock lang="bash" :code="claudeEnv" />

          <h4 class="docs-h4">{{ t('docs.integration.claudeFileTitle') }}</h4>
          <CodeBlock lang="json" :code="claudeFile" />

          <div class="docs-notice">
            <p class="docs-notice-title">{{ t('docs.integration.claudeNoticeTitle') }}</p>
            <i18n-t keypath="docs.integration.claudeNoticeText" tag="p" scope="global" class="docs-notice-text">
              <template #code><code class="docs-code-inline">CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1</code></template>
            </i18n-t>
          </div>
        </section>

        <section id="opencode" class="docs-section">
          <h3 class="docs-h3"><span class="docs-h3-num">5.4</span> {{ t('docs.integration.opencodeTitle') }}</h3>
          <p>{{ t('docs.integration.opencodeDesc') }}</p>

          <h4 class="docs-h4">{{ t('docs.integration.opencodeEnvTitle') }}</h4>
          <CodeBlock lang="bash" :code="opencodeEnv" />

          <h4 class="docs-h4">{{ t('docs.integration.opencodeFileTitle') }}</h4>
          <CodeBlock lang="json" :code="opencodeFile" />

          <div class="docs-notice">
            <p class="docs-notice-title">{{ t('docs.integration.opencodeNoticeTitle') }}</p>
            <i18n-t keypath="docs.integration.opencodeNoticeText" tag="p" scope="global" class="docs-notice-text">
              <template #code><code class="docs-code-inline">baseURL</code></template>
              <template #code2><code class="docs-code-inline">/v1</code></template>
              <template #code3><code class="docs-code-inline">OPENAI_API_KEY</code></template>
              <template #code4><code class="docs-code-inline">apiKey</code></template>
            </i18n-t>
          </div>
        </section>

        <!-- 06 CC-Switch -->
        <section id="ccswitch" class="docs-section">
          <div class="docs-section-head">
            <span class="docs-section-num">06</span>
            <h2 class="docs-h2">{{ t('docs.ccswitch.title') }}</h2>
          </div>
          <p>{{ t('docs.ccswitch.desc') }}</p>

          <h3 class="docs-h3"><span class="docs-h3-num">6.1</span> {{ t('docs.ccswitch.importTitle') }}</h3>
          <p>{{ t('docs.ccswitch.importDesc', { site: siteName }) }}</p>
          <ol class="docs-steps">
            <li>
              <span class="docs-step-num">1</span>
              <i18n-t keypath="docs.ccswitch.importStep1" tag="span" scope="global">
                <template #strong><strong>{{ t('docs.ccswitch.importStep1Strong') }}</strong></template>
              </i18n-t>
            </li>
            <li>
              <span class="docs-step-num">2</span>
              <i18n-t keypath="docs.ccswitch.importStep2" tag="span" scope="global">
                <template #strong><strong>{{ t('docs.ccswitch.importStep2Strong') }}</strong></template>
                <template #code><code class="docs-code-inline">ccswitch://</code></template>
              </i18n-t>
            </li>
            <li><span class="docs-step-num">3</span><span>{{ t('docs.ccswitch.importStep3') }}</span></li>
            <li>
              <span class="docs-step-num">4</span>
              <i18n-t keypath="docs.ccswitch.importStep4" tag="span" scope="global">
                <template #strong><strong>{{ t('docs.ccswitch.importStep4Strong') }}</strong></template>
              </i18n-t>
            </li>
          </ol>

          <h3 class="docs-h3"><span class="docs-h3-num">6.2</span> {{ t('docs.ccswitch.configTitle') }}</h3>
          <p>{{ t('docs.ccswitch.configDesc') }}</p>
          <div class="docs-table-wrap">
            <table class="docs-table">
              <thead>
                <tr>
                  <th>{{ t('docs.ccswitch.colItem') }}</th>
                  <th>{{ t('docs.ccswitch.colValue') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>{{ t('docs.ccswitch.rowApp') }}</td>
                  <td>{{ t('docs.ccswitch.rowAppValue') }}</td>
                </tr>
                <tr>
                  <td>{{ t('docs.ccswitch.rowName') }}</td>
                  <td><code class="docs-code-inline">{{ siteName }}</code></td>
                </tr>
                <tr>
                  <td>{{ t('docs.ccswitch.rowSite') }}</td>
                  <td><code class="docs-code-inline">{{ siteRoot }}</code></td>
                </tr>
                <tr>
                  <td>{{ t('docs.ccswitch.rowEndpoint') }}</td>
                  <td><code class="docs-code-inline">{{ endpoint }}</code></td>
                </tr>
                <tr>
                  <td>{{ t('docs.ccswitch.rowUsage') }}</td>
                  <td>{{ t('docs.ccswitch.rowUsageValue') }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <h3 class="docs-h3"><span class="docs-h3-num">6.3</span> {{ t('docs.ccswitch.useTitle') }}</h3>
          <p>{{ t('docs.ccswitch.useDesc') }}</p>
          <div class="docs-notice">
            <p class="docs-notice-title">{{ t('docs.ccswitch.useNoticeTitle') }}</p>
            <i18n-t keypath="docs.ccswitch.useNoticeText" tag="p" scope="global" class="docs-notice-text">
              <template #code><code class="docs-code-inline">{{ siteRoot }}</code></template>
              <template #code2><code class="docs-code-inline">{{ endpoint }}</code></template>
            </i18n-t>
          </div>

          <h3 class="docs-h3"><span class="docs-h3-num">6.4</span> {{ t('docs.ccswitch.manualTitle') }}</h3>
          <p>{{ t('docs.ccswitch.manualDesc') }}</p>
          <p class="docs-strong-line">{{ t('docs.ccswitch.manualClaude') }}</p>
          <CodeBlock lang="json" :code="ccswitchClaude" />
          <p class="docs-strong-line">{{ t('docs.ccswitch.manualCodex') }}</p>
          <CodeBlock lang="json" :code="ccswitchCodex" />
        </section>

        <!-- 07 用量查询 -->
        <section id="usage" class="docs-section">
          <div class="docs-section-head">
            <span class="docs-section-num">07</span>
            <h2 class="docs-h2">{{ t('docs.usage.title') }}</h2>
          </div>
          <i18n-t keypath="docs.usage.desc" tag="p" scope="global">
            <template #strong><strong>{{ t('docs.usage.descStrong') }}</strong></template>
          </i18n-t>
          <ul class="docs-list">
            <li><strong>{{ t('docs.usage.item1') }}</strong> — {{ t('docs.usage.item1Desc') }}</li>
            <li><strong>{{ t('docs.usage.item2') }}</strong> — {{ t('docs.usage.item2Desc') }}</li>
            <li><strong>{{ t('docs.usage.item3') }}</strong> — {{ t('docs.usage.item3Desc') }}</li>
            <li><strong>{{ t('docs.usage.item4') }}</strong> — {{ t('docs.usage.item4Desc') }}</li>
          </ul>
        </section>

        <!-- 08 常见问题 -->
        <section id="faq" class="docs-section">
          <div class="docs-section-head">
            <span class="docs-section-num">08</span>
            <h2 class="docs-h2">{{ t('docs.faq.title') }}</h2>
          </div>
          <details class="docs-faq">
            <summary>{{ t('docs.faq.q1') }}</summary>
            <i18n-t keypath="docs.faq.a1" tag="p" scope="global">
              <template #code><code class="docs-code-inline">sk-</code></template>
            </i18n-t>
          </details>
          <details class="docs-faq">
            <summary>{{ t('docs.faq.q2') }}</summary>
            <i18n-t keypath="docs.faq.a2" tag="p" scope="global">
              <template #strong><strong>{{ t('docs.faq.q2Concurrency') }}</strong></template>
            </i18n-t>
          </details>
          <details class="docs-faq">
            <summary>{{ t('docs.faq.q3') }}</summary>
            <p>{{ t('docs.faq.a3') }}</p>
          </details>
        </section>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import CodeBlock from '@/components/docs/CodeBlock.vue'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteName = computed(() => appStore.siteName)
const siteLogo = computed(() =>
  sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }),
)

// 站点自己的域名，不写死：换成任何一个部署都能直接用
const siteRoot = computed(() => {
  const configured = appStore.apiBaseUrl?.trim()
  if (configured) return configured.replace(/\/+$/, '')
  if (typeof window !== 'undefined') return window.location.origin.replace(/\/+$/, '')
  return ''
})

// api_base_url 可能已经带 /v1，拼接时要幂等
const endpoint = computed(() => {
  const root = siteRoot.value
  if (!root) return '/v1'
  return /\/v1$/i.test(root) ? root : `${root}/v1`
})

const isAuthenticated = computed(() => authStore.isAuthenticated)
const authTarget = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))
const authLabel = computed(() =>
  isAuthenticated.value ? t('home.dashboard') : t('home.login'),
)

const introFeatures = computed(() => [
  { title: t('docs.intro.featureAgents'), desc: t('docs.intro.featureAgentsDesc') },
  { title: t('docs.intro.featureCompat'), desc: t('docs.intro.featureCompatDesc') },
  { title: t('docs.intro.featureEndpoint'), desc: t('docs.intro.featureEndpointDesc') },
  { title: t('docs.intro.featureBilling'), desc: t('docs.intro.featureBillingDesc') },
])

const envSnippet = computed(
  () => `# OpenAI 兼容模式
export OPENAI_API_KEY="YOUR_API_KEY"
export OPENAI_BASE_URL="${endpoint.value}"

# Anthropic / Claude Code 模式
export ANTHROPIC_AUTH_TOKEN="YOUR_API_KEY"
export ANTHROPIC_BASE_URL="${endpoint.value}"
export CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1`,
)

const curlSnippet = computed(
  () => `# OpenAI 兼容模式
curl ${endpoint.value}/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -d '{
    "model": "gpt-5.5",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'`,
)

const codexConfig = computed(() => `# ~/.codex/config.toml
base_url = "${endpoint.value}"`)
const codexAuth = computed(() => `# ~/.codex/auth.json
{
  "OPENAI_API_KEY": "YOUR_API_KEY"
}`)
const codexWsConfig = computed(
  () => `# ~/.codex/config.toml
base_url = "${endpoint.value}"
supports_websockets = true
responses_websockets_v2 = true`,
)
const claudeEnv = computed(() => `export ANTHROPIC_AUTH_TOKEN="YOUR_API_KEY"
export ANTHROPIC_BASE_URL="${endpoint.value}"
export CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1`)
const claudeFile = computed(() => `# ~/.claude/settings.json
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "YOUR_API_KEY",
    "ANTHROPIC_BASE_URL": "${endpoint.value}",
    "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1"
  }
}`)
const opencodeEnv = computed(() => `export OPENAI_API_KEY="YOUR_API_KEY"`)
const opencodeFile = computed(() => `# ~/.config/opencode/opencode.json
{
  "provider": {
    "baseURL": "${endpoint.value}",
    "apiKey": "YOUR_API_KEY"
  }
}`)
const ccswitchClaude = computed(() => `{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "YOUR_API_KEY",
    "ANTHROPIC_BASE_URL": "${siteRoot.value}"
  }
}`)
const ccswitchCodex = computed(() => `# auth.json
{
  "OPENAI_API_KEY": "YOUR_API_KEY"
}

# config.toml
[model_providers.custom]
name = "${siteName.value}"
base_url = "${endpoint.value}"
wire_api = "responses"
requires_openai_auth = true`)

// ==================== 主题 ====================
// main.ts 已经按 localStorage / 系统偏好设好 dark 类，这里只读结果。
const isDark = ref(document.documentElement.classList.contains('dark'))

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// ==================== 目录高亮 ====================
// 左侧目录跟着滚动位置走：标题滚到视口上方 20% 处就算当前章节。
const activeSection = ref('intro')

function updateActiveSection() {
  const headings = Array.from(
    document.querySelectorAll<HTMLElement>('.docs-main .docs-section[id]'),
  )
  if (!headings.length) return
  const threshold = window.innerHeight * 0.2
  let current = headings[0].id
  for (const el of headings) {
    if (el.getBoundingClientRect().top <= threshold) current = el.id
  }
  activeSection.value = current
}

let scrollHandler: (() => void) | undefined

onMounted(() => {
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
  if (!authStore.isAuthenticated) {
    authStore.checkAuth()
  }
  updateActiveSection()
  let ticking = false
  scrollHandler = () => {
    if (ticking) return
    ticking = true
    window.requestAnimationFrame(() => {
      updateActiveSection()
      ticking = false
    })
  }
  window.addEventListener('scroll', scrollHandler, { passive: true })
})

onBeforeUnmount(() => {
  if (scrollHandler) window.removeEventListener('scroll', scrollHandler)
})
</script>

<style scoped>
.docs-page {
  min-height: 100vh;
  background-color: rgb(var(--n-50));
}

.dark .docs-page {
  background-color: rgb(var(--d-950));
}

/* 顶栏：文档是公共页，导航和后台壳层不一样，所以自己带一条 */
.docs-topbar {
  position: sticky;
  top: 0;
  z-index: 30;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  height: 3.5rem;
  padding-inline: 1.25rem;
  border-bottom: 1px solid rgb(var(--line-light) / 0.1);
  background-color: rgb(248 248 246 / 0.9);
  backdrop-filter: blur(12px);
}

.dark .docs-topbar {
  border-bottom-color: rgb(var(--line) / var(--line-a));
  background-color: rgb(20 20 19 / 0.88);
}

.docs-brand {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 0;
  font-size: 0.9375rem;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: rgb(var(--n-900));
}

.dark .docs-brand {
  color: rgb(var(--n-50));
}

.docs-brand-mark {
  width: 1.75rem;
  height: 1.75rem;
  border-radius: var(--r-md);
  object-fit: contain;
}

.docs-brand-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.docs-topbar-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.docs-icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border: 0;
  border-radius: var(--r-md);
  background: transparent;
  color: rgb(var(--n-600));
  cursor: pointer;
  transition: background-color 120ms ease, color 120ms ease;
}

.docs-icon-btn:hover {
  background-color: rgb(var(--line-light) / 0.07);
  color: rgb(var(--n-900));
}

.dark .docs-icon-btn {
  color: rgb(var(--n-400));
}

.dark .docs-icon-btn:hover {
  background-color: rgb(255 255 255 / 0.06);
  color: rgb(var(--n-50));
}

.docs-login {
  display: inline-flex;
  align-items: center;
  height: 2rem;
  padding-inline: 0.75rem;
  border-radius: var(--r-md);
  background-color: rgb(var(--primary));
  font-size: 0.8125rem;
  font-weight: 600;
  color: #fff;
  transition: background-color 120ms ease;
}

.docs-login:hover {
  background-color: rgb(var(--primary-hover));
}

/* 版式：左侧目录固定，正文限宽到可读区间 */
.docs-layout {
  display: grid;
  grid-template-columns: 14rem minmax(0, 1fr);
  gap: 3.5rem;
  width: 100%;
  max-width: 78rem;
  margin-inline: auto;
  padding: 2.5rem 1.25rem 6rem;
}

.docs-sidebar {
  position: sticky;
  top: 5rem;
  align-self: start;
  max-height: calc(100vh - 6.5rem);
  overflow-y: auto;
}

.docs-nav {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
  font-size: 0.8125rem;
}

.docs-nav-section {
  margin: 1rem 0 0.375rem;
  padding-left: 0.625rem;
  font-size: 0.6875rem;
  font-weight: 600;
  letter-spacing: 0.02em;
  color: rgb(var(--n-400));
}

.docs-nav-section:first-child {
  margin-top: 0;
}

.docs-nav a {
  padding: 0.3125rem 0.625rem;
  border-radius: var(--r-md);
  color: rgb(var(--n-600));
  text-decoration: none;
  transition: background-color 120ms ease, color 120ms ease;
}

.docs-nav a:hover {
  background-color: rgb(var(--line-light) / 0.05);
  color: rgb(var(--n-900));
}

.dark .docs-nav a {
  color: rgb(var(--n-400));
}

.dark .docs-nav a:hover {
  background-color: rgb(255 255 255 / 0.05);
  color: rgb(var(--n-50));
}

/* 正文 */
.docs-main {
  min-width: 0;
  max-width: 46rem;
}

.docs-hero {
  padding-bottom: 2rem;
  border-bottom: 1px solid rgb(var(--line-light) / 0.09);
  margin-bottom: 2.5rem;
}

.dark .docs-hero {
  border-bottom-color: rgb(var(--line) / 0.09);
}

.docs-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  margin: 0 0 1rem;
  padding: 0.1875rem 0.5rem;
  border: 1px solid rgb(var(--p-500) / 0.22);
  border-radius: 999px;
  background-color: rgb(var(--p-500) / 0.1);
  font-size: 0.6875rem;
  font-weight: 600;
  color: rgb(var(--p-700));
}

.dark .docs-badge {
  color: rgb(var(--p-300));
}

.docs-badge-dot {
  width: 0.3125rem;
  height: 0.3125rem;
  border-radius: 999px;
  background-color: rgb(var(--p-500));
}

.docs-h1 {
  margin: 0;
  font-size: 2.25rem;
  font-weight: 800;
  line-height: 1.15;
  letter-spacing: -0.03em;
  color: rgb(var(--n-900));
}

.docs-h1 span {
  display: block;
}

.docs-h1-accent {
  color: rgb(var(--p-600));
}

.dark .docs-h1 {
  color: rgb(var(--n-50));
}

.dark .docs-h1-accent {
  color: rgb(var(--p-400));
}

.docs-hero-desc {
  margin: 1rem 0 0;
  font-size: 0.9375rem;
  line-height: 1.7;
  color: rgb(var(--n-600));
}

.dark .docs-hero-desc {
  color: rgb(var(--n-400));
}

.docs-hero-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 1.25rem;
}

.docs-hero-meta span {
  padding: 0.1875rem 0.5625rem;
  border: 1px solid rgb(var(--line-light) / 0.12);
  border-radius: 999px;
  font-size: 0.6875rem;
  color: rgb(var(--n-600));
}

.dark .docs-hero-meta span {
  border-color: rgb(var(--line) / var(--line-a));
  color: rgb(var(--n-400));
}

.docs-section {
  margin-bottom: 3rem;
  scroll-margin-top: 4.5rem;
}

.docs-section-head {
  display: flex;
  align-items: baseline;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.docs-section-num {
  font-family: var(--font-numeric);
  font-size: 0.75rem;
  font-weight: 600;
  color: rgb(var(--p-600));
}

.dark .docs-section-num {
  color: rgb(var(--p-400));
}

.docs-h2 {
  margin: 0;
  font-size: 1.375rem;
  font-weight: 750;
  letter-spacing: -0.02em;
  color: rgb(var(--n-900));
}

.dark .docs-h2 {
  color: rgb(var(--n-50));
}

.docs-h3 {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  margin: 2rem 0 0.75rem;
  font-size: 1rem;
  font-weight: 650;
  letter-spacing: -0.01em;
  color: rgb(var(--n-800));
}

.dark .docs-h3 {
  color: rgb(var(--n-100));
}

.docs-h3-num {
  font-family: var(--font-numeric);
  font-size: 0.8125rem;
  font-weight: 600;
  color: rgb(var(--n-400));
}

.docs-h4 {
  margin: 1.5rem 0 0.625rem;
  font-size: 0.875rem;
  font-weight: 600;
  color: rgb(var(--n-700));
}

.dark .docs-h4 {
  color: rgb(var(--n-300));
}

.docs-main p {
  margin: 0 0 0.875rem;
  font-size: 0.875rem;
  line-height: 1.75;
  color: rgb(var(--n-700));
}

.dark .docs-main p {
  color: rgb(var(--n-300));
}

.docs-strong-line {
  font-weight: 600;
  color: rgb(var(--n-800));
}

.dark .docs-strong-line {
  color: rgb(var(--n-100));
}

.docs-link {
  color: rgb(var(--p-700));
  text-decoration: underline;
  text-underline-offset: 2px;
}

.dark .docs-link {
  color: rgb(var(--p-400));
}

.docs-code-inline {
  padding: 0.0625rem 0.3125rem;
  border-radius: 5px;
  background-color: rgb(var(--line-light) / 0.06);
  font-family: var(--font-numeric);
  font-size: 0.8125em;
  color: rgb(var(--n-800));
}

.dark .docs-code-inline {
  background-color: rgb(255 255 255 / 0.06);
  color: rgb(var(--n-100));
}

.docs-feature-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(13rem, 1fr));
  gap: 0.75rem;
  margin: 1.25rem 0;
}

.docs-feature {
  padding: 0.875rem 1rem;
  border: 1px solid rgb(var(--line-light) / 0.1);
  border-radius: var(--r-lg);
}

.dark .docs-feature {
  border-color: rgb(var(--line) / var(--line-a));
}

.docs-feature-title {
  margin: 0 0 0.25rem;
  font-size: 0.8125rem;
  font-weight: 650;
  color: rgb(var(--n-900));
}

.dark .docs-feature-title {
  color: rgb(var(--n-50));
}

.docs-feature-desc {
  margin: 0;
  font-size: 0.75rem;
  line-height: 1.6;
  color: rgb(var(--n-500));
}

.docs-steps {
  display: flex;
  flex-direction: column;
  gap: 0.625rem;
  margin: 1rem 0;
  padding: 0;
  list-style: none;
}

.docs-steps li {
  display: flex;
  align-items: flex-start;
  gap: 0.625rem;
  font-size: 0.8125rem;
  line-height: 1.65;
  color: rgb(var(--n-700));
}

.dark .docs-steps li {
  color: rgb(var(--n-300));
}

.docs-step-num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 1.125rem;
  height: 1.125rem;
  margin-top: 0.125rem;
  border-radius: 999px;
  background-color: rgb(var(--p-500) / 0.14);
  font-size: 0.6875rem;
  font-weight: 700;
  color: rgb(var(--p-700));
}

.dark .docs-step-num {
  color: rgb(var(--p-300));
}

.docs-endpoint-card {
  position: relative;
  margin: 1rem 0;
  padding: 1rem 1.125rem;
  border: 1px solid rgb(var(--p-500) / 0.25);
  border-radius: var(--r-lg);
  background-color: rgb(var(--p-500) / 0.06);
}

.docs-endpoint-badge {
  position: absolute;
  top: 0.75rem;
  right: 0.75rem;
  padding: 0.0625rem 0.375rem;
  border-radius: 999px;
  background-color: rgb(var(--p-500) / 0.16);
  font-size: 0.625rem;
  font-weight: 700;
  color: rgb(var(--p-700));
}

.dark .docs-endpoint-badge {
  color: rgb(var(--p-300));
}

.docs-endpoint-label {
  margin: 0 0 0.25rem;
  font-size: 0.6875rem;
  font-weight: 600;
  color: rgb(var(--n-500));
}

.docs-endpoint-url {
  display: block;
  font-family: var(--font-numeric);
  font-size: 0.9375rem;
  font-weight: 600;
  color: rgb(var(--n-900));
  overflow-wrap: anywhere;
}

.dark .docs-endpoint-url {
  color: rgb(var(--n-50));
}

.docs-endpoint-desc {
  margin: 0.5rem 0 0;
  font-size: 0.75rem;
  color: rgb(var(--n-500));
}

.docs-notice {
  margin: 1rem 0;
  padding: 0.75rem 0.875rem;
  border-left: 2px solid rgb(var(--p-500) / 0.5);
  border-radius: var(--r-md);
  background-color: rgb(var(--line-light) / 0.04);
}

.dark .docs-notice {
  background-color: rgb(255 255 255 / 0.04);
}

.docs-notice-title {
  margin: 0 0 0.25rem;
  font-size: 0.75rem;
  font-weight: 700;
  color: rgb(var(--n-800));
}

.dark .docs-notice-title {
  color: rgb(var(--n-100));
}

.docs-notice-text {
  margin: 0;
  font-size: 0.75rem;
  line-height: 1.7;
  color: rgb(var(--n-600));
}

.dark .docs-notice-text {
  color: rgb(var(--n-300));
}

.docs-list {
  margin: 1rem 0;
  padding-left: 1.125rem;
  font-size: 0.8125rem;
  line-height: 1.9;
  color: rgb(var(--n-700));
}

.dark .docs-list {
  color: rgb(var(--n-300));
}

.docs-table-wrap {
  margin: 1rem 0;
  overflow-x: auto;
  border: 1px solid rgb(var(--line-light) / 0.1);
  border-radius: var(--r-md);
}

.dark .docs-table-wrap {
  border-color: rgb(var(--line) / var(--line-a));
}

.docs-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8125rem;
}

.docs-table th,
.docs-table td {
  padding: 0.5rem 0.75rem;
  text-align: left;
  border-bottom: 1px solid rgb(var(--line-light) / 0.07);
}

.dark .docs-table th,
.dark .docs-table td {
  border-bottom-color: rgb(var(--line) / 0.07);
}

.docs-table th {
  font-weight: 600;
  color: rgb(var(--n-500));
}

.docs-table td {
  color: rgb(var(--n-700));
}

.dark .docs-table td {
  color: rgb(var(--n-300));
}

.docs-table tbody tr:last-child td {
  border-bottom: 0;
}

.docs-faq {
  margin-bottom: 0.5rem;
  border-bottom: 1px solid rgb(var(--line-light) / 0.08);
  padding-bottom: 0.5rem;
}

.dark .docs-faq {
  border-bottom-color: rgb(var(--line) / 0.08);
}

.docs-faq summary {
  padding: 0.375rem 0;
  font-size: 0.875rem;
  font-weight: 600;
  color: rgb(var(--n-800));
  cursor: pointer;
}

.dark .docs-faq summary {
  color: rgb(var(--n-100));
}

.docs-faq p {
  margin: 0.25rem 0 0.5rem;
  font-size: 0.8125rem;
}

@media (max-width: 1023px) {
  .docs-layout {
    grid-template-columns: minmax(0, 1fr);
    gap: 0;
    padding: 1.75rem 1rem 4rem;
  }

  .docs-sidebar {
    position: static;
    max-height: none;
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid rgb(var(--line-light) / 0.08);
  }

  .docs-nav {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 0.25rem;
  }

  .docs-nav-section {
    display: none;
  }

  .docs-nav a {
    border: 1px solid rgb(var(--line-light) / 0.12);
  }

  .dark .docs-nav a {
    border-color: rgb(var(--line) / var(--line-a));
  }

  .docs-main {
    max-width: none;
  }

  .docs-h1 {
    font-size: 1.75rem;
  }
}
</style>
