import NextAuth from 'next-auth'
import Credentials from 'next-auth/providers/credentials'
import Google from 'next-auth/providers/google'

import { login } from '@/app/api/auth'
import { isLoginRes } from '@/types/auth'

import { authConfig } from './auth.config'

export const { auth, signIn, signOut, handlers } = NextAuth({
  ...authConfig,
  providers: [
    Google({
      clientId: process.env.AUTH_GOOGLE_ID,
      clientSecret: process.env.AUTH_GOOGLE_SECRET,
      authorization: {
        params: {
          prompt: 'select_account',
          access_type: 'offline',
          response_type: 'code',
        },
      },
    }),
    Credentials({
      async authorize(credentials) {
        const response = await login(credentials.email as string, credentials.password as string)

        if (!response) return null

        if (isLoginRes(response)) {
          return { ...response.user, accessToken: response.accessToken }
        }

        if ('errors' in response) {
          throw new Error(JSON.stringify(response.errors))
        }

        throw new Error(JSON.stringify({ global: 'Unknown authentication error' }))
      },
    }),
  ],
})
