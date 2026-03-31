import { Suspense } from 'react'

import MainLayout from '@/ui/components/layouts/MainLayout'
import CategoryToggle from '@/ui/components/modules/Rooms/GenreToggle'
import RoomCard from '@/ui/components/modules/Rooms/RoomCard'
import RoomPageHeader from '@/ui/components/modules/Rooms/RoomPageHeader'
import GridCardsContainer from '@/ui/components/shared/GridCardsContainer'

export default function RoomsPage() {
  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <RoomPageHeader />
        <div className="flex justify-between">
          <Suspense>
            <CategoryToggle />
          </Suspense>
        </div>

        {/* Room list goes here */}
        <GridCardsContainer>
          {Array.from({ length: 10 }).map((_, index) => (
            <RoomCard key={index} id={`${index}`} />
          ))}
        </GridCardsContainer>
      </div>
    </MainLayout>
  )
}
