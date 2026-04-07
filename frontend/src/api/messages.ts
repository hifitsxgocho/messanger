import type { Message } from '../types/message'
import { apiClient } from './client'

export const messagesApi = {
  list: (conversationId: string, after?: string) =>
    apiClient
      .get<Message[]>(`/conversations/${conversationId}/messages`, {
        params: after ? { after } : undefined,
      })
      .then((r) => r.data),

  send: (conversationId: string, body: string) =>
    apiClient
      .post<Message>(`/conversations/${conversationId}/messages`, { body })
      .then((r) => r.data),

  markRead: (conversationId: string, messageId: string) =>
    apiClient
      .put(`/conversations/${conversationId}/messages/${messageId}/read`)
      .then((r) => r.data),
}
