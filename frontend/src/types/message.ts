export interface Message {
  id: string
  conversationId: string
  senderId: string
  body: string
  createdAt: string
  readAt: string | null
}
