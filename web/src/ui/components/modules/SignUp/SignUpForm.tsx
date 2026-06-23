'use client'
import { startTransition, useActionState, useState } from 'react'

import Link from 'next/link'
import { useForm } from 'react-hook-form'

import MailIcon from '@/assets/auth/ic_mail.svg'
import PasswordIcon from '@/assets/auth/ic_password.svg'
import RepeatPasswordIcon from '@/assets/auth/ic_repeat_password.svg'
import UsernameIcon from '@/assets/auth/ic_username.svg'
import { registerUser } from '@/lib/actions/auth.actions'
import AuthInput from '@/ui/components/modules/Auth/AuthInput'
import Button from '@/ui/components/shared/Button'

import AuthSocialButtons from '../Login/OAuthButtons'

const initialState = {
  success: false,
  errors: {},
}

export default function SignUpForm() {
  const [state, formAction, isPending] = useActionState(registerUser, initialState)
  const [showPassword, setShowPassword] = useState<boolean>(false)
  const togglePassword = () => setShowPassword(prev => !prev)

  const { register, handleSubmit } = useForm({
    defaultValues: {
      email: '',
      username: '',
      password: '',
      repeatPassword: '',
    },
  })

  const onSubmit = (data: any) => {
    const formData = new FormData()
    Object.entries(data).forEach(([key, value]) => formData.append(key, value as string))

    startTransition(() => {
      formAction(formData)
    })
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex w-full flex-col gap-6">
      <div className="flex flex-col gap-1 text-center">
        <h1 className="text-foreground-strong text-2xl font-bold">Create account</h1>
        <p className="text-sm text-foreground-muted">Join Webby to start watching and chatting.</p>
      </div>
      <div className="flex flex-col gap-4">
        <AuthInput Icon={UsernameIcon} {...register('username')} name="username" type="text" placeholder="Username" />
        <AuthInput Icon={MailIcon} {...register('email')} type="email" placeholder="Email" />
        <AuthInput
          Icon={PasswordIcon}
          {...register('password')}
          showPassword={showPassword}
          togglePassword={togglePassword}
          type="password"
          placeholder="Password"
        />
        <AuthInput
          Icon={RepeatPasswordIcon}
          {...register('repeatPassword')}
          showPassword={showPassword}
          togglePassword={togglePassword}
          type="password"
          placeholder="Repeat Password"
        />
      </div>

      {state.errors && (
        <div className="flex flex-col gap-1">
          {Object.entries(state.errors as Record<string, string>).map(([field, message]) => (
            <p key={field} className="text-red-500 text-sm">
              {message}
            </p>
          ))}
        </div>
      )}
      <Button
        disabled={isPending}
        type="submit"
        viewType={isPending ? 'loading' : 'confirm'}
        className="w-full text-[16px]"
        paddingClasses="py-3"
      >
        {isPending ? 'Creating account...' : 'Create account'}
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
        Already have an account?{' '}
        <Link href="/login" className="text-foreground-strong font-semibold hover:underline">
          Log in
        </Link>
      </p>
    </form>
  )
}
