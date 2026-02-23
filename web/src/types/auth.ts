import { User } from 'next-auth'

type Role = 'User' | 'Admin' | 'Moderator'

type GoogleAuthRes = {
  success: boolean
  message: string
  data: {
    user: {
      userId: string
      email: string
      username: string
      about: string
      avatarUrl: string
      role: Role
    }
    accessToken: string
  }
  errors: null
}

type LoginRes = {
  user: User
  accessToken: string
}

type ReturnError = {
  success: boolean
  errors: Record<string, string>
}

export function isLoginRes(data: LoginRes | ReturnError): data is LoginRes {
  return typeof data === 'object' && data !== null && 'user' in data && 'accessToken' in data
}

export { type GoogleAuthRes, type Role, type LoginRes }
