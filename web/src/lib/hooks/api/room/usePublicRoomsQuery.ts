import { useInfiniteQuery } from '@tanstack/react-query'

import { getPublicRooms } from '@/lib/actions/room.actions'

const PAGE_SIZE = 20

export const usePublicRoomsQuery = (searchQuery: string = '', category: string = '') => {
  return useInfiniteQuery({
    queryKey: ['public-rooms', searchQuery, category],

    queryFn: async ({ pageParam = 1 }) => {
      return await getPublicRooms(searchQuery, pageParam, PAGE_SIZE, category)
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
