import { useInfiniteQuery } from '@tanstack/react-query'

import { searchVideos } from '@/lib/actions/video.actions'

const PAGE_SIZE = 20

export const useSearchVideosQuery = (searchQuery: string, searchPlatform?: 'Webby' | 'YouTube') => {
  return useInfiniteQuery({
    queryKey: ['search-videos', searchQuery, searchPlatform],

    queryFn: async ({ pageParam = 1 }) => {
      const page = typeof pageParam === 'number' ? pageParam : 1
      const token = typeof pageParam === 'string' ? pageParam : null

      return await searchVideos(searchQuery, page, PAGE_SIZE, searchPlatform, token)
    },

    initialPageParam: 1 as string | number | undefined,

    getNextPageParam: (lastPage, allPages) => {
      const token = lastPage?.data?.nextPageToken
      if (token) return token

      const items = lastPage?.data?.items || []
      if (items.length === PAGE_SIZE) return allPages.length + 1

      return undefined
    },

    staleTime: 5 * 60 * 1000,
  })
}
