import { useMutation, useQueryClient } from '@tanstack/react-query'

import { deleteRoomAction } from '@/lib/actions/room.actions'
import { extractServerMessage } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useDeleteRoom = (options?: { onSuccess?: () => void }) => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async (roomId: string) => await deleteRoomAction(roomId),
    onSuccess: response => {
      if (response.success) {
        addToast('Room deleted successfully', 'success')
        queryClient.invalidateQueries({ queryKey: ['user-rooms'] })
        queryClient.invalidateQueries({ queryKey: ['public-rooms'] })
        if (options?.onSuccess) options.onSuccess()
      } else {
        addToast(extractServerMessage(response.errors) || 'Failed to delete room', 'error')
      }
    },
    onError: () => addToast('Critical error while deleting room', 'error'),
  })
}
