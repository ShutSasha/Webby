import { Suspense } from 'react'

import VerifyPage from '@/ui/components/modules/SignUp/EmailVerifyPage'

export default async function Page({ searchParams }: { searchParams: Promise<{ email?: string }> }) {
  const { email } = await searchParams

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <VerifyPage email={email} />
    </Suspense>
  )
}
