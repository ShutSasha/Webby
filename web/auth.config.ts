import type { NextAuthConfig } from 'next-auth'

import $api from '@/app/api'
import { clog, serverLog } from '@/lib/utils/utils'
import { AuthRes } from '@/types/auth'

export const authConfig = {
  pages: {
    signIn: '/login',
  },
  callbacks: {
    async signIn({ user, account }) {
      if (account?.provider === 'google') {
        try {
          clog('authConfig.signIn - Google user:', user)

          const { data: googleAuthResponse } = await $api.post<AuthRes>('/auth/google-auth', {
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
          user.accessTokenExpires = googleAuthResponse.data.accessTokenExpiresAt

          return true
        } catch (error) {
          serverLog('Backend sync error:', error)
          return false
        }
      }

      return true
    },

    async jwt({ token, user }) {
      const timeNow = Math.floor(Date.now() / 1000)

      clog('Data now', timeNow)
      clog('Token time', token.accessTokenExpires)

      if (user) {
        token.id = user.userId
        token.username = user.username
        token.image = user.avatarUrl
        token.role = user.role
        token.accessToken = user.accessToken
        token.accessTokenExpires = user.accessTokenExpires
      }

      // seconds
      if (timeNow < token.accessTokenExpires - 120) {
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
        session.error = token.error
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
    clog('entered in refresh req')

    const response = await $api.post(
      '/auth/refresh',
      {},
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    )

    const serverResponse = response.data as AuthRes
    const newToken = serverResponse.data.accessToken
    const expiresAt = serverResponse.data.accessTokenExpiresAt

    clog('REFRESH RESPONSE', serverResponse)

    const updatedData = {
      ...token,
      accessToken: newToken,
      accessTokenExpires: expiresAt,
    }

    clog('PREPARED DATA', updatedData)

    return updatedData
  } catch (error) {
    serverLog('refresh error', error, true)

    return {
      ...token,
      error: 'RefreshAccessTokenError',
    }
  }
}
