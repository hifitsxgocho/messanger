import type { User, UserPublic } from '../types/user'

export const MOCK_CURRENT_USER: User = {
  id: 'user-1',
  email: 'me@example.com',
  username: 'me',
  bio: 'Разработчик и любитель кофе ☕',
  avatarUrl: '',
  createdAt: '2024-01-01T00:00:00Z',
}

export const MOCK_USERS: UserPublic[] = [
  {
    id: 'user-2',
    username: 'alice',
    bio: 'UI/UX дизайнер',
    avatarUrl: '',
  },
  {
    id: 'user-3',
    username: 'bob',
    bio: 'Бэкенд разработчик, Go enthusiast',
    avatarUrl: '',
  },
  {
    id: 'user-4',
    username: 'carol',
    bio: 'Фронтенд разработчик',
    avatarUrl: '',
  },
  {
    id: 'user-5',
    username: 'dave',
    bio: 'DevOps инженер',
    avatarUrl: '',
  },
]
