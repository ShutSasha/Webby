import { useMutation, useQueryClient } from '@tanstack/react-query'

import { cancelVideoUploadAction } from '@/lib/actions/video.actions'
import { extractServerMessage, serverLog } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export function useCancelUploadVideo() {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async (videoId: string) => {
      const response = await cancelVideoUploadAction(videoId)

      if (!response.success) {
        const msg = extractServerMessage(response.errors)
        throw new Error(msg ?? `Something went wrong while canceling video upload`)
      }

      return response
    },
    onSuccess: () => {
      addToast('User video upload has been canceled successfully', 'success')

      queryClient.invalidateQueries({
        predicate: query => query.queryKey[0] === 'user-videos' && query.queryKey[2] === true,
      })
    },
    onError: error => {
      serverLog('Handle submit for cancel upload user video', error, true)
      addToast(error.message, 'error')
    },
  })
}
