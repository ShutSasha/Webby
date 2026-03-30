import { useQuery } from '@tanstack/react-query'

import { getUserFollows } from '@/lib/actions/user.actions'

export const useUserFollowsQuery = (userId: string) => {
  return useQuery({
    queryKey: ['user-follows', userId],

    queryFn: async () => {
      const response = await getUserFollows(userId)

      if (!response.success) {
        throw new Error(response.message || 'Failed to fetch user follows')
      }

      return response.data || []
    },

    enabled: !!userId,

    staleTime: 5 * 60 * 1000,
  })
}
