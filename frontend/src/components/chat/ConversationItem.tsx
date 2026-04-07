import { clsx } from 'clsx'
import type { Conversation } from '../../types/conversation'
import { Avatar } from '../ui/Avatar'
import { formatConversationTime } from '../../utils/formatDate'
import { useAppSelector } from '../../hooks/useStore'

interface Props {
  conversation: Conversation
  isActive: boolean
  onClick: () => void
}

export function ConversationItem({ conversation, isActive, onClick }: Props) {
  const currentUserId = useAppSelector((s) => s.auth.user?.id)
  const { partner, lastMessage, unreadCount } = conversation

  const preview = lastMessage
    ? `${lastMessage.senderId === currentUserId ? 'Вы: ' : ''}${lastMessage.body}`
    : 'Нет сообщений'

  return (
    <button
      onClick={onClick}
      className={clsx(
        'w-full flex items-center gap-3 px-4 py-3 hover:bg-gray-50 transition-colors text-left',
        isActive && 'bg-blue-50 hover:bg-blue-50'
      )}
    >
      <Avatar username={partner.username} avatarUrl={partner.avatarUrl} size="md" />
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between gap-2">
          <span className="font-medium text-gray-900 text-sm truncate">{partner.username}</span>
          {lastMessage && (
            <span className="text-xs text-gray-400 flex-shrink-0">
              {formatConversationTime(lastMessage.createdAt)}
            </span>
          )}
        </div>
        <div className="flex items-center justify-between gap-2 mt-0.5">
          <p className="text-xs text-gray-500 truncate">{preview}</p>
          {unreadCount > 0 && (
            <span className="flex-shrink-0 min-w-[18px] h-[18px] bg-blue-600 text-white text-xs rounded-full flex items-center justify-center px-1">
              {unreadCount > 99 ? '99+' : unreadCount}
            </span>
          )}
        </div>
      </div>
    </button>
  )
}
