import type { NextAuthConfig } from 'next-auth'

import $api from '@/app/api'
import { GoogleAuthRes } from '@/types/auth'

export const authConfig = {
  pages: {
    signIn: '/login',
  },
  callbacks: {
    async signIn({ user, account }) {
      if (account?.provider === 'google') {
        try {
          console.log('authConfig.signIn - Google user:', user)

          const { data: googleAuthResponse } = await $api.post<GoogleAuthRes>('/auth/google-auth', {
            id: user.id,
            name: user.name,
            email: user.email,
            image: user.image,
          })

          console.log('GOOGLE_AUTH RESPONSE', googleAuthResponse)
          if (!googleAuthResponse.success) return false

          user.userId = googleAuthResponse.data.userId
          user.username = googleAuthResponse.data.username
          user.email = googleAuthResponse.data.email
          user.avatarUrl = googleAuthResponse.data.avatarUrl

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
        token.email = user.email
        token.image = user.avatarUrl
      }
      return token
    },

    async session({ session, token }) {
      if (token && session.user) {
        session.user.id = token.id
        session.user.username = token.username
        session.user.email = token.email
        session.user.image = token.image
      }
      return session
    },
    authorized({ auth, request: { nextUrl } }) {
      const isLoggedIn = !!auth?.user
      const isOnChats = nextUrl.pathname.startsWith('/chats')

      if (isOnChats) {
        if (isLoggedIn) return true
        return false
      }
      return true
    },
  },
  providers: [],
} satisfies NextAuthConfig
