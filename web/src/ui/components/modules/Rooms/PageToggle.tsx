'use client'

import { usePathname } from 'next/navigation'

import AnimatedTabs, { TabItem } from '../../shared/AnimatedTabs'

export default function PageToggle() {
  const pathname = usePathname()

  const tabs: TabItem[] = [
    {
      label: 'Public rooms',
      href: '/rooms',
      isActive: pathname === '/rooms',
    },
    {
      label: 'My rooms',
      href: '/rooms/my',
      isActive: pathname === '/rooms/my',
    },
  ]

  return <AnimatedTabs tabs={tabs} layoutId="rooms-tabs" />
}
