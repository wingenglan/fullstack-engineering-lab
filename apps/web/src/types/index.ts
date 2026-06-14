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
