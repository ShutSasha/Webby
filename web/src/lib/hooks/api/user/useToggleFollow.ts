import { useMutation, useQueryClient } from '@tanstack/react-query'

import { toggleFollow } from '@/lib/actions/user.actions'

export const useToggleFollowMutation = (currentUserId?: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (targetId: string) => {
      const response = await toggleFollow(targetId)

      if (!response.success) {
        throw new Error(response.message || 'Toggle follow error')
      }

      return response
    },

    onSettled: (_, __, targetId) => {
      queryClient.invalidateQueries({ queryKey: ['user-followers', targetId] })
      queryClient.invalidateQueries({ queryKey: ['user-profile', targetId] })

      if (currentUserId) {
        queryClient.invalidateQueries({ queryKey: ['user-follows', currentUserId] })
        queryClient.invalidateQueries({ queryKey: ['user-profile', currentUserId] })
      }
    },
  })
}
