'use client'

import { useState } from 'react'

import { OTPInput } from 'input-otp'

import { resendVerifyCode, verifyUser } from '@/lib/actions/auth.actions'
import { useIsClient } from '@/lib/hooks/useIsClient'
import { useResendTimer } from '@/lib/hooks/useResendTimer'
import { parseAxiosError } from '@/lib/utils/general.utils'
import { useCommonStore } from '@/stores/common.store'

import { Slot } from '../../SignUp/EmailVerifyPage'

export default function CodeStep() {
  const forgotPasswordEmail = useCommonStore(state => state.forgotPasswordEmail)
  const setForgotPasswordStep = useCommonStore(state => state.setForgotPasswordStep)

  const [loading, setLoading] = useState<boolean>(false)
  const [errors, setErrors] = useState<Record<string, string> | null>(null)
  const [codeMessage, setCodeMessage] = useState<string>('')

  const isClientReady = useIsClient()
  const { timeLeft, start: startTimer } = useResendTimer(forgotPasswordEmail, 60)

  const handleResendCode = async () => {
    if (loading || timeLeft > 0) return

    setLoading(true)
    setErrors(null)
    setCodeMessage('')

    try {
      await resendVerifyCode({ email: forgotPasswordEmail })
      setCodeMessage('Verification code sent successfully.')
      startTimer()
    } catch (error) {
      const serverErrors = parseAxiosError(error)
      setErrors(serverErrors)
    } finally {
      setLoading(false)
    }
  }

  const handleComplete = async (code: string) => {
    setErrors(null)
    setLoading(true)

    try {
      const res = await verifyUser({ email: forgotPasswordEmail, code })

      if (res?.errors) {
        setErrors(res.errors as Record<string, string>)
      } else {
        setForgotPasswordStep(3)
      }
    } catch (error) {
      setErrors(parseAxiosError(error))
    } finally {
      setLoading(false)
      setCodeMessage('')
    }
  }

  if (!isClientReady) return null

  return (
    <div className="flex flex-col items-center gap-6 w-full animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="text-center space-y-2">
        <h2 className="text-white text-2xl font-bold">Confirm your email</h2>
        <p className="text-neutral-400 text-sm">
          We have sent a recovery code to <br />
          <span className="text-white font-medium">{forgotPasswordEmail}</span>
        </p>
      </div>

      <div className="space-y-4 flex flex-col items-center">
        <OTPInput
          maxLength={6}
          onComplete={handleComplete}
          disabled={loading}
          containerClassName="group flex items-center has-[:disabled]:opacity-50"
          render={({ slots }) => (
            <div className="flex gap-2">
              {slots.map((slot, idx) => (
                <Slot key={idx} {...slot} />
              ))}
            </div>
          )}
        />

        {errors && (
          <div className="space-y-1">
            {Object.entries(errors).map(([field, message]) => (
              <p key={field} className="text-red-500 text-sm text-center">
                {message}
              </p>
            ))}
          </div>
        )}

        <button
          type="button"
          className={`text-sm transition-all ${
            timeLeft > 0 || loading
              ? 'text-neutral-500 cursor-not-allowed'
              : 'text-emerald-500 hover:underline cursor-pointer'
            }`}
          onClick={handleResendCode}
          disabled={timeLeft > 0 || loading}
        >
          {timeLeft > 0 ? `Resend code in ${timeLeft}s` : loading ? 'Sending...' : 'Resend code'}
        </button>

        {codeMessage && <p className="text-green-500 text-sm">{codeMessage}</p>}
      </div>

      <button
        type="button"
        onClick={() => setForgotPasswordStep(1)}
        className="mt-2 text-neutral-500 text-sm hover:text-white transition-colors cursor-pointer"
      >
        Wrong email? <span className="text-emerald-500 hover:underline">Change it</span>
      </button>
    </div>
  )
}
