import { useInfiniteQuery } from '@tanstack/react-query'

import { getComplaintsAction } from '@/lib/actions/admin.actions'
import { unwrapServerAction } from '@/lib/utils/general.utils'

export const useGetComplaintsQuery = (limit: number = 10) => {
  return useInfiniteQuery({
    queryKey: ['complaints', limit],
    queryFn: async ({ pageParam = 1 }) => {
      return unwrapServerAction(await getComplaintsAction(pageParam, limit))
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
