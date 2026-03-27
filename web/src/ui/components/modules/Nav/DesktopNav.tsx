'use client'
import { useState } from 'react'

import Link from 'next/link'

import LogoIcon from '@/assets/icons/ic_logo.svg'
import ExpandIcon from '@/assets/icons/Nav/arrow-right-from-line.svg'
import { navDesktopElements } from '@/lib/placeholder-data/nav-side'
import { cn } from '@/lib/utils/utils'

import { DesktopNavElement } from './NavElements'
import { UserProfile } from './UserProfile'

export default function DesktopNav() {
  const [isExpanded, setIsExpanded] = useState<boolean>(false)

  const toggleSideNav = () => setIsExpanded(prev => !prev)

  return (
    <aside className="hidden lg:block h-screen sticky top-0 shrink-0 z-50">
      <nav
        className={cn(
          `flex flex-col h-full bg-[#0A0A0A] border-r border-neutral-800/50 py-4 px-3 transition-all duration-300
          ease-in-out`,
          isExpanded ? 'w-64 items-stretch' : 'w-20 items-center',
        )}
      >
        <div
          className={cn(
            'flex mb-8 w-full',
            isExpanded ? 'flex-row items-center justify-between' : 'flex-col items-center gap-6',
          )}
        >
          <Link href={'/'} className={cn('flex items-center group', isExpanded ? 'gap-3' : 'justify-center')}>
            <LogoIcon className="size-8 text-emerald-500 transition-transform group-hover:scale-110 shrink-0" />
            <span
              className={cn(
                'text-2xl tracking-tight font-bold text-neutral-100 transition-all duration-300',
                isExpanded ? 'opacity-100 w-auto ml-1' : 'w-0 h-0 opacity-0 overflow-hidden',
              )}
            >
              Webby
            </span>
          </Link>

          <button
            onClick={toggleSideNav}
            className="p-1.5 rounded-lg hover:bg-neutral-800/60 text-neutral-400 hover:text-neutral-100
              transition-colors shrink-0"
          >
            <ExpandIcon
              className={cn(
                'size-5 transition-transform duration-300 stroke-[1.5px]',
                isExpanded ? 'rotate-180' : 'rotate-0',
              )}
            />
          </button>
        </div>

        <ul className="flex flex-col gap-1.5 list-none w-full flex-1">
          {navDesktopElements.map(item => (
            <li key={item.key}>
              <DesktopNavElement
                href={item.href}
                text={item.text}
                isExpanded={isExpanded}
                Icon={item.icon}
                iconSize={item.iconSize || 'size-6'}
              />
            </li>
          ))}
        </ul>

        <div className="mt-auto pt-4 border-t border-neutral-800/50 w-full flex items-center justify-center">
          <UserProfile isExpanded={isExpanded} iconSize="size-10" />
        </div>
      </nav>
    </aside>
  )
}
