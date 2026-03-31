import { RoomTabValue } from '@/lib/constants/room.constants'
import MainLayout from '@/ui/components/layouts/MainLayout'
import AnimatedTabsContainer from '@/ui/components/modules/Rooms/AnimatedTabsContainer'
import RoomCard from '@/ui/components/modules/Rooms/RoomCard'
import RoomPageHeader from '@/ui/components/modules/Rooms/RoomPageHeader'
import GridCardsContainer from '@/ui/components/shared/GridCardsContainer'

type Props = {
  searchParams: Promise<{ tab: RoomTabValue | string }>
}

export default async function RoomsPage({ searchParams }: Props) {
  const { tab } = await searchParams

  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <RoomPageHeader />
        <AnimatedTabsContainer currentTab={tab} />
        <GridCardsContainer>
          {Array.from({ length: 10 }).map((_, index) => (
            <RoomCard key={index} id={`${index}`} />
          ))}
        </GridCardsContainer>
      </div>
    </MainLayout>
  )
}
