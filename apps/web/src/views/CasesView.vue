<template>
  <div class="cases-view py-12 px-6">
    <div class="max-w-7xl mx-auto">
      <h1 class="section-title mb-3">案例列表</h1>
      <p class="section-desc mb-10">探索实战工程案例</p>

      <!-- Filters -->
      <div class="flex flex-wrap gap-3 mb-8">
        <el-button
          v-for="cat in allCategories"
          :key="cat"
          :type="activeCategory === cat ? 'primary' : 'default'"
          size="default"
          round
          @click="activeCategory = cat"
        >
          {{ cat }}
        </el-button>
      </div>

      <!-- Grid -->
      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <CaseCard
          v-for="item in filteredCases"
          :key="item.id"
          v-bind="item"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import CaseCard from '@/components/CaseCard.vue'

const activeCategory = ref('全部')

const allCases = [
  { id: 'jwt', title: 'JWT 认证授权', description: '基于 JWT 的完整认证流程，包含注册、登录、Token 管理和 Redis 黑名单。', tags: ['Go', 'JWT', 'Redis'], difficulty: 'easy' as const, status: 'available' as const, category: '认证授权', icon: 'i-lucide-shield', to: '/cases/jwt-auth' },
  { id: 'ws', title: 'WebSocket 实时通信', description: '基于 WebSocket 协议的实时双向通信聊天室，支持多房间、消息持久化、在线状态。', tags: ['WebSocket', 'Go', 'Vue'], difficulty: 'medium' as const, status: 'available' as const, category: '实时通信', icon: 'i-lucide-message-circle', to: '/cases/websocket' },
  { id: 'redis-lock', title: 'Redis 分布式锁', description: '基于 Redis SET NX EX 实现的分布式锁，演示互斥、自动过期和并发争抢。', tags: ['Redis', 'Go', 'Lua'], difficulty: 'medium' as const, status: 'available' as const, category: '缓存', icon: 'i-lucide-lock', to: '/cases/redis-lock' },
  { id: 'redis-data', title: 'Redis 数据类型实战', description: '覆盖五大核心数据结构实战：String 验证码/计数器、Hash 画像、List 活动流、Set 标签、ZSet 排行榜。', tags: ['Redis', 'Go', 'String', 'Hash', 'List'], difficulty: 'medium' as const, status: 'available' as const, category: '缓存', icon: 'i-lucide-database', to: '/cases/redis-data' },
  { id: 'mqtt', title: 'MQTT IoT 演示', description: '基于 MQTT 协议的 IoT 传感器模拟，通过 Mosquitto Broker 发布数据，SSE 实时推送到前端仪表盘。', tags: ['MQTT', 'Go', 'Vue', 'SSE', 'IoT'], difficulty: 'medium' as const, status: 'available' as const, category: 'IoT通信', icon: 'i-lucide-radio', to: '/cases/mqtt' },
  { id: 'tcp', title: 'TCP 自定义协议', description: 'Go 实现的 TCP Server，演示自定义文本协议（PING/ECHO/TIME 等），前端通过 HTTP 管理会话。', tags: ['TCP', 'Go', 'Vue', 'SSE'], difficulty: 'hard' as const, status: 'available' as const, category: '网络协议', icon: 'i-lucide-terminal', to: '/cases/tcp' },
  { id: 'chunk-upload', title: '大文件分片上传', description: '大文件分片上传，支持断点续传。', tags: ['文件', 'Go', 'Vue'], difficulty: 'medium' as const, status: 'coming-soon' as const, category: '文件存储', icon: 'i-lucide-upload', to: '#' },
  { id: 'scheduler', title: '定时任务', description: '基于 Cron 的定时任务调度系统。', tags: ['Go', 'Cron'], difficulty: 'medium' as const, status: 'coming-soon' as const, category: '定时任务', icon: 'i-lucide-clock', to: '#' },
  { id: 'payment', title: '支付对接', description: '支付网关集成与订单管理。', tags: ['支付', 'Go'], difficulty: 'hard' as const, status: 'coming-soon' as const, category: '支付对接', icon: 'i-lucide-credit-card', to: '#' },
]

const allCategories = computed(() => {
  const cats = new Set(allCases.map((c) => c.category))
  return ['全部', ...cats]
})

const filteredCases = computed(() => {
  if (activeCategory.value === '全部') return allCases
  return allCases.filter((c) => c.category === activeCategory.value)
})
</script>
