<template>
  <div class="glow-card p-6" :class="$attrs.class">
    <div v-if="title || $slots.header" class="flex items-center justify-between mb-3">
      <h3 v-if="title" class="text-lg font-semibold text-text-primary">{{ title }}</h3>
      <slot v-if="$slots.header" name="header" />
    </div>
    <slot />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
}>()
</script>
