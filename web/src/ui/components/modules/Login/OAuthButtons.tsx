'use client'

import { signIn } from 'next-auth/react'

import GoogleIcon from '@/assets/auth/ic_google.svg'

export default function GoogleButton() {
  return (
    <div
      className="flex items-center justify-center gap-1.5 w-full py-1 rounded-xl border border-border
        hover:bg-emerald-500 max-w-30 cursor-pointer hover:border-transparent group transition-all duration-300"
      onClick={() => signIn('google', { callbackUrl: '/' })}
    >
      <GoogleIcon
        className="w-7 h-7 text-neutral-300 transition-colors duration-300 ease-out group-hover:text-neutral-900"
      />
      <p className="text-neutral-300 group-hover:text-neutral-900 text-sm transition-all duration-300 font-medium">
        Google
      </p>
    </div>
  )
}
