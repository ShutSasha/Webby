'use client'

import { usePathname } from 'next/navigation'

import AnimatedTabs, { TabItem } from '../../shared/AnimatedTabs'

export default function PageToggle() {
  const pathname = usePathname()

  const tabs: TabItem[] = [
    {
      label: 'Public playlists',
      href: '/playlists',
      isActive: pathname === '/playlists',
    },
    {
      label: 'My playlists',
      href: '/playlists/my',
      isActive: pathname === '/playlists/my',
    },
  ]

  return <AnimatedTabs tabs={tabs} layoutId="playlists-tabs" />
}
