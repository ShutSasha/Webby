import { useMutation } from '@tanstack/react-query'

import { deleteMessageAction } from '@/lib/actions/chat.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useDeleteMessageMutation = () => {
  const addToast = useToastStore(state => state.addToast)
  return useMutation({
    mutationFn: async ({ chatId, messageId }: { chatId: string; messageId: string }) => {
      return unwrapServerAction(await deleteMessageAction(chatId, messageId))
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message

      addToast(errorMessage, 'error')
    },
  })
}
