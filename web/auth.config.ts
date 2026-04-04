import type { NextAuthConfig } from 'next-auth'

import { mapUserData, refreshAccessToken } from '@/lib/utils/auth.utils'
import { serverLog } from '@/lib/utils/general.utils'
import { AuthRes } from '@/types/auth.types'

const TOKEN_REFRESH_BUFFER = 120
const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:5000/api'

export const authConfig = {
  pages: {
    signIn: '/login',
  },
  callbacks: {
    async signIn({ user, account }) {
      if (account?.provider === 'google') {
        try {
          const response = await fetch(`${API_URL}/auth/google-auth`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              id: user.id,
              name: user.name,
              email: user.email,
              image: user.image,
            }),
          })

          const res: AuthRes = await response.json()

          if (!res.success) return false

          const serverUser = res.data.user

          user.userId = serverUser.userId
          user.username = serverUser.username
          user.avatarUrl = serverUser.avatarUrl
          user.role = serverUser.role
          user.accessToken = res.data.accessToken
          user.accessTokenExpires = res.data.accessTokenExpiresAt

          return true
        } catch (error) {
          serverLog('Backend sync error:', error)
          return false
        }
      }

      return true
    },

    async jwt({ token, user, trigger, session: sessionData }) {
      if (user) {
        mapUserData(token, user)
        return token
      }

      if (trigger === 'update' && sessionData) {
        if (sessionData.image) token.image = sessionData.image

        return token
      }

      // seconds
      const timeNow = Math.floor(Date.now() / 1000)
      if (timeNow < token.accessTokenExpires - TOKEN_REFRESH_BUFFER) {
        return token
      }

      return await refreshAccessToken(token)
    },

    async session({ session, token }) {
      if (token && session.user) {
        mapUserData(session.user, token)
        session.error = token.error
      }
      return session
    },
  },
  providers: [],
} satisfies NextAuthConfig
