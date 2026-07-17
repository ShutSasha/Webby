import { useQuery } from '@tanstack/react-query'

import { getUserFollowers } from '@/lib/actions/user.actions'

export const useUserFollowersQuery = (userId: string) => {
  return useQuery({
    queryKey: ['user-followers', userId],

    queryFn: async () => {
      const response = await getUserFollowers(userId)

      if (!response.success) {
        throw new Error(response.message || 'Failed to fetch user followers')
      }

      return response.data || []
    },

    enabled: !!userId,
    staleTime: 5 * 60 * 1000,
  })
}
