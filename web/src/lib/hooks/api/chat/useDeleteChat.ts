import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'

import { deleteChatAction } from '@/lib/actions/chat.actions'

export const useDeleteChatMutation = () => {
  const queryClient = useQueryClient()
  const router = useRouter()

  return useMutation({
    mutationFn: async (chatId: string) => {
      const response = await deleteChatAction(chatId)

      if (!response.success) {
        throw new Error(response.message || 'Failed to delete chat')
      }

      return response
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['chat-history'] })

      router.push('/chats')
    },
  })
}
