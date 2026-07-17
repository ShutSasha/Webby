import { useMutation } from '@tanstack/react-query'

import { addQueueItemAction } from '@/lib/actions/room.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useAddQueueItemMutation = () => {
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ roomId, videoId }: { roomId: string; videoId: string }) => {
      return unwrapServerAction(await addQueueItemAction(roomId, videoId))
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message

      addToast(errorMessage, 'error')
    },
  })
}
