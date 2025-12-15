'use client'
import { useState } from 'react'

import Link from 'next/link'

import GoogleIcon from '@/assets/auth/ic_google.svg'
import MailIcon from '@/assets/auth/ic_mail.svg'
import PasswordIcon from '@/assets/auth/ic_password.svg'
import RepeatPasswordIcon from '@/assets/auth/ic_repeat_password.svg'
import UsernameIcon from '@/assets/auth/ic_username.svg'
import AuthInput from '@/components/AuthInput'
import Button from '@/components/Button'

export default function SignUpForm() {
  const [email, setEmail] = useState<string>('')
  const [username, setUsername] = useState<string>('')
  const [password, setPassword] = useState<string>('')
  const [repeatPassword, setRepeatPassword] = useState<string>('')

  const handleInputsChange = (key: 'email' | 'username' | 'password' | 'repeatPassword') => (value: string) => {
    if (key === 'email') setEmail(value)
    if (key === 'username') setUsername(value)
    if (key === 'password') setPassword(value)
    if (key === 'repeatPassword') setRepeatPassword(value)
  }

  return (
    <form className="flex w-full flex-col">
      <h1 className="text-white text-xl font-bold text-center mb-4">Sign up</h1>

      <div className="flex flex-col gap-3">
        <AuthInput
          Icon={MailIcon}
          name="email"
          type="email"
          value={email}
          onChange={handleInputsChange('email')}
          placeholder="Email"
        />
        <AuthInput
          Icon={UsernameIcon}
          name="username"
          type="text"
          value={username}
          onChange={handleInputsChange('username')}
          placeholder="Username"
        />
        <AuthInput
          Icon={PasswordIcon}
          name="password"
          type="password"
          value={password}
          onChange={handleInputsChange('password')}
          placeholder="Password"
        />
        <AuthInput
          Icon={RepeatPasswordIcon}
          name="repeatPassword"
          type="password"
          value={repeatPassword}
          onChange={handleInputsChange('repeatPassword')}
          placeholder="Repeat Password"
        />
        <Button
          type="submit"
          viewType="Confirm"
          className="text-[16px] leading-[22px] font-semibold w-fit mx-auto"
          paddingClasses="px-5 py-2"
        >
          Sign Up
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
        <GoogleIcon className="w-11 h-11 hover:text-emerald-500 transition-colors duration-300 ease-out cursor-pointer" />
      </div>
    </form>
  )
}
