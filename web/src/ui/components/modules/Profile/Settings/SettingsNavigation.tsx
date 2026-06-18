'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

import { cn } from '@/lib/utils/general.utils'

type Props = {
  id: string
}

export default function SettingsNavigation({ id }: Props) {
  const pathname = usePathname()

  const navItems = [
    { path: '/profile-settings', label: 'Profile settings' },
    { path: '/secure', label: 'Secure' },
    { path: '/achievements', label: 'Achievements' },
  ]

  return (
    <nav className="flex flex-col gap-2 w-full md:w-64 shrink-0">
      {navItems.map(item => {
        const isActive = pathname.endsWith(item.path)

        return (
          <Link
            key={item.path}
            href={`/profile/${id}${item.path}`}
            className={cn(
              `relative px-5 py-3.5 rounded-2xl transition-all duration-300 text-[15px] flex items-center group
              overflow-hidden`,
              isActive
                ? 'bg-emerald-500/10 text-emerald-400 font-semibold'
                : 'text-neutral-400 font-medium hover:bg-neutral-800/40 hover:text-neutral-200',
            )}
          >
            {isActive && (
              <span
                className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-3/5 bg-emerald-500 rounded-r-full
                  shadow-[0_0_10px_rgba(16,185,129,0.5)]"
              />
            )}

            <span className="relative z-10">{item.label}</span>

            {!isActive && (
              <span
                className="absolute inset-0 bg-linear-to-r from-white/2 to-transparent opacity-0 group-hover:opacity-100
                  transition-opacity duration-300"
              />
            )}
          </Link>
        )
      })}
    </nav>
  )
}
