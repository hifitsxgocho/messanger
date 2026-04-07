import { useState, useEffect } from 'react'
import type { UserPublic } from '../../types/user'
import { Avatar } from '../ui/Avatar'
import { Spinner } from '../ui/Spinner'
import { USE_MOCK } from '../../api/client'
import { usersApi } from '../../api/users'
import { MOCK_USERS } from '../../mocks/users'

interface Props {
  query: string
  onStartChat: (userId: string) => void
}

export function SearchResults({ query, onStartChat }: Props) {
  const [results, setResults] = useState<UserPublic[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!query.trim()) { setResults([]); return }

    const timer = setTimeout(async () => {
      setLoading(true)
      try {
        if (USE_MOCK) {
          await new Promise((r) => setTimeout(r, 300))
          const q = query.toLowerCase()
          setResults(MOCK_USERS.filter((u) => u.username.toLowerCase().includes(q)))
        } else {
          const res = await usersApi.search(query)
          setResults(res)
        }
      } finally {
        setLoading(false)
      }
    }, 300)

    return () => clearTimeout(timer)
  }, [query])

  if (loading) {
    return <div className="flex justify-center py-6"><Spinner /></div>
  }

  if (results.length === 0) {
    return (
      <div className="py-8 text-center">
        <p className="text-sm text-gray-400">Пользователи не найдены</p>
      </div>
    )
  }

  return (
    <div>
      {results.map((user) => (
        <button
          key={user.id}
          onClick={() => onStartChat(user.id)}
          className="w-full flex items-center gap-3 px-4 py-3 hover:bg-gray-50 transition-colors text-left"
        >
          <Avatar username={user.username} avatarUrl={user.avatarUrl} size="md" />
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-gray-900 truncate">{user.username}</p>
            {user.bio && <p className="text-xs text-gray-400 truncate">{user.bio}</p>}
          </div>
          <svg className="w-4 h-4 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
        </button>
      ))}
    </div>
  )
}
