import { useInfiniteQuery } from '@tanstack/react-query'

import { getUnreadNotifications } from '@/lib/actions/notification.actions'

const PAGE_SIZE = 10

export const useGetUnreadNotificationsQuery = () => {
  return useInfiniteQuery({
    queryKey: ['unread-notifications'],

    queryFn: async ({ pageParam = 1 }) => {
      return await getUnreadNotifications(pageParam, PAGE_SIZE)
    },

    initialPageParam: 1,

    getNextPageParam: (lastPage, allPages) => {
      const items = lastPage?.data?.items || []
      if (items.length === PAGE_SIZE) {
        return allPages.length + 1
      }
      return undefined
    },

    staleTime: 1 * 60 * 1000,
  })
}
