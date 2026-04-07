import axios from 'axios'
import { storage } from '../utils/storage'

const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true'

export const apiClient = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
})

apiClient.interceptors.request.use((config) => {
  const token = storage.getToken()
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

apiClient.interceptors.response.use(
  (res) => res,
  (error) => {
    if (error.response?.status === 401) {
      storage.removeToken()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export { USE_MOCK }
