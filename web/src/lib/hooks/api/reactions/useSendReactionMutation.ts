import { useMutation } from '@tanstack/react-query'

import { sendReactionAction } from '@/lib/actions/reaction.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

type SendReactionPayload = {
  roomId: string
  reactionId: string
}

export const useSendReactionMutation = () => {
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ roomId, reactionId }: SendReactionPayload) => {
      return unwrapServerAction(await sendReactionAction(roomId, reactionId))
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message || 'Failed to send reaction'
      addToast(errorMessage, 'error')
    },
  })
}
