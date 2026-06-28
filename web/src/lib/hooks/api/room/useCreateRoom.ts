import { useMutation, useQueryClient } from '@tanstack/react-query'

import { createRoomAction } from '@/lib/actions/room.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useCreateRoom = (options?: { onSuccess?: () => void }) => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async (formData: FormData) => unwrapServerAction(await createRoomAction(formData)),
    onSuccess: () => {
      addToast('Room created successfully!', 'success')

      queryClient.invalidateQueries({ queryKey: ['user-rooms'] })
      queryClient.invalidateQueries({ queryKey: ['public-rooms'] })

      if (options?.onSuccess) options.onSuccess()
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message

      addToast(errorMessage, 'error')
    },
  })
}
