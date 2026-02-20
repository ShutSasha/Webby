'use client'
import { startTransition, useActionState, useEffect, useState } from 'react'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { signIn, useSession } from 'next-auth/react'
import { useForm } from 'react-hook-form'

import { authenticate } from '@/app/api/auth'
import GoogleIcon from '@/assets/auth/ic_google.svg'
import MailIcon from '@/assets/auth/ic_mail.svg'
import PasswordIcon from '@/assets/auth/ic_password.svg'
import AuthInput from '@/components/AuthInput'
import Button from '@/components/Button'
import { serverLog } from '@/lib/utils/utils'

const initialState = {
  success: false,
  errors: null,
}

export default function LoginForm() {
  const [state, formAction, isPending] = useActionState(authenticate, initialState)
  const [showPassword, setShowPassword] = useState<boolean>(false)
  const togglePassword = () => setShowPassword(prev => !prev)
  const router = useRouter()
  const { update } = useSession()

  const { register, handleSubmit, watch } = useForm({
    defaultValues: {
      email: '',
      password: '',
    },
  })

  const watchedEmail = watch('email')

  const onSubmit = (data: any) => {
    const formData = new FormData()
    Object.entries(data).forEach(([key, value]) => formData.append(key, value as string))

    startTransition(() => {
      formAction(formData)
    })
  }

  const callbackUrl = '/'

  useEffect(() => {
    if (!state.success) return

    const handleSuccess = async () => {
      await update()
      router.replace('/')
      router.refresh()
    }

    handleSuccess()
  }, [state.success])

  const resendVerifyCode = async () => {
    const api = process.env.NEXT_PUBLIC_API_URL

    try {
      await fetch(`${api}/auth/resend-verification-code`, {
        method: 'POST',
        body: JSON.stringify({ email: watchedEmail }),
        headers: {
          'Content-Type': 'application/json',
        },
        keepalive: true,
      })
    } catch (error) {
      serverLog('Error resending verification code:', error)
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex w-full flex-col">
      <h1 className="text-white text-xl font-bold text-center mb-4">Login</h1>

      <div className="flex flex-col gap-3">
        <AuthInput Icon={MailIcon} {...register('email')} type="email" placeholder="Email" autoComplete="email" />
        <AuthInput
          Icon={PasswordIcon}
          {...register('password')}
          showPassword={showPassword}
          togglePassword={togglePassword}
          type="password"
          placeholder="Password"
        />
        <div className={`${state?.errors ? 'block' : 'hidden'}`}>
          {state?.errors &&
            Object.entries(state.errors as Record<string, string>).map(([field, message]) => (
              <p key={field} className="text-red-500 text-sm">
                {message}
                {message === `User isn't verified` ? (
                  <span>
                    {'. '}
                    Verify it{' '}
                    <Link
                      href={`/sign-up/email-verify?email=${watchedEmail}`}
                      className="underline"
                      onClick={resendVerifyCode}
                    >
                      here
                    </Link>
                  </span>
                ) : (
                  ''
                )}
              </p>
            ))}
        </div>
        <input type="hidden" name="redirectTo" value={callbackUrl} />
        <Button
          disabled={isPending}
          type="submit"
          viewType="Confirm"
          className="text-[16px] leading-[22px] font-semibold w-fit mx-auto"
          paddingClasses="px-5 py-2"
        >
          {isPending ? 'Sending...' : 'Log in'}
        </Button>
        <p className="text-center text-sm leading-5 text-neutral-300">
          {`Haven't`} an account yet?{' '}
          <Link href="/sign-up" className="text-emerald-500 hover:underline">
            Sign up
          </Link>
        </p>
        <hr className="border-neutral-300" />
        <p className="text-center text-sm leading-5 text-neutral-300">or log in via </p>
      </div>
      <div className="flex flex-row gap-1 items-center justify-center">
        <GoogleIcon
          className="w-11 h-11 hover:text-emerald-500 transition-colors duration-300 ease-out cursor-pointer"
          onClick={() => signIn('google', { callbackUrl: '/' })}
        />
      </div>
    </form>
  )
}
