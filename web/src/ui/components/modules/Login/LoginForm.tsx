'use client'
import { useActionState, useEffect, useState } from 'react'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useSession } from 'next-auth/react'

import { authenticate } from '@/app/api/auth'
import GoogleIcon from '@/assets/auth/ic_google.svg'
import MailIcon from '@/assets/auth/ic_mail.svg'
import PasswordIcon from '@/assets/auth/ic_password.svg'
import AuthInput from '@/components/AuthInput'
import Button from '@/components/Button'

const initialState = {
  success: false,
  errors: null,
}

export default function LoginForm() {
  const [email, setEmail] = useState<string>('')
  const [password, setPassword] = useState<string>('')
  const [state, formAction, isPending] = useActionState(authenticate, initialState)
  const { update } = useSession()
  const router = useRouter()

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

  const handleInputsChange = (key: 'email' | 'password') => (value: string) => {
    if (key === 'email') setEmail(value)
    if (key === 'password') setPassword(value)
  }

  return (
    <form action={formAction} className="flex w-full flex-col">
      <h1 className="text-white text-xl font-bold text-center mb-4">Login</h1>

      <div className="flex flex-col gap-3">
        <AuthInput
          Icon={MailIcon}
          name="email"
          type="email"
          value={email}
          onChange={handleInputsChange('email')}
          placeholder="Email"
          autoComplete="email"
        />

        <AuthInput
          Icon={PasswordIcon}
          name="password"
          type="password"
          value={password}
          onChange={handleInputsChange('password')}
          placeholder="Password"
        />
        <div className={`${state?.errors ? 'block' : 'hidden'}`}>
          {state?.errors &&
            Object.entries(state.errors as Record<string, string>).map(([field, message]) => (
              <p key={field} className="text-red-500 text-sm">
                {message}
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
        <GoogleIcon className="w-11 h-11 hover:text-emerald-500 transition-colors duration-300 ease-out cursor-pointer" />
      </div>
    </form>
  )
}
