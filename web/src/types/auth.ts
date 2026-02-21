type Role = 'User' | 'Admin' | 'Moderator'

type GoogleAuthRes = {
  success: boolean
  message: string
  data: {
    userId: string
    email: string
    username: string
    about: string
    avatarUrl: string
    role: Role
  }
  errors: null
}

export { type GoogleAuthRes, type Role }
