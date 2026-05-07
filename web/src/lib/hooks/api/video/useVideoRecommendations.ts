import { useInfiniteQuery } from '@tanstack/react-query'

import { getVideoRecommendations } from '@/lib/actions/video.actions'

const PAGE_SIZE = 20

type PageParam = {
  page: number
  seed?: number | null
}

export const useVideoRecommendationsQuery = (videoId: string) => {
  return useInfiniteQuery({
    queryKey: ['video-recommendations', videoId],

    queryFn: async ({ pageParam }) => {
      const { page, seed } = pageParam as PageParam

      return await getVideoRecommendations(videoId, page, PAGE_SIZE, seed)
    },

    initialPageParam: { page: 1 } as PageParam,

    getNextPageParam: (lastPage, allPages) => {
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
    enabled: !!videoId,
  })
}
