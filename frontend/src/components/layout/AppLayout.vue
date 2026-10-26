<template>
  <div class="min-h-screen bg-gray-50 dark:bg-dark-950">
    <!-- Header: 通栏，跨整个视口宽度（OpenRouter 的结构） -->
    <AppHeader />

    <!-- Body: 侧栏 + 内容，位于顶栏下方 -->
    <div class="relative flex">
      <!-- Sidebar -->
      <AppSidebar />

      <!-- Main Content Area -->
      <div class="relative min-h-[calc(100vh-4rem)] min-w-0 flex-1 transition-all duration-300">
        <!-- Main Content -->
        <main class="p-4 md:p-6 lg:p-8">
          <!-- 页面主标题：视觉锚点。个别视图自带 h1.page-title，
               由样式表中的 :has() 规则自动隐藏这里的副本，避免重复。 -->
          <div v-if="pageTitle" class="app-page-title">
            <h1>{{ pageTitle }}</h1>
            <p v-if="pageDescription">{{ pageDescription }}</p>
          </div>

          <slot />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import '@/styles/onboarding.css'
import { computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useOnboardingTour } from '@/composables/useOnboardingTour'
import { useOnboardingStore } from '@/stores/onboarding'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
import { usePageTitle } from '@/composables/usePageTitle'

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.user?.role === 'admin')

const { pageTitle, pageDescription } = usePageTitle()

const { replayTour } = useOnboardingTour({
  storageKey: isAdmin.value ? 'admin_guide' : 'user_guide',
  autoStart: true
})

const onboardingStore = useOnboardingStore()

onMounted(() => {
  onboardingStore.setReplayCallback(replayTour)
})

defineExpose({ replayTour })
</script>
