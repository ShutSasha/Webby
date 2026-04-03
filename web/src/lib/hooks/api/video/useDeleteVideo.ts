import { useMutation, useQueryClient } from '@tanstack/react-query'

import { deleteVideo } from '@/lib/actions/video.actions'
import { extractServerMessage, serverLog } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export function useDeleteVideo() {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ videoId }: { videoId: string }) => {
      const response = await deleteVideo(videoId)

      if (!response.success) {
        const msg = extractServerMessage(response.errors)
        throw new Error(msg ?? `Something went wrong while deleting video`)
      }

      return response
    },
    onSuccess: () => {
      addToast('User video has been deleted successfully', 'success')

      queryClient.invalidateQueries({
        queryKey: ['user-videos'],
      })

      queryClient.invalidateQueries({
        queryKey: ['search-videos'],
      })
    },
    onError: error => {
      serverLog('Handle submit for delete user video', error, true)
      addToast(error.message, 'error')
    },
  })
}
