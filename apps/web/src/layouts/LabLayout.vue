<template>
  <div class="lab-layout min-h-screen flex flex-col">
    <AppHeader />
    <main class="flex-1">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
    <footer class="border-t border-border py-8 px-6">
      <div class="max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between gap-4">
        <div class="text-text-muted text-sm">
          &copy; 2026 工程实验室 · MIT License
        </div>
        <div class="flex items-center gap-6 text-sm text-text-muted">
          <a href="https://github.com" target="_blank" class="hover:text-text-primary transition-colors">GitHub</a>
          <router-link to="/docs-link" class="hover:text-text-primary transition-colors">文档</router-link>
          <router-link to="/roadmap" class="hover:text-text-primary transition-colors">路线图</router-link>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.lab-layout {
  background:
    radial-gradient(ellipse 80% 50% at 50% -20%, rgba(99, 102, 241, 0.08), transparent),
    var(--color-bg);
}
</style>
