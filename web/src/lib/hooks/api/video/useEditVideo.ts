import { useMutation, useQueryClient } from '@tanstack/react-query'

import { updateVideoMetadataAction } from '@/lib/actions/video.actions'
import { extractServerMessage } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export const useUpdateVideoMetadata = () => {
  const addToast = useToastStore(state => state.addToast)
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (formData: FormData) => await updateVideoMetadataAction(formData),
    onSuccess: response => {
      if (response.success) {
        addToast('Video successfully updated!', 'success')

        queryClient.invalidateQueries({ queryKey: ['user-videos'] })
        queryClient.invalidateQueries({ queryKey: ['search-videos'] })
      } else {
        const msg = extractServerMessage(response.errors)
        addToast(msg || 'Error updating video', 'error')
      }
    },
    onError: () => {
      addToast('Critical error while updating', 'error')
    },
  })
}
