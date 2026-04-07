import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import type { Conversation } from '../../types/conversation'
import { ConversationList } from '../chat/ConversationList'
import { SearchBar } from '../search/SearchBar'
import { SearchResults } from '../search/SearchResults'
import { Avatar } from '../ui/Avatar'
import { useAppSelector, useAppDispatch } from '../../hooks/useStore'
import { logout } from '../../store/authSlice'
import { USE_MOCK } from '../../api/client'
import { conversationsApi } from '../../api/conversations'
import { MOCK_CONVERSATIONS } from '../../mocks/conversations'
import { MOCK_USERS } from '../../mocks/users'

interface SidebarProps {
  conversations: Conversation[]
  activeConversationId?: string
  onSelectConversation: (id: string) => void
  onNewConversation: (convId: string) => void
}

export function Sidebar({
  conversations,
  activeConversationId,
  onSelectConversation,
  onNewConversation,
}: SidebarProps) {
  const [query, setQuery] = useState('')
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const user = useAppSelector((s) => s.auth.user)

  const isSearching = query.trim().length > 0

  async function handleStartChat(userId: string) {
    setQuery('')
    if (USE_MOCK) {
      const user = MOCK_USERS.find((u) => u.id === userId)
      if (!user) return
      const existing = MOCK_CONVERSATIONS.find((c) => c.partner.id === userId)
      if (existing) { onSelectConversation(existing.id); return }
      const newConv = {
        id: `conv-mock-${Date.now()}`,
        partner: user,
        lastMessage: null,
        unreadCount: 0,
        createdAt: new Date().toISOString(),
      }
      onNewConversation(newConv.id)
      return
    }
    const conv = await conversationsApi.create(userId)
    onNewConversation(conv.id)
  }

  return (
    <aside className="w-80 flex-shrink-0 bg-white border-r border-gray-200 flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-gray-100 flex items-center gap-3">
        <button
          onClick={() => navigate('/profile')}
          className="flex items-center gap-2 flex-1 min-w-0 hover:bg-gray-50 rounded-lg p-1 -m-1 transition-colors"
        >
          {user && <Avatar username={user.username} avatarUrl={user.avatarUrl} size="sm" />}
          <span className="font-semibold text-gray-900 text-sm truncate">{user?.username}</span>
        </button>
        <button
          onClick={() => dispatch(logout())}
          className="text-gray-400 hover:text-red-500 transition-colors p-1 rounded"
          title="Выйти"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
              d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
        </button>
      </div>

      {/* Search */}
      <div className="p-3 border-b border-gray-100">
        <SearchBar value={query} onChange={setQuery} />
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto">
        {isSearching ? (
          <SearchResults query={query} onStartChat={handleStartChat} />
        ) : (
          <ConversationList
            conversations={conversations}
            activeId={activeConversationId}
            onSelect={onSelectConversation}
          />
        )}
      </div>
    </aside>
  )
}
