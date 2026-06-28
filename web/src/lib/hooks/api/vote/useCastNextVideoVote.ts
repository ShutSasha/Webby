import { useMutation } from '@tanstack/react-query'

import { castNextVideoVoteAction } from '@/lib/actions/vote.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

type CastNextVideoVoteVariables = {
  roomId: string
  queueItemId: string
}

export const useCastNextVideoVoteMutation = () => {
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ roomId, queueItemId }: CastNextVideoVoteVariables) => {
      return unwrapServerAction(await castNextVideoVoteAction(roomId, queueItemId))
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message

      addToast(errorMessage, 'error')
    },
  })
}
