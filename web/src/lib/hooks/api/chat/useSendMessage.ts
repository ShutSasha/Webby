import { useMutation } from '@tanstack/react-query'

import { sendMessageAction } from '@/lib/actions/chat.actions'

export const useSendMessageMutation = () => {
  return useMutation({
    mutationFn: async ({ chatId, content }: { chatId: string; content: string }) => {
      const response = await sendMessageAction(chatId, content)

      if (!response.success) {
        throw new Error(response.message || 'Failed to send message')
      }

      return response
    },
  })
}
