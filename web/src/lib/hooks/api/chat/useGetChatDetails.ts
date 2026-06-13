import { useQuery } from '@tanstack/react-query'

import { getChatByIdAction } from '@/lib/actions/chat.actions'

export const useGetChatDetailsQuery = (chatId: string | undefined) => {
  return useQuery({
    queryKey: ['chat-details', chatId],
    queryFn: () => getChatByIdAction(chatId!),
    enabled: !!chatId,
  })
}
