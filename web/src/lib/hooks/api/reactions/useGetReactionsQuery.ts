import { useQuery } from '@tanstack/react-query'

import { getReactionsAction } from '@/lib/actions/reaction.actions'
import { unwrapServerAction } from '@/lib/utils/general.utils'

export const useGetReactionsQuery = () => {
  return useQuery({
    queryKey: ['reactions'],
    queryFn: async () => {
      return unwrapServerAction(await getReactionsAction())
    },
    staleTime: 1000 * 60 * 60,
  })
}
