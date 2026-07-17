import { Suspense } from 'react'

import VerifyPage from '@/ui/components/modules/SignUp/EmailVerifyPage'

type Props = {
  searchParams: Promise<{ email?: string }>
}

export default async function Page({ searchParams }: Props) {
  const { email } = await searchParams

  return (
    <Suspense fallback={<VerifyPageSkeleton />}>
      <VerifyPage email={email} />
    </Suspense>
  )
}

function VerifyPageSkeleton() {
  return (
    <div className="flex flex-1 items-center justify-center">
      <div className="flex w-full bg-surface max-w-[700px] rounded-[20px] p-10 box-border flex-col items-center gap-4">
        {/* Title and subtitle placeholders */}
        <div className="space-y-3 flex flex-col items-center w-full mb-2">
          <div className="h-8 w-40 bg-background rounded-lg animate-pulse" />
          <div className="h-4 w-60 bg-background rounded-md animate-pulse" />
        </div>

        {/* OTP Slots placeholder (6 slots) */}
        <div className="flex gap-2">
          {[...Array(6)].map((_, idx) => (
            <div key={idx} className="w-12 h-14 bg-background border-2 border-border rounded-xl animate-pulse" />
          ))}
        </div>

        {/* Resend button placeholder */}
        <div className="h-5 w-32 bg-background rounded-md animate-pulse mt-2" />
      </div>
    </div>
  )
}
