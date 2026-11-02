<template>
  <AppLayout>
    <!-- 只是占位：iframe 由 App 根上的 WorkbenchHost 持有并盖在这里，
         离开工作台时不会被销毁。见 composables/useWorkbenchSlot.ts -->
    <div ref="frameSlot" class="studio-frame" />
  </AppLayout>
</template>

<script setup lang="ts">
import '@/styles/studio.css'
import { onBeforeUnmount, onMounted, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useWorkbenchSlot } from '@/composables/useWorkbenchSlot'

const frameSlot = ref<HTMLElement | null>(null)
const { attach, detach } = useWorkbenchSlot()

onMounted(() => {
  if (frameSlot.value) attach(frameSlot.value)
})
onBeforeUnmount(() => {
  if (frameSlot.value) detach(frameSlot.value)
})
</script>
