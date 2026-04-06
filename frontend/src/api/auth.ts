import api from './index'
import type { ApiResponse, LoginResult, LoginForm } from '@/types'

export function login(params: LoginForm) {
  return api.post<ApiResponse<LoginResult>>('/user/login', params)
}

export function logout() {
  return api.get<ApiResponse>('/user/logout')
}
