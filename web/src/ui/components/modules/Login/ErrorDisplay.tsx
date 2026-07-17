'use client'

import Link from 'next/link'

import { resendVerifyCode, ReturnAuthError } from '@/lib/actions/auth.actions'

type Props = {
  state: ReturnAuthError
  email: string
  className?: string
}

export default function AuthErrorDisplay({ email, state, className }: Props) {
  return (
    <div className={`${state?.errors ? 'block' : 'hidden'} ${className}`}>
      {state?.errors &&
        Object.entries(state.errors as Record<string, string>).map(([field, message]) => (
          <p key={field} className="text-red-500 text-sm">
            {message}
            {message === `User isn't verified` && (
              <span>
                {'. '}
                Verify it{' '}
                <Link
                  href={`/sign-up/email-verify?email=${email}`}
                  className="underline"
                  onClick={() => resendVerifyCode({ email })}
                >
                  here
                </Link>
              </span>
            )}
          </p>
        ))}
    </div>
  )
}
