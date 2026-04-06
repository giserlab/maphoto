import api from './index'
import type { ApiResponse } from '@/types'

export interface StorageFile {
  name: string
  path: string
  size: number
  lastModified: string
  url: string
}

export interface UploadResult {
  success: StorageFile[]
  failed: string[]
  count: number
}

export function uploadFiles(folder: 'thumbs' | 'photos', files: File[]) {
  const formData = new FormData()
  formData.append('folder', folder)
  files.forEach(file => {
    formData.append('files', file)
  })
  return api.post<ApiResponse<UploadResult>>('/storage/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

export function getFiles(folder?: 'thumbs' | 'photos') {
  const params = folder ? { folder } : {}
  return api.get<ApiResponse<StorageFile[]>>('/storage/files', { params })
}

export function deleteFile(folder: string, filename: string) {
  return api.delete<ApiResponse>(`/storage/files/${folder}/${filename}`)
}

export function renameFile(folder: string, filename: string, newName: string) {
  const formData = new FormData()
  formData.append('newName', newName)
  return api.post<ApiResponse<StorageFile>>(`/storage/files/${folder}/${filename}/rename`, formData)
}
