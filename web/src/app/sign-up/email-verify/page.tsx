'use client'

import { useState } from 'react'

// eslint-disable-next-line import/named
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

        {!success && <button className="text-emerald-500 text-sm hover:underline cursor-pointer">Resend code</button>}

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
