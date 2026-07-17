import { useInfiniteQuery } from '@tanstack/react-query'

import { searchUsers } from '@/lib/actions/user.actions'

const PAGE_SIZE = 20

export const useSearchUsersQuery = (searchQuery: string) => {
  return useInfiniteQuery({
    queryKey: ['search-users', searchQuery],
    queryFn: async ({ pageParam = 1 }) => await searchUsers(searchQuery, pageParam as number, PAGE_SIZE),
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