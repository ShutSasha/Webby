import { useInfiniteQuery } from '@tanstack/react-query'

import { searchStreams } from '@/lib/actions/stream.actions'

const PAGE_SIZE = 20

export const useSearchStreamsQuery = (searchQuery: string) => {
  return useInfiniteQuery({
    queryKey: ['search-streams', searchQuery],
    queryFn: async ({ pageParam = 1 }) => await searchStreams(searchQuery, pageParam, PAGE_SIZE),
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      const items = lastPage?.data?.items || []
      if (items.length === PAGE_SIZE) return allPages.length + 1
      return undefined
    },

    staleTime: 5 * 60 * 1000,
  })
}
