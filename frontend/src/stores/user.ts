import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Config, ConfigUpdateForm } from '@/types'
import * as userApi from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const config = ref<Config | null>(null)

  async function fetchConfig() {
    const res = await userApi.getConfig()
    config.value = res.data
    return res.data
  }

  async function updateConfig(data: ConfigUpdateForm) {
    const res = await userApi.updateConfig(data)
    config.value = { ...config.value, ...res.data }
    return res.data
  }

  return {
    config,
    fetchConfig,
    updateConfig,
  }
})
