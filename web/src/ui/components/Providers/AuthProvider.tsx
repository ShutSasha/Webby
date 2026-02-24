'use client'

import { useEffect } from 'react'

import { Session } from 'next-auth'
import { signOut } from 'next-auth/react'

import { clog } from '@/lib/utils/utils'

type Props = {
  session: Session | null
  children: React.ReactNode
}

export function AuthProvider({ session, children }: Props) {
  useEffect(() => {
    clog('log')
    clog('Current session state:', session)

    if (session?.error === 'RefreshAccessTokenError') {
      clog('Refresh failed, signing out...')
      signOut({ callbackUrl: '/login', redirect: true })
    }
  }, [session])

  return <>{children}</>
}
