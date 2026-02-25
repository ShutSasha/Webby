import type { NextAuthConfig } from 'next-auth'

import $api from '@/app/api'
import { mapUserData, refreshAccessToken } from '@/lib/utils/auth'
import { serverLog } from '@/lib/utils/utils'
import { AuthRes } from '@/types/auth'

const TOKEN_REFRESH_BUFFER = 120

export const authConfig = {
  pages: {
    signIn: '/login',
  },
  callbacks: {
    async signIn({ user, account }) {
      if (account?.provider === 'google') {
        try {
          const { data: res } = await $api.post<AuthRes>('/auth/google-auth', {
            id: user.id,
            name: user.name,
            email: user.email,
            image: user.image,
          })

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

    async jwt({ token, user }) {
      if (user) {
        mapUserData(token, user)
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
