'use client'

import { motion } from 'framer-motion'
import { Route } from 'next'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

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
    <div className="flex rounded-2xl bg-neutral-800 order-2 lg:order-1 relative">
      {TABS.map(tab => {
        const isActive = pathname === tab.href

        return (
          <Link
            key={tab.href}
            href={tab.href}
            className={`relative px-4 py-1.5 font-semibold md:px-5 md:py-2 z-10 block md:font-bold uppercase text-[14px]
            leading-5.5 transition-colors duration-300
            ${isActive ? 'text-neutral-900' : 'text-neutral-600 hover:text-neutral-300'}`}
          >
            {isActive && (
              <motion.div
                layoutId="active-pill"
                className="absolute inset-0 bg-emerald-500 rounded-2xl -z-10"
                transition={{ type: 'spring', bounce: 0.2, duration: 0.45 }}
              />
            )}
            {tab.label}
          </Link>
        )
      })}
    </div>
  )
}
