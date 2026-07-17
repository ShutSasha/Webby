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

      const res = await searchVideos(searchQuery, page, PAGE_SIZE, searchPlatform, token, seed)

      return res
    },

    initialPageParam: { page: 1 } as PageParam,

    getNextPageParam: lastPage => {
      if (!lastPage?.data) return undefined

      const token = lastPage.data.nextPageToken
      if (token) {
        return { page: lastPage.data.page + 1, token }
      }

      const { page, pageSize, totalCount, contentSeed } = lastPage.data
      const hasMore = page * pageSize < totalCount

      if (hasMore) {
        return {
          page: page + 1,
          seed: contentSeed,
        }
      }

      return undefined
    },

    staleTime: 5 * 60 * 1000,
  })
}
