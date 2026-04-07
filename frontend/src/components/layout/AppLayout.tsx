import { useEffect } from 'react'
import { Outlet, useNavigate, useParams } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { useAppDispatch, useAppSelector } from '../../hooks/useStore'
import { fetchConversations } from '../../store/conversationsSlice'
import { fetchMe } from '../../store/authSlice'

export function AppLayout() {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { conversationId } = useParams()
  const conversations = useAppSelector((s) => s.conversations.items)
  const user = useAppSelector((s) => s.auth.user)

  useEffect(() => {
    if (!user) dispatch(fetchMe())
    dispatch(fetchConversations())
  }, [dispatch])

  function handleSelectConversation(id: string) {
    navigate(`/conversations/${id}`)
  }

  function handleNewConversation(convId: string) {
    dispatch(fetchConversations())
    navigate(`/conversations/${convId}`)
  }

  return (
    <div className="flex w-full h-screen bg-gray-100">
      <Sidebar
        conversations={conversations}
        activeConversationId={conversationId}
        onSelectConversation={handleSelectConversation}
        onNewConversation={handleNewConversation}
      />
      <main className="flex-1 flex flex-col min-w-0">
        <Outlet />
      </main>
    </div>
  )
}
