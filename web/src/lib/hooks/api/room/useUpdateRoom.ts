import { useMutation, useQueryClient } from '@tanstack/react-query'

import { updateRoomAction } from '@/lib/actions/room.actions'
import { extractServerMessage } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useUpdateRoom = (options?: { onSuccess?: () => void }) => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ roomId, formData }: { roomId: string; formData: FormData }) =>
      await updateRoomAction(roomId, formData),
    onSuccess: response => {
      if (response.success) {
        addToast('Room updated successfully', 'success')
        queryClient.invalidateQueries({ queryKey: ['user-rooms'] })
        queryClient.invalidateQueries({ queryKey: ['public-rooms'] })
        if (options?.onSuccess) options.onSuccess()
      } else {
        addToast(extractServerMessage(response.errors) || 'Failed to update room', 'error')
      }
    },
    onError: () => addToast('Critical error while updating room', 'error'),
  })
}
