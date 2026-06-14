<template>
  <div class="home-view">
    <!-- Hero Section -->
    <section class="relative py-24 px-6 overflow-hidden">
      <div class="absolute inset-0 bg-gradient-to-b from-brand/5 to-transparent" />
      <div class="max-w-4xl mx-auto text-center relative z-10">
        <div class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full border border-brand/20 bg-brand/5 text-brand text-sm mb-8">
          <i class="i-lucide-sparkles w-4 h-4" />
          开源工程实战实验室
        </div>
        <h1 class="text-4xl md:text-6xl font-bold mb-6 leading-tight">
          <span class="text-gradient">构建、学习、运行</span>
          <br />
          <span class="text-text-primary">真实工程案例</span>
        </h1>
        <p class="text-text-secondary text-lg md:text-xl mb-10 max-w-2xl mx-auto">
          可运行、可体验、可学习的全栈工程实践案例库
        </p>
        <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
          <router-link to="/cases" class="btn-primary text-base px-8 py-3 flex items-center gap-2">
            <i class="i-lucide-rocket w-5 h-5" />
            开始体验
          </router-link>
          <router-link to="/docs-link" class="btn-ghost text-base px-8 py-3 flex items-center gap-2">
            <i class="i-lucide-book-open w-5 h-5" />
            查看文档
          </router-link>
        </div>
      </div>
    </section>

    <!-- Search -->
    <section class="px-6 pb-16">
      <div class="max-w-2xl mx-auto">
        <el-input
          v-model="searchQuery"
          size="large"
          placeholder="搜索案例、技术、关键词..."
          class="search-input"
        >
          <template #prefix>
            <i class="i-lucide-search w-5 h-5 text-text-muted" />
          </template>
        </el-input>
      </div>
    </section>

    <!-- Case Categories -->
    <section class="px-6 pb-20">
      <div class="max-w-7xl mx-auto">
        <h2 class="section-title text-center mb-3">技术分类</h2>
        <p class="section-desc text-center mb-12">按领域探索工程案例</p>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <GlowCard
            v-for="cat in categories"
            :key="cat.name"
            class="text-center cursor-pointer"
          >
            <div class="w-12 h-12 mx-auto mb-3 rounded-xl flex items-center justify-center" :style="{ background: cat.color + '15' }">
              <i :class="cat.icon" class="w-6 h-6" :style="{ color: cat.color }" />
            </div>
            <h3 class="text-text-primary font-medium mb-1">{{ cat.name }}</h3>
            <p class="text-text-muted text-xs">{{ cat.count }} 个案例</p>
          </GlowCard>
        </div>
      </div>
    </section>

    <!-- Featured Cases -->
    <section class="px-6 pb-20">
      <div class="max-w-7xl mx-auto">
        <h2 class="section-title text-center mb-3">热门案例</h2>
        <p class="section-desc text-center mb-12">从这些实战案例开始学习</p>
        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
          <CaseCard
            v-for="item in featuredCases"
            :key="item.id"
            v-bind="item"
          />
        </div>
      </div>
    </section>

    <!-- Learning Roadmap Preview -->
    <section class="px-6 pb-20">
      <div class="max-w-7xl mx-auto">
        <h2 class="section-title text-center mb-3">学习路线</h2>
        <p class="section-desc text-center mb-12">从基础到生产的结构化学习路径</p>
        <div class="grid md:grid-cols-4 gap-6">
          <GlowCard v-for="(phase, i) in roadmapPhases" :key="i" class="text-center">
            <div class="w-10 h-10 mx-auto mb-4 rounded-full flex items-center justify-center font-bold" :class="phase.done ? 'bg-success text-white' : phase.current ? 'bg-brand text-white' : 'bg-brand/10 text-brand'">
              <i v-if="phase.done" class="i-lucide-check w-5 h-5" />
              <span v-else>{{ i + 1 }}</span>
            </div>
            <h3 class="text-text-primary font-semibold mb-2">{{ phase.title }}</h3>
            <p class="text-text-secondary text-sm">{{ phase.desc }}</p>
            <div v-if="phase.current" class="mt-3 px-3 py-1 inline-block text-xs rounded-full bg-brand/10 text-brand">
              进行中
            </div>
            <div v-if="phase.done" class="mt-3 px-3 py-1 inline-block text-xs rounded-full bg-success/10 text-success">
              已完成
            </div>
          </GlowCard>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section class="px-6 pb-20">
      <div class="max-w-7xl mx-auto">
        <div class="grid md:grid-cols-3 gap-6">
          <GlowCard v-for="feature in features" :key="feature.title" class="text-center">
            <div class="w-12 h-12 mx-auto mb-4 rounded-xl bg-brand/10 flex items-center justify-center">
              <i :class="feature.icon" class="w-6 h-6 text-brand" />
            </div>
            <h3 class="text-text-primary font-semibold mb-2">{{ feature.title }}</h3>
            <p class="text-text-secondary text-sm">{{ feature.desc }}</p>
          </GlowCard>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import GlowCard from '@/components/GlowCard.vue'
