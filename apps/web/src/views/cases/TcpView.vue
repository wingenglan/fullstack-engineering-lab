<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { createChatSession, closeChatSession, sendChatMessage, listChatRooms, getChatRoomMessages, getChatSessionSSEUrl } from '@/api/tcp'
import GlowCard from '@/components/GlowCard.vue'
import type { TCPChatRoomInfo, TCPChatSessionEvent, TCPChatMessage } from '@/types'

// 聊天会话状态
const sessionId = ref('')
const nickname = ref('')
const room = ref('')
const joined = ref(false)

// UI 状态
const messages = ref<Array<{ type: string; nickname: string; content: string; time: string }>>([])
const inputText = ref('')
const connecting = ref(false)
const rooms = ref<TCPChatRoomInfo[]>([])

// 加入表单
const joinForm = ref({ nickname: '用户' + Math.floor(Math.random() * 1000), room: '大厅' })

let eventSource: EventSource | null = null
const messagesEndRef = ref<HTMLElement | null>(null)
let pollTimer: ReturnType<typeof setInterval> | null = null

async function handleJoin() {
  if (!joinForm.value.nickname.trim() || !joinForm.value.room.trim()) return

  connecting.value = true
  try {
    const { data } = await createChatSession({
      nickname: joinForm.value.nickname.trim(),
      room: joinForm.value.room.trim(),
    })
    sessionId.value = data.session_id
    nickname.value = joinForm.value.nickname.trim()
    room.value = joinForm.value.room.trim()
    joined.value = true

    // 加载历史消息
    await loadHistory()
    // 启动 SSE
    startSSE()
    // 开始轮询房间列表
    pollTimer = window.setInterval(loadRooms, 3000)
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '加入失败')
  } finally {
    connecting.value = false
  }
}

async function loadHistory() {
  try {
    const { data } = await getChatRoomMessages(room.value, 50)
    const msgs = data || []
    messages.value = msgs.map(m => ({
      type: m.type === 'message' ? 'message' : 'system',
      nickname: m.nickname,
      content: m.content,
      time: m.time,
    }))
  } catch { /* 忽略 */ }
}

async function loadRooms() {
  try {
    const { data } = await listChatRooms()
    rooms.value = data || []
  } catch { /* 忽略 */ }
}

function startSSE() {
  stopSSE()
  const url = getChatSessionSSEUrl(sessionId.value)
  eventSource = new EventSource(url)
  eventSource.addEventListener('chat_event', (e: MessageEvent) => {
    const evt: TCPChatSessionEvent = JSON.parse(e.data)
    if (evt.type === 'message' || evt.type === 'system') {
      messages.value.push({
        type: evt.type,
        nickname: evt.nickname || '',
        content: evt.content || '',
        time: evt.timestamp || '',
      })
      scrollToBottom()
    } else if (evt.type === 'disconnected') {
      ElMessage.warning('连接已断开')
      stopSSE()
    }
  })
  eventSource.onerror = () => {
    // 静默重连
  }
}

function stopSSE() {
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
}

