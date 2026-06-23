'use client'

import { useRouter } from 'next/navigation'

import XCircleIcon from '@/assets/icons/shared/circle-xmark.svg'
import Button from '@/ui/components/shared/Button'

export default function PaymentFailedPage() {
  const router = useRouter()

  return (
    <div className="relative min-h-[80vh] w-full overflow-hidden rounded-2xl flex items-center justify-center">
      <div className="relative z-10 max-w-2xl mx-auto px-6 text-center flex flex-col items-center">
        <XCircleIcon className="size-24 text-red-500 mb-8 drop-shadow-[0_0_15px_rgba(239,68,68,0.3)]" />

        <h1 className="text-3xl md:text-4xl font-bold text-foreground-secondary mb-6 tracking-tight">Payment Failed</h1>

        <p className="text-foreground-muted text-md mb-8 leading-relaxed">
          We couldn&apos;t process your transaction, or the payment was cancelled. Don&apos;t worry, your account has
          not been charged.
        </p>

        <div className="flex gap-4 justify-center">
          <Button
            viewType="confirm"
            onClick={() => router.push('/premium')}
            className="rounded-xl px-8 py-3 text-lg font-bold bg-surface-inverse-secondary
              hover:bg-surface-inverse-muted text-foreground-inverse-subtle shadow-lg transition-all hover:scale-105"
          >
            Try Again
          </Button>

          <Button
            viewType="cancel"
            onClick={() => router.push('/')}
            className="rounded-xl px-8 py-3 text-lg font-semibold"
          >
            Back to Home
          </Button>
        </div>
      </div>
    </div>
  )
}
