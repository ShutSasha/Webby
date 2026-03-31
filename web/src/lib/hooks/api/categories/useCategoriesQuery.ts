import { useInfiniteQuery } from '@tanstack/react-query'

import { getCategories } from '@/lib/actions/category.actions'

const PAGE_SIZE = 20

export const useCategoriesQuery = (searchQuery: string = '') => {
  return useInfiniteQuery({
    queryKey: ['categories', searchQuery],

    queryFn: async ({ pageParam = 1 }) => {
      return await getCategories(searchQuery, pageParam, PAGE_SIZE)
    },

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
