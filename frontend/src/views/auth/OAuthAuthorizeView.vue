<template>
  <AuthLayout>
    <!-- 已授权过的用户只会看到「正在跳转」一闪而过 -->
    <div v-if="phase === 'checking' || phase === 'redirecting'" class="py-6 text-center">
      <div class="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
      <p class="text-sm text-gray-500 dark:text-dark-400">
        {{ phase === 'checking' ? t('oauthProvider.checking') : t('oauthProvider.redirecting') }}
      </p>
    </div>

    <div v-else-if="phase === 'consent'" class="space-y-6">
      <div>
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
          {{ t('oauthProvider.approveTitle', { app: appName }) }}
        </h2>
        <p class="mt-2 text-sm leading-6 text-gray-600 dark:text-dark-300">
          {{ t('oauthProvider.approveIntro', { app: appName, site: siteName }) }}
        </p>
      </div>

      <div>
        <p class="mb-2 text-xs font-medium uppercase tracking-wide text-gray-500 dark:text-dark-400">
          {{ t('oauthProvider.permissionsTitle') }}
        </p>
        <ul class="space-y-2 text-sm text-gray-700 dark:text-dark-200">
          <li class="flex gap-2">
            <Icon name="user" size="sm" class="mt-0.5 shrink-0 text-primary-600 dark:text-primary-400" />
            <span>{{ t('oauthProvider.permProfile') }}</span>
          </li>
          <li class="flex gap-2">
            <Icon name="key" size="sm" class="mt-0.5 shrink-0 text-primary-600 dark:text-primary-400" />
            <span>
              {{ t('oauthProvider.permKey') }}
              <span class="mt-1 block text-xs text-gray-500 dark:text-dark-400">{{ t('oauthProvider.permKeyHint') }}</span>
            </span>
          </li>
        </ul>
      </div>

      <div v-if="user" class="rounded-lg bg-gray-50 px-3 py-2 text-sm dark:bg-dark-800">
        <span class="text-gray-500 dark:text-dark-400">{{ t('oauthProvider.signedInAs') }}</span>
        <span class="ml-2 font-medium text-gray-900 dark:text-white">{{ user.email || user.username }}</span>
      </div>

      <div class="flex gap-3">
        <button type="button" class="btn btn-secondary flex-1" :disabled="submitting" @click="cancel">
          {{ t('oauthProvider.cancel') }}
        </button>
        <button type="button" class="btn btn-primary flex-1" :disabled="submitting" @click="approve">
          {{ t('oauthProvider.approve') }}
        </button>
      </div>
    </div>

    <div v-else class="space-y-4 text-center">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('oauthProvider.failedTitle') }}</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300">{{ errorMessage }}</p>
      <router-link to="/" class="btn btn-secondary">{{ t('oauthProvider.backHome') }}</router-link>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { AuthLayout } from '@/components/layout'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore, useAuthStore } from '@/stores'
import { extractApiErrorMessage } from '@/utils/apiError'
import oauthAPI, { type AuthorizeParams } from '@/api/oauthProvider'

/**
 * OAuth 授权确认页（视频工作台登录桥接）。
 *
 * 路由要求登录：未登录的用户会被路由守卫带去登录，登录后回到这里。
 * 授权过一次的用户不显示确认页，直接签发授权码并跳回工作台。
 */
const { t } = useI18n()
const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()

type Phase = 'checking' | 'consent' | 'redirecting' | 'error'
const phase = ref<Phase>('checking')
const appName = ref('')
const errorMessage = ref('')
const submitting = ref(false)

const user = computed(() => authStore.user)
const siteName = computed(() => appStore.siteName || 'ForAPI')

const params = computed<AuthorizeParams>(() => {
  const q = route.query
  const s = (k: string) => (typeof q[k] === 'string' ? (q[k] as string) : '')
  return {
    client_id: s('client_id'),
    redirect_uri: s('redirect_uri'),
    response_type: s('response_type'),
    state: s('state'),
    code_challenge: s('code_challenge'),
    code_challenge_method: s('code_challenge_method')
  }
})

function fail(err: unknown) {
  errorMessage.value = extractApiErrorMessage(err, t('oauthProvider.failedTitle'))
  phase.value = 'error'
}

async function approve() {
  submitting.value = true
  try {
    const target = await oauthAPI.approveAuthorize(params.value)
    phase.value = 'redirecting'
    // 跳到外部站点：只允许后端返回的、已与配置精确比对过的回调地址
    window.location.replace(target)
  } catch (err) {
    fail(err)
  } finally {
    submitting.value = false
  }
}

function cancel() {
  window.location.replace('/')
}

onMounted(async () => {
  try {
    const check = await oauthAPI.checkAuthorize(params.value)
    appName.value = check.app_name || t('studio.nav')
    if (check.consented) {
      await approve()
      return
    }
    phase.value = 'consent'
  } catch (err) {
    fail(err)
  }
})
</script>
