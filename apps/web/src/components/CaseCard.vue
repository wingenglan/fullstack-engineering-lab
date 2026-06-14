<template>
  <router-link :to="to" class="case-card group block glow-card p-6 cursor-pointer" :class="{ 'opacity-60 pointer-events-none': status === 'coming-soon' }">
    <!-- Header -->
    <div class="flex items-start justify-between mb-4">
      <div class="w-10 h-10 rounded-lg flex items-center justify-center" :class="status === 'available' ? 'bg-brand/10' : 'bg-white/5'">
        <i :class="iconClass" class="w-5 h-5" :style="{ color: status === 'available' ? '#6366F1' : '#71717A' }" />
      </div>
      <div class="flex items-center gap-2">
        <span
          class="px-2 py-0.5 text-xs rounded-full font-medium"
          :class="difficultyClass"
        >
          {{ difficultyLabel }}
        </span>
        <span
          class="w-2 h-2 rounded-full"
          :class="status === 'available' ? 'bg-success' : 'bg-text-muted'"
        />
      </div>
    </div>

    <!-- Content -->
    <h3 class="text-text-primary font-semibold text-lg mb-2 group-hover:text-brand transition-colors">
      {{ title }}
    </h3>
    <p class="text-text-secondary text-sm mb-4 line-clamp-2">
      {{ description }}
    </p>

    <!-- Tags -->
    <div class="flex flex-wrap gap-2 mb-4">
      <span
        v-for="tag in tags"
        :key="tag"
        class="px-2.5 py-0.5 text-xs rounded-full bg-white/5 text-text-secondary border border-border"
      >
        {{ tag }}
      </span>
    </div>

    <!-- Action -->
    <div class="flex items-center justify-between">
      <span class="text-xs text-text-muted uppercase tracking-wider">
        {{ status === 'available' ? '已上线' : '即将推出' }}
      </span>
      <i v-if="status === 'available'" class="i-lucide-arrow-right w-4 h-4 text-text-muted group-hover:text-brand group-hover:translate-x-1 transition-all" />
    </div>
  </router-link>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  title: string
  description: string
  tags: string[]
  difficulty: 'easy' | 'medium' | 'hard'
  status: 'available' | 'coming-soon'
  icon?: string
  to: string
}>()

const iconClass = computed(() => props.icon || 'i-lucide-box')

const difficultyLabel = computed(() => {
  const map = { easy: '入门', medium: '进阶', hard: '高级' }
  return map[props.difficulty]
})

const difficultyClass = computed(() => {
  const map = {
    easy: 'bg-success/10 text-success',
    medium: 'bg-warning/10 text-warning',
    hard: 'bg-danger/10 text-danger',
  }
  return map[props.difficulty]
})
</script>
