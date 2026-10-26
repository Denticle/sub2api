<template>
  <div class="rounded-xl border border-gray-200 bg-white dark:border-dark-800 dark:bg-dark-900">
    <!-- 工具条：搜索在前（模型页的主要操作），筛选在后 -->
    <div class="flex flex-col gap-3 p-3 sm:flex-row sm:items-center">
      <div class="relative w-full sm:max-w-xs">
        <Icon name="search" size="sm" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-dark-600" />
        <input
          :value="search"
          type="text"
          :placeholder="t('modelPlaza.filters.searchPlaceholder')"
          class="input w-full pl-9 pr-9"
          @input="$emit('update:search', ($event.target as HTMLInputElement).value)"
        />
        <button
          v-if="search"
          type="button"
          class="absolute right-2.5 top-1/2 -translate-y-1/2 text-gray-400 transition-colors hover:text-gray-600 dark:text-dark-600 dark:hover:text-dark-300"
          @click="$emit('update:search', '')"
        >
          <Icon name="x" size="xs" />
        </button>
      </div>

      <div class="flex flex-wrap items-center gap-1.5 sm:ml-auto">
        <!-- 平台 -->
        <button
          v-for="p in ['all', ...platforms]"
          :key="`platform-${p}`"
          type="button"
          class="inline-flex h-7 items-center gap-1.5 rounded-md px-2.5 text-[12.5px] font-medium transition disabled:cursor-not-allowed disabled:opacity-35"
          :class="chipClass(platform === p)"
          :disabled="p !== 'all' && !platformEnabled(p)"
          @click="$emit('update:platform', p)"
        >
          <PlatformIcon v-if="p !== 'all'" :platform="p as GroupPlatform" size="xs" />
          {{ p === 'all' ? t('modelPlaza.filters.allPlatforms') : platformLabel(p) }}
        </button>
      </div>
    </div>

    <!-- 第二行：分组与倍率，仅在存在多个可选项时出现，避免单值噪音 -->
    <div
      v-if="groups.length > 1 || rates.length > 1"
      class="flex flex-wrap items-center gap-x-4 gap-y-2 border-t border-gray-200/70 px-3 py-2.5 dark:border-dark-800"
    >
      <div v-if="groups.length > 1" class="flex flex-wrap items-center gap-1.5">
        <span class="text-[11.5px] font-semibold text-gray-400 dark:text-dark-600">{{ t('modelPlaza.filters.groupLabel') }}</span>
        <button
          type="button"
          class="h-6 rounded px-2 text-[12px] font-medium transition"
          :class="chipClass(groupId === 'all')"
          @click="$emit('update:groupId', 'all')"
        >{{ t('modelPlaza.filters.all') }}</button>
        <button
          v-for="g in groups"
          :key="`group-${g.id}`"
          type="button"
          class="h-6 rounded px-2 text-[12px] font-medium transition disabled:cursor-not-allowed disabled:opacity-35"
          :class="chipClass(groupId === g.id)"
          :disabled="!groupEnabled(g)"
          @click="$emit('update:groupId', g.id)"
        >{{ g.name }}</button>
      </div>

      <div v-if="rates.length > 1" class="flex flex-wrap items-center gap-1.5">
        <span class="text-[11.5px] font-semibold text-gray-400 dark:text-dark-600">{{ t('modelPlaza.filters.rateLabel') }}</span>
        <button
          type="button"
          class="h-6 rounded px-2 font-mono text-[12px] font-medium transition"
          :class="chipClass(rate === 'all')"
          @click="$emit('update:rate', 'all')"
        >{{ t('modelPlaza.filters.all') }}</button>
        <button
          v-for="r in rates"
          :key="`rate-${r}`"
          type="button"
          class="h-6 rounded px-2 font-mono text-[12px] font-medium transition disabled:cursor-not-allowed disabled:opacity-35"
          :class="chipClass(rate === r)"
          :disabled="!rateEnabled(r)"
          @click="$emit('update:rate', r)"
        >{{ r }}x</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import PlatformIcon from '@/components/common/PlatformIcon.vue'
import { platformLabel } from '@/utils/platformColors'
import type { GroupPlatform } from '@/types'

const props = defineProps<{
  /** 数据中出现的平台(去重排序后)。 */
  platforms: string[]
  /** 全量分组(含平台与生效倍率),三个维度的置灰联动由此推导。 */
  groups: Array<{ id: number; name: string; platform: string; rate: number }>
  /** 全量生效倍率去重升序。 */
  rates: number[]
  platform: string
  groupId: number | 'all'
  rate: number | 'all'
  /** 模型名搜索词(纯前端过滤)。 */
  search: string
}>()

defineEmits<{
  'update:platform': [value: string]
  'update:groupId': [value: number | 'all']
  'update:rate': [value: number | 'all']
  'update:search': [value: string]
}>()

const { t } = useI18n()

/**
 * 三个维度互为约束(faceted):某选项可点 ⟺ 在「其他两维」当前选择下仍有分组命中。
 * 「全部」永远可点,作为解除本维约束的出口;可点项组合恒有结果,无需选择修正。
 */
function platformEnabled(p: string): boolean {
  return props.groups.some(
    (g) =>
      g.platform === p &&
      (props.groupId === 'all' || g.id === props.groupId) &&
      (props.rate === 'all' || g.rate === props.rate)
  )
}

function groupEnabled(g: { platform: string; rate: number }): boolean {
  return (
    (props.platform === 'all' || g.platform === props.platform) &&
    (props.rate === 'all' || g.rate === props.rate)
  )
}

function rateEnabled(r: number): boolean {
  return props.groups.some(
    (g) =>
      g.rate === r &&
      (props.platform === 'all' || g.platform === props.platform) &&
      (props.groupId === 'all' || g.id === props.groupId)
  )
}

function chipClass(active: boolean): string {
  return active
    ? 'bg-primary-500/12 text-primary-700 dark:bg-primary-500/15 dark:text-primary-400'
    : 'text-gray-600 enabled:hover:bg-gray-100 enabled:hover:text-gray-900 dark:text-dark-400 dark:enabled:hover:bg-dark-800 dark:enabled:hover:text-white'
}
</script>

<style scoped>
/* 筛选控件统一走中性色 + 主色选中态。
   上游按平台主色给 chip 上色，一屏会出现七八种颜色，像色卡。 */
</style>
