import { useInfiniteQuery } from '@tanstack/react-query'

import { getChatMessagesAction } from '@/lib/actions/chat.actions'

export const useGetChatMessagesQuery = (chatId: string | undefined) => {
  return useInfiniteQuery({
    queryKey: ['chat-messages', chatId],
    queryFn: ({ pageParam = 1 }) => getChatMessagesAction(chatId!, pageParam, 20),
    getNextPageParam: lastPage => {
      if (!lastPage.data) return undefined
      const { page, pageSize, totalCount } = lastPage.data

      if (page * pageSize < totalCount) return page + 1
      return undefined
    },
    enabled: !!chatId,
    initialPageParam: 1,
  })
}
