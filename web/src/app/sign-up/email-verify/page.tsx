'use client'

import { useEffect, useState } from 'react'

import { OTPInput, SlotProps } from 'input-otp'
import Link from 'next/link'
import { useSearchParams } from 'next/navigation'

import $api from '@/app/api'
import { BaseServerResponse, serverLog } from '@/lib/utils/utils'

export default function VerifyPage() {
  const searchParams = useSearchParams()
  const email = searchParams.get('email')

  const [errors, setErrors] = useState<Record<string, string> | null>(null)
  const [success, setSuccess] = useState<boolean>(false)
  const [codeMessage, setCodeMessage] = useState<string>('')

  const [timeLeft, setTimeLeft] = useState<number>(0)
  const TIMER_KEY = `verify_timer_${email}`

  useEffect(() => {
    const savedExpiry = localStorage.getItem(TIMER_KEY)
    if (savedExpiry) {
      const remaining = Math.floor((parseInt(savedExpiry) - Date.now()) / 1000)
      if (remaining > 0) {
        setTimeLeft(remaining)
      } else {
        localStorage.removeItem(TIMER_KEY)
      }
    }
  }, [email, TIMER_KEY])

  useEffect(() => {
    if (timeLeft <= 0) return

    const interval = setInterval(() => {
      setTimeLeft(prev => {
        if (prev <= 1) {
          localStorage.removeItem(TIMER_KEY)
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => clearInterval(interval)
  }, [timeLeft, TIMER_KEY])

  const startTimer = () => {
    const expiryDate = Date.now() + 60 * 1000 // 60 seconds from now
    localStorage.setItem(TIMER_KEY, expiryDate.toString())
    setTimeLeft(60)
  }

  const handleComplete = async (code: string) => {
    setSuccess(false)
    setErrors(null)
    try {
      const response = await $api.post<BaseServerResponse>('/auth/verify-user', {
        email: email || '',
        verificationCode: code,
      })
      setSuccess(response.data.success)
    } catch (error: any) {
      serverLog('VERIFY_USER_ERROR', error)

      const serverErrors = error.response?.data?.errors || null
      setErrors(serverErrors)
    } finally {
      setCodeMessage('')
    }
  }

  const resendVerifyCode = async () => {
    try {
      if (timeLeft > 0) return
      await $api.post('/auth/resend-verification-code', { email })
      setCodeMessage('Verification code resent successfully.')
      startTimer()
    } catch (error) {
      serverLog('Error resending verification code:', error)
    }
  }

  return (
    <div className="flex flex-1 items-center justify-center">
      <div
        className="flex w-full bg-neutral-900 max-w-[700px] rounded-[20px] p-10 box-border flex-col items-center gap-4"
      >
        <div className="space-y-2 text-center">
          <h2 className="text-white text-2xl font-bold">Confirm email</h2>
          <p className="text-neutral-400 text-sm">We have sent a code to your email</p>
        </div>

        <OTPInput
          maxLength={6}
          onComplete={handleComplete}
          containerClassName="group flex items-center has-[:disabled]:opacity-50"
          render={({ slots }) => (
            <div className="flex gap-2">
              {slots.map((slot, idx) => (
                <Slot key={idx} {...slot} />
              ))}
            </div>
          )}
        />

        {errors &&
          Object.entries(errors as Record<string, string>).map(([field, message]) => (
            <p key={field} className="text-red-500 text-sm">
              {message}
            </p>
          ))}

        {!success && (
          <button
            className={`text-sm transition-all ${
              timeLeft > 0 ? 'text-neutral-500 cursor-not-allowed' : 'text-emerald-500 hover:underline cursor-pointer'
            }`}
            onClick={resendVerifyCode}
            disabled={timeLeft > 0}
          >
            {timeLeft > 0 ? `Resend code in ${timeLeft}s` : 'Resend code'}
          </button>
        )}

        {codeMessage && <p className="text-green-500 text-sm">{codeMessage}</p>}

        {success && (
          <>
            <hr className="border-neutral-300 w-full" />
            <p className="text-emerald-500 text-sm font-medium">
              User has been successfully verified. Please Log in to your account now{' '}
              <Link href="/login" className="underline">
                here
              </Link>
            </p>
          </>
        )}
      </div>
    </div>
  )
}

function Slot(props: SlotProps) {
  return (
    <div
      className={` relative w-12 h-14 text-[20px] flex items-center justify-center transition-all duration-300 border-2
        rounded-xl
        ${props.isActive ? 'border-emerald-500 shadow-[0_0_10px_rgba(16,185,129,0.3)]' : 'border-neutral-700'}
        ${props.char ? 'text-white' : 'text-neutral-500'} `}
    >
      {props.char !== null ? <div>{props.char}</div> : null}

      {/* Caret (blinking bar) */}
      {props.isActive && props.char === null && (
        <div className="absolute inset-0 flex items-center justify-center animate-caret-blink">
          <div className="w-px h-8 bg-white" />
        </div>
      )}
    </div>
  )
}
