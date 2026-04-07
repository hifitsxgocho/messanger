export interface User {
  id: string
  email: string
  username: string
  bio: string
  avatarUrl: string
  createdAt: string
}

export interface UserPublic {
  id: string
  username: string
  bio: string
  avatarUrl: string
}
