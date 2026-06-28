import { useQuery } from '@tanstack/react-query'

import { getRoomVotesAction } from '@/lib/actions/vote.actions'
import { unwrapServerAction } from '@/lib/utils/general.utils'

export const useGetRoomVotesQuery = (roomId: string) => {
  return useQuery({
    queryKey: ['room-votes', roomId],
    queryFn: async () => {
      return unwrapServerAction(await getRoomVotesAction(roomId))
    },
    enabled: !!roomId,
  })
}
