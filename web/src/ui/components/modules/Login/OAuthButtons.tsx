'use client'

import { FC, SVGProps } from 'react'

import { signIn } from 'next-auth/react'

import GithubIcon from '@/assets/auth/ic_github.svg' // Переконайся, що файл існує
import GoogleIcon from '@/assets/auth/ic_google.svg'

const BUTTON_CLASSES = `
  flex items-center justify-center gap-1.5 py-1 px-5
  rounded-xl border border-border  cursor-pointer 
  hover:bg-emerald-500 hover:border-transparent 
  group transition-all duration-300
`

const ICON_CLASSES = 'w-7 h-7 text-neutral-300 transition-colors duration-300 ease-out group-hover:text-neutral-900'
const TEXT_CLASSES = 'text-neutral-300 group-hover:text-neutral-900 text-sm transition-all duration-300 font-medium'

interface SocialButtonProps {
  provider: 'google' | 'github'
  label: string
  Icon: FC<SVGProps<SVGSVGElement>>
}

const SocialButton = ({ provider, label, Icon }: SocialButtonProps) => (
  <div className={BUTTON_CLASSES} onClick={() => signIn(provider, { callbackUrl: '/' })}>
    <Icon className={ICON_CLASSES} />
    <p className={TEXT_CLASSES}>{label}</p>
  </div>
)

export default function AuthSocialButtons() {
  return (
    <div className="flex gap-4 w-full justify-center items-center">
      <SocialButton provider="google" label="Google" Icon={GoogleIcon} />
      <SocialButton provider="github" label="Github" Icon={GithubIcon} />
    </div>
  )
}
