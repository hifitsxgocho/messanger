import type { Conversation } from '../../types/conversation'
import { ConversationItem } from './ConversationItem'
import { Spinner } from '../ui/Spinner'
import { useAppSelector } from '../../hooks/useStore'

interface Props {
  conversations: Conversation[]
  activeId?: string
  onSelect: (id: string) => void
}

export function ConversationList({ conversations, activeId, onSelect }: Props) {
  const loading = useAppSelector((s) => s.conversations.loading)

  if (loading && conversations.length === 0) {
    return (
      <div className="flex justify-center py-8">
        <Spinner />
      </div>
    )
  }

  if (conversations.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
        <div className="w-12 h-12 bg-gray-100 rounded-full flex items-center justify-center mb-3">
          <svg className="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
              d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
        </div>
        <p className="text-sm text-gray-500">Нет чатов</p>
        <p className="text-xs text-gray-400 mt-1">Найдите пользователя через поиск</p>
      </div>
    )
  }

  return (
    <div>
      {conversations.map((conv) => (
        <ConversationItem
          key={conv.id}
          conversation={conv}
          isActive={conv.id === activeId}
          onClick={() => onSelect(conv.id)}
        />
      ))}
    </div>
  )
}
