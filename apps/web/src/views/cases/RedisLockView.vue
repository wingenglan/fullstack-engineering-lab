<template>
  <div class="redis-lock-view py-8 px-6">
    <div class="max-w-7xl mx-auto">
      <!-- Breadcrumb -->
      <div class="flex items-center gap-2 text-sm text-text-muted mb-6">
        <router-link to="/cases" class="hover:text-text-primary transition-colors">案例列表</router-link>
        <i class="i-lucide-chevron-right w-4 h-4" />
        <span class="text-text-primary">Redis 分布式锁</span>
      </div>

      <div class="grid lg:grid-cols-3 gap-6">
        <!-- Left Column -->
        <div class="lg:col-span-1 space-y-6">
          <!-- Case Info -->
          <GlowCard>
            <div class="flex items-center gap-3 mb-4">
              <div class="w-10 h-10 rounded-lg bg-brand/10 flex items-center justify-center">
                <i class="i-lucide-lock w-5 h-5 text-brand" />
              </div>
              <div>
                <h2 class="text-text-primary font-bold text-lg">Redis 分布式锁</h2>
                <span class="text-xs text-success">已上线</span>
              </div>
            </div>
            <p class="text-text-secondary text-sm mb-4">
              基于 Redis SET NX EX 和 Lua 脚本实现的分布式锁，演示互斥访问、自动过期和并发争抢场景。
            </p>
            <div class="flex flex-wrap gap-2 mb-4">
              <span v-for="tag in ['Redis', 'Go', 'Lua', 'SET NX', '分布式']" :key="tag" class="px-2.5 py-0.5 text-xs rounded-full bg-white/5 text-text-secondary border border-border">
                {{ tag }}
              </span>
            </div>
            <div class="text-xs text-text-muted">
              难度：<span class="text-warning font-medium">进阶</span>
            </div>
          </GlowCard>

          <!-- Steps Navigation -->
          <GlowCard title="操作步骤">
            <div class="space-y-2">
              <div
                v-for="(step, i) in steps"
                :key="i"
                class="flex items-center gap-3 px-3 py-2.5 rounded-lg cursor-pointer transition-all"
                :class="activeStep === i ? 'bg-brand/10 text-brand' : 'text-text-secondary hover:bg-white/5'"
                @click="activeStep = i"
              >
                <div
                  class="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold shrink-0"
                  :class="step.completed ? 'bg-success text-white' : activeStep === i ? 'bg-brand text-white' : 'bg-white/10 text-text-muted'"
                >
                  <i v-if="step.completed" class="i-lucide-check w-3 h-3" />
                  <span v-else>{{ i + 1 }}</span>
                </div>
                <span class="text-sm font-medium">{{ step.label }}</span>
              </div>
            </div>
          </GlowCard>

          <!-- Lock Concepts -->
          <GlowCard title="核心原理">
            <div class="space-y-3 text-sm">
              <div v-for="(item, i) in concepts" :key="i" class="flex items-start gap-3">
                <div class="w-7 h-7 rounded-full bg-brand/10 flex items-center justify-center shrink-0 mt-0.5">
                  <i :class="item.icon" class="w-3.5 h-3.5 text-brand" />
                </div>
                <div>
                  <span class="text-text-primary font-medium block">{{ item.title }}</span>
                  <span class="text-text-muted text-xs">{{ item.desc }}</span>
                </div>
              </div>
            </div>
          </GlowCard>

          <!-- Redis Key Viewer -->
          <GlowCard title="Redis Key 监控">
            <div class="space-y-2">
              <div
                v-for="key in redisKeys"
                :key="key.name"
                class="flex items-center justify-between px-3 py-2 rounded-lg bg-bg-elevated border border-border"
              >
                <div class="flex items-center gap-2">
                  <span class="w-2 h-2 rounded-full" :class="key.locked ? 'bg-danger' : 'bg-text-muted'" />
                  <span class="text-xs font-mono text-text-secondary">{{ key.name }}</span>
                </div>
                <span class="text-xs font-mono" :class="key.locked ? 'text-danger' : 'text-text-muted'">
                  {{ key.locked ? `TTL: ${key.ttl}s` : '无锁' }}
                </span>
              </div>
              <div v-if="redisKeys.length === 0" class="text-center py-4 text-text-muted text-xs">
                暂无锁记录
              </div>
            </div>
          </GlowCard>
        </div>

        <!-- Right Column -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Lock Status Bar -->
          <div class="flex items-center justify-between p-4 rounded-xl border border-border bg-bg-card">
            <div class="flex items-center gap-3">
              <div
                class="w-3 h-3 rounded-full"
                :class="currentLock.locked ? 'bg-danger shadow-[0_0_8px_rgba(239,68,68,0.5)]' : 'bg-text-muted'"
              />
              <span class="text-sm text-text-primary">
                {{ currentLock.locked ? `资源 "${currentLock.resource}" 已锁定 (TTL: ${Math.ceil(currentLock.ttl_ms / 1000)}s)` : '无活跃锁' }}
              </span>
            </div>
            <el-button
              v-if="currentLock.locked"
              size="small"
              type="danger"
              plain
              @click="handleRelease(currentLock.resource)"
            >
              强制释放
            </el-button>
          </div>

          <!-- Tabs -->
          <div class="flex gap-1 p-1 bg-bg-elevated rounded-xl border border-border">
            <button
              v-for="tab in tabs"
              :key="tab.key"
              class="flex-1 px-4 py-2.5 rounded-lg text-sm font-medium transition-all"
              :class="activeTab === tab.key ? 'bg-brand text-white shadow-lg' : 'text-text-secondary hover:text-text-primary'"
              @click="activeTab = tab.key"
            >
              <i :class="tab.icon" class="w-4 h-4 mr-1.5 inline-block" />
              {{ tab.label }}
            </button>
          </div>

          <!-- Acquire / Release Tab -->
          <GlowCard v-if="activeTab === 'lock'">
            <h3 class="text-text-primary font-semibold text-lg mb-6">获取 / 释放锁</h3>

            <div class="space-y-5">
              <!-- Resource Name -->
              <div>
                <label class="text-sm text-text-secondary mb-2 block">资源名称</label>
                <el-input v-model="lockForm.resource" placeholder="例如: order:12345" />
                <p class="text-xs text-text-muted mt-1">资源名将作为 Redis key 的一部分，不同资源的锁互不影响</p>
              </div>

              <!-- TTL Slider -->
              <div>
                <div class="flex items-center justify-between mb-2">
                  <label class="text-sm text-text-secondary">锁过期时间 (TTL)</label>
                  <span class="text-sm font-mono text-brand">{{ lockForm.ttl }}s</span>
                </div>
                <el-slider v-model="lockForm.ttl" :min="1" :max="60" :step="1" />
                <p class="text-xs text-text-muted mt-1">TTL 保证即使持有者崩溃，锁也会自动释放，防止死锁</p>
              </div>

              <!-- Action Buttons -->
              <div class="flex gap-3">
                <el-button
                  type="primary"
                  :loading="acquiring"
                  class="flex-1"
                  @click="handleAcquire"
                >
                  <i class="i-lucide-lock w-4 h-4 mr-2" />
                  获取锁
                </el-button>
                <el-button
                  type="danger"
                  plain
                  :disabled="!currentLock.locked"
                  class="flex-1"
                  @click="handleRelease(lockForm.resource)"
                >
                  <i class="i-lucide-unlock w-4 h-4 mr-2" />
                  释放锁
                </el-button>
              </div>

              <!-- Result -->
              <div v-if="lockResult" class="p-4 rounded-lg border" :class="lockResult.success ? 'bg-success/5 border-success/20' : 'bg-danger/5 border-danger/20'">
                <div class="flex items-center gap-2 mb-2">
                  <i :class="lockResult.success ? 'i-lucide-check-circle text-success' : 'i-lucide-x-circle text-danger'" class="w-4 h-4" />
                  <span class="text-sm font-medium" :class="lockResult.success ? 'text-success' : 'text-danger'">
                    {{ lockResult.success ? '操作成功' : '操作失败' }}
                  </span>
                </div>
                <pre class="text-xs font-mono text-text-secondary bg-bg-elevated rounded p-3 overflow-x-auto">{{ lockResult.message }}</pre>
              </div>
            </div>
          </GlowCard>

          <!-- Status Query Tab -->
          <GlowCard v-if="activeTab === 'status'">
            <h3 class="text-text-primary font-semibold text-lg mb-6">锁状态查询</h3>

            <div class="space-y-5">
              <div>
                <label class="text-sm text-text-secondary mb-2 block">资源名称</label>
                <div class="flex gap-3">
                  <el-input v-model="statusResource" placeholder="输入资源名称查询锁状态" class="flex-1" />
                  <el-button type="primary" :loading="querying" @click="handleQueryStatus">
                    <i class="i-lucide-search w-4 h-4 mr-1" />
                    查询
                  </el-button>
                </div>
              </div>

              <!-- Quick Query Buttons -->
              <div class="flex flex-wrap gap-2">
                <el-button
                  v-for="res in quickResources"
                  :key="res"
                  size="small"
                  plain
                  @click="statusResource = res; handleQueryStatus()"
                >
                  {{ res }}
                </el-button>
              </div>

              <!-- Status Result -->
              <div v-if="statusResult" class="grid grid-cols-2 gap-4">
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <span class="text-xs text-text-muted block mb-1">资源</span>
                  <span class="text-text-primary font-mono text-sm">{{ statusResult.resource }}</span>
                </div>
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <span class="text-xs text-text-muted block mb-1">状态</span>
                  <span class="text-sm flex items-center gap-1.5" :class="statusResult.locked ? 'text-danger' : 'text-success'">
                    <span class="w-2 h-2 rounded-full" :class="statusResult.locked ? 'bg-danger' : 'bg-success'" />
                    {{ statusResult.locked ? '已锁定' : '未锁定' }}
                  </span>
                </div>
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <span class="text-xs text-text-muted block mb-1">持有者</span>
                  <span class="text-text-primary font-mono text-xs break-all">
                    {{ statusResult.owner || '-' }}
                  </span>
                </div>
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <span class="text-xs text-text-muted block mb-1">剩余 TTL</span>
                  <span class="text-text-primary font-mono">
                    {{ statusResult.locked ? `${Math.ceil(statusResult.ttl_ms / 1000)}s` : '-' }}
                  </span>
                </div>
              </div>
            </div>
          </GlowCard>

          <!-- Contention Demo Tab -->
          <GlowCard v-if="activeTab === 'contention'">
            <h3 class="text-text-primary font-semibold text-lg mb-6">并发争抢演示</h3>
            <p class="text-text-secondary text-sm mb-6">
              模拟多个协程同时争抢同一把分布式锁，演示 Redis 锁的互斥特性。同一时刻只有一个协程能获取到锁。
            </p>

            <div class="space-y-5">
              <!-- Config -->
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="text-sm text-text-secondary mb-2 block">资源名称</label>
                  <el-input v-model="contentionForm.resource" placeholder="demo-resource" />
                </div>
                <div>
                  <label class="text-sm text-text-secondary mb-2 block">锁 TTL (秒)</label>
                  <el-input-number v-model="contentionForm.ttl" :min="5" :max="60" class="w-full" />
                </div>
                <div>
                  <label class="text-sm text-text-secondary mb-2 block">并发协程数</label>
                  <el-slider v-model="contentionForm.goroutines" :min="2" :max="10" :step="1" show-stops />
                  <span class="text-xs text-text-muted">{{ contentionForm.goroutines }} 个协程</span>
                </div>
                <div>
                  <label class="text-sm text-text-secondary mb-2 block">持有时间 (ms)</label>
                  <el-input-number v-model="contentionForm.hold_ms" :min="100" :max="5000" :step="100" class="w-full" />
                </div>
              </div>

              <!-- Run Button -->
              <el-button
                type="primary"
                :loading="contending"
                class="w-full"
                @click="handleContention"
              >
                <i class="i-lucide-zap w-4 h-4 mr-2" />
                启动并发争抢
              </el-button>

              <!-- Contention Results -->
              <div v-if="contentionResult" class="space-y-4">
                <!-- Summary -->
                <div class="grid grid-cols-3 gap-3">
                  <div class="p-4 rounded-lg bg-bg-elevated border border-border text-center">
                    <span class="text-2xl font-bold text-text-primary block">{{ contentionResult.summary.total }}</span>
                    <span class="text-xs text-text-muted">总协程</span>
                  </div>
                  <div class="p-4 rounded-lg bg-success/5 border border-success/20 text-center">
                    <span class="text-2xl font-bold text-success block">{{ contentionResult.summary.succeeded }}</span>
                    <span class="text-xs text-success/70">获取成功</span>
                  </div>
                  <div class="p-4 rounded-lg bg-danger/5 border border-danger/20 text-center">
                    <span class="text-2xl font-bold text-danger block">{{ contentionResult.summary.failed }}</span>
                    <span class="text-xs text-danger/70">获取失败</span>
                  </div>
                </div>

                <!-- Timeline -->
                <div class="p-4 rounded-lg bg-[#0D0D0F] border border-border">
                  <h4 class="text-sm font-medium text-text-secondary mb-3">执行详情</h4>
                  <div class="space-y-2">
                    <div
                      v-for="r in contentionResult.results"
                      :key="r.goroutine_id"
                      class="flex items-center gap-3 text-sm font-mono"
                    >
                      <span class="text-text-muted text-xs w-8 shrink-0">#{{ r.goroutine_id }}</span>
                      <span class="w-4 h-4 rounded-full flex items-center justify-center shrink-0" :class="r.acquired ? 'bg-success' : 'bg-danger'">
                        <i :class="r.acquired ? 'i-lucide-check w-2.5 h-2.5 text-white' : 'i-lucide-x w-2.5 h-2.5 text-white'" />
                      </span>
                      <span class="text-text-secondary flex-1">{{ r.message }}</span>
                      <span class="text-xs text-text-muted">{{ r.wait_ms }}ms</span>
                    </div>
                  </div>
                </div>

                <!-- Explanation -->
                <div class="p-4 rounded-lg bg-brand/5 border border-brand/20">
                  <div class="flex items-start gap-2">
                    <i class="i-lucide-info w-4 h-4 text-brand mt-0.5 shrink-0" />
                    <div class="text-sm text-text-secondary">
                      <p class="mb-1"><strong class="text-text-primary">为什么只有 1 个成功？</strong></p>
                      <p>Redis 分布式锁使用 <code class="text-brand">SET key value NX EX ttl</code> 命令，NX 保证了原子性——同一时刻只有一个客户端能设置成功。其他客户端会立即返回失败，这就是分布式锁的互斥语义。</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </GlowCard>

          <!-- Request Log -->
          <div class="rounded-lg border border-border overflow-hidden">
            <div class="flex items-center justify-between px-4 py-2.5 bg-bg-elevated border-b border-border">
              <div class="flex items-center gap-2">
                <i class="i-lucide-terminal w-4 h-4 text-brand" />
                <span class="text-sm font-medium text-text-primary">请求日志</span>
              </div>
              <span class="text-xs text-text-muted">{{ requestLogs.length }} 条请求</span>
            </div>
            <RequestLog :logs="requestLogs" />
          </div>

          <!-- API Reference -->
          <GlowCard title="接口列表">
            <div class="overflow-x-auto">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-border">
                    <th class="text-left py-3 text-text-muted font-medium">方法</th>
                    <th class="text-left py-3 text-text-muted font-medium">接口</th>
                    <th class="text-left py-3 text-text-muted font-medium">说明</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="api in apiList" :key="api.endpoint" class="border-b border-border/50">
                    <td class="py-3">
                      <span class="px-2 py-0.5 rounded text-xs font-bold bg-brand/20 text-brand">
                        {{ api.method }}
                      </span>
                    </td>
                    <td class="py-3 font-mono text-text-secondary text-xs">{{ api.endpoint }}</td>
                    <td class="py-3 text-text-secondary">{{ api.desc }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </GlowCard>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import * as lockApi from '@/api/lock'
import { requestLogs } from '@/api/request'
import GlowCard from '@/components/GlowCard.vue'
import RequestLog from '@/components/RequestLog.vue'
import type { LockStatusResponse, ContentionDemoResponse } from '@/types'

const activeStep = ref(0)
const activeTab = ref('lock')
const acquiring = ref(false)
const querying = ref(false)
const contending = ref(false)

const lockForm = reactive({ resource: 'order:10086', ttl: 10 })
const statusResource = ref('')
const contentionForm = reactive({
  resource: 'demo-resource',
  ttl: 10,
  goroutines: 5,
  hold_ms: 500,
})

const lockResult = ref<{ success: boolean; message: string } | null>(null)
const statusResult = ref<LockStatusResponse | null>(null)
const contentionResult = ref<ContentionDemoResponse | null>(null)

const currentLock = ref<{ locked: boolean; resource: string; ttl_ms: number; owner: string }>({
  locked: false, resource: '', ttl_ms: 0, owner: '',
})

const redisKeys = ref<{ name: string; locked: boolean; ttl: number }[]>([])

const quickResources = ['order:10086', 'user:login', 'payment:9527']

const tabs = [
  { key: 'lock', label: '获取 / 释放', icon: 'i-lucide-lock' },
  { key: 'status', label: '状态查询', icon: 'i-lucide-search' },
  { key: 'contention', label: '并发争抢', icon: 'i-lucide-zap' },
]

const steps = computed(() => [
  { label: '获取分布式锁', completed: lockResult.value?.success === true },
  { label: '查询锁状态', completed: statusResult.value !== null },
  { label: '演示并发争抢', completed: contentionResult.value !== null },
  { label: '释放锁', completed: false },
])

const concepts = [
  { title: 'SET NX EX', desc: '原子操作：不存在才设置，同时设置过期时间', icon: 'i-lucide-key' },
  { title: 'Lua 原子释放', desc: '比较 value 后再删除，防止误删其他客户端的锁', icon: 'i-lucide-code' },
  { title: 'TTL 自动过期', desc: '即使持有者崩溃，锁也会自动释放', icon: 'i-lucide-timer' },
  { title: 'Owner 标识', desc: '每个锁持有者有唯一 ID，只有持有者能释放锁', icon: 'i-lucide-fingerprint' },
]

const apiList = [
  { method: 'POST', endpoint: '/api/v1/lock/acquire', desc: '获取分布式锁' },
  { method: 'POST', endpoint: '/api/v1/lock/release', desc: '释放分布式锁' },
  { method: 'POST', endpoint: '/api/v1/lock/status', desc: '查询锁状态' },
  { method: 'POST', endpoint: '/api/v1/lock/contention', desc: '并发争抢演示' },
]

async function handleAcquire() {
  if (!lockForm.resource) {
    ElMessage.warning('请输入资源名称')
    return
  }
  acquiring.value = true
  lockResult.value = null
  try {
    const res = await lockApi.acquireLock({ resource: lockForm.resource, ttl: lockForm.ttl })
    lockResult.value = {
      success: true,
      message: JSON.stringify(res.data, null, 2),
    }
    currentLock.value = {
      locked: true,
      resource: lockForm.resource,
      ttl_ms: lockForm.ttl * 1000,
      owner: res.data.owner,
    }
    updateRedisKeys()
    activeStep.value = 0
    ElMessage.success('锁获取成功')
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : '获取锁失败'
    lockResult.value = { success: false, message: msg }
    ElMessage.error(msg)
  } finally {
    acquiring.value = false
  }
}

async function handleRelease(resource: string) {
  try {
    await lockApi.releaseLock({ resource })
    if (currentLock.value.resource === resource) {
      currentLock.value = { locked: false, resource: '', ttl_ms: 0, owner: '' }
    }
    removeRedisKey(resource)
    ElMessage.success('锁已释放')
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '释放锁失败')
  }
}

