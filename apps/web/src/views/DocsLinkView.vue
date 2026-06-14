<template>
  <div class="docs-link-view py-12 px-6">
    <div class="max-w-3xl mx-auto text-center">
      <div class="w-16 h-16 mx-auto mb-6 rounded-2xl bg-brand/10 flex items-center justify-center">
        <i class="i-lucide-book-open w-8 h-8 text-brand" />
      </div>
      <h1 class="section-title mb-3">项目文档</h1>
      <p class="section-desc mb-10">
        全面的指南、API 参考和架构设计文档
      </p>

      <div class="grid md:grid-cols-2 gap-6 mb-10">
        <GlowCard v-for="link in docLinks" :key="link.title" class="text-left">
          <div class="flex items-start gap-4">
            <div class="w-10 h-10 shrink-0 rounded-lg bg-brand/10 flex items-center justify-center">
              <i :class="link.icon" class="w-5 h-5 text-brand" />
            </div>
            <div>
              <h3 class="text-text-primary font-semibold mb-1">{{ link.title }}</h3>
              <p class="text-text-secondary text-sm mb-3">{{ link.desc }}</p>
              <a
                :href="link.href"
                target="_blank"
                class="inline-flex items-center gap-1 text-sm text-brand hover:text-brand-light transition-colors"
              >
                阅读更多
                <i class="i-lucide-external-link w-3.5 h-3.5" />
              </a>
            </div>
          </div>
        </GlowCard>
      </div>

      <a
        :href="docsUrl"
        target="_blank"
        class="btn-primary inline-flex items-center gap-2 text-base px-8 py-3"
      >
        <i class="i-lucide-book-open w-5 h-5" />
        打开文档站
      </a>
      <p class="text-text-muted text-sm mt-4">
        文档站地址：<code class="text-brand">{{ docsUrl }}</code>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import GlowCard from '@/components/GlowCard.vue'

const docsUrl = 'http://localhost:8888/docs/'

const docLinks = [
  { title: '快速开始', desc: '使用 Docker Compose 几分钟内启动完整项目。', icon: 'i-lucide-rocket', href: '#quick-start' },
  { title: '项目架构', desc: '了解 monorepo 结构和各组件设计。', icon: 'i-lucide-layout-grid', href: '#architecture' },
  { title: '技术栈', desc: '每个技术的详细说明和版本信息。', icon: 'i-lucide-stack', href: '#tech-stack' },
  { title: 'JWT 认证案例', desc: '深入了解 JWT 认证 Demo 的实现细节。', icon: 'i-lucide-shield', href: '#jwt-auth' },
]
</script>
