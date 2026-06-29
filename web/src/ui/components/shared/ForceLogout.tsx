'use client'

import { useEffect } from 'react'

import { signOut } from 'next-auth/react'

export default function ForceLogout() {
  useEffect(() => {
    signOut({ callbackUrl: '/login', redirect: true })
  }, [])

  return (
    <div className="flex h-screen w-full flex-col items-center justify-center gap-4 bg-background">
      <div className="size-8 border-4 border-emerald-500/30 border-t-emerald-500 rounded-full animate-spin" />
      <p className="text-foreground-muted text-sm">Access denied. Signing out...</p>
    </div>
  )
}
