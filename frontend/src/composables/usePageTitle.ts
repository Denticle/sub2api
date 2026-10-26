import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useAdminSettingsStore } from '@/stores/adminSettings'
import { resolveRouteMetaKeys } from '@/router/title'
import { resolveSiteBillingMode } from '@/utils/siteBillingMode'

/**
 * 页面标题 / 描述的单一来源。
 *
 * 顶栏（面包屑层级）与内容区（页面主标题）都要用同一份值，
 * 因此抽到这里，避免两处各自维护 route.meta 解析逻辑而漂移。
 */
export function usePageTitle() {
  const route = useRoute()
  const { t } = useI18n()
  const appStore = useAppStore()
  const authStore = useAuthStore()
  const adminSettingsStore = useAdminSettingsStore()

  // /purchase 的标题/描述随站点计费模式切换，与 document.title 共用同一解析。
  const routeMetaKeys = computed(() =>
    resolveRouteMetaKeys(route, {
      billingMode: resolveSiteBillingMode(appStore.cachedPublicSettings),
    }),
  )

  const pageTitle = computed(() => {
    // 自定义页面优先用菜单项的标签，而不是通用的「自定义页面」
    if (route.name === 'CustomPage') {
      const id = route.params.id as string
      const publicItems = appStore.cachedPublicSettings?.custom_menu_items ?? []
      const menuItem =
        publicItems.find((item) => item.id === id) ??
        (authStore.isAdmin
          ? adminSettingsStore.customMenuItems.find((item) => item.id === id)
          : undefined)
      if (menuItem?.label) return menuItem.label
    }
    const titleKey = routeMetaKeys.value.titleKey
    if (titleKey) return t(titleKey)
    return (route.meta.title as string) || ''
  })

  const pageDescription = computed(() => {
    const descKey = routeMetaKeys.value.descriptionKey
    if (descKey) return t(descKey)
    return (route.meta.description as string) || ''
  })

  return { pageTitle, pageDescription }
}
