import { useEffect, useRef } from 'react'
import { useParams } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../hooks/useStore'
import { fetchMessages, sendMessage } from '../store/messagesSlice'
import { updateLastMessage } from '../store/conversationsSlice'
import { MessageThread } from '../components/chat/MessageThread'
import { MessageInput } from '../components/chat/MessageInput'
import { EmptyChat } from '../components/chat/EmptyChat'
import { Avatar } from '../components/ui/Avatar'

const POLL_INTERVAL = 3000

export function ChatPage() {
  const { conversationId } = useParams<{ conversationId: string }>()
  const dispatch = useAppDispatch()
  const conversations = useAppSelector((s) => s.conversations.items)
  const messages = useAppSelector((s) =>
    conversationId ? (s.messages.byConversation[conversationId] ?? []) : []
  )
  const loading = useAppSelector((s) => s.messages.loading)
  const pollingRef = useRef<ReturnType<typeof setInterval> | null>(null)

  const conversation = conversations.find((c) => c.id === conversationId)

  useEffect(() => {
    if (!conversationId) return

    dispatch(fetchMessages({ conversationId }))

    pollingRef.current = setInterval(() => {
      dispatch(fetchMessages({ conversationId }))
    }, POLL_INTERVAL)

    return () => {
      if (pollingRef.current) clearInterval(pollingRef.current)
    }
  }, [conversationId, dispatch])

  async function handleSend(body: string) {
    if (!conversationId) return
    const result = await dispatch(sendMessage({ conversationId, body }))
    if (sendMessage.fulfilled.match(result)) {
      dispatch(updateLastMessage({ conversationId, message: result.payload.message }))
    }
  }

  if (!conversationId) {
    return <EmptyChat />
  }

  return (
    <div className="flex flex-col h-full bg-gray-50">
      {/* Chat header */}
      {conversation && (
        <div className="flex items-center gap-3 px-4 py-3 bg-white border-b border-gray-100 shadow-sm">
          <Avatar
            username={conversation.partner.username}
            avatarUrl={conversation.partner.avatarUrl}
            size="md"
          />
          <div>
            <p className="font-semibold text-gray-900 text-sm">{conversation.partner.username}</p>
            {conversation.partner.bio && (
              <p className="text-xs text-gray-400 truncate max-w-xs">{conversation.partner.bio}</p>
            )}
          </div>
        </div>
      )}

      <MessageThread messages={messages} loading={loading} />
      <MessageInput onSend={handleSend} />
    </div>
  )
}
