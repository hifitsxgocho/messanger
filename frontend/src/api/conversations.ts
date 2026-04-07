import type { Conversation } from '../types/conversation'
import { apiClient } from './client'

export const conversationsApi = {
  list: () => apiClient.get<Conversation[]>('/conversations').then((r) => r.data),

  create: (userId: string) =>
    apiClient.post<Conversation>('/conversations', { userId }).then((r) => r.data),

  getById: (id: string) => apiClient.get<Conversation>(`/conversations/${id}`).then((r) => r.data),
}
