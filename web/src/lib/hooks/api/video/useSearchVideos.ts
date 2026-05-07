import { useInfiniteQuery } from '@tanstack/react-query'

import { searchVideos } from '@/lib/actions/video.actions'

const PAGE_SIZE = 20

type PageParam = {
  page: number
  token?: string | null
  seed?: number | null
}

export const useSearchVideosQuery = (searchQuery: string, searchPlatform?: 'Webby' | 'YouTube') => {
  return useInfiniteQuery({
    queryKey: ['search-videos', searchQuery, searchPlatform],

    queryFn: async ({ pageParam }) => {
      const { page, token, seed } = pageParam as PageParam

      return await searchVideos(searchQuery, page, PAGE_SIZE, searchPlatform, token, seed)
    },

    initialPageParam: { page: 1 } as PageParam,

    getNextPageParam: (lastPage, allPages) => {
      const token = lastPage?.data?.nextPageToken
      if (token) {
        return { page: allPages.length + 1, token }
      }

      const items = lastPage?.data?.items || []
      if (items.length === PAGE_SIZE) {
        return {
          page: allPages.length + 1,
          seed: lastPage?.data?.contentSeed,
        }
      }

      return undefined
    },

    staleTime: 5 * 60 * 1000,
  })
}
