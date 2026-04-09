import { useInfiniteQuery } from '@tanstack/react-query'

import { getUserNotifications } from '@/lib/actions/notification.actions'

const PAGE_SIZE = 10

export const useGetUserNotificationsQuery = () => {
  return useInfiniteQuery({
    queryKey: ['user-notifications'],

    queryFn: async ({ pageParam = 1 }) => {
      return await getUserNotifications(pageParam, PAGE_SIZE)
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
