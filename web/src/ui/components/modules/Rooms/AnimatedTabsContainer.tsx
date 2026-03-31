'use client'

import { useSearchParams } from 'next/navigation'

import { ROOM_TABS_CONFIG, RoomTabValue } from '@/lib/constants/room.constants'
import { useCategoriesQuery } from '@/lib/hooks/api/category/useCategoriesQuery'
import AnimatedTabs, { AnimatedTabsSkeleton } from '@/ui/components/shared/AnimatedTabs'

type Props = {
  currentTab: RoomTabValue | string
}

export default function AnimatedTabsContainer({ currentTab }: Props) {
  const { data, isPending } = useCategoriesQuery()
  const searchParams = useSearchParams()

  const createQueryString = (tabValue: string) => {
    const params = new URLSearchParams(searchParams.toString())
    params.set('tab', tabValue)
    return `?${params.toString()}`
  }

  if (isPending) {
    return (
      <div className="flex justify-between overflow-hidden">
        <AnimatedTabsSkeleton count={6} />
      </div>
    )
  }

  const categories = data?.pages.flatMap(page => page?.data?.items || []) || []

  const roomsTabs =
    categories.length > 0
      ? [
          {
            label: 'All',
            href: createQueryString('all'),
            isActive: currentTab === 'all',
          },
          ...categories
            .filter(category => category.name.toLowerCase() !== 'all')
            .map(category => ({
              label: category.name,
              href: createQueryString(category.name.toLowerCase()),
              isActive: currentTab === category.name.toLowerCase(),
            })),
        ]
      : ROOM_TABS_CONFIG.map(config => ({
          label: config.label,
          href: createQueryString(config.value),
          isActive: currentTab === config.value,
        }))

  return (
    <div className="flex justify-between w-full overflow-x-auto no-scrollbar">
      <AnimatedTabs tabs={roomsTabs} layoutId="rooms-page-tabs" />
    </div>
  )
}
