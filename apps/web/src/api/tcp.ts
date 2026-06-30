import request from './request'
import type { ApiResponse, TCPSessionResponse, TCPSessionInfo, TCPCommandResponse, TCPStatsResponse, TCPChatSessionRequest, TCPChatMessageRequest, TCPChatRoomInfo, TCPChatMessage } from '@/types'

// --- 基础协议 ---

export const createSession = (): Promise<ApiResponse<TCPSessionResponse>> =>
  request.post('/tcp/sessions')

export const closeSession = (id: string): Promise<ApiResponse<{ session_id: string; closed: boolean }>> =>
  request.delete(`/tcp/sessions/${id}`)

export const listSessions = (): Promise<ApiResponse<TCPSessionInfo[]>> =>
  request.get('/tcp/sessions')

export const sendCommand = (id: string, command: string): Promise<ApiResponse<TCPCommandResponse>> =>
  request.post(`/tcp/sessions/${id}/send`, { command })

export const getStats = (): Promise<ApiResponse<TCPStatsResponse>> =>
  request.get('/tcp/stats')

export const getSessionSSEUrl = (sessionId: string) => `${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/tcp/sessions/${sessionId}/stream`

// --- 聊天室 ---

export const createChatSession = (data: TCPChatSessionRequest): Promise<ApiResponse<{ session_id: string; nickname: string; room: string; created_at: string }>> =>
  request.post('/tcp/chat/sessions', data)

export const closeChatSession = (id: string): Promise<ApiResponse<{ session_id: string; closed: boolean }>> =>
  request.delete(`/tcp/chat/sessions/${id}`)

export const sendChatMessage = (id: string, content: string): Promise<ApiResponse<{ session_id: string; sent: boolean; timestamp: string }>> =>
  request.post(`/tcp/chat/sessions/${id}/msg`, { content })

export const listChatRooms = (): Promise<ApiResponse<TCPChatRoomInfo[]>> =>
  request.get('/tcp/chat/rooms')

export const getChatRoomInfo = (name: string): Promise<ApiResponse<TCPChatRoomInfo>> =>
  request.get(`/tcp/chat/rooms/${encodeURIComponent(name)}`)

export const getChatRoomMessages = (name: string, limit?: number): Promise<ApiResponse<TCPChatMessage[]>> =>
  request.get(`/tcp/chat/rooms/${encodeURIComponent(name)}/messages`, { params: { limit } })

export const getChatSessionSSEUrl = (sessionId: string) => `${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/tcp/chat/sessions/${sessionId}/stream`
