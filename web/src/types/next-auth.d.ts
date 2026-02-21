/* eslint-disable */

import NextAuth, { DefaultSession } from 'next-auth'
import { JWT } from 'next-auth/jwt'
import { Role } from './auth'

declare module 'next-auth' {
  interface User {
    userId: string
    username: string
    avatarUrl: string
    email: string
    role: Role
  }

  interface Session {
    user: {
      id: string
      username: string
      email: string
      image: string
      role: Role
    }
  }
}

// jwt token
declare module 'next-auth/jwt' {
  interface JWT {
    id: string
    username: string
    email: string
    image: string
    role: Role
  }
}
