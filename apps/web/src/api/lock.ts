import request from './request'
import type {
  ApiResponse,
  LockAcquireRequest,
  LockAcquireResponse,
  LockReleaseRequest,
  LockStatusRequest,
  LockStatusResponse,
  ContentionDemoRequest,
  ContentionDemoResponse,
} from '@/types'

export function acquireLock(data: LockAcquireRequest): Promise<ApiResponse<LockAcquireResponse>> {
  return request.post('/lock/acquire', data)
}

export function releaseLock(data: LockReleaseRequest): Promise<ApiResponse<null>> {
  return request.post('/lock/release', data)
}

export function getLockStatus(data: LockStatusRequest): Promise<ApiResponse<LockStatusResponse>> {
  return request.post('/lock/status', data)
}

export function contentionDemo(data: ContentionDemoRequest): Promise<ApiResponse<ContentionDemoResponse>> {
  return request.post('/lock/contention', data)
}
