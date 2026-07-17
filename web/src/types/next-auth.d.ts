/* eslint-disable */

import NextAuth, { DefaultSession } from 'next-auth'
import { JWT } from 'next-auth/jwt'
import { Role } from './auth.types'

declare module 'next-auth' {
  interface User {
    userId: string
    username: string
    avatarUrl: string
    role: Role
    accessToken: string
    accessTokenExpires: number
  }

  interface Session {
    user: {
      id: string
      username: string
      image: string
      role: Role
      accessToken: string
      accessTokenExpires: number
    }
    error?: 'RefreshAccessTokenError'
  }
}

// jwt token
declare module 'next-auth/jwt' {
  interface JWT {
    id: string
    username: string
    image: string
    role: Role
    accessToken: string
    accessTokenExpires: number
    error?: 'RefreshAccessTokenError'
  }
}
