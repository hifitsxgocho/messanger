import { useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { RegisterForm } from '../components/auth/RegisterForm'
import { useAppDispatch, useAppSelector } from '../hooks/useStore'
import { register as registerAction, clearError } from '../store/authSlice'
import type { RegisterRequest } from '../types/auth'

export function RegisterPage() {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { loading, error, token } = useAppSelector((s) => s.auth)

  useEffect(() => {
    if (token) navigate('/', { replace: true })
    return () => { dispatch(clearError()) }
  }, [token, navigate, dispatch])

  async function handleSubmit(data: RegisterRequest) {
    await dispatch(registerAction(data))
  }

  return (
    <div className="w-full min-h-screen flex items-center justify-center bg-gray-50 p-4">
      <div className="w-full max-w-sm">
        <div className="text-center mb-8">
          <div className="w-16 h-16 bg-blue-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
            <svg className="w-9 h-9 text-white" fill="currentColor" viewBox="0 0 24 24">
              <path d="M20 2H4c-1.1 0-2 .9-2 2v18l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2z" />
            </svg>
          </div>
          <h1 className="text-2xl font-bold text-gray-900">Регистрация</h1>
          <p className="text-gray-500 text-sm mt-1">Создайте аккаунт</p>
        </div>
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
          <RegisterForm onSubmit={handleSubmit} loading={loading} error={error} />
          <p className="mt-4 text-center text-sm text-gray-500">
            Уже есть аккаунт?{' '}
            <Link to="/login" className="text-blue-600 hover:underline font-medium">
              Войти
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}
