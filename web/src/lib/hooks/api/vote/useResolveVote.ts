import { useMutation } from '@tanstack/react-query'

import { resolveVoteAction } from '@/lib/actions/vote.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

type ResolveVoteVariables = {
  roomId: string
  voteId: string
  rightChoice: string
}

export const useResolveVoteMutation = () => {
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ roomId, voteId, rightChoice }: ResolveVoteVariables) => {
      return unwrapServerAction(await resolveVoteAction(roomId, voteId, rightChoice))
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message
      
      addToast(errorMessage, 'error')
    },
  })
}