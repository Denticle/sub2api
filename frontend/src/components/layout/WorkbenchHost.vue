<template>
  <div
    v-if="started && authStore.isAuthenticated"
    ref="host"
    class="workbench-host"
    :class="{ 'workbench-host--hidden': !visible }"
    :style="hostStyle"
    :aria-hidden="visible ? undefined : 'true'"
  >
    <iframe
      ref="frame"
      :src="initialSrc"
      class="studio-frame-iframe"
      :title="t('nav.workbench')"
      allow="clipboard-read; clipboard-write; fullscreen"
      referrerpolicy="same-origin"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch, watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'
import { isNavigationFailure, NavigationFailureType, useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useWorkbenchSlot } from '@/composables/useWorkbenchSlot'

/**
 * 工作台 iframe 的宿主，挂在 App 根上，整个会话只建一次。
 *
 * 为什么不放在 WorkbenchView 里：离开工作台时视图卸载，iframe 被销毁，
 * 画布里正在跑的生成任务就丢了。这里让 iframe 一直活着，盖在 WorkbenchView
 * 的占位元素上；离开工作台时隐藏（visibility + inert），任务在后台继续跑完，
 * 回来时结果已经在页面上。
 *
 * 和画布之间的消息：
 *   主站 -> 画布：{ type: "studio:navigate", path }，侧栏切功能时让画布在内部换页
 *   画布 -> 主站：{ type: "studio:route", path }，画布里自己换页时报回来，侧栏高亮跟上
 */
const TOOLS = ['canvas', 'image', 'video', 'prompts', 'assets'] as const
type Tool = (typeof TOOLS)[number]

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { slot, started, reset } = useWorkbenchSlot()

const host = ref<HTMLElement | null>(null)
const frame = ref<HTMLIFrameElement | null>(null)

function toolOf(value: unknown): Tool | '' {
  const v = Array.isArray(value) ? value[0] : value
  return typeof v === 'string' && (TOOLS as readonly string[]).includes(v) ? (v as Tool) : ''
}

const onWorkbench = computed(() => route.name === 'Workbench')
const routeTool = computed<Tool>(() => toolOf(route.params.tool) || 'canvas')
const visible = computed(() => onWorkbench.value && slot.value !== null)

// 只在创建 iframe 那一刻取一次，之后换功能走消息，不改 src（改 src 就是整页重载）。
// started 和 slot 在同一次 attach 里设好，这里赶在 iframe 渲染前同步取值。
const initialSrc = ref('/studio/canvas')
watch(
  started,
  (on) => {
    // ?v= 让浏览器每次建 iframe 都重新取 index.html：以前 /studio/ 没有 Cache-Control，
    // 浏览器会拿启发式缓存的旧页面，部署了新画布也不生效
    if (on) initialSrc.value = `/studio/${routeTool.value}?v=${Date.now()}`
  },
  { immediate: true, flush: 'sync' }
)

// ---- 盖在占位元素上 ----
const rect = ref({ top: 0, left: 0, width: 0, height: 0 })
function measure() {
  const el = slot.value
  if (!el) return
  const r = el.getBoundingClientRect()
  rect.value = { top: r.top, left: r.left, width: r.width, height: r.height }
}
const hostStyle = computed(() => ({
  top: `${rect.value.top}px`,
  left: `${rect.value.left}px`,
  width: `${rect.value.width}px`,
  height: `${rect.value.height}px`
}))

// 侧栏收起 / 展开、移动端抽屉都会挪动占位元素的位置却不一定改它的尺寸，
// ResizeObserver 接不住；显示期间每帧读一次位置，变了才写回，隐藏时停掉。
let raf = 0
function track() {
  const el = slot.value
  if (!el) return
  const r = el.getBoundingClientRect()
  const c = rect.value
  if (r.top !== c.top || r.left !== c.left || r.width !== c.width || r.height !== c.height) measure()
  raf = requestAnimationFrame(track)
}
watch(
  visible,
  (on) => {
    cancelAnimationFrame(raf)
    if (on) {
      measure()
      raf = requestAnimationFrame(track)
    }
  },
  { flush: 'post', immediate: true }
)

// 隐藏时用 inert 拦住焦点和点击；Vue 对 inert 的模板绑定在各版本间不一致，直接设属性
watchEffect(() => {
  if (host.value) host.value.inert = !visible.value
})

// ---- 和画布对话 ----
let childReady = false
let childTool: Tool | '' = ''

function navigateChild(tool: Tool) {
  frame.value?.contentWindow?.postMessage({ type: 'studio:navigate', path: `/${tool}` }, window.location.origin)
}

// 侧栏切功能 -> 画布内部换页。离开工作台时不动画布
watch(routeTool, (tool) => {
  if (onWorkbench.value && childReady && tool !== childTool) navigateChild(tool)
})

// 已在「我的画布」里某个画布中，再点一次「我的画布」：路由不变，watch 不触发，在这里让画布回列表
const removeAfterEach = router.afterEach((to, _from, failure) => {
  if (to.name !== 'Workbench' || !childReady) return
  if (isNavigationFailure(failure, NavigationFailureType.duplicated)) navigateChild(routeTool.value)
})

function onMessage(e: MessageEvent) {
  if (e.origin !== window.location.origin) return
  if (!frame.value || e.source !== frame.value.contentWindow) return
  const data = e.data
  if (!data || data.type !== 'studio:route' || typeof data.path !== 'string') return

  const tool = toolOf(data.path.replace(/^\/+/, '').split('/')[0])
  if (!childReady) {
    // 画布刚启动：以侧栏当前选中的为准
    childReady = true
    childTool = tool
    if (onWorkbench.value && tool !== routeTool.value) navigateChild(routeTool.value)
    return
  }
  childTool = tool
  // 画布里自己换了功能：侧栏高亮跟过去。人不在工作台时别把他拽回来
  if (onWorkbench.value && tool && tool !== routeTool.value) {
    router.replace({ name: 'Workbench', params: { tool } })
  }
}

// 退出登录：销毁 iframe，换账号进来重新加载，不沿用上一个人的画布状态
watch(
  () => authStore.isAuthenticated,
  (authed) => {
    if (authed) return
    childReady = false
    childTool = ''
    reset()
  }
)

onMounted(() => window.addEventListener('message', onMessage))
onBeforeUnmount(() => {
  window.removeEventListener('message', onMessage)
  cancelAnimationFrame(raf)
  removeAfterEach()
})
</script>

<style scoped>
.workbench-host {
  position: fixed;
  /* 低于顶栏下拉（z-40）、移动端侧栏（z-30 / z-50）、Toast，高于普通页面内容 */
  z-index: 5;
}

.workbench-host--hidden {
  visibility: hidden;
  pointer-events: none;
}
</style>
