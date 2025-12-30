'use client'
import { useState } from 'react'

import { navMobileItems } from '@/lib/placeholder-data/nav-side'

import NavElement from './NavElement'

export default function MobileNav() {
  const [isExpanded, setIsExpanded] = useState<boolean>(false)

  const toggleSideNav = () => setIsExpanded(prev => !prev)

  return (
    <nav
      className={`flex w-full md:hidden bg-neutral-900 transition-all duration-300 ease-out text-sm leading-5 fixed
        bottom-0 border-t border-border`}
    >
      <ul className="flex justify-between items-center list-none w-full gap-2 max-[360px]:px-1.5 px-4">
        {navMobileItems.map(item => (
          <li key={item.key} className="min-w-0">
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
