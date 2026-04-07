import type { User, UserPublic } from '../types/user'
import { apiClient } from './client'

export const usersApi = {
  getMe: () => apiClient.get<User>('/users/me').then((r) => r.data),

  updateMe: (data: Partial<Pick<User, 'username' | 'bio'>>) =>
    apiClient.put<User>('/users/me', data).then((r) => r.data),

  uploadAvatar: (file: File) => {
    const form = new FormData()
    form.append('avatar', file)
    return apiClient
      .post<{ avatarUrl: string }>('/users/me/avatar', form, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      .then((r) => r.data)
  },

  search: (q: string) =>
    apiClient.get<UserPublic[]>('/users/search', { params: { q, limit: 20 } }).then((r) => r.data),

  getById: (id: string) => apiClient.get<UserPublic>(`/users/${id}`).then((r) => r.data),
}
