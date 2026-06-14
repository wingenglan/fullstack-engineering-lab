import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '@/api/auth'
import type { User } from '@/types'
import { isTokenExpired } from '@/utils/jwt'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const displayName = computed(() => user.value?.nickname || user.value?.username || '')

  async function login(username: string, password: string) {
    loading.value = true
    try {
      const res = await authApi.login({ username, password })
      token.value = res.data.access_token
      localStorage.setItem('token', res.data.access_token)
      await fetchProfile()
    } finally {
      loading.value = false
    }
  }

  async function register(username: string, email: string, password: string) {
    loading.value = true
    try {
      await authApi.register({ username, email, password })
    } finally {
      loading.value = false
    }
  }

  async function fetchProfile() {
    try {
      const res = await authApi.getProfile()
      user.value = res.data
    } catch {
      user.value = null
      token.value = null
      localStorage.removeItem('token')
    }
  }

  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // ignore logout error
    } finally {
      user.value = null
      token.value = null
      localStorage.removeItem('token')
    }
  }

  async function initAuth() {
    const saved = localStorage.getItem('token')
    if (saved && !isTokenExpired(saved)) {
      token.value = saved
      await fetchProfile()
    } else {
      localStorage.removeItem('token')
      token.value = null
    }
  }

  return {
    user,
    token,
    loading,
    isAuthenticated,
    displayName,
    login,
    register,
    fetchProfile,
    logout,
    initAuth,
  }
})
