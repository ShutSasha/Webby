'use client'
import { useActionState, useState } from 'react'

import Link from 'next/link'
import { signIn } from 'next-auth/react'

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
  const [email, setEmail] = useState<string>('')
  const [username, setUsername] = useState<string>('')
  const [password, setPassword] = useState<string>('')
  const [repeatPassword, seRepeatPassword] = useState<string>('')

  const [state, formAction, isPending] = useActionState(registerUser, initialState)
  const [showPassword, setShowPassword] = useState<boolean>(false)

  const togglePassword = () => setShowPassword(prev => !prev)

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value)
  }

  const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUsername(e.target.value)
  }

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value)
  }

  const handleRepeatPasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    seRepeatPassword(e.target.value)
  }

  return (
    <form action={formAction} className="flex w-full flex-col">
      <h1 className="text-white text-xl font-bold text-center mb-4">Sign up</h1>

      <div className="flex flex-col gap-3">
        <AuthInput
          value={email}
          onChange={handleEmailChange}
          Icon={MailIcon}
          name="email"
          type="email"
          placeholder="Email"
          autoComplete="email"
        />
        <AuthInput
          value={username}
          onChange={handleUsernameChange}
          Icon={UsernameIcon}
          name="username"
          type="text"
          placeholder="Username"
          autoComplete="username"
        />
        <AuthInput
          value={password}
          onChange={handlePasswordChange}
          Icon={PasswordIcon}
          showPassword={showPassword}
          togglePassword={togglePassword}
          name="password"
          type="password"
          placeholder="Password"
        />
        <AuthInput
          value={repeatPassword}
          onChange={handleRepeatPasswordChange}
          Icon={RepeatPasswordIcon}
          showPassword={showPassword}
          togglePassword={togglePassword}
          name="repeatPassword"
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
