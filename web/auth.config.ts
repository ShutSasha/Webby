import type { NextAuthConfig } from 'next-auth'

import $api from '@/app/api'
import { clog, serverLog } from '@/lib/utils/utils'
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

          user.userId = googleAuthResponse.data.user.userId
          user.username = googleAuthResponse.data.user.username
          user.avatarUrl = googleAuthResponse.data.user.avatarUrl
          user.role = googleAuthResponse.data.user.role
          user.accessToken = googleAuthResponse.data.accessToken

          return true
        } catch (error) {
          console.error('Backend sync error:', error)
          return false
        }
      }
      return true
    },

    async jwt({ token, user }) {
      console.log(1)
      if (user) {
        token.id = user.userId
        token.username = user.username
        token.image = user.avatarUrl
        token.role = user.role
        token.accessToken = user.accessToken
        token.accessTokenExpires = Date.now() + 30 * 1000 // 30 seconds
      }

      if (Date.now() < token.accessTokenExpires) {
        return token
      }

      return await refreshAccessToken(token, token.accessToken)
    },

    async session({ session, token }) {
      if (token && session.user) {
        session.user.id = token.id
        session.user.username = token.username
        session.user.image = token.image
        session.user.role = token.role
        session.user.accessToken = token.accessToken
      }

      return session
    },
    authorized({ auth, request: { nextUrl } }) {
      return true
    },
  },
  providers: [],
} satisfies NextAuthConfig

async function refreshAccessToken(token: any, accessToken: string) {
  try {
    console.log('AHHAHAHA')
    const { data } = await $api.post(
      '/auth/refresh',
      {},
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    )

    console.log(data)

    return {
      ...token,
      accessToken: data.accessToken,
      accessTokenExpires: Date.now() + 30 * 1000, // 30 seconds
    }
  } catch (error) {
    serverLog('refresh error', error, true)

    return {
      ...token,
      error: 'RefreshAccessTokenError',
    }
  }
}
