'use client'
import { useState } from 'react'

import Link from 'next/link'

import LogoIcon from '@/assets/icons/ic_logo.svg'
import ExpandIcon from '@/assets/icons/Nav/arrow-right-from-line.svg'
import { navDesktopElements } from '@/lib/placeholder-data/nav-side'
import { cn } from '@/lib/utils/general.utils'
import { useNotificationPopupStore } from '@/stores/notification-popup.store'

import { DesktopNavElement } from './DesktopNavElement'
import { UserProfile } from './UserProfile'
import GlobalSearchModal from '../GlobalSearch/GlobalSearchModal'

type Props = {
  isAdmin?: boolean
}

export default function DesktopNav({ isAdmin }: Props) {
  const [isExpanded, setIsExpanded] = useState<boolean>(false)
  const [isSearchOpen, setIsSearchOpen] = useState<boolean>(false)
  const unreadCount = useNotificationPopupStore(state => state.unreadCount)

  const visibleNavElements = navDesktopElements.filter(item => !item.requireAdmin || isAdmin)

  const toggleSideNav = () => setIsExpanded(prev => !prev)

  const handleNavClick = (action?: string) => {
    if (action === 'search') {
      setIsSearchOpen(true)
    }
  }

  return (
    <>
      <aside className="block h-screen sticky top-0 shrink-0 z-50">
        <nav
          className={cn(
            `flex flex-col h-full bg-[#0A0A0A] border-r border-neutral-800/50 py-3 px-3 xl:py-3 xl:px-3 2xl:py-4
            2xl:px-3 transition-all duration-300 ease-in-out`,
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
              <LogoIcon
                className={cn(
                  'size-8 text-neutral-50 transition-transform group-hover:scale-110 shrink-0',
                  'bg-linear-to-br from-emerald-400 to-emerald-600 rounded-[10px]',
                )}
              />
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
            {visibleNavElements.map(item => (
              <li key={item.key}>
                <DesktopNavElement
                  key={item.key}
                  href={item.href}
                  onClick={() => handleNavClick(item.action)}
                  text={item.text}
                  isExpanded={isExpanded}
                  Icon={item.icon}
                  iconSize={item.iconSize || 'size-6'}
                  badgeCount={item.key === 'notifications' ? unreadCount : undefined}
                />
              </li>
            ))}
          </ul>

          <div className="mt-auto pt-4 border-t border-neutral-800/50 w-full flex items-center justify-center">
            <UserProfile isExpanded={isExpanded} iconSize="size-10" />
          </div>
        </nav>
      </aside>
      <GlobalSearchModal isOpen={isSearchOpen} onClose={() => setIsSearchOpen(false)} />
    </>
  )
}
