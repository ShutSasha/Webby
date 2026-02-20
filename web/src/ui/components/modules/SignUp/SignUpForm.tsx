'use client'
import { startTransition, useActionState, useState } from 'react'

import Link from 'next/link'
import { signIn } from 'next-auth/react'
import { useForm } from 'react-hook-form'

import { registerUser } from '@/app/api/auth'
import GoogleIcon from '@/assets/auth/ic_google.svg'
import MailIcon from '@/assets/auth/ic_mail.svg'
import PasswordIcon from '@/assets/auth/ic_password.svg'
import RepeatPasswordIcon from '@/assets/auth/ic_repeat_password.svg'
import UsernameIcon from '@/assets/auth/ic_username.svg'
import AuthInput from '@/components/AuthInput'
import Button from '@/components/Button'

const initialState = {
  success: false,
  errors: null,
}

export default function SignUpForm() {
  const [state, formAction, isPending] = useActionState(registerUser, initialState)
  const { register, handleSubmit } = useForm({
    defaultValues: {
      email: '',
      username: '',
      password: '',
      repeatPassword: '',
    },
  })

  const [showPassword, setShowPassword] = useState<boolean>(false)

  const togglePassword = () => setShowPassword(prev => !prev)

  const onSubmit = (data: any) => {
    const formData = new FormData()
    Object.entries(data).forEach(([key, value]) => formData.append(key, value as string))

    startTransition(() => {
      formAction(formData)
    })
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex w-full flex-col">
      <h1 className="text-white text-xl font-bold text-center mb-4">Sign up</h1>

      <div className="flex flex-col gap-3">
        <AuthInput Icon={MailIcon} {...register('email')} type="email" placeholder="Email" />
        <AuthInput Icon={UsernameIcon} {...register('username')} name="username" type="text" placeholder="Username" />
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
        <div>
          {state.errors &&
            Object.entries(state.errors as Record<string, string>).map(([field, message]) => (
              <p key={field} className="text-red-500 text-sm">
                {message}
              </p>
            ))}
        </div>
        <Button
          disabled={isPending}
          type="submit"
          viewType="Confirm"
          className="text-[16px] leading-[22px] font-semibold w-fit mx-auto"
          paddingClasses="px-5 py-2"
        >
          {isPending ? 'Sending...' : 'Sign Up'}
        </Button>
        <p className="text-center text-sm leading-5 text-neutral-300">
          Already have an account?{' '}
          <Link href="/login" className="text-emerald-500 hover:underline">
            Log in
          </Link>
        </p>
        <hr className="border-neutral-300" />
        <p className="text-center text-sm leading-5 text-neutral-300">or sign up via </p>
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
