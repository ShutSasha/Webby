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
      if (!lastPage) return undefined

      const totalPages = Math.ceil(lastPage.totalCount / limit)

      if (allPages.length < totalPages) {
        return allPages.length + 1
      }

      return undefined
    },
    initialPageParam: 1,
    staleTime: 5 * 60 * 1000,
  })
}
