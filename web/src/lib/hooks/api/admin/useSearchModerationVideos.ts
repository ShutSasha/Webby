import { useInfiniteQuery } from '@tanstack/react-query'

import { searchModerationVideosAction } from '@/lib/actions/admin.actions'
import { unwrapServerAction } from '@/lib/utils/general.utils'

export const useSearchModerationVideosQuery = (searchText: string = '', limit: number = 20) => {
  return useInfiniteQuery({
    queryKey: ['moderation-videos-search', searchText, limit],
    queryFn: async ({ pageParam = 1 }) => {
      return unwrapServerAction(await searchModerationVideosAction(searchText, pageParam, limit))
    },
    getNextPageParam: (lastPage, allPages) => {
      const items = lastPage?.items || []

      if (items.length === limit) {
        return allPages.length + 1
      }

      return undefined
    },
    initialPageParam: 1,
    staleTime: 5 * 60 * 1000,
  })
}
