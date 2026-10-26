<template>
  <header class="sticky top-0 z-30 border-b border-gray-200/70 bg-gray-50/85 backdrop-blur-xl dark:border-dark-800/80 dark:bg-dark-950/85">
    <div class="flex h-14 w-full items-center justify-between gap-4 px-6 lg:px-10">
      <!-- 左:站点 logo + 名称 -->
      <div class="flex min-w-0 items-center gap-2">
        <template v-if="settings">
          <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-7 w-7 rounded-lg object-contain" />
          <span class="truncate text-[15px] font-bold tracking-tight text-gray-900 dark:text-white">
            {{ siteName }}
          </span>
        </template>
        <template v-else>
          <span class="h-7 w-7 flex-shrink-0 animate-pulse rounded-lg bg-gray-200 dark:bg-dark-800" aria-hidden="true"></span>
          <span class="h-4 w-28 animate-pulse rounded bg-gray-200 dark:bg-dark-800" aria-hidden="true"></span>
        </template>
      </div>

      <!-- 右:登录 / 回到后台 -->
      <RouterLink
        v-if="isAuthenticated"
        :to="backTarget"
        class="btn btn-primary flex-shrink-0"
      >
        {{ t('modelPlaza.nav.backToDashboard') }}
      </RouterLink>
      <RouterLink
        v-else
        :to="{ path: '/login', query: { redirect: '/model-plaza' } }"
        class="btn btn-primary flex-shrink-0"
      >
        {{ t('modelPlaza.nav.login') }}
      </RouterLink>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { sanitizeUrl } from '@/utils/url'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const settings = computed(() => appStore.cachedPublicSettings)
const siteName = computed(() => settings.value?.site_name || 'Sub2API')
const siteLogo = computed(() =>
  sanitizeUrl(settings.value?.site_logo || '', { allowRelative: true, allowDataUrl: true })
)
const isAuthenticated = computed(() => authStore.isAuthenticated)
const backTarget = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))
</script>
