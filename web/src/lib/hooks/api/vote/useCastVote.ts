import { useMutation } from '@tanstack/react-query'

import { castVoteAction } from '@/lib/actions/vote.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

type CastVoteVariables = {
  roomId: string
  voteId: string
  choice: string
}

export const useCastVoteMutation = () => {
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ roomId, voteId, choice }: CastVoteVariables) => {
      return unwrapServerAction(await castVoteAction(roomId, voteId, choice))
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message

      addToast(errorMessage, 'error')
    },
  })
}
