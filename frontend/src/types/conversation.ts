import type { UserPublic } from './user'

export interface Conversation {
  id: string
  partner: UserPublic
  lastMessage: LastMessage | null
  unreadCount: number
  createdAt: string
}

export interface LastMessage {
  body: string
  senderId: string
  createdAt: string
}
