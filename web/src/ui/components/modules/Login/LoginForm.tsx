'use client'
import { startTransition, useActionState, useEffect, useState } from 'react'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useSession } from 'next-auth/react'
import { useForm } from 'react-hook-form'

import MailIcon from '@/assets/auth/ic_mail.svg'
import PasswordIcon from '@/assets/auth/ic_password.svg'
import { authenticate } from '@/lib/actions/auth.actions'
import AuthInput from '@/ui/components/modules/Auth/AuthInput'
import Button from '@/ui/components/shared/Button'

import AuthErrorDisplay from './ErrorDisplay'
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
    <form onSubmit={handleSubmit(onSubmit)} className="flex w-full flex-col gap-6">
      <div className="flex flex-col gap-1 text-center">
        <h1 className="text-foreground-secondary text-2xl font-bold">Log in to Webby</h1>
        <p className="text-sm text-foreground-muted">Welcome back! Please enter your details.</p>
      </div>

      <div className="flex flex-col gap-4">
        <AuthInput Icon={MailIcon} {...register('email')} type="email" placeholder="Email" autoComplete="email" />

        <div className="flex flex-col gap-2">
          <AuthInput
            Icon={PasswordIcon}
            {...register('password')}
            showPassword={showPassword}
            togglePassword={togglePassword}
            type="password"
            placeholder="Password"
          />

          <Link
            href={'/forgot-password'}
            className="text-xs text-foreground-muted hover:text-foreground-tertiary transition-colors self-end"
          >
            Forgot password?
          </Link>
        </div>
      </div>

      <AuthErrorDisplay email={watchedEmail} state={state} />
      <input type="hidden" name="redirectTo" value={callbackUrl} />

      <Button
        disabled={isPending}
        type="submit"
        viewType={isPending ? 'loading' : 'confirm'}
        className="w-full text-[16px]"
        paddingClasses="py-3"
      >
        {isPending ? 'Logging in...' : 'Log in'}
      </Button>

      <div className="flex items-center gap-3">
        <div className="h-px w-full bg-neutral-800" />
        <span className="text-xs text-foreground-faint uppercase tracking-wider whitespace-nowrap">
          or continue with
        </span>
        <div className="h-px w-full bg-neutral-800" />
      </div>

      <AuthSocialButtons />

      <p className="text-center text-sm text-foreground-muted mt-2">
        Don&apos;t have an account?{' '}
        <Link href="/sign-up" className="text-foreground-strong font-semibold hover:underline">
          Sign up
        </Link>
      </p>
    </form>
  )
}
