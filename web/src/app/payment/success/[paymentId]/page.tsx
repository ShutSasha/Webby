'use client'

import { useParams, useRouter } from 'next/navigation'

import CheckCircleIcon from '@/assets/icons/shared/check-circle.svg'
import XCircleIcon from '@/assets/icons/shared/circle-xmark.svg'
import { useSubscriptionDetailsQuery } from '@/lib/hooks/api/payment/useSubscriptionDetails'
import Button from '@/ui/components/shared/Button'

export default function PaymentSuccessPage() {
  const params = useParams()
  const router = useRouter()
  const paymentId = params.paymentId as string

  const { data: subscription, isLoading, isError } = useSubscriptionDetailsQuery(paymentId)

  const formattedDate = subscription?.expirationDate
    ? new Date(subscription.expirationDate).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
    : ''

  return (
    <div className="relative min-h-[80vh] w-full overflow-hidden rounded-2xl flex items-center justify-center">
      <div className="relative z-10 max-w-2xl mx-auto px-6 text-center">
        {isLoading ? (
          <div className="flex flex-col items-center gap-6 animate-pulse">
            <div className="size-24 rounded-full bg-neutral-800/80 mb-2" />
            <div className="h-8 w-64 bg-neutral-800/80 rounded-md" />
            <div className="h-4 w-48 bg-neutral-800/60 rounded-md" />
          </div>
        ) : isError ? (
          <div className="flex flex-col items-center">
            <XCircleIcon className="size-28 text-red-500 mb-8 drop-shadow-[0_0_15px_rgba(239,68,68,0.3)]" />
            <h1 className="text-4xl font-bold text-neutral-100 mb-4 tracking-tight">Payment Verification Failed</h1>
            <p className="text-neutral-400 text-lg mb-10 leading-relaxed">
              We couldn&apos;t verify your premium status. If you were charged, please contact support.
            </p>
            <Button
              viewType="cancel"
              onClick={() => router.push('/premium')}
              className="rounded-xl px-10 py-3 font-semibold"
            >
              Back to Premium
            </Button>
          </div>
        ) : (
          <div className="flex flex-col items-center">
            <CheckCircleIcon className="size-28 text-emerald-500 mb-8 drop-shadow-[0_0_15px_rgba(16,185,129,0.3)]" />
            <h1 className="text-4xl md:text-5xl font-bold text-neutral-100 mb-6 tracking-tight">
              Welcome to Premium,
              <br /> {subscription?.username}!
            </h1>
            <p className="text-neutral-400 text-lg mb-4 leading-relaxed">
              Your payment was successful. You now have access to all limitless viewing features.
            </p>
            <div className="bg-neutral-900/50 border border-emerald-500/20 rounded-xl px-6 py-4 mb-10">
              <p className="text-emerald-100 font-medium">
                Your subscription is active until <span className="font-bold text-emerald-400">{formattedDate}</span>
              </p>
            </div>

            <div className="flex gap-4">
              <Button
                viewType="confirm"
                onClick={() => router.push('/')}
                className="rounded-xl px-10 py-3.5 text-lg font-bold bg-emerald-500 hover:bg-emerald-400
                  text-neutral-900 shadow-lg shadow-emerald-500/20 transition-all hover:scale-105"
              >
                Go to Dashboard
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
