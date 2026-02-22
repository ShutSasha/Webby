import { Suspense } from 'react'

import VerifyPage from '@/ui/components/modules/SignUp/EmailVerifyPage'

type Props = {
  searchParams: Promise<{ email?: string }>
}

export default async function Page({ searchParams }: Props) {
  const { email } = await searchParams

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <VerifyPage email={email} />
    </Suspense>
  )
}
