import { useQuery } from '@tanstack/react-query'

import { getUnreadNotificationsCount } from '@/lib/actions/notification.actions'

export const useUnreadNotificationsCountQuery = (enabled: boolean = true) => {
  return useQuery({
    queryKey: ['unread-notifications-count'],
    
    queryFn: async () => {
      const response = await getUnreadNotificationsCount()
      if (!response.success) throw new Error(response.message)
      return response.data ?? 0
    },

    enabled,
    
    refetchInterval: 60 * 1000, 
    staleTime: 60 * 1000,
  })
}