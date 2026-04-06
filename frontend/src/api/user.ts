import api from './index'
import type { ApiResponse, Config, ConfigUpdateForm } from '@/types'

export function getConfig() {
  return api.get<ApiResponse<Config>>('/user/config')
}

export function updateConfig(data: ConfigUpdateForm) {
  return api.post<ApiResponse<Config>>('/user/config/update', data)
}
