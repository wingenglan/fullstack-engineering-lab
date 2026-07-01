<template>
  <header class="app-header sticky top-0 z-50 backdrop-blur-xl border-b border-border">
    <div class="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
      <!-- Logo -->
      <router-link to="/" class="flex items-center gap-3 group">
        <div class="w-9 h-9 rounded-lg bg-gradient-to-br from-brand to-accent flex items-center justify-center">
          <i class="i-lucide-flask-conical w-5 h-5 text-white" />
        </div>
        <span class="text-text-primary font-semibold text-lg hidden sm:block group-hover:text-brand transition-colors">
          工程实验室
        </span>
      </router-link>

      <!-- Navigation -->
      <nav class="hidden md:flex items-center gap-1">
        <router-link
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="nav-link px-4 py-2 rounded-lg text-sm text-text-secondary hover:text-text-primary hover:bg-white/5 transition-all"
          active-class="!text-brand !bg-brand/10"
        >
          {{ item.label }}
        </router-link>
        <a
          href="https://github.com"
          target="_blank"
          class="nav-link px-3 py-2 rounded-lg text-text-secondary hover:text-text-primary hover:bg-white/5 transition-all ml-1"
        >
          <i class="i-lucide-github w-5 h-5" />
        </a>
      </nav>

      <!-- Auth -->
      <div class="flex items-center gap-3">
        <template v-if="authStore.isAuthenticated">
          <span class="text-sm text-text-secondary hidden sm:block">
            {{ authStore.displayName }}
          </span>
          <button
            class="px-3 py-1.5 text-sm text-text-secondary border border-border rounded-lg hover:text-text-primary hover:border-brand/50 transition-all"
            @click="handleLogout"
          >
            退出登录
          </button>
        </template>
        <template v-else>
          <router-link
            to="/cases/jwt-auth"
            class="px-4 py-1.5 text-sm bg-brand text-white rounded-lg hover:bg-brand-dark transition-all"
          >
            开始体验
          </router-link>
        </template>

        <!-- Mobile menu -->
        <button class="md:hidden p-2 text-text-secondary" @click="mobileOpen = !mobileOpen">
          <i :class="mobileOpen ? 'i-lucide-x' : 'i-lucide-menu'" class="w-5 h-5" />
        </button>
      </div>
    </div>

    <!-- Mobile nav -->
    <transition name="slide">
      <div v-if="mobileOpen" class="md:hidden border-t border-border bg-bg-card">
        <nav class="px-6 py-4 flex flex-col gap-2">
          <router-link
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="px-4 py-2.5 rounded-lg text-text-secondary hover:text-text-primary hover:bg-white/5 transition-all"
            active-class="!text-brand !bg-brand/10"
            @click="mobileOpen = false"
          >
            {{ item.label }}
          </router-link>
        </nav>
      </div>
    </transition>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const mobileOpen = ref(false)

const navItems = [
  { label: '案例', to: '/cases' },
  { label: '路线图', to: '/roadmap' },
  { label: '文档', to: '/docs-link' },
]

async function handleLogout() {
  await authStore.logout()
}
</script>

<style scoped>
.app-header {
  background: rgba(9, 9, 11, 0.8);
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.2s ease;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
