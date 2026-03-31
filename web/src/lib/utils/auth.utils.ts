import { JWT } from 'next-auth/jwt'

import { AuthRes } from '@/types/auth.types'

import { serverLog } from './general.utils'

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:5000/api'

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
    const response = await fetch(`${API_URL}/auth/refresh`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token.accessToken}`,
        'Content-Type': 'application/json',
      },
    })

    if (!response.ok) throw new Error('Refresh request failed')

    const res: AuthRes = await response.json()

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
