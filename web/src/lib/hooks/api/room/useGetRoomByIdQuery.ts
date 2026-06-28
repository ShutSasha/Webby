import { useQuery } from '@tanstack/react-query'

import { getRoomById } from '@/lib/actions/room.actions'
import { unwrapServerAction } from '@/lib/utils/general.utils'

export const useGetRoomByIdQuery = (roomId: string) => {
  return useQuery({
    queryKey: ['room', roomId],
    queryFn: async () => {
      return unwrapServerAction(await getRoomById(roomId))
    },
    enabled: !!roomId,
  })
}
