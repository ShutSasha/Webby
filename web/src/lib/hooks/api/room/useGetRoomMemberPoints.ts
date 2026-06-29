import { useQuery } from '@tanstack/react-query'

import { getRoomMemberPointsAction } from '@/lib/actions/room.actions'
import { unwrapServerAction } from '@/lib/utils/general.utils'

export const useGetRoomMemberPointsQuery = (roomId: string | undefined) => {
  return useQuery({
    queryKey: ['room-member-points', roomId],
    queryFn: async () => unwrapServerAction(await getRoomMemberPointsAction(roomId!)),
    enabled: !!roomId,
    staleTime: 1000 * 60 * 5,
  })
}
