'use client'
import { useState } from 'react'

import Link from 'next/link'

import ExpandIcon from '@/assets/icons/ic_expand.svg'
import LogoIcon from '@/assets/icons/ic_logo.svg'
import UserProfileIcon from '@/assets/icons/ic_user_profile.svg'
import { navItems } from '@/lib/placeholder-data/nav-side'

import NavElement from './NavElement'

export default function SideNav() {
  const [isExpanded, setIsExpanded] = useState<boolean>(false)

  const toggleSideNav = () => setIsExpanded(prev => !prev)

  return (
    <nav
      className={`hidden md:flex flex-col h-screen shrink-0 bg-neutral-900 p-2 transition-all duration-300 ease-out
        ${isExpanded ? 'w-60 items-start' : 'w-12 items-center'} text-sm leading-5 sticky top-0`}
    >
      <div className={`flex w-full items-center ${isExpanded ? 'mb-3 flex-row justify-between' : 'flex-col'}`}>
        <Link href={'/'}>
          <LogoIcon className="h-8 w-8 text-white" />
        </Link>
        <button onClick={toggleSideNav}>
          <ExpandIcon
            className={`h-6 w-6 cursor-pointer transition-all duration-300 hover:text-emerald-500
              ${isExpanded ? 'my-0 rotate-180' : 'my-3 rotate-0'} `}
          />
        </button>
      </div>

      <ul className={`flex flex-col gap-2 list-none w-full ${isExpanded ? '' : 'items-center'}`}>
        {navItems.map(item => (
          <li key={item.key}>
            <NavElement
              href={item.href}
              text={item.text}
              isExpanded={isExpanded}
              Icon={item.icon}
              iconSize={item.iconSize}
            />
          </li>
        ))}
      </ul>

      <NavElement
        href={'/sign-up'}
        text={'Log in | Sign up'}
        isExpanded={isExpanded}
        Icon={UserProfileIcon}
        className={`mt-auto ${isExpanded ? 'w-full' : ''}`}
        iconSize="size-6"
      />
    </nav>
  )
}
