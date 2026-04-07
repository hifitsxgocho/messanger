import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../hooks/useStore'
import { setUser } from '../store/authSlice'
import { Avatar } from '../components/ui/Avatar'
import { Button } from '../components/ui/Button'
import { EditProfileModal } from '../components/profile/EditProfileModal'
import { USE_MOCK } from '../api/client'
import { usersApi } from '../api/users'

export function ProfilePage() {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const user = useAppSelector((s) => s.auth.user)
  const [editOpen, setEditOpen] = useState(false)
  const [saving, setSaving] = useState(false)

  if (!user) return null

  async function handleSave(data: { username: string; bio?: string }) {
    setSaving(true)
    try {
      if (USE_MOCK) {
        await new Promise((r) => setTimeout(r, 400))
        dispatch(setUser({ ...user!, ...data, bio: data.bio ?? '' }))
      } else {
        const updated = await usersApi.updateMe(data)
        dispatch(setUser(updated))
      }
      setEditOpen(false)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="flex-1 flex items-center justify-center p-6 bg-gray-50">
      <div className="w-full max-w-sm">
        <button
          onClick={() => navigate(-1)}
          className="flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 mb-6 transition-colors"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Назад
        </button>

        <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 text-center">
          <Avatar username={user.username} avatarUrl={user.avatarUrl} size="lg" className="mx-auto mb-4" />
          <h1 className="text-xl font-bold text-gray-900">{user.username}</h1>
          <p className="text-sm text-gray-500 mt-1">{user.email}</p>
          {user.bio && <p className="text-sm text-gray-600 mt-3 leading-relaxed">{user.bio}</p>}
          <Button className="mt-6 w-full" variant="ghost" onClick={() => setEditOpen(true)}>
            Редактировать профиль
          </Button>
        </div>
      </div>

      {editOpen && (
        <EditProfileModal
          user={user}
          onSave={handleSave}
          onClose={() => setEditOpen(false)}
          saving={saving}
        />
      )}
    </div>
  )
}
