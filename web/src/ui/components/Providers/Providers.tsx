import { SessionProvider } from 'next-auth/react'

import { auth } from '@/workspace/auth'

import { AuthProvider } from './AuthProvider'

export default async function Providers({ children }: { children: React.ReactNode }) {
  const session = await auth()

  return (
    <SessionProvider session={session} refetchOnWindowFocus={true} refetchInterval={1500}>
      <AuthProvider>{children}</AuthProvider>
    </SessionProvider>
  )
}
