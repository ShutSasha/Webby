import NextAuth from 'next-auth'
import Credentials from 'next-auth/providers/credentials'
import Google from 'next-auth/providers/google'
import { authConfig } from './auth.config'
import $api from '@/app/lib/api'
import { z } from 'zod'
import type { User } from 'next-auth'

const CredentialsSchema = z.object({
  username: z.string().min(4),
  password: z.string().min(8),
})

async function login(username: string, password: string): Promise<User | undefined> {
  try {
    const { data: user } = await $api.post('/auth/login', { username, password })
    return user
  } catch (error) {
    console.error('Failed to fetch user:', error)
    return undefined
  }
}

export const { auth, signIn, signOut, handlers } = NextAuth({
  ...authConfig,
  providers: [
    Google({
      clientId: process.env.AUTH_GOOGLE_ID,
      clientSecret: process.env.AUTH_GOOGLE_SECRET,
    }),
    Credentials({
      async authorize(credentials) {
        const parsedCredentials = CredentialsSchema.safeParse(credentials)

        if (parsedCredentials.success) {
          const { username, password } = parsedCredentials.data
          const user = await login(username, password)

          if (!user) return null

          return {
            ...user,
            id: user._id,
          }
        }
        return null
      },
    }),
  ],
})
