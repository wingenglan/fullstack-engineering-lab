<template>
  <div class="jwt-auth-view py-8 px-6">
    <div class="max-w-7xl mx-auto">
      <!-- Breadcrumb -->
      <div class="flex items-center gap-2 text-sm text-text-muted mb-6">
        <router-link to="/cases" class="hover:text-text-primary transition-colors">案例列表</router-link>
        <i class="i-lucide-chevron-right w-4 h-4" />
        <span class="text-text-primary">JWT 认证授权</span>
      </div>

      <div class="grid lg:grid-cols-3 gap-6">
        <!-- Left Column -->
        <div class="lg:col-span-1 space-y-6">
          <!-- Case Info -->
          <GlowCard>
            <div class="flex items-center gap-3 mb-4">
              <div class="w-10 h-10 rounded-lg bg-brand/10 flex items-center justify-center">
                <i class="i-lucide-shield w-5 h-5 text-brand" />
              </div>
              <div>
                <h2 class="text-text-primary font-bold text-lg">JWT 认证授权</h2>
                <span class="text-xs text-success">已上线</span>
              </div>
            </div>
            <p class="text-text-secondary text-sm mb-4">
              基于 JWT 的完整认证授权案例，包含注册、登录、Token 管理和基于 Redis 的黑名单机制。
            </p>
            <div class="flex flex-wrap gap-2 mb-4">
              <span v-for="tag in ['Go', 'Gin', 'JWT', 'Redis', 'bcrypt']" :key="tag" class="px-2.5 py-0.5 text-xs rounded-full bg-white/5 text-text-secondary border border-border">
                {{ tag }}
              </span>
            </div>
            <div class="text-xs text-text-muted">
              难度：<span class="text-success font-medium">入门</span>
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
                <i v-if="step.completed" class="i-lucide-check-circle w-4 h-4 text-success ml-auto" />
              </div>
            </div>
          </GlowCard>

          <!-- Flow Diagram -->
          <GlowCard title="认证流程">
            <div class="space-y-2 text-sm">
              <div v-for="(flow, i) in authFlow" :key="i" class="flex items-center gap-3">
                <div class="w-7 h-7 rounded-full bg-brand/10 flex items-center justify-center shrink-0">
                  <i :class="flow.icon" class="w-3.5 h-3.5 text-brand" />
                </div>
                <span class="text-text-secondary">{{ flow.label }}</span>
              </div>
            </div>
          </GlowCard>
        </div>

        <!-- Right Column -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Auth Status Bar -->
          <div class="flex items-center justify-between p-4 rounded-xl border border-border bg-bg-card">
            <div class="flex items-center gap-3">
              <div
                class="w-3 h-3 rounded-full"
                :class="authStore.isAuthenticated ? 'bg-success shadow-[0_0_8px_rgba(16,185,129,0.5)]' : 'bg-text-muted'"
              />
              <span class="text-sm text-text-primary">
                {{ authStore.isAuthenticated ? `已登录：${authStore.displayName}` : '未登录' }}
              </span>
            </div>
            <el-button v-if="authStore.isAuthenticated" size="small" type="danger" plain @click="handleLogout">
              退出登录
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

          <!-- Register Tab -->
          <GlowCard v-if="activeTab === 'register'">
            <h3 class="text-text-primary font-semibold text-lg mb-6">注册新用户</h3>
            <el-form :model="registerForm" label-position="top" @submit.prevent="handleRegister">
              <el-form-item label="用户名">
                <el-input v-model="registerForm.username" placeholder="请输入用户名（3-64 个字符）" />
              </el-form-item>
              <el-form-item label="邮箱">
                <el-input v-model="registerForm.email" type="email" placeholder="请输入邮箱" />
              </el-form-item>
              <el-form-item label="密码">
                <el-input v-model="registerForm.password" type="password" show-password placeholder="请输入密码（6 位以上）" />
              </el-form-item>
              <el-button type="primary" :loading="registering" class="w-full" @click="handleRegister">
                <i class="i-lucide-user-plus w-4 h-4 mr-2" />
                注册
              </el-button>
            </el-form>
          </GlowCard>

          <!-- Login Tab -->
          <GlowCard v-if="activeTab === 'login'">
            <h3 class="text-text-primary font-semibold text-lg mb-6">用户登录</h3>
            <el-form :model="loginForm" label-position="top" @submit.prevent="handleLogin">
              <el-form-item label="用户名">
                <el-input v-model="loginForm.username" placeholder="请输入用户名" />
              </el-form-item>
              <el-form-item label="密码">
                <el-input v-model="loginForm.password" type="password" show-password placeholder="请输入密码" />
              </el-form-item>
              <el-button type="primary" :loading="loggingIn" class="w-full" @click="handleLogin">
                <i class="i-lucide-log-in w-4 h-4 mr-2" />
                登录
              </el-button>
            </el-form>
          </GlowCard>

          <!-- Profile Tab -->
          <GlowCard v-if="activeTab === 'profile'">
            <div v-if="authStore.isAuthenticated && authStore.user" class="space-y-6">
              <h3 class="text-text-primary font-semibold text-lg">用户信息</h3>
              <div class="grid grid-cols-2 gap-4">
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <span class="text-xs text-text-muted block mb-1">用户 ID</span>
                  <span class="text-text-primary font-mono">{{ authStore.user.id }}</span>
                </div>
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <span class="text-xs text-text-muted block mb-1">用户名</span>
                  <span class="text-text-primary">{{ authStore.user.username }}</span>
                </div>
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <span class="text-xs text-text-muted block mb-1">邮箱</span>
                  <span class="text-text-primary">{{ authStore.user.email }}</span>
                </div>
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <span class="text-xs text-text-muted block mb-1">状态</span>
                  <span class="text-success text-sm flex items-center gap-1.5">
                    <span class="w-2 h-2 rounded-full bg-success" />
                    正常
                  </span>
                </div>
                <div class="col-span-2 p-4 rounded-lg bg-bg-elevated border border-border">
                  <span class="text-xs text-text-muted block mb-1">注册时间</span>
                  <span class="text-text-primary font-mono text-sm">{{ authStore.user.created_at }}</span>
                </div>
              </div>
            </div>
            <div v-else class="text-center py-12">
              <i class="i-lucide-lock w-12 h-12 text-text-muted mx-auto mb-4 block" />
              <p class="text-text-secondary mb-4">请先登录以查看用户信息</p>
              <el-button type="primary" @click="activeTab = 'login'">前往登录</el-button>
            </div>
          </GlowCard>

          <!-- Token Tab -->
          <GlowCard v-if="activeTab === 'token'">
            <div v-if="authStore.token" class="space-y-6">
              <h3 class="text-text-primary font-semibold text-lg">Token 详情</h3>

              <!-- Raw Token -->
              <div>
                <div class="flex items-center justify-between mb-2">
                  <span class="text-sm text-text-secondary font-medium">Access Token</span>
                  <button class="flex items-center gap-1.5 text-xs text-text-muted hover:text-text-primary" @click="copyToken">
                    <i class="i-lucide-copy w-3.5 h-3.5" />
                    复制
                  </button>
                </div>
                <pre class="p-4 rounded-lg bg-[#0D0D0F] border border-border text-xs font-mono text-text-secondary break-all whitespace-pre-wrap overflow-x-auto max-h-32">{{ authStore.token }}</pre>
              </div>

              <!-- Token Header -->
              <div>
                <span class="text-sm text-text-secondary font-medium mb-2 block">Header（头部）</span>
                <pre class="p-4 rounded-lg bg-[#0D0D0F] border border-border text-xs font-mono overflow-x-auto"><span class="text-blue-400">{{ formatJson(tokenHeader) }}</span></pre>
              </div>

              <!-- Token Payload -->
              <div>
                <span class="text-sm text-text-secondary font-medium mb-2 block">Payload（载荷）</span>
                <pre class="p-4 rounded-lg bg-[#0D0D0F] border border-border text-xs font-mono overflow-x-auto"><span class="text-purple-400">{{ formatJson(tokenPayload) }}</span></pre>
              </div>

              <!-- Expiry Info -->
              <div v-if="tokenPayload?.exp" class="flex items-center gap-3 p-3 rounded-lg bg-bg-elevated border border-border">
                <i class="i-lucide-clock w-4 h-4 text-text-muted" />
                <span class="text-sm text-text-secondary">
                  过期时间：<span class="text-text-primary font-mono">{{ formatDate(tokenPayload.exp) }}</span>
                </span>
                <span
                  class="ml-auto px-2 py-0.5 text-xs rounded-full"
                  :class="isExpired ? 'bg-danger/10 text-danger' : 'bg-success/10 text-success'"
                >
                  {{ isExpired ? '已过期' : '有效' }}
                </span>
              </div>
            </div>
            <div v-else class="text-center py-12">
              <i class="i-lucide-key w-12 h-12 text-text-muted mx-auto mb-4 block" />
              <p class="text-text-secondary mb-4">暂无 Token，请先登录</p>
              <el-button type="primary" @click="activeTab = 'login'">前往登录</el-button>
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
                    <th class="text-left py-3 text-text-muted font-medium">认证</th>
                    <th class="text-left py-3 text-text-muted font-medium">说明</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="api in apiList" :key="api.endpoint" class="border-b border-border/50">
                    <td class="py-3">
                      <span class="px-2 py-0.5 rounded text-xs font-bold" :class="methodBadge(api.method)">
                        {{ api.method }}
                      </span>
                    </td>
                    <td class="py-3 font-mono text-text-secondary text-xs">{{ api.endpoint }}</td>
                    <td class="py-3">
                      <i v-if="api.auth" class="i-lucide-lock w-3.5 h-3.5 text-warning" />
                      <i v-else class="i-lucide-globe w-3.5 h-3.5 text-text-muted" />
                    </td>
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
import { useAuthStore } from '@/stores/auth'
import { requestLogs } from '@/api/request'
import { decodeTokenHeader, getTokenPayload, isTokenExpired } from '@/utils/jwt'
import GlowCard from '@/components/GlowCard.vue'
import RequestLog from '@/components/RequestLog.vue'

