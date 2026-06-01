import { useInfiniteQuery } from '@tanstack/react-query'

import { getRoomQueueAction } from '@/lib/actions/room.actions'

const PAGE_SIZE = 20

export const useRoomQueueQuery = (roomId: string) => {
  return useInfiniteQuery({
    queryKey: ['room-queue', roomId],

    queryFn: async ({ pageParam = 1 }) => {
      return await getRoomQueueAction(roomId, pageParam, PAGE_SIZE)
    },

    initialPageParam: 1,

    getNextPageParam: lastPage => {
      if (!lastPage?.data) return undefined

      const { page, pageSize, totalCount } = lastPage.data
      const hasMore = page * pageSize < totalCount

      return hasMore ? page + 1 : undefined
    },

    staleTime: 5 * 60 * 1000,
    enabled: !!roomId,
  })
}
