<template>
  <div class="code-block rounded-lg overflow-hidden border border-border">
    <div class="flex items-center justify-between px-4 py-2 bg-bg-elevated border-b border-border">
      <span class="text-xs text-text-muted font-mono">{{ language || 'text' }}</span>
      <button
        class="flex items-center gap-1.5 text-xs text-text-muted hover:text-text-primary transition-colors"
        @click="handleCopy"
      >
        <i :class="copied ? 'i-lucide-check' : 'i-lucide-copy'" class="w-3.5 h-3.5" />
        {{ copied ? 'Copied' : 'Copy' }}
      </button>
    </div>
    <pre class="p-4 overflow-x-auto text-sm leading-6 bg-[#0D0D0F]"><code class="text-text-secondary font-mono">{{ code }}</code></pre>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  code: string
  language?: string
}>()

const copied = ref(false)

function handleCopy() {
  navigator.clipboard.writeText(props.code)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}
</script>
