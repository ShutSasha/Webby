import { useInfiniteQuery } from '@tanstack/react-query'

import { searchStreams } from '@/lib/actions/stream.actions'

const PAGE_SIZE = 20

export const useSearchStreamsQuery = (searchQuery: string) => {
  return useInfiniteQuery({
    queryKey: ['search-streams', searchQuery],

    queryFn: async ({ pageParam }) => await searchStreams(searchQuery, pageParam as string | undefined, PAGE_SIZE),

    initialPageParam: undefined as string | undefined,

    getNextPageParam: lastPage => {
      const token = lastPage?.data?.nextPageToken

      return token ? token : undefined
    },

    staleTime: 5 * 60 * 1000,
  })
}
