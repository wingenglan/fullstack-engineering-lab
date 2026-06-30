import request from './request'
import type { ApiResponse, MQTTDeviceInfo, AddDeviceRequest, DeviceCommandRequest } from '@/types'

// 设备管理
export const getDevices = (): Promise<ApiResponse<MQTTDeviceInfo[]>> =>
  request.get('/mqtt/devices')

export const getDevice = (id: string): Promise<ApiResponse<MQTTDeviceInfo>> =>
  request.get(`/mqtt/devices/${id}`)

export const addDevice = (data: AddDeviceRequest): Promise<ApiResponse<MQTTDeviceInfo>> =>
  request.post('/mqtt/devices', data)

export const removeDevice = (id: string): Promise<ApiResponse<{ device_id: string; removed: boolean }>> =>
  request.delete(`/mqtt/devices/${id}`)

export const sendDeviceCommand = (id: string, data: DeviceCommandRequest): Promise<ApiResponse<{ device_id: string; command: string; timestamp: string }>> =>
  request.post(`/mqtt/devices/${id}/commands`, data)

// SSE
export const getGlobalSSEUrl = () => `${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/mqtt/subscribe`

export const getDeviceSSEUrl = (deviceId: string) => `${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/mqtt/devices/${deviceId}/subscribe`

// 历史消息
export const getMessages = (limit?: number): Promise<ApiResponse<any[]>> =>
  request.get('/mqtt/messages', { params: { limit } })
