'use client'
import { startTransition, useActionState, useEffect, useState } from 'react'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useSession } from 'next-auth/react'
import { useForm } from 'react-hook-form'

import { authenticate } from '@/app/api/auth'
import MailIcon from '@/assets/auth/ic_mail.svg'
import PasswordIcon from '@/assets/auth/ic_password.svg'
import AuthInput from '@/components/AuthInput'
import Button from '@/components/Button'

import ErrorDisplay from './ErrorDisplay'
import AuthSocialButtons from './OAuthButtons'

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

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex w-full flex-col">
      <h1 className="text-white text-xl font-bold text-center mb-4">Login</h1>

      <div className="flex flex-col gap-3 mb-2">
        <AuthInput Icon={MailIcon} {...register('email')} type="email" placeholder="Email" autoComplete="email" />
        <AuthInput
          Icon={PasswordIcon}
          {...register('password')}
          showPassword={showPassword}
          togglePassword={togglePassword}
          type="password"
          placeholder="Password"
        />
        <ErrorDisplay email={watchedEmail} state={state} />
        <input type="hidden" name="redirectTo" value={callbackUrl} />
        <Button
          disabled={isPending}
          type="submit"
          viewType="Confirm"
          className={`text-[16px] leading-[22px] font-semibold w-fit mx-auto ${
            isPending ? 'bg-neutral-700 hover:bg-neutral-700 cursor-not-allowed' : 'cursor-pointer'
          }`}
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

      <AuthSocialButtons />
    </form>
  )
}
