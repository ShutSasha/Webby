import { useInfiniteQuery } from '@tanstack/react-query'

import { searchVideos } from '@/lib/actions/video.actions'

const PAGE_SIZE = 20

export const useSearchVideosQuery = (searchQuery: string) => {
  return useInfiniteQuery({
    queryKey: ['search-videos', searchQuery],

    queryFn: async ({ pageParam = 1 }) => {
      return await searchVideos(searchQuery, pageParam, PAGE_SIZE)
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
