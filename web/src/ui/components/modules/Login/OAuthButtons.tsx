'use client'
import { MouseEvent } from 'react'

import { signIn } from 'next-auth/react'

// import GithubIcon from '@/assets/auth/ic_github.svg'
import GoogleIcon from '@/assets/auth/ic_google.svg'
import { cn } from '@/lib/utils/general.utils'

const OAUTH_PROVIDERS = [
  { id: 'google', label: 'Google', Icon: GoogleIcon },
  // { id: 'github', label: 'GitHub', Icon: GithubIcon },
] as const

type OAuthProviderId = (typeof OAUTH_PROVIDERS)[number]['id']

export default function AuthSocialButtons() {
  const handleSignIn = (provider: OAuthProviderId, e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()
    signIn(provider, { redirectTo: '/' })
  }

  return (
    <div className="flex w-full flex-col gap-3 items-center">
      {OAUTH_PROVIDERS.map(({ id, label, Icon }) => (
        <button
          key={id}
          onClick={e => handleSignIn(id, e)}
          className={cn(
            'flex w-full items-center justify-center gap-3',
            'px-5 py-3 rounded-full transition-all duration-300 ease-out cursor-pointer',
            'bg-neutral-800 border border-neutral-700/60 text-foreground-subtle',
            'hover:bg-neutral-700/60 hover:border-neutral-500 hover:text-foreground-secondary',
          )}
        >
          <Icon className="size-5 shrink-0" />
          <span className="text-sm font-semibold whitespace-nowrap">Continue with {label}</span>
        </button>
      ))}
    </div>
  )
}
