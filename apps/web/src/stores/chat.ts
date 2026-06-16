import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useWebSocket } from '@/composables/useWebSocket'
import * as chatApi from '@/api/chat'
import type {
  ChatRoom,
  ChatMessage,
  OnlineUser,
  WSMessage,
  WSHistoryPayload,
  WSNewMessagePayload,
  WSUserEventPayload,
  WSOnlineUsersPayload,
  WSUserTypingPayload,
  WSErrorPayload,
} from '@/types'
import { useAuthStore } from './auth'

export const useChatStore = defineStore('chat', () => {
  // WebSocket 实例（整个 store 生命周期只创建一次）
  const ws = useWebSocket()

  // 状态
  const rooms = ref<ChatRoom[]>([])
  const currentRoomId = ref<number | null>(null)
  const messages = ref<ChatMessage[]>([])
  const onlineUsers = ref<OnlineUser[]>([])
  const typingUsers = ref<Map<number, string>>(new Map())
  const loadingRooms = ref(false)
  const loadingMessages = ref(false)
  const hasMoreMessages = ref(true)
  const systemNotices = ref<string[]>([])

  // 计算属性
  const currentRoom = computed(() => rooms.value.find((r) => r.id === currentRoomId.value))
  const currentMessages = computed(() =>
    messages.value.filter((m) => m.room_id === currentRoomId.value)
  )
  const onlineCount = computed(() => onlineUsers.value.length)
  const connected = computed(() => ws.status.value === 'connected')

  // 输入指示器去抖动定时器
  const typingTimers = new Map<number, ReturnType<typeof setTimeout>>()

  // ========== WebSocket 消息处理 ==========

  function handleWSMessage(msg: WSMessage) {
    const payload = msg.payload as Record<string, unknown>

    switch (msg.type) {
      case 'room_history':
        handleRoomHistory(payload as unknown as WSHistoryPayload)
        break
      case 'new_message':
        handleNewMessage(payload as unknown as WSNewMessagePayload)
        break
      case 'user_joined':
        handleUserJoined(payload as unknown as WSUserEventPayload)
        break
      case 'user_left':
        handleUserLeft(payload as unknown as WSUserEventPayload)
        break
      case 'online_users':
        handleOnlineUsers(payload as unknown as WSOnlineUsersPayload)
        break
      case 'user_typing':
        handleUserTyping(payload as unknown as WSUserTypingPayload)
        break
      case 'error':
        console.error('WS 错误:', (payload as unknown as WSErrorPayload).message)
        break
      case 'pong':
        break
    }
  }

  function handleRoomHistory(payload: WSHistoryPayload) {
    if (payload.room_id !== currentRoomId.value) return
    // 避免重复：过滤掉已有的消息
    const existingIds = new Set(messages.value.map((m) => m.id))
    const newMsgs = payload.messages.filter((m) => !existingIds.has(m.id))
    messages.value = [...newMsgs, ...messages.value]
    if (payload.messages.length < 50) {
      hasMoreMessages.value = false
    }
  }

  function handleNewMessage(payload: WSNewMessagePayload) {
    // 避免重复
    if (messages.value.some((m) => m.id === payload.id)) return

    const msg: ChatMessage = {
      id: payload.id,
      room_id: payload.room_id,
      user_id: payload.user_id,
      username: payload.username,
      content: payload.content,
      msg_type: payload.msg_type,
      created_at: payload.created_at,
    }
    messages.value.push(msg)
    clearTypingUser(payload.user_id)
  }

  function handleUserJoined(payload: WSUserEventPayload) {
    if (payload.room_id === currentRoomId.value) {
      systemNotices.value.push(`${payload.username} 加入了聊天室`)
      setTimeout(() => systemNotices.value.shift(), 5000)
    }
    refreshOnlineUsers()
  }

  function handleUserLeft(payload: WSUserEventPayload) {
    if (payload.room_id === currentRoomId.value) {
      systemNotices.value.push(`${payload.username} 离开了聊天室`)
      setTimeout(() => systemNotices.value.shift(), 5000)
    }
    refreshOnlineUsers()
  }

  function handleOnlineUsers(payload: WSOnlineUsersPayload) {
    if (payload.room_id === currentRoomId.value) {
      onlineUsers.value = payload.users
    }
  }

  function handleUserTyping(payload: WSUserTypingPayload) {
    if (payload.room_id !== currentRoomId.value) return
    typingUsers.value.set(payload.user_id, payload.username)

    const existingTimer = typingTimers.get(payload.user_id)
    if (existingTimer) clearTimeout(existingTimer)

    const timer = setTimeout(() => clearTypingUser(payload.user_id), 3000)
    typingTimers.set(payload.user_id, timer)
  }

  function clearTypingUser(userId: number) {
    typingUsers.value.delete(userId)
    const timer = typingTimers.get(userId)
    if (timer) {
      clearTimeout(timer)
      typingTimers.delete(userId)
    }
  }

  // 注册消息处理器（只执行一次）
  ws.onMessage(handleWSMessage)

  // ========== 公开方法 ==========

  async function loadRooms() {
    loadingRooms.value = true
    try {
      const res = await chatApi.getRooms()
      rooms.value = res.data
    } catch {
      console.error('加载聊天室列表失败')
    } finally {
      loadingRooms.value = false
    }
  }

  async function joinRoom(roomId: number) {
    // 离开当前房间
    if (currentRoomId.value !== null && currentRoomId.value !== roomId) {
      leaveCurrentRoom()
    }

    currentRoomId.value = roomId
    messages.value = []
    onlineUsers.value = []
    systemNotices.value = []
    hasMoreMessages.value = true

    // 通过 REST API 加载历史消息
    loadingMessages.value = true
    try {
      const res = await chatApi.getMessageHistory(roomId, 50)
      messages.value = res.data.messages
      hasMoreMessages.value = res.data.has_more
    } catch {
      console.error('加载消息历史失败')
    } finally {
      loadingMessages.value = false
    }

    // 通过 WebSocket 加入房间
    ws.send('join_room', { room_id: roomId })

    // 刷新在线用户
    refreshOnlineUsers()
  }

  function leaveCurrentRoom() {
    if (currentRoomId.value !== null) {
      ws.send('leave_room', { room_id: currentRoomId.value })
      currentRoomId.value = null
      messages.value = []
      onlineUsers.value = []
    }
  }

  async function loadMoreHistory() {
    if (!currentRoomId.value || !hasMoreMessages.value || loadingMessages.value) return

    const oldestMsg = currentMessages.value[0]
    if (!oldestMsg) return

    loadingMessages.value = true
    try {
      const res = await chatApi.getMessageHistory(currentRoomId.value, 50, oldestMsg.id)
      const existingIds = new Set(messages.value.map((m) => m.id))
      const newMsgs = res.data.messages.filter((m) => !existingIds.has(m.id))
      messages.value = [...newMsgs, ...messages.value]
      hasMoreMessages.value = res.data.has_more
    } catch {
      console.error('加载更多历史消息失败')
    } finally {
      loadingMessages.value = false
    }
  }

  function sendMessage(content: string) {
    if (!currentRoomId.value || !content.trim()) return
    ws.send('send_message', {
      room_id: currentRoomId.value,
      content: content.trim(),
      msg_type: 1,
    })
  }

  function sendTyping() {
    if (!currentRoomId.value) return
    ws.send('typing', { room_id: currentRoomId.value })
  }

  async function refreshOnlineUsers() {
    if (!currentRoomId.value) return
    try {
      const res = await chatApi.getOnlineUsers(currentRoomId.value)
      onlineUsers.value = res.data
    } catch {
      // 静默失败
    }
  }

  function connectWS() {
    const authStore = useAuthStore()
    if (!authStore.token) return
    ws.connect(authStore.token)
    loadRooms()
  }

  function disconnectWS() {
    leaveCurrentRoom()
    ws.disconnect()
    rooms.value = []
  }

  async function createRoom(name: string, description: string, type = 1) {
    const res = await chatApi.createRoom({ name, description, type })
    rooms.value.push(res.data)
    return res.data
  }

  return {
    // 状态
    rooms,
    currentRoomId,
    currentRoom,
    messages,
    currentMessages,
    onlineUsers,
    onlineCount,
    typingUsers,
    connected,
    loadingRooms,
    loadingMessages,
    hasMoreMessages,
    systemNotices,
    status: ws.status,

    // 方法
    loadRooms,
    joinRoom,
    leaveCurrentRoom,
    loadMoreHistory,
    sendMessage,
    sendTyping,
    refreshOnlineUsers,
    connectWS,
    disconnectWS,
    createRoom,
  }
})
