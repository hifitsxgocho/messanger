import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import type { Message } from '../types/message'
import { messagesApi } from '../api/messages'
import { USE_MOCK } from '../api/client'
import { MOCK_MESSAGES } from '../mocks/messages'

interface MessagesState {
  byConversation: Record<string, Message[]>
  loading: boolean
}

const initialState: MessagesState = { byConversation: {}, loading: false }

export const fetchMessages = createAsyncThunk(
  'messages/fetch',
  async ({ conversationId, after }: { conversationId: string; after?: string }) => {
    if (USE_MOCK) {
      await new Promise((r) => setTimeout(r, 200))
      return { conversationId, messages: MOCK_MESSAGES[conversationId] ?? [] }
    }
    const messages = await messagesApi.list(conversationId, after)
    return { conversationId, messages }
  }
)

export const sendMessage = createAsyncThunk(
  'messages/send',
  async ({ conversationId, body }: { conversationId: string; body: string }) => {
    if (USE_MOCK) {
      await new Promise((r) => setTimeout(r, 150))
      const msg: Message = {
        id: `mock-${Date.now()}`,
        conversationId,
        senderId: 'user-1',
        body,
        createdAt: new Date().toISOString(),
        readAt: null,
      }
      return { conversationId, message: msg }
    }
    const message = await messagesApi.send(conversationId, body)
    return { conversationId, message }
  }
)

const messagesSlice = createSlice({
  name: 'messages',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchMessages.pending, (state) => { state.loading = true })
      .addCase(fetchMessages.fulfilled, (state, action) => {
        state.loading = false
        const { conversationId, messages } = action.payload
        if (!action.meta.arg.after) {
          state.byConversation[conversationId] = messages
        } else {
          const existing = state.byConversation[conversationId] ?? []
          const existingIds = new Set(existing.map((m) => m.id))
          const newMsgs = messages.filter((m) => !existingIds.has(m.id))
          state.byConversation[conversationId] = [...existing, ...newMsgs]
        }
      })
      .addCase(fetchMessages.rejected, (state) => { state.loading = false })
      .addCase(sendMessage.fulfilled, (state, action) => {
        const { conversationId, message } = action.payload
        if (!state.byConversation[conversationId]) {
          state.byConversation[conversationId] = []
        }
        state.byConversation[conversationId].push(message)
      })
  },
})

export default messagesSlice.reducer
