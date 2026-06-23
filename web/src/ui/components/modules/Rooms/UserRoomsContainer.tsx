'use client'

import { useUserRoomsQuery } from '@/lib/hooks/api/room/useUserRoomsQuery'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import GridCardsContainer from '@/ui/components/shared/GridCardsContainer'

import { RoomCardSkeleton } from './RoomCard'
import UserRoomItem from './UserRoomItem'

type Props = {
  query: string
}

export default function UserRoomsContainer({ query }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useUserRoomsQuery(query)

  const rooms = data?.pages.flatMap(page => page?.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  if (isLoading && rooms.length === 0) {
    return (
      <GridCardsContainer>
        {[...new Array(12)].map((_, idx) => (
          <RoomCardSkeleton key={idx} />
        ))}
      </GridCardsContainer>
    )
  }

  if (!isLoading && rooms.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center py-20">
        {query.length > 0 ? (
          <p className="text-foreground0 text-center">
            Rooms by query <span className="text-foreground-subtle">{`'${query}'`}</span> not found
          </p>
        ) : (
          <p className="text-foreground0 text-center">You haven&apos;t created any rooms yet</p>
        )}
      </div>
    )
  }

  return (
    <>
      <GridCardsContainer>
        {rooms.map((room, index) => {
          const isLast = rooms.length === index + 1

          const item = <UserRoomItem key={room.id} room={room} />

          if (isLast) {
            return (
              <div ref={lastElementRef} key={`last-${room.id}`}>
                {item}
              </div>
            )
          }

          return item
        })}
      </GridCardsContainer>

      {isFetchingNextPage && (
        <div className="w-full flex justify-center py-8">
          <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </>
  )
}