async function handleQueryStatus() {
  if (!statusResource.value) {
    ElMessage.warning('请输入资源名称')
    return
  }
  querying.value = true
  try {
    const res = await lockApi.getLockStatus({ resource: statusResource.value })
    statusResult.value = res.data
    activeStep.value = 1
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '查询失败')
  } finally {
    querying.value = false
  }
}

async function handleContention() {
  contending.value = true
  contentionResult.value = null
  try {
    const res = await lockApi.contentionDemo({
      resource: contentionForm.resource,
      ttl: contentionForm.ttl,
      goroutines: contentionForm.goroutines,
      hold_ms: contentionForm.hold_ms,
    })
    contentionResult.value = res.data
    activeStep.value = 2
    ElMessage.success('并发争抢演示完成')
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '演示失败')
  } finally {
    contending.value = false
  }
}

function updateRedisKeys() {
  const existing = redisKeys.value.find(k => k.name === `lock:${lockForm.resource}`)
  if (existing) {
    existing.locked = true
    existing.ttl = lockForm.ttl
  } else {
    redisKeys.value.push({ name: `lock:${lockForm.resource}`, locked: true, ttl: lockForm.ttl })
  }
}

function removeRedisKey(resource: string) {
  const key = `lock:${resource}`
  const existing = redisKeys.value.find(k => k.name === key)
  if (existing) {
    existing.locked = false
    existing.ttl = 0
  }
}
</script>
