import request from './request'
import type { ApiResponse, LoginRequest, LoginResponse, RegisterRequest, User } from '@/types'

export function register(data: RegisterRequest): Promise<ApiResponse<null>> {
  return request.post('/auth/register', data)
}

export function login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
  return request.post('/auth/login', data)
}

export function getProfile(): Promise<ApiResponse<User>> {
  return request.get('/auth/profile')
}

export function logout(): Promise<ApiResponse<null>> {
  return request.post('/auth/logout')
}

export function healthCheck(): Promise<ApiResponse<{ status: string; db: string; redis: string }>> {
  return request.get('/health')
}
