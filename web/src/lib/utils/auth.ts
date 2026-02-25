import { JWT } from 'next-auth/jwt'

import $api from '@/app/api'
import { AuthRes } from '@/types/auth'

import { clog, serverLog } from './utils'

export const mapUserData = (target: any, source: any) => {
  target.id = source.userId ?? source.id
  target.username = source.username
  target.image = source.avatarUrl ?? source.image
  target.role = source.role
  target.accessToken = source.accessToken
  target.accessTokenExpires = source.accessTokenExpires
}

export async function refreshAccessToken(token: JWT): Promise<JWT> {
  try {
    clog('Refreshing access token...')

    const { data: res } = await $api.post<AuthRes>(
      '/auth/refresh',
      {},
      {
        headers: { Authorization: `Bearer ${token.accessToken}` },
      },
    )

    return {
      ...token,
      accessToken: res.data.accessToken,
      accessTokenExpires: res.data.accessTokenExpiresAt,
      error: undefined,
    }
  } catch (error) {
    serverLog('refresh error', error, true)

    return {
      ...token,
      error: 'RefreshAccessTokenError',
    }
  }
}
