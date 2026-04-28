import { useMutation, useQueryClient } from '@tanstack/react-query'

import { createRoomAction } from '@/lib/actions/room.actions'
import { extractServerMessage } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useCreateRoom = (options?: { onSuccess?: () => void }) => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async (formData: FormData) => await createRoomAction(formData),
    onSuccess: response => {
      if (response.success) {
        addToast('Room created successfully!', 'success')

        queryClient.invalidateQueries({ queryKey: ['user-rooms'] })
        queryClient.invalidateQueries({ queryKey: ['public-rooms'] })

        if (options?.onSuccess) options.onSuccess()
      } else {
        const msg = extractServerMessage(response.errors)
        addToast(msg || 'Failed to create room', 'error')
      }
    },
    onError: () => {
      addToast('Critical error while creating room', 'error')
    },
  })
}
