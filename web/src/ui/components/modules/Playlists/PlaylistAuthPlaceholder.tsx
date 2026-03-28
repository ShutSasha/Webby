import Link from 'next/link'

import LockIcon from '@/assets/icons/shared/lock.svg'

export default function PlaylistAuthPlaceholder() {
  return (
    <div className="flex flex-col flex-1 items-center justify-center py-32 px-4 animate-in fade-in duration-500">
      <div
        className="size-20 bg-neutral-900/80 rounded-full flex items-center justify-center mb-6 border
          border-neutral-800"
      >
        <LockIcon className="size-10 text-neutral-500 stroke-1" />
      </div>

      <h2 className="text-2xl font-bold text-neutral-200 mb-3 text-center">Sign in to view your playlists</h2>
      <p className="text-neutral-400 text-center max-w-md mb-8 leading-relaxed">
        Keep track of your favorite videos, create custom collections, and manage your saved content by logging into
        your account.
      </p>

      <Link
        href="/login"
        className="px-8 py-3 bg-emerald-500 hover:bg-emerald-400 text-neutral-900 font-bold rounded-xl transition-colors
          duration-300"
      >
        Sign In
      </Link>
    </div>
  )
}
