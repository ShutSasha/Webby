'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

type Props = {
  id: string
}

export default function SettingsNavigation({ id }: Props) {
  const pathname = usePathname()

  const baseClasses = 'border px-4 py-2 transition-colors rounded-xl font-medium'
  const activeClasses = 'bg-emerald-500 border-transparent text-neutral-900'
  const inactiveClasses = 'hover:bg-neutral-800 border-border'

  const getClasses = (path: string) => {
    const isActive = pathname.endsWith(path)
    return `${baseClasses} ${isActive ? activeClasses : inactiveClasses}`
  }

  return (
    <div className="flex flex-col gap-2 w-60">
      <Link href={`/profile/${id}/profile-settings`} className={getClasses('/profile-settings')}>
        Profile settings
      </Link>

      <Link href={`/profile/${id}/secure`} className={getClasses('/secure')}>
        Secure
      </Link>

      <Link href={`/profile/${id}/badges`} className={getClasses('/badges')}>
        Badges
      </Link>
    </div>
  )
}
