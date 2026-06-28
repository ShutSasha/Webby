import { useMutation } from '@tanstack/react-query'

import { createRightChoiceVoteAction } from '@/lib/actions/vote.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

type CreateVoteVariables = {
  roomId: string
  voteText: string
  duration: number
  choices: string[]
}

export const useCreateRightChoiceVoteMutation = () => {
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ roomId, voteText, duration, choices }: CreateVoteVariables) => {
      return unwrapServerAction(await createRightChoiceVoteAction(roomId, { voteText, duration, choices }))
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message

      addToast(errorMessage, 'error')
    },
  })
}
