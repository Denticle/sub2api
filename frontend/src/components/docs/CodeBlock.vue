<template>
  <div class="docs-code">
    <div class="docs-code-bar">
      <span class="docs-code-lang">{{ lang }}</span>
      <button type="button" class="docs-code-copy" :aria-label="t('docs.copy')" @click="copy">
        <Icon :name="copied ? 'check' : 'copy'" size="xs" />
        <span>{{ copied ? t('docs.copied') : t('docs.copy') }}</span>
      </button>
    </div>
    <pre class="docs-code-pre"><code>{{ code }}</code></pre>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { useClipboard } from '@/composables/useClipboard'

const props = defineProps<{ code: string; lang?: string }>()

const { t } = useI18n()
const { copyToClipboard } = useClipboard()
const copied = ref(false)
let timer: ReturnType<typeof setTimeout> | undefined

async function copy() {
  const ok = await copyToClipboard(props.code)
  if (!ok) return
  copied.value = true
  clearTimeout(timer)
  timer = setTimeout(() => {
    copied.value = false
  }, 1600)
}
</script>

<style scoped>
.docs-code {
  margin: 1rem 0;
  overflow: hidden;
  border: 1px solid rgb(var(--line-light) / 0.1);
  border-radius: var(--r-md);
  background-color: rgb(var(--line-light) / 0.03);
}

.dark .docs-code {
  border-color: rgb(var(--line) / var(--line-a));
  background-color: rgb(255 255 255 / 0.03);
}

.docs-code-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  padding: 0.25rem 0.5rem 0.25rem 0.75rem;
  border-bottom: 1px solid rgb(var(--line-light) / 0.07);
}

.dark .docs-code-bar {
  border-bottom-color: rgb(var(--line) / 0.07);
}

.docs-code-lang {
  font-family: var(--font-numeric);
  font-size: 0.6875rem;
  color: rgb(var(--n-400));
}

.docs-code-copy {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.125rem 0.375rem;
  border: 0;
  border-radius: 5px;
  background: transparent;
  font-size: 0.6875rem;
  color: rgb(var(--n-500));
  cursor: pointer;
  transition: background-color 120ms ease, color 120ms ease;
}

.docs-code-copy:hover {
  background-color: rgb(var(--line-light) / 0.07);
  color: rgb(var(--n-800));
}

.dark .docs-code-copy:hover {
  background-color: rgb(255 255 255 / 0.07);
  color: rgb(var(--n-100));
}

.docs-code-copy:focus-visible {
  outline: 2px solid rgb(var(--p-500) / 0.55);
  outline-offset: 1px;
}

.docs-code-pre {
  margin: 0;
  padding: 0.75rem;
  overflow-x: auto;
  font-family: var(--font-numeric);
  font-size: 0.75rem;
  line-height: 1.7;
  color: rgb(var(--n-800));
  tab-size: 2;
}

.dark .docs-code-pre {
  color: rgb(var(--n-200));
}

.docs-code-pre code {
  font-family: inherit;
}
</style>
