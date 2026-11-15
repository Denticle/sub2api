<template>
  <div class="lp">
    <!-- ── 公告条：后台「公告」里面向全部用户的最新一条 ── -->
    <div v-if="announcement && !announceDismissed" class="nv-announce" data-testid="home-announce">
      <div class="lp-wrap nv-announce-inner">
        <component
          :is="isAuthenticated ? 'router-link' : 'span'"
          :to="isAuthenticated ? '/dashboard' : undefined"
          class="nv-announce-link"
        >
          <span class="nv-announce-tag">{{ t('home.lp.announceTag') }}</span>
          <span class="nv-announce-text">{{ announcement.title }}</span>
          <span v-if="isAuthenticated" class="nv-announce-go">{{ t('home.lp.announceMore') }}</span>
        </component>
        <button type="button" class="nv-announce-x" :aria-label="t('home.lp.announceClose')" @click="dismissAnnouncement">
          <svg viewBox="0 0 24 24"><path d="M6 6l12 12M18 6L6 18" /></svg>
        </button>
      </div>
    </div>

    <!-- ── 导航 ── -->
    <header class="nv">
      <div class="lp-wrap nv-inner">
        <router-link class="nv-brand" to="/">
          <img :src="siteLogo || '/logo.svg'" alt="" width="28" height="28" />
          <span class="lp-brand-name">{{ siteName }}</span>
        </router-link>

        <nav class="nv-links" :aria-label="t('home.lp.mainNav')">
          <div ref="dropRef" class="nv-drop">
            <button
              type="button"
              class="nv-link nv-trigger"
              :aria-expanded="menuOpen ? 'true' : 'false'"
              aria-controls="nv-menu"
              @click="menuOpen = !menuOpen"
            >
              {{ t('home.lp.navProduct') }}
              <svg class="nv-caret" viewBox="0 0 24 24"><path d="M7 10l5 5 5-5" /></svg>
            </button>
            <div id="nv-menu" class="nv-menu">
              <div class="nv-col">
                <p class="nv-col-h">{{ t('home.lp.menuStudio') }}</p>
                <component
                  :is="item.external ? 'a' : 'router-link'"
                  v-for="item in studioMenu"
                  :key="item.key"
                  class="nv-item"
                  v-bind="linkAttrs(item.to, item.external)"
                  @click="menuOpen = false"
                >
                  <span class="nv-ico"><svg viewBox="0 0 24 24" v-html="ICONS[item.icon]" /></span>
                  <span><b>{{ item.title }}</b><small>{{ item.desc }}</small></span>
                </component>
              </div>
              <div class="nv-col">
                <p class="nv-col-h">{{ t('home.lp.menuDev') }}</p>
                <component
                  :is="item.external ? 'a' : 'router-link'"
                  v-for="item in devMenu"
                  :key="item.key"
                  class="nv-item"
                  v-bind="linkAttrs(item.to, item.external)"
                  @click="menuOpen = false"
                >
                  <span class="nv-ico"><svg viewBox="0 0 24 24" v-html="ICONS[item.icon]" /></span>
                  <span><b>{{ item.title }}</b><small>{{ item.desc }}</small></span>
                </component>
              </div>
              <a class="nv-feature" href="#lp-gallery" @click="menuOpen = false">
                <img :src="featureImage" alt="" loading="lazy" />
                <span class="nv-feature-txt">
                  <b>{{ t('home.lp.menuFeature') }}</b>
                  <small>{{ t('home.lp.menuFeatureDesc') }}</small>
                </span>
              </a>
            </div>
          </div>
          <router-link v-if="showModelPlazaEntry" class="nv-link" to="/model-plaza">{{ t('home.lp.navPlaza') }}</router-link>
          <a class="nv-link" href="#lp-gallery">{{ t('home.lp.navGallery') }}</a>
          <component :is="docsIsRoute ? 'router-link' : 'a'" v-if="docsTarget" class="nv-link" v-bind="linkAttrs(docsTarget, !docsIsRoute)">
            {{ t('home.lp.navDocs') }}
          </component>
        </nav>

        <div class="nv-actions">
          <span class="nv-locale"><LocaleSwitcher /></span>
          <button
            type="button"
            class="nv-icon"
            :aria-label="isDark ? t('home.lp.themeToLight') : t('home.lp.themeToDark')"
            :title="isDark ? t('home.lp.themeToLight') : t('home.lp.themeToDark')"
            @click="$emit('toggle-theme')"
          >
            <svg v-if="isDark" viewBox="0 0 24 24"><circle cx="12" cy="12" r="4" /><path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" /></svg>
            <svg v-else viewBox="0 0 24 24"><path d="M20 14.5A8 8 0 019.5 4 8 8 0 1020 14.5z" /></svg>
          </button>
          <template v-if="isAuthenticated">
            <router-link class="nv-cta" :to="dashboardPath">{{ t('home.lp.navConsole') }}</router-link>
          </template>
          <template v-else>
            <router-link class="nv-login" to="/login">{{ t('home.lp.navLogin') }}</router-link>
            <router-link class="nv-cta" :to="signupTarget">{{ t('home.lp.navStart') }}</router-link>
          </template>
        </div>
      </div>
    </header>

    <main>
      <!-- ── ① Hero ── -->
      <section class="lp-hero">
        <div class="lp-wrap lp-hero-grid">
          <div>
            <p class="lp-status">
              <span class="lp-status-dot"></span>
              {{ t('home.lp.status', { models: modelCountLabel, providers: providerCount }) }}
            </p>
            <h1 class="lp-h1">
              <span>{{ t('home.lp.headlineOne') }}</span>
              <span>{{ t('home.lp.headlineTwo') }}<span class="accent">{{ t('home.lp.headlineAccent') }}</span></span>
            </h1>
            <p class="lp-lede">
              {{ t('home.lp.ledeA') }}<span class="hl">{{ t('home.lp.ledeHl') }}</span>{{ t('home.lp.ledeB') }}
            </p>
            <div class="lp-cta-row">
              <router-link class="lp-btn lp-btn-primary" :to="createTarget">
                {{ t('home.lp.ctaPrimary') }}
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14M13 6l6 6-6 6" /></svg>
              </router-link>
              <a class="lp-btn lp-btn-quiet" href="#lp-gallery">{{ t('home.lp.ctaSecondary') }}</a>
            </div>
            <p class="lp-note">{{ t('home.lp.note') }}</p>
          </div>

          <!-- 画布示意：短剧 / 广告 / 代码三种场景 -->
          <div ref="cvFitRef" class="cv-fit" aria-hidden="true" :style="{ '--cv-scale': cvScale }">
            <div class="cv-canvas">
              <svg width="100%" height="100%" viewBox="0 0 560 378" preserveAspectRatio="none">
                <path class="cv-wire" d="M140 60 C180 60, 186 100, 224 100" />
                <path class="cv-wire" d="M344 110 C388 110, 392 66, 430 66" />
                <path class="cv-wire" d="M140 240 C180 240, 186 300, 224 300" />
                <path class="cv-wire" d="M344 100 C392 100, 400 180, 434 180" />
                <circle class="cv-packet" cx="182" cy="80" r="3" />
                <circle class="cv-packet" cx="400" cy="88" r="3" />
                <circle class="cv-packet" cx="182" cy="270" r="3" />
              </svg>

              <div class="cv-node prompt" style="left: 16px; top: 34px">
                <div class="body"><em>{{ t('home.lp.cvDramaTitle') }}</em><br />{{ t('home.lp.cvDramaBody') }}</div>
                <div class="lbl">{{ t('home.lp.cvDramaLabel') }}</div>
              </div>
              <div class="cv-node img" style="left: 224px; top: 74px">
                <div class="thumb" :style="bg(MEDIA.boards)"></div>
                <div class="lbl">{{ t('home.lp.cvBoards') }}</div>
              </div>
              <div class="cv-node vid" style="left: 436px; top: 156px">
                <div class="thumb" :style="bg(MEDIA.beast)"></div>
                <div class="lbl">{{ t('home.lp.cvDramaOut') }}</div>
              </div>
              <div class="cv-node vid" style="left: 432px; top: 40px">
                <div class="thumb" :style="bg(MEDIA.ad)"></div>
                <div class="lbl">{{ t('home.lp.cvAd') }}</div>
              </div>
              <div class="cv-node prompt" style="left: 16px; top: 214px">
                <div class="body"><em>{{ t('home.lp.cvCodeTitle') }}</em><br />{{ t('home.lp.cvCodeBody') }}</div>
                <div class="lbl">coding · claude-fable-5</div>
              </div>
              <div class="cv-node out" style="left: 224px; top: 274px">
                <div class="thumb cv-code">
                  <span class="k">function</span> <span class="f">CardGrid</span>() {<br />
                  &nbsp;&nbsp;<span class="k">return</span> (<br />
                  &nbsp;&nbsp;&nbsp;&nbsp;&lt;<span class="tg">div</span> <span class="a">className</span>=<span class="s">"grid"</span>&gt;
                </div>
                <div class="lbl">{{ t('home.lp.cvCodeOut') }}</div>
              </div>
              <span class="cv-badge">{{ t('home.lp.cvBadge') }}</span>
            </div>
          </div>
        </div>

        <div class="lp-wrap">
          <dl class="lp-stats" data-testid="home-providers">
            <div><dt>{{ t('home.lp.statModels') }}</dt><dd>{{ modelCountLabel }}<small v-if="unit">{{ unit }}</small></dd></div>
            <div><dt>{{ t('home.lp.statProviders') }}</dt><dd>{{ providerCount }}<small v-if="unit">{{ unit }}</small></dd></div>
            <div><dt>{{ t('home.lp.statSkills') }}</dt><dd>{{ SKILL_COUNT }}<small v-if="unit">{{ unit }}</small></dd></div>
            <div><dt>{{ t('home.lp.statBilling') }}</dt><dd>{{ t('home.lp.billingValue') }}<small>token</small></dd></div>
          </dl>
        </div>
      </section>

      <!-- ── ② 工作台：Agent + 画布 ── -->
      <section v-if="showWorkbench" id="lp-workbench" class="lp-section ws-sec">
        <div class="lp-wrap">
          <header class="ws-head">
            <h2>{{ t('home.lp.wsTitle') }}</h2>
            <p>{{ t('home.lp.wsSub') }}</p>
          </header>

          <div class="ws-app" aria-hidden="true">
            <div class="ws-bar">
              <span class="ws-dots"><i></i><i></i><i></i></span>
              <span class="ws-crumb">{{ t('home.lp.wsCrumbRoot') }} <b>/</b> {{ t('home.lp.wsCrumbFile') }}</span>
              <span class="ws-bar-r"><span class="ws-live"></span>{{ t('home.lp.wsRunning') }}</span>
            </div>
            <div class="ws-body">
              <aside class="ws-agent">
                <div class="ws-msg ws-msg-user">{{ t('home.lp.wsUserMsg') }}</div>
                <div class="ws-msg ws-msg-ai">
                  <p>{{ t('home.lp.wsReply') }}</p>
                  <ol class="ws-steps">
                    <li class="done"><span>{{ t('home.lp.wsStep1') }}</span><em>{{ t('home.lp.wsStep1Meta') }}</em></li>
                    <li class="done"><span>{{ t('home.lp.wsStep2') }}</span><em>{{ t('home.lp.wsStep2Meta') }}</em></li>
                    <li class="done"><span>{{ t('home.lp.wsStep3') }}</span><em>{{ t('home.lp.wsStep3Meta') }}</em></li>
                    <li class="run"><span>{{ t('home.lp.wsStep4') }}</span><em>72%</em><i style="--p: 72%"></i></li>
                    <li><span>{{ t('home.lp.wsStep5') }}</span><em>1080p</em></li>
                  </ol>
                </div>
                <div class="ws-input">{{ t('home.lp.wsInput') }}<span class="ws-send"></span></div>
              </aside>

              <div class="ws-canvas">
                <div class="ws-shots">
                  <figure v-for="shot in shots" :key="shot.n" class="ws-shot">
                    <img :src="shot.src" alt="" loading="lazy" />
                    <figcaption>
                      <span>{{ t('home.lp.wsShot', { n: shot.n, name: shot.name }) }}</span>
                      <span class="tc">{{ shot.dur }}</span>
                    </figcaption>
                  </figure>
                </div>
                <div class="ws-player">
                  <div class="ws-screen">
                    <img :src="MEDIA.beast" alt="" loading="lazy" />
                    <span class="ws-play"></span>
                    <span class="ws-res">1080p · 16:9</span>
                  </div>
                  <div class="ws-timeline">
                    <span class="ws-time">00:22 / 00:30</span>
                    <div class="ws-strip">
                      <span v-for="shot in shots" :key="shot.n" :style="bg(shot.src)"></span>
                      <b class="ws-head-line" style="left: 73%"></b>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <p class="ws-note">{{ t('home.lp.wsNote', { n: SKILL_COUNT }) }}</p>
        </div>
      </section>

      <!-- ── ③ 两种用法 ── -->
      <section id="lp-work" class="lp-section lp-band">
        <div class="lp-wrap">
          <header class="lp-sec-head">
            <p class="lp-kicker">{{ t('home.lp.waysKicker') }}</p>
            <h2>{{ t('home.lp.waysTitle') }}</h2>
            <p>{{ t('home.lp.waysSub') }}</p>
          </header>
          <div class="lp-rails">
            <div v-if="showWorkbench" class="lp-rail">
              <span class="lp-rail-tag"><span class="dot"></span>{{ t('home.lp.waysMakeTag') }}</span>
              <h3>{{ t('home.lp.waysMakeTitle') }}</h3>
              <p>{{ t('home.lp.waysMakeDesc') }}</p>
              <ul class="lp-rail-list">
                <li v-for="(line, i) in makeLines" :key="i"><span class="n">{{ pad(i + 1) }}</span><span>{{ line }}</span></li>
              </ul>
              <router-link class="lp-btn lp-btn-quiet lp-btn-sm" :to="WORKBENCH">{{ t('home.lp.waysMakeCta') }}</router-link>
            </div>
            <div class="lp-rail">
              <span class="lp-rail-tag"><span class="dot"></span>{{ t('home.lp.waysApiTag') }}</span>
              <h3>{{ t('home.lp.waysApiTitle') }}</h3>
              <p>{{ t('home.lp.waysApiDesc') }}</p>
              <ul class="lp-rail-list">
                <li v-for="(line, i) in apiLines" :key="i"><span class="n">{{ pad(i + 1) }}</span><span>{{ line }}</span></li>
              </ul>
              <component
                :is="docsIsRoute ? 'router-link' : 'a'"
                v-if="docsTarget"
                class="lp-btn lp-btn-quiet lp-btn-sm"
                v-bind="linkAttrs(docsTarget, !docsIsRoute)"
              >
                {{ t('home.lp.waysApiCta') }}
              </component>
            </div>
          </div>
        </div>
      </section>

      <!-- ── ④ 作品墙：视频工作台灵感库 ── -->
      <section v-if="gallery.length" id="lp-gallery" class="lp-section" data-testid="home-gallery">
        <div class="lp-wrap">
          <header class="lp-sec-head lp-sec-head-row">
            <div>
              <p class="lp-kicker">{{ t('home.lp.galKicker') }}</p>
              <h2>{{ t('home.lp.galTitle') }}</h2>
              <p>{{ t('home.lp.galSub') }}</p>
            </div>
            <router-link v-if="galleryTotal" class="lp-link" :to="createTarget">{{ t('home.lp.galAll', { n: galleryTotal }) }} →</router-link>
          </header>
          <div class="lp-filter" role="group">
            <button
              v-for="f in filters"
              :key="f.key"
              type="button"
              :class="{ on: filter === f.key }"
              :aria-pressed="filter === f.key"
              @click="filter = f.key"
            >
              {{ f.label }}
            </button>
          </div>
          <div class="lp-gallery lp-gallery-wide">
            <router-link
              v-for="item in visibleGallery"
              :key="item.id"
              class="lp-shot"
              :class="{ 'lp-shot-video': item.mode === 'video' }"
              :to="createTarget"
              :title="item.title"
            >
              <div class="lp-shot-art">
                <img :src="item.image" :alt="item.title" loading="lazy" />
                <span v-if="item.mode === 'video'" class="lp-shot-badge">▶</span>
              </div>
              <span class="lp-shot-cap">
                <span v-if="item.credit" class="m">{{ item.credit }}</span>
                <span class="t">{{ item.title }}</span>
              </span>
            </router-link>
          </div>
        </div>
      </section>

      <!-- ── ⑤ 收尾 ── -->
      <section class="cl">
        <div class="lp-wrap">
          <div class="cl-panel">
            <div class="cl-copy">
              <h2>{{ t('home.lp.clTitle') }}</h2>
              <p>{{ t('home.lp.clDesc') }}</p>
              <div class="cl-cta">
                <router-link class="cl-btn cl-btn-main" :to="createTarget">{{ t('home.lp.clPrimary') }}</router-link>
                <a v-if="contactHref" class="cl-btn cl-btn-ghost" :href="contactHref">{{ t('home.lp.clSecondary') }}</a>
              </div>
              <ul class="cl-points">
                <li>{{ t('home.lp.clPoint1') }}</li>
                <li>{{ t('home.lp.clPoint2') }}</li>
              </ul>
            </div>
            <div class="cl-visual" aria-hidden="true">
              <img :src="MEDIA.sea" alt="" loading="lazy" />
              <img :src="MEDIA.temple" alt="" loading="lazy" />
              <img :src="MEDIA.bridge" alt="" loading="lazy" />
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- ── 页脚 ── -->
    <footer class="ft">
      <div class="lp-wrap">
        <div class="ft-top">
          <div class="ft-brand">
            <router-link class="nv-brand" to="/">
              <img :src="siteLogo || '/logo.svg'" alt="" width="28" height="28" />
              <span>{{ siteName }}</span>
            </router-link>
            <p>{{ t('home.lp.ftTagline') }}</p>
          </div>
          <nav class="ft-cols" :aria-label="t('home.lp.footerNav')">
            <div v-for="col in footerCols" :key="col.key">
              <h4>{{ col.title }}</h4>
              <component
                :is="link.external ? 'a' : 'router-link'"
                v-for="link in col.links"
                :key="link.key"
                v-bind="linkAttrs(link.to, link.external)"
              >
                {{ link.label }}
              </component>
            </div>
          </nav>
        </div>
        <div class="ft-bottom">
          <span>© {{ currentYear }} {{ siteName }}</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import { sanitizeUrl } from '@/utils/url'
import { isLocalDocsLink, resolveDocsLink } from '@/utils/docsLink'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'

defineProps<{ isDark: boolean }>()
defineEmits<{ 'toggle-theme': [] }>()

const { t, locale } = useI18n()
// 中文数字后带量词「个」，英文不需要。
const unit = computed(() => (String(locale.value).startsWith('zh') ? '个' : ''))
const authStore = useAuthStore()
const appStore = useAppStore()

// 视频工作台内置技能数量；跟着工作台版本走，升级工作台时一起核对。
const SKILL_COUNT = 73
const WORKBENCH = '/workbench/canvas'

// 示意图用到的画面，全部取自工作台灵感库里的真实作品封面。
const MEDIA = {
  boards: 'https://cms-assets.youmind.com/media/1792482820813_mtdoc7_HU854tgbIAE9SlB.jpg',
  ad: 'https://media.goodcase.ai/media/poster/case-96017ad7706f.jpg',
  temple: 'https://media.goodcase.ai/cases/177801c0a33f.jpg',
  bridge: 'https://media.goodcase.ai/cases/b4c4f9bdf65c.jpg',
  sea: 'https://media.goodcase.ai/cases/a002c9dccc74.jpg',
  beast: 'https://media.goodcase.ai/cases/c5e8fdb6d69a.jpg',
}

const ICONS: Record<string, string> = {
  canvas: '<path d="M4 6h6v5H4zM14 13h6v5h-6zM10 8.5h2a2 2 0 012 2V15"/>',
  agent: '<path d="M12 3l1.6 4.4L18 9l-4.4 1.6L12 15l-1.6-4.4L6 9l4.4-1.6z"/><path d="M18 15l.8 2.2L21 18l-2.2.8L18 21l-.8-2.2L15 18l2.2-.8z"/>',
  film: '<rect x="3" y="5" width="18" height="14" rx="2"/><path d="M7 5v14M17 5v14M3 9.5h4M17 9.5h4M3 14.5h4M17 14.5h4"/>',
  image: '<rect x="3" y="4" width="18" height="16" rx="2"/><circle cx="9" cy="10" r="1.6"/><path d="M21 16l-5-5-9 9"/>',
  skill: '<path d="M4 5h7v7H4zM13 5h7v7h-7zM4 14h7v5H4zM13 14h7v5h-7z"/>',
  code: '<path d="M9 8l-4 4 4 4M15 8l4 4-4 4"/>',
  grid: '<circle cx="7" cy="7" r="2.5"/><circle cx="17" cy="7" r="2.5"/><circle cx="7" cy="17" r="2.5"/><circle cx="17" cy="17" r="2.5"/>',
  book: '<path d="M5 4h10a3 3 0 013 3v13H8a3 3 0 01-3-3z"/><path d="M5 17a3 3 0 013-3h10"/>',
  meter: '<path d="M4 18a8 8 0 1116 0"/><path d="M12 18l4-6"/>',
}

// ── 站点设置 ──
const settings = computed(() => appStore.cachedPublicSettings)
const siteName = computed(() => settings.value?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() =>
  sanitizeUrl(settings.value?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }),
)
const docsTarget = computed(() => resolveDocsLink(sanitizeUrl(settings.value?.doc_url || appStore.docUrl || '')))
const docsIsRoute = computed(() => isLocalDocsLink(docsTarget.value))

const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))
const registrationEnabled = computed(() => settings.value?.registration_enabled !== false)
const signupTarget = computed(() => (registrationEnabled.value ? '/register' : '/login'))

