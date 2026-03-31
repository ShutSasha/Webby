import { RoomTabValue } from '@/lib/constants/room.constants'
import MainLayout from '@/ui/components/layouts/MainLayout'
import AnimatedTabsContainer from '@/ui/components/modules/Rooms/AnimatedTabsContainer'
import RoomsContainer from '@/ui/components/modules/Rooms/Public/RoomsContainer'
import RoomPageHeader from '@/ui/components/modules/Rooms/RoomPageHeader'

type Props = {
  searchParams: Promise<{ tab: RoomTabValue | string; query?: string }>
}

export default async function RoomsPage({ searchParams }: Props) {
  const { tab, query } = await searchParams
  const safeQuery = query || ''

  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <RoomPageHeader />
        <AnimatedTabsContainer currentTab={tab} />
        <RoomsContainer query={safeQuery} category={tab || ''} />
      </div>
    </MainLayout>
  )
}
