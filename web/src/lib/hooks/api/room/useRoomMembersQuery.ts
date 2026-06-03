import { useInfiniteQuery } from '@tanstack/react-query'

import { getRoomMembersAction } from '@/lib/actions/room.actions'

const PAGE_SIZE = 10

export const useRoomMembersQuery = (roomId: string, searchQuery: string = '') => {
  return useInfiniteQuery({
    queryKey: ['room-members', roomId, searchQuery],

    queryFn: async ({ pageParam = 1 }) => {
      return await getRoomMembersAction(roomId, pageParam, PAGE_SIZE, searchQuery)
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
