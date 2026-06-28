import { useMutation, useQueryClient } from '@tanstack/react-query'

import { ChatMessage, ChatUser, sendMessageAction } from '@/lib/actions/chat.actions'

type SendMessagePayload = {
  chatId: string
  content: string
  sender: ChatUser
}

export const useSendMessageMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ chatId, content }: SendMessagePayload) => {
      const response = await sendMessageAction(chatId, content)

      if (!response.success) {
        throw new Error(response.message || 'Failed to send message')
      }

      return response
    },
    onMutate: async newMsg => {
      await queryClient.cancelQueries({ queryKey: ['chat-messages', newMsg.chatId] })

      const previousMessages = queryClient.getQueryData(['chat-messages', newMsg.chatId])

      const optimisticMessage: ChatMessage = {
        id: `temp-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`,
        content: newMsg.content,
        createdAt: new Date().toISOString(),
        isEdited: false,
        sender: newMsg.sender,
      }

      queryClient.setQueryData(['chat-messages', newMsg.chatId], (oldData: any) => {
        if (!oldData || !oldData.pages || oldData.pages.length === 0) return oldData

        const newPages = [...oldData.pages]
        const firstPage = { ...newPages[0] }

        if (firstPage.data && firstPage.data.items) {
          firstPage.data = {
            ...firstPage.data,
            items: [optimisticMessage, ...firstPage.data.items],
          }
          newPages[0] = firstPage
        }

        return { ...oldData, pages: newPages }
      })

      return { previousMessages }
    },
    onError: (err, newMsg, context) => {
      if (context?.previousMessages) {
        queryClient.setQueryData(['chat-messages', newMsg.chatId], context.previousMessages)
      }
    },
  })
}