const authStore = useAuthStore()

const activeStep = ref(0)
const activeTab = ref('register')
const registering = ref(false)
const loggingIn = ref(false)

const registerForm = reactive({ username: '', email: '', password: '' })
const loginForm = reactive({ username: '', password: '' })

const tabs = [
  { key: 'register', label: '注册', icon: 'i-lucide-user-plus' },
  { key: 'login', label: '登录', icon: 'i-lucide-log-in' },
  { key: 'profile', label: '用户信息', icon: 'i-lucide-user' },
  { key: 'token', label: 'Token', icon: 'i-lucide-key' },
]

const steps = computed(() => [
  { label: '注册用户', completed: !!authStore.user },
  { label: '登录获取 Token', completed: !!authStore.token },
  { label: '获取用户信息', completed: !!authStore.user },
  { label: '查看 Token 详情', completed: !!authStore.token },
  { label: '退出登录', completed: false },
])

const authFlow = [
  { label: '注册（bcrypt 密码加密）', icon: 'i-lucide-user-plus' },
  { label: '登录（验证密码）', icon: 'i-lucide-log-in' },
  { label: '签发 JWT Token', icon: 'i-lucide-key' },
  { label: '访问受保护接口', icon: 'i-lucide-shield-check' },
  { label: '退出（Token 加入黑名单）', icon: 'i-lucide-ban' },
]

