'use client'

import { useEffect } from 'react'

import { SessionProvider, useSession, signOut } from 'next-auth/react'

import { clog } from '@/lib/utils/utils'

export default function Providers({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <RefreshTokenProvider>{children}</RefreshTokenProvider>
    </SessionProvider>
  )
}

export function RefreshTokenProvider({ children }: { children: React.ReactNode }) {
  const { data: session } = useSession()

  useEffect(() => {
    clog('Current session state:', session)

    if (session?.error === 'RefreshAccessTokenError') {
      clog('Refresh failed, signing out...')
      signOut({ callbackUrl: '/login', redirect: true })
    }
  }, [session?.error])

  return <>{children}</>
}
