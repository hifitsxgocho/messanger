import { useEffect, useRef } from 'react'
import type { Message } from '../../types/message'
import { MessageBubble } from './MessageBubble'
import { Spinner } from '../ui/Spinner'
import { useAppSelector } from '../../hooks/useStore'

interface Props {
  messages: Message[]
  loading: boolean
}

export function MessageThread({ messages, loading }: Props) {
  const currentUserId = useAppSelector((s) => s.auth.user?.id)
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages.length])

  if (loading && messages.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <Spinner />
      </div>
    )
  }

  if (messages.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <p className="text-sm text-gray-400">Начните переписку</p>
      </div>
    )
  }

  return (
    <div className="flex-1 overflow-y-auto px-4 py-4 flex flex-col gap-2">
      {messages.map((msg) => (
        <MessageBubble key={msg.id} message={msg} isOwn={msg.senderId === currentUserId} />
      ))}
      <div ref={bottomRef} />
    </div>
  )
}
