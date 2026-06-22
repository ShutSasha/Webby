import { useMutation, useQueryClient } from '@tanstack/react-query'

import { updateRoomAction } from '@/lib/actions/room.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useUpdateRoom = (options?: { onSuccess?: () => void }) => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ roomId, formData }: { roomId: string; formData: FormData }) => {
      return unwrapServerAction(await updateRoomAction(roomId, formData))
    },
    onSuccess: (_, variables) => {
      addToast('Room updated successfully', 'success')
      queryClient.invalidateQueries({ queryKey: ['user-rooms'] })
      queryClient.invalidateQueries({ queryKey: ['public-rooms'] })

      queryClient.invalidateQueries({ queryKey: ['room', variables.roomId] })

      if (options?.onSuccess) options.onSuccess()
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message || 'Failed to update room'
      addToast(errorMessage, 'error')
    },
  })
}
