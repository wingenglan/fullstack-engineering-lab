import request from './request'
import type {
  ApiResponse,
  ChatRoom,
  MessageHistoryResponse,
  CreateRoomRequest,
  OnlineUser,
} from '@/types'

// 获取所有聊天室
export function getRooms(): Promise<ApiResponse<ChatRoom[]>> {
  return request.get('/chat/rooms')
}

// 获取聊天室详情
export function getRoomInfo(id: number): Promise<ApiResponse<ChatRoom>> {
  return request.get(`/chat/rooms/${id}`)
}

// 创建聊天室
export function createRoom(data: CreateRoomRequest): Promise<ApiResponse<ChatRoom>> {
  return request.post('/chat/rooms', data)
}

// 获取消息历史
export function getMessageHistory(
  roomId: number,
  limit = 50,
  before?: number
): Promise<ApiResponse<MessageHistoryResponse>> {
  const params: Record<string, string | number> = { room_id: roomId, limit }
  if (before) params.before = before
  return request.get('/chat/messages', { params })
}

// 获取房间在线用户
export function getOnlineUsers(roomId: number): Promise<ApiResponse<OnlineUser[]>> {
  return request.get(`/chat/rooms/${roomId}/online`)
}
