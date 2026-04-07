import type { Conversation } from '../types/conversation'

export const MOCK_CONVERSATIONS: Conversation[] = [
  {
    id: 'conv-1',
    partner: { id: 'user-2', username: 'alice', bio: 'UI/UX дизайнер', avatarUrl: '' },
    lastMessage: {
      body: 'Привет! Как дела?',
      senderId: 'user-2',
      createdAt: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
    },
    unreadCount: 2,
    createdAt: '2024-01-10T00:00:00Z',
  },
  {
    id: 'conv-2',
    partner: { id: 'user-3', username: 'bob', bio: 'Бэкенд разработчик', avatarUrl: '' },
    lastMessage: {
      body: 'Смотри, нашёл крутую библиотеку для Go',
      senderId: 'user-3',
      createdAt: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
    },
    unreadCount: 0,
    createdAt: '2024-01-08T00:00:00Z',
  },
  {
    id: 'conv-3',
    partner: { id: 'user-4', username: 'carol', bio: 'Фронтенд разработчик', avatarUrl: '' },
    lastMessage: {
      body: 'Спасибо за помощь!',
      senderId: 'user-1',
      createdAt: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
    },
    unreadCount: 0,
    createdAt: '2024-01-05T00:00:00Z',
  },
]
