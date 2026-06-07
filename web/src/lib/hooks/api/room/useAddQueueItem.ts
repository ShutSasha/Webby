import { useMutation } from '@tanstack/react-query'

import { addQueueItemAction } from '@/lib/actions/room.actions'

export const useAddQueueItemMutation = () => {
  return useMutation({
    mutationFn: async ({ roomId, videoId }: { roomId: string; videoId: string }) => {
      return await addQueueItemAction(roomId, videoId)
    },
  })
}
