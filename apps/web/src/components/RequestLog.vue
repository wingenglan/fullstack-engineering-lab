<template>
  <div class="request-log max-h-80 overflow-y-auto bg-[#0D0D0F]">
    <div v-if="logs.length === 0" class="p-6 text-center text-text-muted text-sm">
      暂无请求记录，请尝试注册或登录。
    </div>
    <div
      v-for="log in logs"
      :key="log.id"
      class="px-4 py-2 border-b border-border/50 hover:bg-white/[0.02] transition-colors cursor-pointer"
      @click="toggleExpand(log.id)"
    >
      <div class="flex items-center gap-3 text-sm font-mono">
        <span class="text-text-muted text-xs w-16 shrink-0">
          {{ formatTime(log.timestamp) }}
        </span>
        <span
          class="px-1.5 py-0.5 rounded text-xs font-bold"
          :class="methodClass(log.method)"
        >
          {{ log.method }}
        </span>
        <span class="text-text-secondary flex-1 truncate">{{ log.url }}</span>
        <span
          v-if="log.status"
          class="text-xs font-mono"
          :class="statusClass(log.status)"
        >
          {{ log.status }}
        </span>
        <span v-if="log.duration !== null" class="text-xs text-text-muted">
          {{ log.duration }}ms
        </span>
        <span v-if="!log.status" class="text-xs text-text-muted animate-pulse">请求中...</span>
      </div>

      <!-- Expanded details -->
      <transition name="expand">
        <div v-if="expanded.has(log.id)" class="mt-2 pt-2 border-t border-border/50">
          <div v-if="log.requestData" class="mb-2">
            <span class="text-xs text-text-muted block mb-1">请求体：</span>
            <pre class="text-xs text-text-secondary font-mono bg-bg-card rounded p-2 overflow-x-auto">{{ formatJson(log.requestData) }}</pre>
          </div>
          <div v-if="log.responseData">
            <span class="text-xs text-text-muted block mb-1">响应体：</span>
            <pre class="text-xs text-text-secondary font-mono bg-bg-card rounded p-2 overflow-x-auto">{{ formatJson(log.responseData) }}</pre>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import type { RequestLogEntry } from '@/types'

const props = defineProps<{
  logs: RequestLogEntry[]
}>()

const expanded = ref(new Set<number>())

function toggleExpand(id: number) {
  if (expanded.value.has(id)) {
    expanded.value.delete(id)
  } else {
    expanded.value.add(id)
  }
}

function formatTime(timestamp: string): string {
  const d = new Date(timestamp)
  return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}:${d.getSeconds().toString().padStart(2, '0')}`
}

function methodClass(method: string): string {
  const map: Record<string, string> = {
    GET: 'bg-success/20 text-success',
    POST: 'bg-brand/20 text-brand',
    PUT: 'bg-warning/20 text-warning',
    DELETE: 'bg-danger/20 text-danger',
  }
  return map[method] || 'bg-white/10 text-text-secondary'
}

function statusClass(status: number): string {
  if (status >= 200 && status < 300) return 'text-success'
  if (status >= 400 && status < 500) return 'text-warning'
  return 'text-danger'
}

function formatJson(str: string): string {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

watch(() => props.logs.length, async () => {
  await nextTick()
  // 自动滚动由父组件处理
})
</script>

<style scoped>
.expand-enter-active,
.expand-leave-active {
  transition: all 0.2s ease;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
}
</style>
