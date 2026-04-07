import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import type { Conversation } from '../types/conversation'
import { conversationsApi } from '../api/conversations'
import { USE_MOCK } from '../api/client'
import { MOCK_CONVERSATIONS } from '../mocks/conversations'

interface ConversationsState {
  items: Conversation[]
  loading: boolean
}

const initialState: ConversationsState = { items: [], loading: false }

export const fetchConversations = createAsyncThunk('conversations/fetch', async () => {
  if (USE_MOCK) {
    await new Promise((r) => setTimeout(r, 300))
    return MOCK_CONVERSATIONS
  }
  return conversationsApi.list()
})

export const createConversation = createAsyncThunk(
  'conversations/create',
  async (userId: string) => {
    if (USE_MOCK) {
      await new Promise((r) => setTimeout(r, 300))
      return null
    }
    return conversationsApi.create(userId)
  }
)

const conversationsSlice = createSlice({
  name: 'conversations',
  initialState,
  reducers: {
    updateLastMessage(state, action) {
      const { conversationId, message } = action.payload
      const conv = state.items.find((c) => c.id === conversationId)
      if (conv) {
        conv.lastMessage = {
          body: message.body,
          senderId: message.senderId,
          createdAt: message.createdAt,
        }
      }
    },
    prependConversation(state, action) {
      const exists = state.items.find((c) => c.id === action.payload.id)
      if (!exists) state.items.unshift(action.payload)
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchConversations.pending, (state) => { state.loading = true })
      .addCase(fetchConversations.fulfilled, (state, action) => {
        state.loading = false
        state.items = action.payload
      })
      .addCase(fetchConversations.rejected, (state) => { state.loading = false })
  },
})

export const { updateLastMessage, prependConversation } = conversationsSlice.actions
export default conversationsSlice.reducer
