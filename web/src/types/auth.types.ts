import { User } from 'next-auth'

import { BaseServerResponse } from './general.types'
import { Role } from './user.types'

type AuthPayload = {
  user: {
    userId: string
    email: string
    username: string
    about: string
    avatarUrl: string
    role: Role
  }
  accessToken: string
  accessTokenExpiresAt: number
}

export type AuthServerResponse = BaseServerResponse<AuthPayload>

export type LoginResponse = {
  user: User
  accessToken: string
}

export type ReturnError = {
  success: boolean
  errors: Record<string, string>
}

export function isLoginRes(data: LoginResponse | ReturnError): data is LoginResponse {
  return typeof data === 'object' && data !== null && 'user' in data && 'accessToken' in data
}
