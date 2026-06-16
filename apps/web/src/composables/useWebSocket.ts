import { ref } from 'vue'
import type { WSMessage } from '@/types'

// WebSocket 连接状态
export type WSStatus = 'connecting' | 'connected' | 'disconnected' | 'error'

// 消息处理器类型
type MessageHandler = (msg: WSMessage) => void

/**
 * WebSocket 连接 composable（单例模式）
 * 封装连接管理、自动重连、心跳检测
 * 由 Pinia store 调用一次，不要在组件中重复调用
 */
export function useWebSocket() {
  const status = ref<WSStatus>('disconnected')
  const reconnectCount = ref(0)

  let ws: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let heartbeatTimer: ReturnType<typeof setInterval> | null = null
  const messageHandlers: MessageHandler[] = []
  let currentToken = ''

  // 重连配置
  const maxReconnectAttempts = 10
  const baseReconnectDelay = 1000
  const maxReconnectDelay = 30000
  const heartbeatInterval = 30000

  function getReconnectDelay(): number {
    const delay = Math.min(
      baseReconnectDelay * Math.pow(2, reconnectCount.value),
      maxReconnectDelay
    )
    return delay + Math.random() * 1000
  }

  function connect(token: string) {
    // 已连接则跳过
    if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
      return
    }

    currentToken = token
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const wsUrl = `${protocol}//${host}/api/v1/chat/ws?token=${encodeURIComponent(token)}`

    status.value = 'connecting'

    try {
      ws = new WebSocket(wsUrl)
    } catch {
      status.value = 'error'
      scheduleReconnect()
      return
    }

    ws.onopen = () => {
      status.value = 'connected'
      reconnectCount.value = 0
      startHeartbeat()
    }

    ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        for (const handler of messageHandlers) {
          handler(msg)
        }
      } catch {
        console.warn('WebSocket 消息解析失败:', event.data)
      }
    }

    ws.onclose = (event) => {
      stopHeartbeat()
      status.value = 'disconnected'
      ws = null

      if (event.code !== 1000 && currentToken) {
        scheduleReconnect()
      }
    }

    ws.onerror = () => {
      status.value = 'error'
    }
  }

  function send(type: string, payload: unknown): boolean {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
      console.warn('WebSocket 未连接，无法发送消息')
      return false
    }
    ws.send(JSON.stringify({ type, payload }))
    return true
  }

  function scheduleReconnect() {
    if (reconnectCount.value >= maxReconnectAttempts) {
      status.value = 'error'
      return
    }
    if (reconnectTimer) clearTimeout(reconnectTimer)

    const delay = getReconnectDelay()
    reconnectCount.value++
    reconnectTimer = setTimeout(() => connect(currentToken), delay)
  }

  function startHeartbeat() {
    stopHeartbeat()
    heartbeatTimer = setInterval(() => send('ping', {}), heartbeatInterval)
  }

  function stopHeartbeat() {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
  }

  function onMessage(handler: MessageHandler) {
    messageHandlers.push(handler)
  }

  function offMessage(handler: MessageHandler) {
    const idx = messageHandlers.indexOf(handler)
    if (idx >= 0) messageHandlers.splice(idx, 1)
  }

  function disconnect() {
    stopHeartbeat()
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    currentToken = ''
    reconnectCount.value = 0

    if (ws) {
      ws.close(1000)
      ws = null
    }
    status.value = 'disconnected'
  }

  return {
    status,
    reconnectCount,
    connect,
    disconnect,
    send,
    onMessage,
    offMessage,
  }
}
