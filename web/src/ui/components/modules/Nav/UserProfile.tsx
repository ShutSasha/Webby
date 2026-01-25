'use client'

import Image from 'next/image'
import Link from 'next/link'
import { useSession } from 'next-auth/react'

import UserProfileIcon from '@/assets/icons/ic_user_profile.svg'

import { DesktopNavElement } from './NavElements'

type UserProfileProps = {
  isExpanded: boolean
  iconSize: string
}

export function UserProfile({ isExpanded, iconSize }: UserProfileProps) {
  const { data: session, status } = useSession()

  if (status === 'loading') {
    return (
      <div
        className={`${iconSize} bg-neutral-800 mt-auto animate-pulse rounded-full ${isExpanded ? 'w-full' : ''}`}
      ></div>
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
      />
    )
  }

  return (
    <Link
      href={'/'}
      className={`flex h-fit cursor-pointer flex-col items-center rounded-sm mt-auto ${
        isExpanded
          ? 'flex-row transition-colors duration-300 ease-out hover:bg-emerald-500 hover:text-neutral-800 w-full'
          : 'flex-col hover:text-emerald-500'
        }`}
    >
      <div className={`w-full flex-1 flex flex-row items-center ${isExpanded ? 'gap-2 p-1 ' : 'gap-0 p-0 '}`}>
        <Image
          src={
            session.user.image
              ? session.user.image
              : 'https://i.pinimg.com/originals/44/64/20/4464203a781eed3650f1fdd624c4d02a.jpg'
          }
          alt="user profile picture"
          width={24}
          height={24}
          className={`${iconSize} object-cover rounded-full`}
        />
        <p
          className={`overflow-hidden whitespace-nowrap transition-opacity font-medium
            ${isExpanded ? 'opacity-100 ' : 'w-0 h-0 opacity-0'}`}
        >
          {session.user.username}
        </p>
      </div>
    </Link>
  )
}
