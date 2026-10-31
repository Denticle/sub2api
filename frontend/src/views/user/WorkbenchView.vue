<template>
  <AppLayout>
    <div class="studio-frame">
      <iframe
        ref="frame"
        :src="initialSrc"
        class="studio-frame-iframe"
        :title="t('nav.workbench')"
        allow="clipboard-read; clipboard-write; fullscreen"
        referrerpolicy="same-origin"
      />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import '@/styles/studio.css'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { isNavigationFailure, NavigationFailureType, useRoute, useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'

/**
 * 工作台。
 *
 * 画布是独立构建的 React SPA（挂在 /studio/），这里用同源 iframe 装进主站外壳。
 * 侧栏「工作台」分组的五项都落到这个视图，靠 :tool 区分：
 *
 *   /workbench/canvas  -> /studio/canvas
 *   /workbench/image   -> /studio/image
 *   ...
 *
 * iframe 只在首次挂载时设 src；之后切换功能用 postMessage 让画布在内部换页，
 * 不整页重载（画布的 bundle 有 3.6MB，重载一次就是一次白屏）。
 * 画布内部换页（比如从画布列表进到某个画布，再点返回）也会报回来，
 * 主站只同步顶层功能名，让侧栏高亮跟得上。
 */
const TOOLS = ['canvas', 'image', 'video', 'prompts', 'assets'] as const
type Tool = (typeof TOOLS)[number]

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const frame = ref<HTMLIFrameElement | null>(null)

function routeTool(): Tool {
  const raw = route.params.tool
  const v = Array.isArray(raw) ? raw[0] : raw
  return (TOOLS as readonly string[]).includes(v) ? (v as Tool) : 'canvas'
}

const initialSrc = `/studio/${routeTool()}`

/** 画布已经启动、能接消息 */
let childReady = false
/** 画布报回来的当前顶层功能 */
let childTool: Tool | '' = ''

function navigateChild(tool: Tool) {
  frame.value?.contentWindow?.postMessage({ type: 'studio:navigate', path: `/${tool}` }, window.location.origin)
}

// 侧栏切功能 -> 画布内部换页
watch(
  () => route.params.tool,
  () => {
    const tool = routeTool()
    if (childReady && tool !== childTool) navigateChild(tool)
  }
)

/**
 * 已经在「我的画布」但人在某个画布里时，再点一次侧栏的「我的画布」，
 * 路由参数不变、watch 不触发——用户会以为点了没反应。这里接住同路由的重复点击，
 * 让画布回到列表。
 */
const removeAfterEach = router.afterEach((to, _from, failure) => {
  if (to.name !== 'Workbench' || !childReady) return
  if (isNavigationFailure(failure, NavigationFailureType.duplicated)) navigateChild(routeTool())
})

function onMessage(e: MessageEvent) {
  if (e.origin !== window.location.origin) return
  if (e.source !== frame.value?.contentWindow) return
  const data = e.data
  if (!data || data.type !== 'studio:route' || typeof data.path !== 'string') return

  const first = data.path.replace(/^\/+/, '').split('/')[0]
  const tool = (TOOLS as readonly string[]).includes(first) ? (first as Tool) : ''

  if (!childReady) {
    // 画布刚启动：它落在哪不重要，以侧栏当前选中的为准
    childReady = true
    childTool = tool
    if (tool !== routeTool()) navigateChild(routeTool())
    return
  }

  childTool = tool
  // 画布里自己换了功能：侧栏高亮跟过去，不在主站历史里多压一条
  if (tool && tool !== routeTool()) router.replace({ name: 'Workbench', params: { tool } })
}

onMounted(() => window.addEventListener('message', onMessage))
onBeforeUnmount(() => {
  window.removeEventListener('message', onMessage)
  removeAfterEach()
})
</script>
