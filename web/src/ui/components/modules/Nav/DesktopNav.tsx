'use client'
import { useState } from 'react'

import Link from 'next/link'

import ExpandIcon from '@/assets/icons/ic_expand.svg'
import LogoIcon from '@/assets/icons/ic_logo.svg'
import UserProfileIcon from '@/assets/icons/ic_user_profile.svg'
import { navDesktopElements } from '@/lib/placeholder-data/nav-side'

import { DesktopNavElement } from './NavElements'

export default function DesktopNav() {
  const [isExpanded, setIsExpanded] = useState<boolean>(false)

  const toggleSideNav = () => setIsExpanded(prev => !prev)

  return (
    <aside className={'hidden md:block h-screen p-4 transition-all duration-300 ease-in-out sticky top-0'}>
      <nav
        className={`flex flex-col h-full shrink-0 bg-neutral-900 p-2.5 transition-all duration-300 ease-in-out
          rounded-lg ${isExpanded ? 'w-60 items-start' : 'w-14 items-center'} text-sm leading-5`}
      >
        <div className={`flex w-full items-center ${isExpanded ? 'mb-3 flex-row justify-between' : 'flex-col'}`}>
          <Link href={'/'}>
            <LogoIcon className="h-8 w-8 text-zinc-950" />
          </Link>
          <button onClick={toggleSideNav}>
            <ExpandIcon
              className={`h-6 w-6 cursor-pointer transition-all duration-300 hover:text-emerald-500
                ${isExpanded ? 'my-0 rotate-180' : 'my-3 rotate-0'} `}
            />
          </button>
        </div>

        <ul className={`flex flex-col gap-2 list-none w-full ${isExpanded ? '' : 'items-center'}`}>
          {navDesktopElements.map(item => (
            <li key={item.key}>
              <DesktopNavElement
                href={item.href}
                text={item.text}
                isExpanded={isExpanded}
                Icon={item.icon}
                iconSize={item.iconSize}
              />
            </li>
          ))}
        </ul>

        <DesktopNavElement
          href={'/sign-up'}
          text={'Log in | Sign up'}
          isExpanded={isExpanded}
          Icon={UserProfileIcon}
          className={`mt-auto ${isExpanded ? 'w-full' : ''}`}
          iconSize="size-6"
        />
      </nav>
    </aside>
  )
}
