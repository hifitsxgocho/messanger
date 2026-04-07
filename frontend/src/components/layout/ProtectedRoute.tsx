import { Navigate } from 'react-router-dom'
import { useAppSelector } from '../../hooks/useStore'

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const token = useAppSelector((s) => s.auth.token)
  if (!token) return <Navigate to="/login" replace />
  return <>{children}</>
}
