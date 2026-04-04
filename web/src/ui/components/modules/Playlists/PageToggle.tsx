'use client'

import { motion } from 'framer-motion'
import { Route } from 'next'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

import { cn } from '@/lib/utils/general.utils'

type Tab = {
  label: string
  href: Route
}

const TABS: Tab[] = [
  { label: 'public playlists', href: '/playlists' },
  { label: 'my playlists', href: '/playlists/my' },
]

export default function PageToggle() {
  const pathname = usePathname()

  return (
    <div className="flex rounded-xl bg-neutral-900/60 p-1 border border-neutral-800 order-2 lg:order-1 relative">
      {TABS.map(tab => {
        const isActive = pathname === tab.href

        return (
          <Link
            key={tab.href}
            href={tab.href}
            className={cn(
              `relative px-4 py-1.5 md:px-5 md:py-2 z-10 block font-semibold md:font-bold uppercase text-[13px]
              tracking-wide`,
              'transition-colors duration-300 rounded-lg',
              isActive ? 'text-neutral-100' : 'text-neutral-500 hover:text-neutral-300',
            )}
          >
            {isActive && (
              <motion.div
                layoutId="active-pill"
                className="absolute inset-0 bg-neutral-800 rounded-lg -z-10 shadow-sm border border-neutral-700/50"
                transition={{ type: 'spring', stiffness: 400, damping: 30 }}
              />
            )}
            {tab.label}
          </Link>
        )
      })}
    </div>
  )
}
