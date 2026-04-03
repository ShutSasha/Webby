import { useInfiniteQuery } from '@tanstack/react-query'

import { getUserVideos } from '@/lib/actions/video.actions'

const PAGE_SIZE = 15

export const useGetUserVideosQuery = (userId: string) => {
  return useInfiniteQuery({
    queryKey: ['user-videos', userId],

    queryFn: async ({ pageParam = 1 }) => {
      return await getUserVideos(userId, pageParam, PAGE_SIZE)
    },

    initialPageParam: 1,

    getNextPageParam: (lastPage, allPages) => {
      const items = lastPage?.data?.items || []

      if (items.length === PAGE_SIZE) {
        return allPages.length + 1
      }

      return undefined
    },

    enabled: !!userId,
    staleTime: 5 * 60 * 1000,
  })
}