// 后端模式下普通用户进不了用户端，simple 模式也不显示工作台：这两种情况首页不宣传工作台。
// 没登录也直接链到工作台：路由守卫会先带去登录，登录完再回来。
const showWorkbench = computed(() => !appStore.backendModeEnabled && !authStore.isSimpleMode)
const createTarget = computed(() => (showWorkbench.value ? WORKBENCH : isAuthenticated.value ? dashboardPath.value : signupTarget.value))

const showModelPlazaEntry = computed(
  () =>
    isFeatureFlagEnabled(FeatureFlags.modelPlaza) &&
    (isAuthenticated.value || settings.value?.model_plaza_require_auth !== true),
)

// 「联系销售」只在后台填了联系方式且能解析出可点的地址时出现。
const contactHref = computed(() => {
  const raw = (settings.value?.contact_info || '').trim()
  const mail = raw.match(/[\w.+-]+@[\w-]+\.[\w.-]+/)
  if (mail) return `mailto:${mail[0]}`
  const url = raw.match(/https?:\/\/\S+/)
  return url ? sanitizeUrl(url[0]) : ''
})

function linkAttrs(to: string, external?: boolean) {
  if (external || /^https?:\/\//.test(to)) return { href: to, target: '_blank', rel: 'noopener noreferrer' }
  if (to.startsWith('#')) return { href: to }
  return { to }
}

function bg(src: string) {
  return { backgroundImage: `url("${src}")`, backgroundSize: 'cover', backgroundPosition: 'center' }
}

const pad = (n: number) => String(n).padStart(2, '0')

// ── 模型广场（公开接口）──
interface PlazaGroup { platform: string; models?: { name: string }[] }
const plazaGroups = ref<PlazaGroup[]>([])

async function loadPlaza() {
  try {
    const res = await fetch('/api/v1/model-plaza', { headers: { Accept: 'application/json' } })
    if (!res.ok) return
    const body = await res.json()
    plazaGroups.value = body?.data?.groups ?? []
  } catch {
    // 拿不到就显示占位，不让首页报错
  }
}

const modelCount = computed(() => {
  const seen = new Set<string>()
  for (const g of plazaGroups.value) for (const m of g.models ?? []) seen.add(m.name)
  return seen.size
})
const modelCountLabel = computed(() => (modelCount.value ? String(modelCount.value) : '—'))
const providerCount = computed(
  () => new Set(plazaGroups.value.filter((g) => (g.models ?? []).length > 0).map((g) => g.platform)).size,
)

// ── 公告 ──
interface PublicAnnouncement { id: number; title: string }
const announcement = ref<PublicAnnouncement | null>(null)
const DISMISS_KEY = 'home_announce_dismissed'
const announceDismissed = ref(false)

async function loadAnnouncement() {
  try {
    const res = await fetch('/api/v1/public/announcements', { headers: { Accept: 'application/json' } })
    if (!res.ok) return
    const body = await res.json()
    const first = (body?.data ?? [])[0]
    if (!first?.title) return
    announcement.value = first
    announceDismissed.value = localStorage.getItem(DISMISS_KEY) === String(first.id)
  } catch {
    // 没有公告就不显示公告条
  }
}

function dismissAnnouncement() {
  announceDismissed.value = true
  if (announcement.value) localStorage.setItem(DISMISS_KEY, String(announcement.value.id))
}

// ── 作品墙 ──
interface Inspiration { id: string; title: string; image: string; mode: 'image' | 'video'; credit?: string }
const gallery = ref<Inspiration[]>([])
const galleryTotal = ref(0)
const filter = ref<'all' | 'image' | 'video'>('all')
const filters = computed(() => [
  { key: 'all' as const, label: t('home.lp.galFilterAll') },
  { key: 'image' as const, label: t('home.lp.galFilterImage') },
  { key: 'video' as const, label: t('home.lp.galFilterVideo') },
])
// 「全部」里视频排在前面，和设计稿一致；同类保持接口给的顺序。
const visibleGallery = computed(() => {
  const list = gallery.value.filter((g) => filter.value === 'all' || g.mode === filter.value)
  if (filter.value === 'all') list.sort((a, b) => Number(b.mode === 'video') - Number(a.mode === 'video'))
  return list.slice(0, 24)
})
const featureImage = computed(() => gallery.value.find((g) => g.mode === 'video')?.image || MEDIA.beast)

async function loadGallery() {
  try {
    const res = await fetch('/api/v1/public/inspirations', { headers: { Accept: 'application/json' } })
    if (!res.ok) return
    const body = await res.json()
    gallery.value = body?.data?.items ?? []
    galleryTotal.value = body?.data?.total ?? 0
  } catch {
    // 工作台不可用时整段作品墙不显示
  }
}

const shots = computed(() => [
  { n: 1, name: t('home.lp.wsShot1'), dur: '7s', src: MEDIA.temple },
  { n: 2, name: t('home.lp.wsShot2'), dur: '8s', src: MEDIA.bridge },
  { n: 3, name: t('home.lp.wsShot3'), dur: '7s', src: MEDIA.sea },
  { n: 4, name: t('home.lp.wsShot4'), dur: '8s', src: MEDIA.beast },
])

const makeLines = computed(() => [1, 2, 3, 4, 5].map((i) => t(`home.lp.waysMake${i}`)))
const apiLines = computed(() => [1, 2, 3, 4, 5].map((i) => t(`home.lp.waysApi${i}`)))

const studioMenu = computed(() =>
  showWorkbench.value
    ? [
        { key: 'canvas', icon: 'canvas', to: '#lp-workbench', title: t('home.lp.menuCanvas'), desc: t('home.lp.menuCanvasDesc') },
        { key: 'agent', icon: 'agent', to: '#lp-workbench', title: t('home.lp.menuAgent'), desc: t('home.lp.menuAgentDesc') },
        { key: 'drama', icon: 'film', to: '/workbench/video', title: t('home.lp.menuDrama'), desc: t('home.lp.menuDramaDesc') },
        { key: 'image', icon: 'image', to: '/workbench/image', title: t('home.lp.menuImage'), desc: t('home.lp.menuImageDesc') },
        { key: 'skills', icon: 'skill', to: '/workbench/prompts', title: t('home.lp.menuSkills'), desc: t('home.lp.menuSkillsDesc', { n: SKILL_COUNT }) },
      ].map((i) => ({ ...i, external: false }))
    : [],
)

const devMenu = computed(() => {
  const items: { key: string; icon: string; to: string; title: string; desc: string; external: boolean }[] = [
    { key: 'api', icon: 'code', to: '#lp-work', title: t('home.lp.menuApi'), desc: t('home.lp.menuApiDesc'), external: false },
  ]
  if (showModelPlazaEntry.value) {
    items.push({ key: 'plaza', icon: 'grid', to: '/model-plaza', title: t('home.lp.menuPlaza'), desc: t('home.lp.menuPlazaDesc', { n: modelCountLabel.value }), external: false })
  }
  if (docsTarget.value) {
    items.push({ key: 'docs', icon: 'book', to: docsTarget.value, title: t('home.lp.menuDocs'), desc: t('home.lp.menuDocsDesc'), external: !docsIsRoute.value })
  }
  items.push({ key: 'usage', icon: 'meter', to: isAuthenticated.value ? '/usage' : '/key-usage', title: t('home.lp.menuUsage'), desc: t('home.lp.menuUsageDesc'), external: false })
  return items
})

const footerCols = computed(() => {
  type L = { key: string; label: string; to: string; external?: boolean }
  const product: L[] = showWorkbench.value
    ? [
        { key: 'studio', label: t('home.lp.ftStudio'), to: WORKBENCH },
        { key: 'canvas', label: t('home.lp.ftCanvas'), to: '#lp-workbench' },
        { key: 'skills', label: t('home.lp.ftSkills'), to: '/workbench/prompts' },
        { key: 'gallery', label: t('home.lp.ftGallery'), to: '#lp-gallery' },
      ]
    : [{ key: 'gallery', label: t('home.lp.ftGallery'), to: '#lp-gallery' }]
  const dev: L[] = [{ key: 'api', label: t('home.lp.ftApi'), to: '#lp-work' }]
  if (showModelPlazaEntry.value) dev.push({ key: 'plaza', label: t('home.lp.ftPlaza'), to: '/model-plaza' })
  if (docsTarget.value) dev.push({ key: 'docs', label: t('home.lp.ftDocs'), to: docsTarget.value, external: !docsIsRoute.value })
  dev.push({ key: 'keyusage', label: t('home.lp.ftKeyUsage'), to: '/key-usage' })
  const legal: L[] = (settings.value?.login_agreement_documents ?? []).map((d) => ({
    key: d.id,
    label: d.title,
    to: `/legal/${encodeURIComponent(d.id)}`,
  }))
  const cols = [
    { key: 'product', title: t('home.lp.ftProduct'), links: product },
    { key: 'dev', title: t('home.lp.ftDev'), links: dev },
  ]
  if (legal.length) cols.push({ key: 'legal', title: t('home.lp.ftLegal'), links: legal })
  return cols
})

const currentYear = new Date().getFullYear()

// 画布示意图按 560px 宽设计，窄屏时整体等比缩小。
const cvFitRef = ref<HTMLElement | null>(null)
const cvScale = ref(1)
let cvObserver: ResizeObserver | null = null
function fitCanvas() {
  const w = cvFitRef.value?.clientWidth ?? 560
  cvScale.value = Math.min(1, w / 560)
}

// ── 下拉菜单：点外面或按 Esc 收起 ──
const menuOpen = ref(false)
const dropRef = ref<HTMLElement | null>(null)
function onDocClick(e: MouseEvent) {
  if (menuOpen.value && dropRef.value && !dropRef.value.contains(e.target as Node)) menuOpen.value = false
}
function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape') menuOpen.value = false
}

onMounted(() => {
  document.addEventListener('click', onDocClick)
  document.addEventListener('keydown', onKey)
  if (typeof ResizeObserver !== 'undefined' && cvFitRef.value) {
    cvObserver = new ResizeObserver(fitCanvas)
    cvObserver.observe(cvFitRef.value)
  }
  fitCanvas()
  loadPlaza()
  loadAnnouncement()
  loadGallery()
})
onBeforeUnmount(() => {
  cvObserver?.disconnect()
  document.removeEventListener('click', onDocClick)
  document.removeEventListener('keydown', onKey)
})
</script>

<style scoped src="./landing-page.css"></style>
