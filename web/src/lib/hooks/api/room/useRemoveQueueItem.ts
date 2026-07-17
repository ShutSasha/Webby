import { useMutation } from '@tanstack/react-query'

import { removeQueueItemAction } from '@/lib/actions/room.actions'

export const useRemoveQueueItemMutation = () => {
  return useMutation({
    mutationFn: async ({ roomId, itemId }: { roomId: string; itemId: string }) => {
      return await removeQueueItemAction(roomId, itemId)
    },
  })
}
