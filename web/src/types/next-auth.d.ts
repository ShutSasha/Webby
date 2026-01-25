/* eslint-disable */

import NextAuth, { DefaultSession } from 'next-auth'
import { JWT } from 'next-auth/jwt'

declare module 'next-auth' {
  interface User {
    userId: string
    username: string
    avatarUrl: string
    email: string
  }

  interface Session {
    user: {
      id: string
      username: string
      email: string
      image: string
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
  }
}
