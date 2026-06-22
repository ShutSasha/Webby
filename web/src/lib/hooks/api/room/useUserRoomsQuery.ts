import { useInfiniteQuery } from '@tanstack/react-query'

import { getUserRooms } from '@/lib/actions/room.actions'

const PAGE_SIZE = 20

export const useUserRoomsQuery = (searchQuery: string = '') => {
  return useInfiniteQuery({
    queryKey: ['user-rooms', searchQuery],

    queryFn: async ({ pageParam = 1 }) => {
      return await getUserRooms(searchQuery, pageParam, PAGE_SIZE)
    },

    initialPageParam: 1,

    getNextPageParam: (lastPage, allPages) => {
      const items = lastPage?.data?.items || []

      if (items.length === PAGE_SIZE) {
        return allPages.length + 1
      }

      return undefined
    },

    staleTime: 5 * 60 * 1000,
  })
}
