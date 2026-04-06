import api from './index'
import type { ApiResponse, Place, PlaceAddForm, PlaceUpdateForm } from '@/types'

export function getAll() {
  return api.get<ApiResponse<Place[]>>('/place/all')
}

export function create(data: PlaceAddForm) {
  return api.post<ApiResponse<Place>>('/place/add', data)
}

export function update(id: number, data: PlaceUpdateForm) {
  return api.post<ApiResponse<Place>>(`/place/update/${id}`, data)
}

export function remove(id: number) {
  return api.get<ApiResponse>(`/place/del/${id}`)
}

export function updateCover(id: number, url: string) {
  return api.post<ApiResponse>('/place/cover', { id, url })
}

export function addPhoto(id: number, url: string) {
  return api.post<ApiResponse>('/place/pic/add', { id, url })
}

export function removePhoto(id: number, url: string) {
  return api.post<ApiResponse>('/place/pic/del', { id, url })
}
