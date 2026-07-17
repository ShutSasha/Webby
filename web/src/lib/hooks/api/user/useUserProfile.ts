import { useQuery } from '@tanstack/react-query'

import { getUser, GetUserResponse } from '@/lib/actions/user.actions'

export const useUserProfileQuery = (userId: string, initialData?: GetUserResponse) => {
  return useQuery({
    queryKey: ['user-profile', userId],
    queryFn: async () => {
      const data = await getUser(userId)
      if (!data) throw new Error('User not found')
      return data
    },

    initialData,
    staleTime: 5 * 60 * 1000,
  })
}
