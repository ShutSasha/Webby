import { RoomTabValue } from '@/lib/constants/room.constants'
import AnimatedTabsContainer from '@/ui/components/modules/Rooms/AnimatedTabsContainer'
import RoomsContainer from '@/ui/components/modules/Rooms/Public/RoomsContainer'

type Props = {
  searchParams: Promise<{ tab: RoomTabValue | string; query?: string }>
}

export default async function RoomsPage({ searchParams }: Props) {
  const { tab, query } = await searchParams
  const safeQuery = query || ''

  return (
    <>
      <AnimatedTabsContainer currentTab={tab || 'All'} />
      <RoomsContainer query={safeQuery} category={tab || ''} />
    </>
  )
}
