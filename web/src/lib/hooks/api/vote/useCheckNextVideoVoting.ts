import { useQuery } from '@tanstack/react-query'

import { checkNextVideoVotingAction } from '@/lib/actions/vote.actions'
import { unwrapServerAction } from '@/lib/utils/general.utils'

export const useCheckNextVideoVotingQuery = (roomId: string) => {
  return useQuery({
    queryKey: ['has-next-video-voting', roomId],
    queryFn: async () => {
      return unwrapServerAction(await checkNextVideoVotingAction(roomId))
    },
    enabled: !!roomId,
  })
}