async function handleSend() {
  const text = inputText.value.trim()
  if (!text || !sessionId.value) return

  inputText.value = ''
  try {
    await sendChatMessage(sessionId.value, text)
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '发送失败')
    inputText.value = text
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

async function handleLeave() {
  if (sessionId.value) {
    try {
      await closeChatSession(sessionId.value)
    } catch { /* 忽略 */ }
  }
  stopSSE()
  if (pollTimer) clearInterval(pollTimer)
  sessionId.value = ''
  joined.value = false
  messages.value = []
  rooms.value = []
}

function scrollToBottom() {
  nextTick(() => {
    messagesEndRef.value?.scrollIntoView({ behavior: 'smooth' })
  })
}

onMounted(() => {
  loadRooms()
  pollTimer = window.setInterval(loadRooms, 5000)
})

onUnmounted(() => {
  stopSSE()
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<template>
  <div class="tcp-container">
    <div class="page-header">
      <h2>TCP 聊天室</h2>
      <p class="subtitle">基于原始 TCP Socket 的多房间聊天系统，通过 SSE 实时推送消息</p>
    </div>

    <div class="content-layout">
      <!-- 左侧：房间列表 -->
      <div class="left-panel">
        <GlowCard>
          <template #header>
            <span class="panel-title">聊天室列表</span>
          </template>

          <div class="room-list">
            <div
              v-for="r in rooms"
              :key="r.name"
              class="room-item"
              :class="{ active: room === r.name }"
            >
              <span class="room-name"># {{ r.name }}</span>
              <span class="room-count">{{ r.user_count }} 人</span>
            </div>
            <el-empty v-if="rooms.length === 0" description="暂无活跃聊天室" :image-size="40" />
          </div>
        </GlowCard>
      </div>

      <!-- 右侧：聊天区 -->
      <div class="right-panel">
        <template v-if="!joined">
          <!-- 加入表单 -->
          <GlowCard>
            <template #header>
              <span class="panel-title">加入聊天室</span>
            </template>

            <div class="join-form">
              <div class="form-item">
                <label>昵称</label>
                <el-input v-model="joinForm.nickname" placeholder="输入你的昵称" maxlength="20" />
              </div>
              <div class="form-item">
                <label>聊天室</label>
                <el-input v-model="joinForm.room" placeholder="输入聊天室名称" maxlength="30" />
              </div>
              <el-button type="primary" :loading="connecting" @click="handleJoin" style="width: 100%">
                加入聊天
              </el-button>
            </div>
          </GlowCard>
        </template>

        <template v-else>
          <!-- 聊天区 -->
          <GlowCard class="chat-card">
            <template #header>
              <div class="chat-header">
                <div>
                  <span class="chat-room-name"># {{ room }}</span>
                  <span class="chat-nickname">以 {{ nickname }} 的身份</span>
                </div>
                <el-button type="danger" size="small" plain @click="handleLeave">离开</el-button>
              </div>
            </template>

            <!-- 消息区 -->
            <div class="messages-container">
              <div
                v-for="(msg, idx) in messages"
                :key="idx"
                class="message-item"
                :class="{
                  'msg-system': msg.type === 'system',
                  'msg-self': msg.nickname === nickname
                }"
              >
                <div v-if="msg.type === 'system'" class="system-msg">
                  {{ msg.content }}
                </div>
                <div v-else class="user-msg">
                  <span class="msg-nickname">{{ msg.nickname }}</span>
                  <span class="msg-content">{{ msg.content }}</span>
                </div>
              </div>
              <div ref="messagesEndRef" />
            </div>

            <!-- 输入区 -->
            <div class="input-area">
              <el-input
                v-model="inputText"
                placeholder="输入消息，回车发送..."
                @keydown="handleKeydown"
                maxlength="500"
              >
                <template #append>
                  <el-button @click="handleSend" :disabled="!inputText.trim()">发送</el-button>
                </template>
              </el-input>
            </div>
          </GlowCard>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tcp-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px 16px;
}

.page-header {
  margin-bottom: 24px;
}
.page-header h2 {
  margin: 0 0 4px 0;
  font-size: 22px;
  color: #e2e8f0;
}
.subtitle {
  color: #94a3b8;
  margin: 0;
  font-size: 14px;
}

.content-layout {
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 20px;
  align-items: start;
}

.panel-title {
  font-size: 14px;
  font-weight: 600;
  color: #e2e8f0;
}

/* 房间列表 */
.room-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.room-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  border-radius: 8px;
  transition: background 0.15s;
}
.room-item:hover {
  background: rgba(255,255,255,0.04);
}
.room-item.active {
  background: rgba(59, 130, 246, 0.1);
}
.room-name {
  color: #93c5fd;
  font-size: 14px;
}
.room-count {
  font-size: 12px;
  color: #64748b;
}

/* 加入表单 */
.join-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.form-item label {
  display: block;
  margin-bottom: 6px;
  font-size: 13px;
  color: #94a3b8;
}

/* 聊天卡 */
.chat-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.chat-card :deep(.card-body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0 !important;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}
.chat-room-name {
  color: #93c5fd;
  font-weight: 600;
}
.chat-nickname {
  font-size: 12px;
  color: #64748b;
  margin-left: 12px;
}

/* 消息区 */
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  min-height: 400px;
  max-height: 55vh;
}
.message-item {
  margin-bottom: 8px;
}
.msg-system {
  text-align: center;
  color: #64748b;
  font-size: 12px;
  padding: 4px 0;
}
.user-msg {
  display: flex;
  gap: 8px;
  align-items: baseline;
}
.msg-nickname {
  color: #60a5fa;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}
.msg-self .msg-nickname {
  color: #34d399;
}
.msg-content {
  color: #e2e8f0;
  font-size: 14px;
  word-break: break-word;
}

/* 输入区 */
.input-area {
  padding: 12px 16px;
  border-top: 1px solid rgba(255,255,255,0.06);
}

@media (max-width: 768px) {
  .content-layout {
    grid-template-columns: 1fr;
  }
}
</style>
