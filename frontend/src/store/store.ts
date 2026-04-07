import { configureStore } from '@reduxjs/toolkit'
import authReducer from './authSlice'
import conversationsReducer from './conversationsSlice'
import messagesReducer from './messagesSlice'

export const store = configureStore({
  reducer: {
    auth: authReducer,
    conversations: conversationsReducer,
    messages: messagesReducer,
  },
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