const apiList = [
  { method: 'POST', endpoint: '/api/v1/auth/register', auth: false, desc: '用户注册' },
  { method: 'POST', endpoint: '/api/v1/auth/login', auth: false, desc: '用户登录，获取 Token' },
  { method: 'GET', endpoint: '/api/v1/auth/profile', auth: true, desc: '获取当前用户信息' },
  { method: 'POST', endpoint: '/api/v1/auth/logout', auth: true, desc: '退出登录（撤销 Token）' },
  { method: 'GET', endpoint: '/api/v1/health', auth: false, desc: '健康检查' },
]

const tokenHeader = computed(() => {
  if (!authStore.token) return null
  return decodeTokenHeader(authStore.token)
})

const tokenPayload = computed(() => {
  if (!authStore.token) return null
  return getTokenPayload(authStore.token)
})

const isExpired = computed(() => {
  if (!authStore.token) return true
  return isTokenExpired(authStore.token)
})

async function handleRegister() {
  if (!registerForm.username || !registerForm.email || !registerForm.password) {
    ElMessage.warning('请填写所有字段')
    return
  }
  registering.value = true
  try {
    await authStore.register(registerForm.username, registerForm.email, registerForm.password)
    ElMessage.success('注册成功！请登录')
    activeTab.value = 'login'
    loginForm.username = registerForm.username
    activeStep.value = 1
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '注册失败')
  } finally {
    registering.value = false
  }
}

async function handleLogin() {
  if (!loginForm.username || !loginForm.password) {
    ElMessage.warning('请填写所有字段')
    return
  }
  loggingIn.value = true
  try {
    await authStore.login(loginForm.username, loginForm.password)
    ElMessage.success('登录成功！')
    activeTab.value = 'profile'
    activeStep.value = 2
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '登录失败')
  } finally {
    loggingIn.value = false
  }
}

async function handleLogout() {
  await authStore.logout()
  ElMessage.success('已退出登录，Token 已加入黑名单')
  activeTab.value = 'login'
  activeStep.value = 0
}

function copyToken() {
  if (authStore.token) {
    navigator.clipboard.writeText(authStore.token)
    ElMessage.success('Token 已复制！')
  }
}

function formatJson(obj: unknown): string {
  if (!obj) return ''
  return JSON.stringify(obj, null, 2)
}

function formatDate(timestamp: number): string {
  return new Date(timestamp * 1000).toLocaleString()
}

function methodBadge(method: string): string {
  const map: Record<string, string> = {
    GET: 'bg-success/20 text-success',
    POST: 'bg-brand/20 text-brand',
    PUT: 'bg-warning/20 text-warning',
    DELETE: 'bg-danger/20 text-danger',
  }
  return map[method] || ''
}
</script>
