'use client'
import { useState } from 'react'

import Link from 'next/link'

import ExpandIcon from '@/assets/icons/ic_expand.svg'
import LogoIcon from '@/assets/icons/ic_logo.svg'
import UserProfileIcon from '@/assets/icons/ic_user_profile.svg'
import { navItems, navMobileItems } from '@/lib/placeholder-data/nav-side'

import NavElement from './NavElement'

export default function MobileNav() {
  const [isExpanded, setIsExpanded] = useState<boolean>(false)

  const toggleSideNav = () => setIsExpanded(prev => !prev)

  return (
    <nav
      className={`flex w-full md:hidden shrink-0 bg-neutral-900 transition-all duration-300 ease-out text-sm leading-5
        fixed bottom-0 border-t border-border`}
    >
      <ul className="flex justify-between items-center list-none w-full gap-2 px-2">
        {navMobileItems.map(item => (
          <li key={item.key}>
            <NavElement
              href={item.href}
              text={item.text}
              isExpanded={isExpanded}
              Icon={item.icon}
              iconSize={item.iconSize}
              isLogo={item.key === 'logo'}
            />
          </li>
        ))}
      </ul>
    </nav>
  )
}
