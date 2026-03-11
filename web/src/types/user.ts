type Role = 'User' | 'Admin' | 'Moderator'

interface User {
  userId: string
  email: string
  username: string
  about: string
  avatarUrl: string
  role: Role
}

export { type Role, type User }
