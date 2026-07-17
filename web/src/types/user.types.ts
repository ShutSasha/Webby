type Role = 'User' | 'Admin' | 'Moderator'

interface User {
  userId: string
  email: string
  username: string
  about: string
  avatarUrl: string
  role: Role
}

type UserFolowStats = {
  followers: number
  following: number
}

export { type Role, type User, type UserFolowStats }
