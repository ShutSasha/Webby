import { useMutation } from '@tanstack/react-query'

import { editMessageAction } from '@/lib/actions/chat.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

type EditMessageVariables = {
  chatId: string
  messageId: string
  content: string
}

export const useEditMessageMutation = () => {
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ chatId, messageId, content }: EditMessageVariables) => {
      return unwrapServerAction(await editMessageAction(chatId, messageId, content))
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message

      addToast(errorMessage, 'error')
    },
  })
}
