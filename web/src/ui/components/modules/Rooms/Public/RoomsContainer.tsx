'use client'

import { usePublicRoomsQuery } from '@/lib/hooks/api/room/usePublicRoomsQuery'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import GridCardsContainer from '@/ui/components/shared/GridCardsContainer'

import RoomCard, { RoomCardSkeleton } from '../RoomCard'

type Props = {
  query: string
  category: string
  limit?: number
  gridClassName?: string
}

export default function RoomsContainer({ query, category, limit, gridClassName }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = usePublicRoomsQuery(query, category)

  const allRooms = data?.pages.flatMap(page => page?.data?.items || []) || []

  const rooms = limit ? allRooms.slice(0, limit) : allRooms

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  if (isLoading && rooms.length === 0) {
    return (
      <GridCardsContainer className={gridClassName}>
        {[...new Array(limit || 12)].map((_, idx) => (
          <RoomCardSkeleton key={idx} />
        ))}
      </GridCardsContainer>
    )
  }

  if (!isLoading && rooms.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center py-20">
        {query.length > 0 && (
          <p className="text-foreground0 text-center">
            Rooms by query <span className="text-foreground-subtle">{`'${query}'`}</span> not found
          </p>
        )}
        {query.length === 0 && <p className="text-foreground0 text-center">No rooms found</p>}
      </div>
    )
  }

  return (
    <>
      <GridCardsContainer className={gridClassName}>
        {rooms.map((room, index) => {
          const isLast = rooms.length === index + 1

          const item = (
            <RoomCard
              key={room.id}
              id={room.id}
              name={room.name}
              thumbnail={room.thumbnail}
              category={room.categoryName || 'not found'}
              hostUsername={room.hostUsername}
              hostAvatarUrl={room.hostAvatarUrl}
            />
          )

          if (isLast && !limit) {
            return (
              <div ref={lastElementRef} key={`last-${room.id}`}>
                {item}
              </div>
            )
          }

          return item
        })}
      </GridCardsContainer>

      {isFetchingNextPage && !limit && (
        <div className="w-full flex justify-center py-8">
          <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </>
  )
}
