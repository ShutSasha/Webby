'use client'

import Image from 'next/image'
import Link from 'next/link'
import { signOut, useSession } from 'next-auth/react'

import LogoutIcon from '@/assets/auth/ic_logout.svg'
import UserProfileIcon from '@/assets/icons/ic_user_profile.svg'
import { cn } from '@/lib/utils/general.utils'
import { BLUR_DATA_URLS } from '@/ui/images'

import { DesktopNavElement } from './NavElements'

type UserProfileProps = {
  isExpanded: boolean
  iconSize: string
}

export function UserProfile({ isExpanded, iconSize }: UserProfileProps) {
  const { data: session, status } = useSession()

  const handleLogout = async () => {
    await signOut({ redirectTo: '/login' })
  }

  if (status === 'loading') {
    return (
      <div
        className={cn(
          iconSize,
          'bg-neutral-800 animate-pulse rounded-full shrink-0',
          isExpanded && 'w-full rounded-xl h-12',
        )}
      />
    )
  }

  if (!session?.user) {
    return (
      <DesktopNavElement
        href={'/sign-up'}
        text={'Log in | Sign up'}
        isExpanded={isExpanded}
        Icon={UserProfileIcon}
        className={`mt-auto ${isExpanded ? 'w-full' : ''}`}
        iconSize={iconSize}
        isUserProfile
      />
    )
  }

  return (
    <div
      className={cn(
        'flex items-center justify-center w-full rounded-xl transition-colors duration-300',
        isExpanded ? 'justify-between px-2 py-1.5 hover:bg-neutral-800/60' : 'justify-center',
      )}
    >
      <Link
        href={`/profile/${session.user.id}`}
        className={cn('flex items-center justify-center', !isExpanded && 'group', isExpanded && 'gap-3')}
      >
        <Image
          src={session.user.image ?? 'https://i.pinimg.com/originals/44/64/20/4464203a781eed3650f1fdd624c4d02a.jpg'}
          alt="user profile picture"
          width={100}
          height={100}
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral800']}
          className={cn(
            iconSize,
            'object-cover rounded-full shrink-0 transition-all duration-300',
            !isExpanded && 'group-hover:ring-2 group-hover:ring-emerald-500/50',
          )}
        />

        <div
          className={cn(
            'flex flex-col transition-all duration-300 ease-in-out',
            isExpanded ? 'w-auto opacity-100' : 'w-0 opacity-0',
          )}
        >
          <span className="text-sm font-semibold text-neutral-200 truncate max-w-[120px]">{session.user.username}</span>
          <span className="text-xs text-neutral-500">View profile</span>
        </div>
      </Link>

      <button
        onClick={handleLogout}
        className={cn(
          `p-2 rounded-lg text-neutral-500 hover:text-red-500 hover:bg-red-500/10 transition-colors duration-300
          shrink-0`,
          !isExpanded && 'hidden',
        )}
        title="Log out"
      >
        <LogoutIcon className="size-5" />
      </button>
    </div>
  )
}
