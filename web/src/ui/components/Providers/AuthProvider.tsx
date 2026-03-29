'use client'

import { useEffect } from 'react'

import { signOut, useSession } from 'next-auth/react'

import { setGlobalToken, setUpdateSession } from '@/lib/utils/auth-token'

type Props = {
  children: React.ReactNode
}

export function AuthProvider({ children }: Props) {
  const { data: session, update } = useSession()

  useEffect(() => {
    setGlobalToken(session?.user?.accessToken || null)
    setUpdateSession(update)

    if (session?.error === 'RefreshAccessTokenError') {
      signOut({ callbackUrl: '/login', redirect: true })
    }
  }, [session, update])

  return <>{children}</>
}
