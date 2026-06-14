<template>
  <div class="glow-card p-6" :class="$attrs.class">
    <h3 v-if="title" class="text-lg font-semibold text-text-primary mb-3">{{ title }}</h3>
    <slot />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
}>()
</script>
