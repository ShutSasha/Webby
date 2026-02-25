'use client'

import Link from 'next/link'

import { resendVerifyCode, ReturnAuthError } from '@/app/api/auth'

type Props = {
  state: ReturnAuthError
  email: string
}

export default function ErrorDisplay({ email, state }: Props) {
  return (
    <div className={`${state?.errors ? 'block' : 'hidden'}`}>
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
