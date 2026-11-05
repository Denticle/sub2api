import { computed, ref } from 'vue'
import { getPublicConfig, studioLoginUrl, type OAuthPublicConfig } from '@/api/oauthProvider'

// 顶栏和首页共用一份；一次页面会话内只请求一次。
const config = ref<OAuthPublicConfig | null>(null)
let pending: Promise<void> | null = null

function load(): Promise<void> {
  if (!pending) {
    pending = getPublicConfig()
      .then((value) => {
        config.value = value
      })
      .catch(() => {
        // 拿不到就不显示入口；允许下次挂载时重试。
        pending = null
      })
  }
  return pending
}

/** 视频工作台入口：后台启用且配置了地址才显示，链接直接走工作台的 forapi 登录。 */
export function useStudioEntry() {
  void load()
  const studioUrl = computed(() => {
    const value = config.value
    if (!value?.enabled || !value.app_url) return ''
    return studioLoginUrl(value.app_url)
  })
  return { studioUrl }
}
