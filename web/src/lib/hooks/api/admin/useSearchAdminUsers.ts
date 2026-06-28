import { useInfiniteQuery } from '@tanstack/react-query'

import { searchAdminUsersAction } from '@/lib/actions/admin.actions'
import { unwrapServerAction } from '@/lib/utils/general.utils'

export const useSearchAdminUsersQuery = (searchText: string = '', limit: number = 10) => {
  return useInfiniteQuery({
    queryKey: ['admin-users-search', searchText, limit],
    queryFn: async ({ pageParam = 1 }) => {
      return unwrapServerAction(await searchAdminUsersAction(searchText, pageParam, limit))
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
