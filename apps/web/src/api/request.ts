import axios, { type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import type { RequestLogEntry } from '@/types'

let logId = 0
const requestTimestamps = new Map<number, number>()

export const requestLogs: RequestLogEntry[] = []

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // Log request
    const id = ++logId
    requestTimestamps.set(id, Date.now())
    ;(config as InternalAxiosRequestConfig & { _logId?: number })._logId = id

    requestLogs.push({
      id,
      timestamp: new Date().toISOString(),
      method: (config.method || 'GET').toUpperCase(),
      url: config.url || '',
      status: null,
      statusText: 'pending',
      requestData: config.data ? JSON.stringify(config.data) : '',
      responseData: '',
      duration: null,
    })

    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const config = response.config as InternalAxiosRequestConfig & { _logId?: number }
    const logId = config._logId
    if (logId) {
      const startTime = requestTimestamps.get(logId)
      const logEntry = requestLogs.find((l) => l.id === logId)
      if (logEntry && startTime) {
        logEntry.status = response.status
        logEntry.statusText = 'OK'
        logEntry.responseData = JSON.stringify(response.data)
        logEntry.duration = Date.now() - startTime
        requestTimestamps.delete(logId)
      }
    }

    const data = response.data
    if (data.code !== 0) {
      return Promise.reject(new Error(data.message || 'Request failed'))
    }
    return data
  },
  (error) => {
    if (error.response) {
      const config = error.response.config as InternalAxiosRequestConfig & { _logId?: number }
      const logId = config?._logId
      if (logId) {
        const startTime = requestTimestamps.get(logId)
        const logEntry = requestLogs.find((l) => l.id === logId)
        if (logEntry && startTime) {
          logEntry.status = error.response.status
          logEntry.statusText = error.response.statusText || 'Error'
          logEntry.responseData = JSON.stringify(error.response.data)
          logEntry.duration = Date.now() - startTime
          requestTimestamps.delete(logId)
        }
      }

      if (error.response.status === 401) {
        localStorage.removeItem('token')
      }
    }
    return Promise.reject(error)
  }
)

export default request
