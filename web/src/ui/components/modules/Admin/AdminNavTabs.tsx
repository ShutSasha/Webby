'use client'

import { Route } from 'next'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

import { cn } from '@/lib/utils/general.utils'

const adminNav = [
  { name: 'Overview', href: '/admin' },
  { name: 'Users Management', href: '/admin/users' },
  { name: 'Reports', href: '/admin/reports' },
]

export default function AdminNavTabs() {
  const pathname = usePathname()

  return (
    <nav className="flex items-center gap-2 border-b border-neutral-800/60 pb-px">
      {adminNav.map(tab => {
        const isActive = pathname === tab.href
        return (
          <Link
            key={tab.name}
            href={tab.href as Route}
            className={cn(
              'px-5 py-2.5 text-sm font-medium transition-all relative rounded-t-lg',
              isActive
                ? 'text-emerald-400 bg-emerald-400/10'
                : 'text-foreground-muted hover:text-foreground-tertiary hover:bg-neutral-800/50',
            )}
          >
            {tab.name}

            {isActive && <span className="absolute -bottom-px left-0 w-full h-0.5 bg-emerald-500 rounded-t-full" />}
          </Link>
        )
      })}
    </nav>
  )
}
