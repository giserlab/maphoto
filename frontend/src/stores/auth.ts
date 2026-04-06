import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginResult, LoginForm } from '@/types'
import * as authApi from '@/api/auth'
import { storage } from '@/utils/storage'
import { STORAGE_KEYS } from '@/utils/constants'

export const useAuthStore = defineStore('auth', () => {
  // State
  const token = ref<string>(storage.get(STORAGE_KEYS.TOKEN) || '')
  const user = ref<User | null>(storage.get<User>(STORAGE_KEYS.USER))

  // Getters
  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.admin || false)
  const username = computed(() => user.value?.username || '')

  // Actions
  async function login(form: LoginForm) {
    const res = await authApi.login(form)
    const data = res.data as LoginResult

    token.value = data.token
    user.value = {
      username: data.username,
      admin: data.admin,
    } as User

    storage.set(STORAGE_KEYS.TOKEN, data.token)
    storage.set(STORAGE_KEYS.USER, user.value)

    return data
  }

  async function logout() {
    try {
      await authApi.logout()
    } finally {
      token.value = ''
      user.value = null
      storage.remove(STORAGE_KEYS.TOKEN)
      storage.remove(STORAGE_KEYS.USER)
    }
  }

  return {
    token,
    user,
    isAuthenticated,
    isAdmin,
    username,
    login,
    logout,
  }
})
