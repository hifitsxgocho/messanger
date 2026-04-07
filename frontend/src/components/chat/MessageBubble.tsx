import { clsx } from 'clsx'
import type { Message } from '../../types/message'
import { formatMessageTime } from '../../utils/formatDate'

interface Props {
  message: Message
  isOwn: boolean
  showTime?: boolean
}

export function MessageBubble({ message, isOwn, showTime = true }: Props) {
  return (
    <div className={clsx('flex', isOwn ? 'justify-end' : 'justify-start')}>
      <div className={clsx('max-w-[70%] flex flex-col gap-0.5', isOwn ? 'items-end' : 'items-start')}>
        <div
          className={clsx(
            'px-3 py-2 rounded-2xl text-sm leading-relaxed',
            isOwn
              ? 'bg-blue-600 text-white rounded-br-sm'
              : 'bg-white text-gray-900 border border-gray-100 shadow-xs rounded-bl-sm'
          )}
        >
          {message.body}
        </div>
        {showTime && (
          <span className="text-[10px] text-gray-400 px-1">
            {formatMessageTime(message.createdAt)}
          </span>
        )}
      </div>
    </div>
  )
}
