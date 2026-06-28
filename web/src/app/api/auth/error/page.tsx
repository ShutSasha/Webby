import Link from 'next/link'

import LockIcon from '@/assets/icons/shared/lock.svg'

export default async function AuthErrorPage({
  searchParams,
}: {
  searchParams: Promise<{ error?: string; message?: string }>
}) {
  const { error, message } = await searchParams

  let errorTitle = 'Authentication Error'
  let errorMessage = 'An unknown error occurred during sign-in. Please try again later.'

  if (error === 'AccessDenied') {
    errorTitle = 'Access Denied'
    errorMessage =
      message || 'You do not have permission to sign in. Your account might not be whitelisted or has been blocked.'
  } else if (error === 'Configuration') {
    errorTitle = 'Configuration Error'
    errorMessage = 'There is a problem with the server authentication configuration.'
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background p-4">
      <div
        className="w-full max-w-md flex flex-col items-center p-8 rounded-2xl bg-surface border border-border shadow-2xl
          shadow-black/5 text-center transition-all"
      >
        <div className="w-16 h-16 flex items-center justify-center rounded-full bg-red-500/10 text-red-500 mb-6">
          <LockIcon className="w-8 h-8 stroke-2" />
        </div>
        <h1 className="mb-3 text-2xl font-bold text-foreground-secondary">{errorTitle}</h1>
        <p className="mb-8 text-sm text-foreground-muted leading-relaxed">{errorMessage}</p>
        <Link
          href="/login"
          className="px-5 py-2.5 rounded-xl bg-linear-to-r w-full from-emerald-500/10 to-emerald-500/5 border
            border-emerald-500/20 text-emerald-400 font-semibold text-sm hover:from-emerald-500/20
            hover:to-emerald-500/10 hover:border-emerald-500/40 transition-all duration-300"
        >
          Return to Sign In
        </Link>
      </div>
    </div>
  )
}
