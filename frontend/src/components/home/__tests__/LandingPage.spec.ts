import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount, RouterLinkStub } from '@vue/test-utils'

import LandingPage from '../LandingPage.vue'

const { appStore, authStore } = vi.hoisted(() => ({
  appStore: {
    cachedPublicSettings: {} as Record<string, unknown>,
    siteName: 'Fallback site',
    siteLogo: '',
    docUrl: '',
    backendModeEnabled: false,
  },
  authStore: {
    isAuthenticated: false,
    isAdmin: false,
    isSimpleMode: false,
  },
}))

vi.mock('@/stores', () => ({
  useAppStore: () => appStore,
  useAuthStore: () => authStore,
}))

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-i18n')>()
  return { ...actual, useI18n: () => ({ t: (key: string) => key, locale: { value: 'zh' } }) }
})

vi.mock('@/utils/featureFlags', () => ({
  FeatureFlags: { modelPlaza: 'modelPlaza' },
  isFeatureFlagEnabled: () => true,
}))

type Routes = Record<string, unknown>

function stubFetch(routes: Routes) {
  const fetchMock = vi.fn(async (url: string) => {
    const body = routes[url]
    if (body === undefined) return { ok: false, json: async () => ({}) }
    return { ok: true, json: async () => ({ code: 0, data: body }) }
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

const plaza = {
  groups: [
    { platform: 'anthropic', models: [{ name: 'claude-a' }, { name: 'claude-b' }] },
    { platform: 'openai', models: [{ name: 'gpt-a' }, { name: 'claude-a' }] },
    { platform: 'gemini', models: [] },
  ],
}

const inspirations = {
  total: 320,
  items: [
    { id: '1', title: '仙侠短片', image: 'https://cdn.example/1.jpg', mode: 'video', credit: 'Seedance 2.5' },
    { id: '2', title: '产品海报', image: 'https://cdn.example/2.jpg', mode: 'image', credit: '蛋黄堡' },
  ],
}

async function mountLanding(routes: Routes = {}, settings: Record<string, unknown> = {}) {
  appStore.cachedPublicSettings = { site_name: 'ForApi', registration_enabled: true, ...settings }
  stubFetch(routes)
  const wrapper = mount(LandingPage, {
    props: { isDark: false },
    global: {
      stubs: {
        RouterLink: RouterLinkStub,
        LocaleSwitcher: { template: '<div data-testid="locale-switcher" />' },
      },
    },
  })
  await flushPromises()
  return wrapper
}

const linkTargets = (wrapper: Awaited<ReturnType<typeof mountLanding>>) =>
  wrapper.findAllComponents(RouterLinkStub).map((l) => l.props('to'))

describe('LandingPage', () => {
  beforeEach(() => {
    authStore.isAuthenticated = false
    authStore.isAdmin = false
    authStore.isSimpleMode = false
    appStore.backendModeEnabled = false
    localStorage.clear()
  })
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('counts distinct models and providers from the model plaza', async () => {
    const wrapper = await mountLanding({ '/api/v1/model-plaza': plaza })
    const stats = wrapper.get('[data-testid="home-providers"]').findAll('dd').map((d) => d.text())
    expect(stats[0]).toContain('3')
    expect(stats[1]).toContain('2')
    expect(wrapper.get('.lp-brand-name').text()).toBe('ForApi')
  })

  it('shows the latest public announcement and remembers dismissal', async () => {
    const wrapper = await mountLanding({
      '/api/v1/public/announcements': [{ id: 7, title: 'Seedance 2.5 已接入' }],
    })
    expect(wrapper.get('[data-testid="home-announce"]').text()).toContain('Seedance 2.5 已接入')

    await wrapper.get('.nv-announce-x').trigger('click')
    expect(wrapper.find('[data-testid="home-announce"]').exists()).toBe(false)
    expect(localStorage.getItem('home_announce_dismissed')).toBe('7')

    const again = await mountLanding({
      '/api/v1/public/announcements': [{ id: 7, title: 'Seedance 2.5 已接入' }],
    })
    expect(again.find('[data-testid="home-announce"]').exists()).toBe(false)
  })

  it('hides the announcement bar when there are no announcements', async () => {
    const wrapper = await mountLanding({ '/api/v1/public/announcements': [] })
    expect(wrapper.find('[data-testid="home-announce"]').exists()).toBe(false)
  })

  it('renders the gallery and filters by media type', async () => {
    const wrapper = await mountLanding({ '/api/v1/public/inspirations': inspirations })
    const gallery = wrapper.get('[data-testid="home-gallery"]')
    expect(gallery.findAll('.lp-shot')).toHaveLength(2)

    const videoBtn = gallery.findAll('.lp-filter button')[2]
    await videoBtn.trigger('click')
    const shots = gallery.findAll('.lp-shot')
    expect(shots).toHaveLength(1)
    expect(shots[0].attributes('title')).toBe('仙侠短片')
  })

  it('omits the gallery section when the studio feed is unavailable', async () => {
    const wrapper = await mountLanding()
    expect(wrapper.find('[data-testid="home-gallery"]').exists()).toBe(false)
  })

  it('sends visitors to register and the workbench', async () => {
    const wrapper = await mountLanding()
    const targets = linkTargets(wrapper)
    expect(targets).toContain('/register')
    expect(targets).toContain('/login')
    expect(targets).toContain('/workbench/canvas')
    expect(wrapper.find('#lp-workbench').exists()).toBe(true)
  })

  it('falls back to login when registration is closed', async () => {
    const wrapper = await mountLanding({}, { registration_enabled: false })
    expect(linkTargets(wrapper)).not.toContain('/register')
  })

  it('shows the console link to signed-in admins', async () => {
    authStore.isAuthenticated = true
    authStore.isAdmin = true
    const wrapper = await mountLanding()
    const targets = linkTargets(wrapper)
    expect(targets).toContain('/admin/dashboard')
    expect(targets).not.toContain('/login')
  })

  it('does not advertise the workbench in backend mode', async () => {
    appStore.backendModeEnabled = true
    const wrapper = await mountLanding()
    expect(wrapper.find('#lp-workbench').exists()).toBe(false)
    expect(linkTargets(wrapper)).not.toContain('/workbench/canvas')
  })

  it('lists legal documents from settings in the footer', async () => {
    const wrapper = await mountLanding({}, {
      login_agreement_documents: [{ id: 'terms', title: '服务条款', content_md: '' }],
    })
    expect(linkTargets(wrapper)).toContain('/legal/terms')
  })

  it('emits toggle-theme from the nav button', async () => {
    const wrapper = await mountLanding()
    await wrapper.get('.nv-icon').trigger('click')
    expect(wrapper.emitted('toggle-theme')).toHaveLength(1)
  })

  it('closes the product menu on Escape', async () => {
    const wrapper = await mountLanding()
    const trigger = wrapper.get('.nv-trigger')
    await trigger.trigger('click')
    expect(trigger.attributes('aria-expanded')).toBe('true')
    document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    await flushPromises()
    expect(trigger.attributes('aria-expanded')).toBe('false')
  })
})
