import api from './index'
import type { ApiResponse, User } from '@/types'

export function getAll() {
  return api.get<ApiResponse<User[]>>('/users')
}

export function create(data: { username: string; password: string; admin: boolean }) {
  return api.post<ApiResponse<User>>('/user/add', data)
}

export function update(id: number, data: Partial<User>) {
  return api.post<ApiResponse<User>>(`/user/update/${id}`, data)
}

export function remove(id: number) {
  return api.get<ApiResponse>(`/user/del/${id}`)
}