import CaseCard from '@/components/CaseCard.vue'

const searchQuery = ref('')

const categories = [
  { name: '认证授权', icon: 'i-lucide-shield', color: '#6366F1', count: 1 },
  { name: '实时通信', icon: 'i-lucide-zap', color: '#F59E0B', count: 0 },
  { name: '缓存', icon: 'i-lucide-database', color: '#EF4444', count: 1 },
  { name: '消息队列', icon: 'i-lucide-list-ordered', color: '#8B5CF6', count: 0 },
  { name: '文件存储', icon: 'i-lucide-hard-drive', color: '#10B981', count: 0 },
  { name: '定时任务', icon: 'i-lucide-clock', color: '#F97316', count: 0 },
  { name: '支付对接', icon: 'i-lucide-credit-card', color: '#EC4899', count: 0 },
  { name: '运维部署', icon: 'i-lucide-container', color: '#06B6D4', count: 0 },
]

const featuredCases = [
  { id: 'jwt', title: 'JWT 认证授权', description: '基于 JWT 的完整认证流程，包含注册、登录、Token 管理和 Redis 黑名单。', tags: ['Go', 'JWT', 'Redis'], difficulty: 'easy' as const, status: 'available' as const, icon: 'i-lucide-shield', to: '/cases/jwt-auth' },
  { id: 'redis-lock', title: 'Redis 分布式锁', description: '基于 Redis SET NX EX 实现的分布式锁，演示互斥、自动过期和并发争抢。', tags: ['Redis', 'Go', 'Lua'], difficulty: 'medium' as const, status: 'available' as const, icon: 'i-lucide-lock', to: '/cases/redis-lock' },
  { id: 'ws', title: 'WebSocket 实时通信', description: '基于 WebSocket 的实时双向通信，支持群聊和在线状态。', tags: ['WebSocket', 'Go', 'Vue'], difficulty: 'medium' as const, status: 'coming-soon' as const, icon: 'i-lucide-message-circle', to: '#' },
  { id: 'chunk', title: '大文件分片上传', description: '大文件分片上传，支持断点续传和上传进度追踪。', tags: ['文件', 'Go', 'Vue'], difficulty: 'medium' as const, status: 'coming-soon' as const, icon: 'i-lucide-upload', to: '#' },
]

const roadmapPhases = [
  { title: '认证与基础', desc: 'JWT、RBAC、基础工程化搭建', done: true, current: false },
  { title: '实时与缓存', desc: 'WebSocket、Redis 分布式锁', done: false, current: true },
  { title: '队列与调度', desc: '消息队列、定时任务、异步处理', done: false, current: false },
  { title: '文件与支付', desc: '文件上传、分片上传、支付对接', done: false, current: false },
]

const features = [
  { title: '一键启动', icon: 'i-lucide-rocket', desc: 'Docker Compose 一条命令启动完整技术栈，包含数据库、缓存和反向代理。' },
  { title: '真实案例', icon: 'i-lucide-code-2', desc: '每个案例都是真实的工程实践，不是玩具示例，包含生产级最佳实践。' },
  { title: '易于扩展', icon: 'i-lucide-puzzle', desc: '模块化 monorepo 架构，新增案例只需添加目录，无需大规模迁移。' },
]
</script>

<style scoped>
.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  padding: 8px 16px;
  background: #18181B;
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.08) inset;
}
</style>
