import type { NextAuthConfig } from 'next-auth'

import $api from '@/app/api'
import { clog } from '@/lib/utils/utils'
import { GoogleAuthRes } from '@/types/auth'

export const authConfig = {
  pages: {
    signIn: '/login',
  },
  callbacks: {
    async signIn({ user, account }) {
      if (account?.provider === 'google') {
        try {
          clog('authConfig.signIn - Google user:', user)

          const { data: googleAuthResponse } = await $api.post<GoogleAuthRes>('/auth/google-auth', {
            id: user.id,
            name: user.name,
            email: user.email,
            image: user.image,
          })

          clog('GOOGLE_AUTH RESPONSE', googleAuthResponse)
          if (!googleAuthResponse.success) return false

          user.userId = googleAuthResponse.data.userId
          user.username = googleAuthResponse.data.username
          user.avatarUrl = googleAuthResponse.data.avatarUrl
          user.role = googleAuthResponse.data.role

          return true
        } catch (error) {
          console.error('Backend sync error:', error)
          return false
        }
      }
      return true
    },

    async jwt({ token, user }) {
      if (user) {
        token.id = user.userId
        token.username = user.username
        token.image = user.avatarUrl
        token.role = user.role
      }
      return token
    },

    async session({ session, token }) {
      if (token && session.user) {
        session.user.id = token.id
        session.user.username = token.username
        session.user.image = token.image
        session.user.role = token.role
      }

      return session
    },
    authorized({ auth, request: { nextUrl } }) {
      return true
    },
  },
  providers: [],
} satisfies NextAuthConfig
