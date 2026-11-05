<template>
  <div class="space-y-2" data-test="key-group-list">
    <ol v-if="modelValue.length > 0" class="space-y-1.5">
      <li
        v-for="(groupId, index) in modelValue"
        :key="groupId"
        class="flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-2.5 py-2 dark:border-dark-600 dark:bg-dark-800"
        data-test="key-group-row"
      >
        <span
          class="w-14 shrink-0 text-xs font-medium text-gray-500 dark:text-gray-400"
          data-test="key-group-route-label"
        >{{ t('keys.groupList.route', { n: index + 2 }) }}</span>
        <span class="min-w-0 flex-1">
          <GroupBadge
            v-if="optionById(groupId)"
            :name="optionById(groupId)!.label"
            :platform="optionById(groupId)!.platform"
            :subscription-type="optionById(groupId)!.subscriptionType"
            :rate-multiplier="optionById(groupId)!.rate"
            :user-rate-multiplier="optionById(groupId)!.userRate"
            :peak-rate-enabled="optionById(groupId)!.peakRateEnabled"
            :peak-start="optionById(groupId)!.peakStart"
            :peak-end="optionById(groupId)!.peakEnd"
            :peak-rate-multiplier="optionById(groupId)!.peakRateMultiplier"
          />
          <span v-else class="text-sm text-gray-400 dark:text-dark-500" data-test="key-group-unavailable">
            {{ t('keys.groupList.unavailable', { id: groupId }) }}
          </span>
        </span>
        <span class="flex shrink-0 items-center gap-0.5">
          <button
            type="button"
            class="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-700 disabled:cursor-not-allowed disabled:opacity-30 dark:hover:bg-dark-700 dark:hover:text-gray-200"
            :disabled="index === 0 && !canPromote"
            :aria-label="index === 0 ? t('keys.groupList.setPrimary') : t('keys.groupList.moveUp')"
            :title="index === 0 ? t('keys.groupList.setPrimary') : t('keys.groupList.moveUp')"
            data-test="key-group-move-up"
            @click="moveUp(index)"
          >
            <Icon name="chevronUp" size="sm" />
          </button>
          <button
            type="button"
            class="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-700 disabled:cursor-not-allowed disabled:opacity-30 dark:hover:bg-dark-700 dark:hover:text-gray-200"
            :disabled="index === modelValue.length - 1"
            :aria-label="t('keys.groupList.moveDown')"
            :title="t('keys.groupList.moveDown')"
            data-test="key-group-move-down"
            @click="moveDown(index)"
          >
            <Icon name="chevronDown" size="sm" />
          </button>
          <button
            type="button"
            class="rounded p-1 text-gray-400 hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20 dark:hover:text-red-400"
            :aria-label="t('keys.groupList.remove')"
            :title="t('keys.groupList.remove')"
            data-test="key-group-remove"
            @click="remove(index)"
          >
            <Icon name="x" size="sm" />
          </button>
        </span>
      </li>
    </ol>

    <Select
      v-if="modelValue.length < maxExtra"
      :key="addSelectKey"
      :model-value="null"
      :options="addableOptions"
      :placeholder="t('keys.groupList.add')"
      :empty-text="t('keys.groupList.noMoreGroups')"
      :searchable="true"
      :search-placeholder="t('keys.searchGroup')"
      :aria-label="t('keys.groupList.add')"
      data-test="key-group-add"
      @update:model-value="add"
    >
      <template #option="{ option, selected }">
        <GroupOptionItem
          :name="(option as unknown as KeyGroupOption).label"
          :platform="(option as unknown as KeyGroupOption).platform"
          :subscription-type="(option as unknown as KeyGroupOption).subscriptionType"
          :rate-multiplier="(option as unknown as KeyGroupOption).rate"
          :user-rate-multiplier="(option as unknown as KeyGroupOption).userRate"
          :peak-rate-enabled="(option as unknown as KeyGroupOption).peakRateEnabled"
          :peak-start="(option as unknown as KeyGroupOption).peakStart"
          :peak-end="(option as unknown as KeyGroupOption).peakEnd"
          :peak-rate-multiplier="(option as unknown as KeyGroupOption).peakRateMultiplier"
          :description="(option as unknown as KeyGroupOption).description"
          :selected="selected"
        />
      </template>
    </Select>
    <p v-else class="text-xs text-gray-500 dark:text-gray-400" data-test="key-group-limit">
      {{ t('keys.groupList.limitReached', { max: MAX_API_KEY_GROUPS }) }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Select from '@/components/common/Select.vue'
import GroupBadge from '@/components/common/GroupBadge.vue'
import GroupOptionItem from '@/components/common/GroupOptionItem.vue'
import Icon from '@/components/icons/Icon.vue'
import { MAX_API_KEY_GROUPS, type KeyGroupOption } from './keyGroups'

// The primary route is edited by the dialog's own group select; this editor owns
// routes 2..N. Moving route 2 up asks the parent to swap it with the primary.
const props = defineProps<{
  /** Ordered extra group IDs (routes 2..N). */
  modelValue: number[]
  /** Groups the user may bind (already filtered by permission). */
  options: KeyGroupOption[]
  /** Current primary group; excluded from the add list. */
  primaryId: number | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: number[]): void
  /** Swap extra route `index` with the primary route. */
  (e: 'promote', index: number): void
}>()

const { t } = useI18n()

const maxExtra = MAX_API_KEY_GROUPS - 1
const canPromote = computed(() => props.primaryId !== null)

// Remount the add-select after each pick so it returns to its placeholder.
const addSelectKey = ref(0)

const optionMap = computed(() => new Map(props.options.map((option) => [option.value, option])))
const optionById = (id: number) => optionMap.value.get(id)

const addableOptions = computed(() => {
  const chosen = new Set(props.modelValue)
  if (props.primaryId !== null) chosen.add(props.primaryId)
  return props.options.filter((option) => !chosen.has(option.value))
})

function add(value: string | number | boolean | null) {
  addSelectKey.value++
  const id = typeof value === 'number' ? value : Number(value)
  if (value === null || !Number.isFinite(id)) return
  if (id === props.primaryId || props.modelValue.includes(id) || props.modelValue.length >= maxExtra) return
  emit('update:modelValue', [...props.modelValue, id])
}

function remove(index: number) {
  emit('update:modelValue', props.modelValue.filter((_, i) => i !== index))
}

function swap(a: number, b: number) {
  const next = [...props.modelValue]
  ;[next[a], next[b]] = [next[b], next[a]]
  emit('update:modelValue', next)
}

function moveUp(index: number) {
  if (index === 0) {
    if (canPromote.value) emit('promote', 0)
    return
  }
  swap(index, index - 1)
}

function moveDown(index: number) {
  if (index < props.modelValue.length - 1) swap(index, index + 1)
}
</script>
