import { ALLOWED_ROOM_TABS, ROOM_TABS_CONFIG, RoomTabValue } from '@/lib/constants/room.constants'
import MainLayout from '@/ui/components/layouts/MainLayout'
import RoomCard from '@/ui/components/modules/Rooms/RoomCard'
import RoomPageHeader from '@/ui/components/modules/Rooms/RoomPageHeader'
import AnimatedTabs from '@/ui/components/shared/AnimatedTabs'
import GridCardsContainer from '@/ui/components/shared/GridCardsContainer'

type Props = {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

export default async function RoomsPage({ searchParams }: Props) {
  const { tab } = await searchParams
  const currentTab = ALLOWED_ROOM_TABS.includes(tab as RoomTabValue) ? (tab as RoomTabValue) : 'all'

  const roomsTabs = ROOM_TABS_CONFIG.map(config => ({
    label: config.label,
    href: `?tab=${config.value}`,
    isActive: currentTab === config.value,
  }))

  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <RoomPageHeader />
        <div className="flex justify-between">
          <AnimatedTabs tabs={roomsTabs} layoutId="rooms-page-tabs" />
        </div>

        <GridCardsContainer>
          {Array.from({ length: 10 }).map((_, index) => (
            <RoomCard key={index} id={`${index}`} />
          ))}
        </GridCardsContainer>
      </div>
    </MainLayout>
  )
}
