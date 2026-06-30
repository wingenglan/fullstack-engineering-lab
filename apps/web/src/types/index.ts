export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface User {
  id: number
  username: string
  email: string
  nickname: string
  status: number
  created_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginResponse {
  access_token: string
  expires_in: number
}

export interface RequestLogEntry {
  id: number
  timestamp: string
  method: string
  url: string
  status: number | null
  statusText: string
  requestData: string
  responseData: string
  duration: number | null
}

export interface CaseItem {
  id: string
  title: string
  description: string
  tags: string[]
  difficulty: 'easy' | 'medium' | 'hard'
  status: 'available' | 'coming-soon'
  category: string
  icon: string
  to: string
}

// ===== Redis Lock Types =====

export interface LockAcquireRequest {
  resource: string
  ttl: number
}

export interface LockAcquireResponse {
  resource: string
  owner: string
  ttl: number
}

export interface LockReleaseRequest {
  resource: string
}

export interface LockStatusRequest {
  resource: string
}

export interface LockStatusResponse {
  resource: string
  locked: boolean
  owner: string
  ttl_ms: number
}

export interface ContentionDemoRequest {
  resource: string
  ttl: number
  goroutines: number
  hold_ms: number
}

export interface ContentionResult {
  goroutine_id: number
  acquired: boolean
  wait_ms: number
  message: string
}

export interface ContentionDemoResponse {
  resource: string
  results: ContentionResult[]
  summary: {
    total: number
    succeeded: number
    failed: number
  }
}

// ===== WebSocket Chat Types =====

export interface ChatRoom {
  id: number
  name: string
  description: string
  type: number // 1=群聊, 2=私聊
  creator_id: number
  member_count: number
  status: number
  created_at: string
}

export interface ChatMessage {
  id: number
  room_id: number
  user_id: number
  username: string
  content: string
  msg_type: number // 1=文本, 2=图片, 3=系统
  created_at: string
}

export interface OnlineUser {
  user_id: number
  username: string
}

export interface MessageHistoryResponse {
  messages: ChatMessage[]
  has_more: boolean
}

export interface CreateRoomRequest {
  name: string
  description: string
  type: number
}

// WebSocket 消息协议
export interface WSMessage<T = unknown> {
  type: string
  payload: T
}

export interface WSHistoryPayload {
  room_id: number
  messages: ChatMessage[]
}

export interface WSNewMessagePayload {
  id: number
  room_id: number
  user_id: number
  username: string
  content: string
  msg_type: number
  created_at: string
}

export interface WSUserEventPayload {
  room_id: number
  user_id: number
  username: string
}

export interface WSOnlineUsersPayload {
  room_id: number
  users: OnlineUser[]
  count: number
}

export interface WSUserTypingPayload {
  room_id: number
  user_id: number
  username: string
}

export interface WSErrorPayload {
  code: number
  message: string
}

// ===== Redis 数据类型 Types =====

// ----- Hash: 用户画像缓存 -----

export interface HashFieldUpdateRequest {
  key: string
  field: string
  value: string
}

export interface HashMultiSetRequest {
  key: string
  fields: Record<string, string>
}

export interface HashProfileResponse {
  key: string
  fields: Record<string, string>
  num_of_fld: number
}

// ----- Set: 标签/收藏夹管理 -----

export interface SetAddMembersRequest {
  key: string
  members: string[]
}

export interface SetRemoveMembersRequest {
  key: string
  members: string[]
}

export interface SetListResponse {
  key: string
  members: string[]
  count: number
}

export interface SetOperationRequest {
  keys: string[]
}

export interface SetOperationResponse {
  keys: string[]
  op: 'intersect' | 'union' | 'diff'
  members: string[]
  count: number
}

// ----- ZSet: 实时排行榜 -----

export interface ZSetAddScoreRequest {
  key: string
  member: string
  score: number
}

export interface ZSetRankRequest {
  key: string
  member: string
}

export interface ZSetTopNRequest {
  key: string
  n: number
}

export interface ZSetMemberResponse {
  key: string
  member: string
  score: number
  rank: number
}

export interface ZSetRankEntry {
  rank: number
  member: string
  score: number
}

export interface ZSetTopNResponse {
  key: string
  total: number
  members: ZSetRankEntry[]
}

// ----- List: 最新活动流 / 简易消息队列 -----

export interface ListPushRequest {
  key: string
  value: string
  pos: 'left' | 'right'
}

export interface ListPopRequest {
  key: string
  pos: 'left' | 'right'
}

export interface ListRangeRequest {
  key: string
  start: number
  stop: number
}

export interface ListPopResponse {
  key: string
  value: string
  pos: string
}

export interface ListRangeResponse {
  key: string
  values: string[]
  length: number
}

// ----- String: 验证码存储 / 计数器 -----

export interface StringSetRequest {
  key: string
  value: string
  ttl: number
}

export interface StringSetResponse {
  key: string
  value: string
  ttl: number
}

export interface StringGetResponse {
  key: string
  value: string
  ttl: number
  exists: boolean
}

export interface StringIncrRequest {
  key: string
  delta: number
}

export interface StringIncrResponse {
  key: string
  before: number
  after: number
  delta: number
}

// ============================================
// MQTT IoT 设备管理演示类型
// ============================================

export interface MQTTDeviceInfo {
  device_id: string
  type: string
  name: string
  online: boolean
  last_seen: string
  properties: Record<string, any>
}

export interface MQTTSSEEvent {
  type: string // device_online / device_offline / telemetry / command_response
  device_id?: string
  topic?: string
  payload?: string
  timestamp: string
}

export interface AddDeviceRequest {
  type: string
  name: string
}

export interface DeviceCommandRequest {
  command: string
  params?: Record<string, any>
}

// ============================================
// TCP 基础协议演示类型
// ============================================

export interface TCPSessionResponse {
  session_id: string
  created_at: string
}

export interface TCPSessionInfo {
  session_id: string
  remote_addr: string
  created_at: string
  last_act_at: string
  duration_sec: number
  command_count: number
  is_alive: boolean
}

export interface TCPCommandResponse {
  session_id: string
  command: string
  response: string
  duration_ms: number
  timestamp: string
}

export interface TCPStatsResponse {
  server_addr: string
  active_sessions: number
  max_sessions: number
  uptime_sec: number
}

// ============================================
// TCP 聊天室演示类型
// ============================================

export interface TCPChatSessionRequest {
  nickname: string
  room: string
}

export interface TCPChatMessageRequest {
  content: string
}

export interface TCPChatSessionEvent {
  type: string // connected / message / system / disconnected
  session_id: string
  nickname?: string
  room?: string
  content?: string
  timestamp: string
}

export interface TCPChatRoomInfo {
  name: string
  user_count: number
  users?: string[]
  created_at: string
}

export interface TCPChatMessage {
  room: string
  nickname: string
  content: string
  type: string // message / system_join / system_leave
  time: string
}
