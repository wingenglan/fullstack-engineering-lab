<template>
  <div class="websocket-chat-view py-8 px-6">
    <div class="max-w-7xl mx-auto">
      <!-- Breadcrumb -->
      <div class="flex items-center gap-2 text-sm text-text-muted mb-6">
        <router-link to="/cases" class="hover:text-text-primary transition-colors">案例列表</router-link>
        <i class="i-lucide-chevron-right w-4 h-4" />
        <span class="text-text-primary">WebSocket 实时通信</span>
      </div>

      <!-- 未登录提示 -->
      <div v-if="!authStore.isAuthenticated" class="mb-6">
        <GlowCard>
          <div class="flex items-center gap-3">
            <i class="i-lucide-alert-triangle w-5 h-5 text-warning" />
            <div>
              <span class="text-text-primary font-medium">请先登录</span>
              <p class="text-text-muted text-sm mt-1">实时通讯需要认证，请先登录或注册账号。</p>
            </div>
          </div>
        </GlowCard>
      </div>

      <div class="grid lg:grid-cols-3 gap-6">
        <!-- Left Column -->
        <div class="lg:col-span-1 space-y-6">
          <!-- Case Info -->
          <GlowCard>
            <div class="flex items-center gap-3 mb-4">
              <div class="w-10 h-10 rounded-lg bg-brand/10 flex items-center justify-center">
                <i class="i-lucide-message-circle w-5 h-5 text-brand" />
              </div>
              <div>
                <h2 class="text-text-primary font-bold text-lg">WebSocket 实时通讯</h2>
                <span class="text-xs text-success">已上线</span>
              </div>
            </div>
            <p class="text-text-secondary text-sm mb-4">
              基于 WebSocket 协议的实时双向通信聊天室，支持多房间、消息持久化、在线状态和输入指示器。
            </p>
            <div class="flex flex-wrap gap-2 mb-4">
              <span v-for="tag in ['WebSocket', 'Go', 'Vue', '实时通信', 'Hub 模型']" :key="tag" class="px-2.5 py-0.5 text-xs rounded-full bg-white/5 text-text-secondary border border-border">
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

          <!-- WebSocket Concepts -->
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

          <!-- Room List -->
          <GlowCard title="聊天室列表">
            <div class="space-y-2">
              <div
                v-for="room in chatStore.rooms"
                :key="room.id"
                class="flex items-center justify-between px-3 py-2.5 rounded-lg cursor-pointer transition-all border"
                :class="chatStore.currentRoomId === room.id
                  ? 'bg-brand/10 border-brand/30 text-brand'
                  : 'bg-bg-elevated border-border text-text-secondary hover:border-brand/20'"
                @click="handleJoinRoom(room.id)"
              >
                <div class="flex items-center gap-2 min-w-0">
                  <i class="i-lucide-hash w-4 h-4 shrink-0" />
                  <span class="text-sm font-medium truncate">{{ room.name }}</span>
                </div>
                <div class="flex items-center gap-1.5 shrink-0">
                  <span class="w-2 h-2 rounded-full bg-success" />
                  <span class="text-xs text-text-muted">{{ room.member_count }}</span>
                </div>
              </div>
              <div v-if="chatStore.rooms.length === 0" class="text-center py-4 text-text-muted text-xs">
                暂无聊天室
              </div>
            </div>
          </GlowCard>

          <!-- Online Users -->
          <GlowCard title="在线用户">
            <div class="space-y-1.5">
              <div
                v-for="user in chatStore.onlineUsers"
                :key="user.user_id"
                class="flex items-center gap-2.5 px-3 py-2 rounded-lg bg-bg-elevated"
              >
                <span class="w-2 h-2 rounded-full bg-success shrink-0" />
                <span class="text-sm text-text-secondary truncate">{{ user.username }}</span>
              </div>
              <div v-if="chatStore.onlineUsers.length === 0" class="text-center py-4 text-text-muted text-xs">
                {{ chatStore.currentRoomId ? '暂无在线用户' : '请先加入聊天室' }}
              </div>
            </div>
          </GlowCard>
        </div>

        <!-- Right Column -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Connection Status Bar -->
          <div class="flex items-center justify-between p-4 rounded-xl border border-border bg-bg-card">
            <div class="flex items-center gap-3">
              <div
                class="w-3 h-3 rounded-full"
                :class="connectionStatusClass"
              />
              <span class="text-sm text-text-primary">
                {{ connectionStatusText }}
              </span>
            </div>
            <div class="flex items-center gap-2">
              <el-button
                v-if="!chatStore.connected && authStore.isAuthenticated"
                size="small"
                type="primary"
                @click="handleConnect"
              >
                连接
              </el-button>
              <el-button
                v-if="chatStore.connected"
                size="small"
                type="danger"
                plain
                @click="handleDisconnect"
              >
                断开
              </el-button>
            </div>
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

          <!-- Chat Room Tab -->
          <GlowCard v-if="activeTab === 'chat'" class="!p-0 overflow-hidden">
            <!-- Room Header -->
            <div class="flex items-center justify-between px-6 py-4 border-b border-border">
              <div class="flex items-center gap-3">
                <i class="i-lucide-hash w-5 h-5 text-brand" />
                <div>
                  <h3 class="text-text-primary font-semibold">
                    {{ chatStore.currentRoom?.name || '请选择聊天室' }}
                  </h3>
                  <p v-if="chatStore.currentRoom" class="text-xs text-text-muted">
                    {{ chatStore.currentRoom.description }} · {{ chatStore.onlineCount }} 人在线
                  </p>
                </div>
              </div>
              <el-button size="small" plain @click="chatStore.refreshOnlineUsers()">
                <i class="i-lucide-refresh-cw w-3.5 h-3.5 mr-1" />
                刷新
              </el-button>
            </div>

            <!-- Messages Area -->
            <div
              ref="messagesContainer"
              class="h-96 overflow-y-auto px-6 py-4 space-y-4"
              @scroll="handleScroll"
            >
              <!-- Load More -->
              <div v-if="chatStore.hasMoreMessages && chatStore.currentRoomId" class="text-center">
                <el-button
                  size="small"
                  text
                  :loading="chatStore.loadingMessages"
                  @click="chatStore.loadMoreHistory()"
                >
                  加载更多历史消息
                </el-button>
              </div>

              <!-- System Notices -->
              <div
                v-for="(notice, i) in chatStore.systemNotices"
                :key="'notice-' + i"
                class="text-center"
              >
                <span class="inline-block px-3 py-1 rounded-full bg-white/5 text-xs text-text-muted">
                  {{ notice }}
                </span>
              </div>

              <!-- Messages -->
              <div
                v-for="msg in chatStore.currentMessages"
                :key="msg.id"
                class="flex gap-3"
                :class="isMyMessage(msg) ? 'flex-row-reverse' : ''"
              >
                <!-- Avatar -->
                <div
                  class="w-8 h-8 rounded-full flex items-center justify-center shrink-0 text-xs font-bold"
                  :class="isMyMessage(msg) ? 'bg-brand text-white' : 'bg-white/10 text-text-secondary'"
                >
                  {{ getAvatarLetter(msg.username) }}
                </div>

                <!-- Message Bubble -->
                <div
                  class="max-w-[70%]"
                  :class="isMyMessage(msg) ? 'items-end' : 'items-start'"
                >
                  <div class="flex items-center gap-2 mb-1" :class="isMyMessage(msg) ? 'flex-row-reverse' : ''">
                    <span class="text-xs font-medium text-text-secondary">{{ msg.username }}</span>
                    <span class="text-xs text-text-muted">{{ formatTime(msg.created_at) }}</span>
                  </div>
                  <div
                    class="px-4 py-2.5 rounded-2xl text-sm"
                    :class="isMyMessage(msg)
                      ? 'bg-brand text-white rounded-tr-sm'
                      : 'bg-white/5 text-text-primary rounded-tl-sm border border-border'"
                  >
                    {{ msg.content }}
                  </div>
                </div>
              </div>

              <!-- Empty State -->
              <div
                v-if="!chatStore.currentRoomId"
                class="flex flex-col items-center justify-center h-full text-text-muted"
              >
                <i class="i-lucide-message-circle w-12 h-12 mb-3 opacity-30" />
                <p class="text-sm">请从左侧选择一个聊天室加入</p>
              </div>
              <div
                v-else-if="chatStore.currentMessages.length === 0 && !chatStore.loadingMessages"
                class="flex flex-col items-center justify-center h-full text-text-muted"
              >
                <i class="i-lucide-messages-square w-12 h-12 mb-3 opacity-30" />
                <p class="text-sm">还没有消息，发送第一条消息吧！</p>
              </div>
            </div>

            <!-- Typing Indicator -->
            <div v-if="chatStore.typingUsers.size > 0" class="px-6 py-1.5 border-t border-border/50">
              <span class="text-xs text-text-muted">
                {{ typingText }} 正在输入...
              </span>
            </div>

            <!-- Input Area -->
            <div v-if="chatStore.currentRoomId" class="px-6 py-4 border-t border-border">
              <div class="flex gap-3">
                <el-input
                  v-model="messageInput"
                  placeholder="输入消息..."
                  :disabled="!chatStore.connected"
                  @keyup.enter="handleSendMessage"
                  @input="handleInput"
                />
                <el-button
                  type="primary"
                  :disabled="!chatStore.connected || !messageInput.trim()"
                  @click="handleSendMessage"
                >
                  <i class="i-lucide-send w-4 h-4" />
                </el-button>
              </div>
            </div>
          </GlowCard>

          <!-- Create Room Tab -->
          <GlowCard v-if="activeTab === 'create'">
            <h3 class="text-text-primary font-semibold text-lg mb-6">创建聊天室</h3>

            <div class="space-y-5">
              <div>
                <label class="text-sm text-text-secondary mb-2 block">房间名称</label>
                <el-input v-model="createForm.name" placeholder="例如：前端技术交流" maxlength="128" show-word-limit />
              </div>

              <div>
                <label class="text-sm text-text-secondary mb-2 block">房间描述</label>
                <el-input
                  v-model="createForm.description"
                  type="textarea"
                  :rows="3"
                  placeholder="简单描述一下这个聊天室的用途..."
                  maxlength="512"
                  show-word-limit
                />
              </div>

              <el-button
                type="primary"
                class="w-full"
                :loading="creating"
                :disabled="!createForm.name.trim() || !chatStore.connected"
                @click="handleCreateRoom"
              >
                <i class="i-lucide-plus w-4 h-4 mr-2" />
                创建聊天室
              </el-button>
            </div>
          </GlowCard>

          <!-- Message Log Tab -->
          <GlowCard v-if="activeTab === 'protocol'">
            <h3 class="text-text-primary font-semibold text-lg mb-6">WebSocket 协议说明</h3>

            <div class="space-y-6">
              <!-- Protocol Format -->
              <div>
                <h4 class="text-text-primary font-medium mb-3">消息格式</h4>
                <div class="p-4 rounded-lg bg-[#0D0D0F] border border-border">
                  <pre class="text-xs font-mono text-text-secondary">{{ protocolFormat }}</pre>
                </div>
              </div>

              <!-- Message Types -->
              <div>
                <h4 class="text-text-primary font-medium mb-3">消息类型</h4>
                <div class="overflow-x-auto">
                  <table class="w-full text-sm">
                    <thead>
                      <tr class="border-b border-border">
                        <th class="text-left py-3 text-text-muted font-medium">方向</th>
                        <th class="text-left py-3 text-text-muted font-medium">类型</th>
                        <th class="text-left py-3 text-text-muted font-medium">说明</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="proto in protocolTypes" :key="proto.type" class="border-b border-border/50">
                        <td class="py-3">
                          <span class="px-2 py-0.5 rounded text-xs font-bold"
                            :class="proto.dir === 'C→S' ? 'bg-brand/20 text-brand' : 'bg-success/20 text-success'">
                            {{ proto.dir }}
                          </span>
                        </td>
                        <td class="py-3 font-mono text-text-secondary text-xs">{{ proto.type }}</td>
                        <td class="py-3 text-text-secondary">{{ proto.desc }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- Flow Diagram -->
              <div>
                <h4 class="text-text-primary font-medium mb-3">连接流程</h4>
                <div class="p-4 rounded-lg bg-[#0D0D0F] border border-border">
                  <pre class="text-xs font-mono text-text-secondary whitespace-pre-wrap">{{ flowDiagram }}</pre>
                </div>
              </div>
            </div>
          </GlowCard>

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
                      <span class="px-2 py-0.5 rounded text-xs font-bold" :class="methodClass(api.method)">
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
import { ref, reactive, computed, watch, nextTick, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useChatStore } from '@/stores/chat'
import { useAuthStore } from '@/stores/auth'
import GlowCard from '@/components/GlowCard.vue'
import type { ChatMessage } from '@/types'

const chatStore = useChatStore()
const authStore = useAuthStore()

const activeStep = ref(0)
const activeTab = ref('chat')
const creating = ref(false)
const messageInput = ref('')
const messagesContainer = ref<HTMLElement | null>(null)

const createForm = reactive({
  name: '',
  description: '',
})

// 连接状态
const connectionStatusClass = computed(() => {
  switch (chatStore.status) {
    case 'connected': return 'bg-success shadow-[0_0_8px_rgba(34,197,94,0.5)]'
    case 'connecting': return 'bg-warning shadow-[0_0_8px_rgba(234,179,8,0.5)] animate-pulse'
    case 'error': return 'bg-danger shadow-[0_0_8px_rgba(239,68,68,0.5)]'
    default: return 'bg-text-muted'
  }
})

const connectionStatusText = computed(() => {
  switch (chatStore.status) {
    case 'connected': return '已连接到 WebSocket 服务'
    case 'connecting': return '正在连接...'
    case 'error': return `连接失败，已重试 ${chatStore.status === 'error' ? '多次' : ''}`
    default: return '未连接'
  }
})

// 输入指示器文本
const typingText = computed(() => {
  const names = Array.from(chatStore.typingUsers.values())
  if (names.length === 1) return names[0]
  if (names.length === 2) return `${names[0]} 和 ${names[1]}`
  return `${names[0]} 等 ${names.length} 人`
})

// 步骤
const steps = computed(() => [
  { label: '连接 WebSocket', completed: chatStore.connected },
  { label: '加入聊天室', completed: chatStore.currentRoomId !== null },
  { label: '发送消息', completed: chatStore.currentMessages.some(m => isMyMessage(m)) },
  { label: '查看在线用户', completed: chatStore.onlineCount > 0 },
])

const tabs = [
  { key: 'chat', label: '聊天室', icon: 'i-lucide-message-circle' },
  { key: 'create', label: '创建房间', icon: 'i-lucide-plus-circle' },
  { key: 'protocol', label: '协议说明', icon: 'i-lucide-file-text' },
]

const concepts = [
  { title: 'WebSocket 协议', desc: 'HTTP 升级为全双工通信，服务端可主动推送', icon: 'i-lucide-plug' },
  { title: 'Hub 模型', desc: '一个 Hub 管理所有连接，按房间分组广播消息', icon: 'i-lucide-network' },
  { title: '消息持久化', desc: '消息同时写入 MySQL，新加入房间可查看历史', icon: 'i-lucide-database' },
  { title: '心跳保活', desc: '客户端定时 Ping，服务端 Pong，检测死连接', icon: 'i-lucide-heart-pulse' },
]

const apiList = [
  { method: 'GET', endpoint: '/api/v1/chat/ws', desc: 'WebSocket 连接（需 Bearer Token）' },
  { method: 'GET', endpoint: '/api/v1/chat/rooms', desc: '获取聊天室列表' },
  { method: 'POST', endpoint: '/api/v1/chat/rooms', desc: '创建聊天室（需认证）' },
  { method: 'GET', endpoint: '/api/v1/chat/rooms/:id', desc: '获取聊天室详情' },
  { method: 'GET', endpoint: '/api/v1/chat/rooms/:id/online', desc: '获取在线用户列表' },
  { method: 'GET', endpoint: '/api/v1/chat/messages', desc: '获取消息历史' },
]

const protocolFormat = `{
  "type": "send_message",
  "payload": {
    "room_id": 1,
    "content": "大家好！",
    "msg_type": 1
  }
}`

const protocolTypes = [
  { dir: 'C→S', type: 'join_room', desc: '加入指定房间，接收该房间消息' },
  { dir: 'C→S', type: 'leave_room', desc: '离开指定房间，停止接收消息' },
  { dir: 'C→S', type: 'send_message', desc: '向当前房间发送消息' },
  { dir: 'C→S', type: 'typing', desc: '通知服务端用户正在输入' },
  { dir: 'C→S', type: 'ping', desc: '心跳检测' },
  { dir: 'S→C', type: 'new_message', desc: '广播新消息到房间所有成员' },
  { dir: 'S→C', type: 'user_joined', desc: '通知有人加入房间' },
  { dir: 'S→C', type: 'user_left', desc: '通知有人离开房间' },
  { dir: 'S→C', type: 'online_users', desc: '推送当前房间在线用户列表' },
  { dir: 'S→C', type: 'user_typing', desc: '广播输入指示器（排除发送者）' },
  { dir: 'S→C', type: 'room_history', desc: '加入房间后推送历史消息' },
  { dir: 'S→C', type: 'pong', desc: '心跳响应' },
  { dir: 'S→C', type: 'error', desc: '错误消息' },
]

const flowDiagram = `Client                        Server
  |                              |
  |-- HTTP GET /chat/ws -------->|  1. WebSocket 握手（带 Token）
  |<----- 101 Switching Protocols|
  |                              |
  |-- join_room {room_id} ------>|  2. 加入房间
  |<-- room_history {msgs} ------|  3. 推送历史消息
  |<-- user_joined {user} -------|  4. 广播加入通知
  |<-- online_users {list} ------|  5. 推送在线用户
  |                              |
  |-- send_message {content} --->|  6. 发送消息
  |<-- new_message {msg} --------|  7. 广播到房间所有人
  |                              |
  |-- typing {room_id} --------->|  8. 输入指示器
  |<-- user_typing {user} -------|  9. 广播给其他人（排除自己）
  |                              |
  |-- ping {} ------------------>|  10. 心跳
  |<-- pong {} ------------------|  11. 心跳响应
  |                              |
  |-- leave_room {room_id} ----->|  12. 离开房间
  |<-- user_left {user} ---------|  13. 广播离开通知`

function methodClass(method: string): string {
  const map: Record<string, string> = {
    GET: 'bg-success/20 text-success',
    POST: 'bg-brand/20 text-brand',
    WS: 'bg-warning/20 text-warning',
  }
  return map[method] || 'bg-white/10 text-text-secondary'
}

// 判断是否是自己的消息
function isMyMessage(msg: ChatMessage): boolean {
  return authStore.user?.id === msg.user_id
}

// 获取头像字母
function getAvatarLetter(username: string): string {
  return username.charAt(0).toUpperCase()
}

// 格式化时间
function formatTime(dateStr: string): string {
  try {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return ''
    return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
  } catch {
    return ''
  }
}

// 连接 WebSocket
function handleConnect() {
  chatStore.connectWS()
  activeStep.value = 0
}

// 断开连接
function handleDisconnect() {
  chatStore.disconnectWS()
}

// 加入房间
async function handleJoinRoom(roomId: number) {
  if (!chatStore.connected) {
    ElMessage.warning('请先连接 WebSocket')
    return
  }
  await chatStore.joinRoom(roomId)
  activeStep.value = 1
}

// 发送消息
function handleSendMessage() {
  if (!messageInput.value.trim()) return
  chatStore.sendMessage(messageInput.value)
  messageInput.value = ''
  activeStep.value = 2
  scrollToBottom()
}

// 输入指示器（节流）
let typingThrottle: ReturnType<typeof setTimeout> | null = null
function handleInput() {
  if (typingThrottle) return
  chatStore.sendTyping()
  typingThrottle = setTimeout(() => {
    typingThrottle = null
  }, 2000)
}

// 滚动到底部
function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// 上拉加载更多
function handleScroll() {
  if (messagesContainer.value && messagesContainer.value.scrollTop < 50) {
    chatStore.loadMoreHistory()
  }
}

// 创建聊天室
async function handleCreateRoom() {
  if (!createForm.name.trim()) {
    ElMessage.warning('请输入房间名称')
    return
  }
  creating.value = true
  try {
    const room = await chatStore.createRoom(createForm.name, createForm.description || '')
    ElMessage.success(`聊天室「${room.name}」创建成功`)
    createForm.name = ''
    createForm.description = ''
    activeTab.value = 'chat'
  } catch {
    ElMessage.error('创建聊天室失败')
  } finally {
    creating.value = false
  }
}

// 监听新消息自动滚动
watch(
  () => chatStore.currentMessages.length,
  () => scrollToBottom()
)

// 初始化：确保 WebSocket 已连接
onMounted(async () => {
  if (authStore.isAuthenticated && !chatStore.connected) {
    chatStore.connectWS()
  }
  await chatStore.loadRooms()
})
</script>
