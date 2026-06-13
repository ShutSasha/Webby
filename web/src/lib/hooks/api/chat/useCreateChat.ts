import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'

import { createPrivateChatAction } from '@/lib/actions/chat.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useCreateChatMutation = () => {
  const router = useRouter()
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async (targetId: string) => {
      return unwrapServerAction(await createPrivateChatAction(targetId))
    },
    onSuccess: data => {
      queryClient.invalidateQueries({ queryKey: ['chat-history'] })
      router.push(`/chats/${data.id}`)
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message

      addToast(errorMessage, 'error')
    },
  })
}
