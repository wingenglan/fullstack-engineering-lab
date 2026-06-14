<template>
  <div class="roadmap-view py-12 px-6">
    <div class="max-w-4xl mx-auto">
      <h1 class="section-title text-center mb-3">开发路线图</h1>
      <p class="section-desc text-center mb-16">从基础工程到生产级实践的演进之路</p>

      <!-- Timeline -->
      <div class="relative">
        <div class="absolute left-6 md:left-1/2 top-0 bottom-0 w-px bg-gradient-to-b from-brand via-accent to-transparent" />

        <div
          v-for="(phase, i) in phases"
          :key="i"
          class="relative mb-12"
          :class="i % 2 === 0 ? 'md:pr-1/2' : 'md:pl-1/2 md:ml-auto'"
        >
          <div
            class="relative ml-16 md:ml-0 glow-card p-6"
            :class="phase.current ? 'ring-1 ring-brand/30' : ''"
          >
            <!-- Dot -->
            <div
              class="absolute -left-[2.75rem] md:-left-4 top-6 w-6 h-6 rounded-full flex items-center justify-center z-10"
              :class="phase.current ? 'bg-brand shadow-[0_0_12px_rgba(99,102,241,0.5)]' : phase.done ? 'bg-success' : 'bg-bg-elevated border border-border'"
            >
              <i v-if="phase.done" class="i-lucide-check w-3 h-3 text-white" />
              <i v-else-if="phase.current" class="i-lucide-loader w-3 h-3 text-white animate-spin" />
            </div>

            <div class="flex items-center gap-3 mb-3">
              <h3 class="text-lg font-semibold text-text-primary">{{ phase.title }}</h3>
              <span
                class="px-2 py-0.5 text-xs rounded-full font-medium"
                :class="phase.current ? 'bg-brand/10 text-brand' : phase.done ? 'bg-success/10 text-success' : 'bg-white/5 text-text-muted'"
              >
                {{ phase.current ? '进行中' : phase.done ? '已完成' : '计划中' }}
              </span>
            </div>
            <p class="text-text-secondary text-sm mb-4">{{ phase.desc }}</p>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="item in phase.items"
                :key="item"
                class="px-3 py-1 text-xs rounded-full border border-border text-text-secondary"
              >
                {{ item }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const phases = [
  {
    title: '阶段一：认证与基础工程',
    desc: 'JWT 认证授权、RBAC 权限控制、Docker 部署、CI/CD 流水线、项目脚手架。',
    items: ['JWT 认证案例', 'Docker Compose', 'CI/CD', 'Monorepo'],
    current: true,
    done: false,
  },
  {
    title: '阶段二：实时通信与缓存',
    desc: 'WebSocket 通信、Redis 分布式锁、实时数据同步。',
    items: ['WebSocket 聊天', 'Redis 分布式锁', '发布/订阅', '在线状态'],
    current: false,
    done: false,
  },
  {
    title: '阶段三：消息队列与调度',
    desc: '消息队列集成、定时任务调度、异步任务处理。',
    items: ['RabbitMQ', 'Cron 定时任务', '任务队列', '死信队列'],
    current: false,
    done: false,
  },
  {
    title: '阶段四：文件存储与支付',
    desc: '文件上传、大文件分片、支付网关集成。',
    items: ['分片上传', 'MinIO 存储', '支付网关', '订单系统'],
    current: false,
    done: false,
  },
  {
    title: '阶段五：搜索与 AI',
    desc: '全文搜索集成、AI 功能、向量搜索。',
    items: ['Elasticsearch', 'AI 对话', '向量数据库', 'RAG'],
    current: false,
    done: false,
  },
  {
    title: '阶段六：运维与可观测性',
    desc: '生产部署、监控告警、日志收集、链路追踪。',
    items: ['K8s 部署', 'Prometheus', 'Grafana', 'Jaeger'],
    current: false,
    done: false,
  },
]
</script>
