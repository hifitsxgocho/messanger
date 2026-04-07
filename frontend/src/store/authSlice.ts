import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import type { PayloadAction } from '@reduxjs/toolkit'
import type { User } from '../types/user'
import type { LoginRequest, RegisterRequest } from '../types/auth'
import { authApi } from '../api/auth'
import { usersApi } from '../api/users'
import { storage } from '../utils/storage'
import { USE_MOCK } from '../api/client'
import { MOCK_CURRENT_USER } from '../mocks/users'

interface AuthState {
  user: User | null
  token: string | null
  loading: boolean
  error: string | null
}

const initialState: AuthState = {
  user: null,
  token: storage.getToken(),
  loading: false,
  error: null,
}

export const fetchMe = createAsyncThunk('auth/fetchMe', async () => {
  if (USE_MOCK) return MOCK_CURRENT_USER
  return usersApi.getMe()
})

export const login = createAsyncThunk('auth/login', async (data: LoginRequest) => {
  if (USE_MOCK) {
    await new Promise((r) => setTimeout(r, 500))
    const token = 'mock-jwt-token'
    storage.setToken(token)
    return { token, user: MOCK_CURRENT_USER }
  }
  const res = await authApi.login(data)
  storage.setToken(res.token)
  return res
})

export const register = createAsyncThunk('auth/register', async (data: RegisterRequest) => {
  if (USE_MOCK) {
    await new Promise((r) => setTimeout(r, 500))
    const token = 'mock-jwt-token'
    storage.setToken(token)
    return { token, user: { ...MOCK_CURRENT_USER, email: data.email, username: data.username } }
  }
  const res = await authApi.register(data)
  storage.setToken(res.token)
  return res
})

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    logout(state) {
      state.user = null
      state.token = null
      storage.removeToken()
    },
    setUser(state, action: PayloadAction<User>) {
      state.user = action.payload
    },
    clearError(state) {
      state.error = null
    },
  },
  extraReducers: (builder) => {
    const pending = (state: AuthState) => { state.loading = true; state.error = null }
    const rejected = (state: AuthState, action: { error: { message?: string } }) => {
      state.loading = false
      state.error = action.error.message ?? 'Ошибка'
    }
    builder
      .addCase(fetchMe.fulfilled, (state, action) => {
        state.user = action.payload
      })
      .addCase(fetchMe.rejected, (state) => {
        state.token = null
        storage.removeToken()
      })
      .addCase(login.pending, pending)
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false
        state.token = action.payload.token
        state.user = action.payload.user
      })
      .addCase(login.rejected, rejected)
      .addCase(register.pending, pending)
      .addCase(register.fulfilled, (state, action) => {
        state.loading = false
        state.token = action.payload.token
        state.user = action.payload.user
      })
      .addCase(register.rejected, rejected)
  },
})

export const { logout, setUser, clearError } = authSlice.actions
export default authSlice.reducer
