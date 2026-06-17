import { useMutation } from '@tanstack/react-query'

import { createNextVideoVoteAction } from '@/lib/actions/vote.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useCreateNextVideoVoteMutation = () => {
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async (roomId: string) => {
      return unwrapServerAction(await createNextVideoVoteAction(roomId))
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message

      addToast(errorMessage, 'error')
    },
  })
}
