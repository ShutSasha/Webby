import { useMutation, useQueryClient } from '@tanstack/react-query'

import { toggleFollow } from '@/lib/actions/user.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useToggleFollowMutation = (currentUserId?: string) => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async (targetId: string) => {
      return unwrapServerAction(await toggleFollow(targetId))
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message
      addToast(errorMessage, 'error')
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
