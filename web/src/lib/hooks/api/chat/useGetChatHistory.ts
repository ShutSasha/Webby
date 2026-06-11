import { useInfiniteQuery } from '@tanstack/react-query'

import { getChatHistoryAction } from '@/lib/actions/chat.actions'

export const useGetChatHistoryQuery = (search?: string) => {
  return useInfiniteQuery({
    queryKey: ['chat-history', search],
    queryFn: ({ pageParam = 1 }) => getChatHistoryAction(pageParam, 10, search),
    getNextPageParam: lastPage => {
      if (!lastPage.data) return undefined
      const { page, pageSize, totalCount } = lastPage.data

      if (page * pageSize < totalCount) return page + 1
      return undefined
    },
    initialPageParam: 1,
  })
}
