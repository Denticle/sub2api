import { ref, shallowRef } from 'vue'

/**
 * 工作台 iframe 的「占位」与「宿主」之间共享的状态。
 *
 * iframe 不能放进 WorkbenchView 里：离开工作台（去 API 密钥、使用记录……）时视图被卸载，
 * iframe 跟着销毁，画布里正在跑的生成任务也就没了。所以 iframe 由挂在 App 根上的
 * WorkbenchHost 持有，整个会话只建一次；WorkbenchView 只渲染一个占位元素，
 * 宿主把 iframe 盖在占位的位置上。离开工作台时占位消失，iframe 隐藏但继续运行。
 */
const slot = shallowRef<HTMLElement | null>(null)
/** 用户进过工作台之后才创建 iframe，没进过就不加载 3MB 的画布 */
const started = ref(false)

export function useWorkbenchSlot() {
  return {
    slot,
    started,
    attach(el: HTMLElement) {
      slot.value = el
      started.value = true
    },
    detach(el: HTMLElement) {
      if (slot.value === el) slot.value = null
    },
    /** 退出登录时销毁，下一个账号进来重新加载 */
    reset() {
      slot.value = null
      started.value = false
    }
  }
}
