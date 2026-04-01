import { useMutation } from '@tanstack/react-query'

import { createVideoMetadataAction } from '@/lib/actions/video.actions'
import $api from '@/lib/config/api.config'
import { extractServerMessage, serverLog } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import { BaseServerResponse } from '@/types/general.types'
import { UploadVideoResponse } from '@/types/video.types'

export const useUploadVideoFile = () => {
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async (formData: FormData) => {
      try {
        const { data: response } = await $api.post<BaseServerResponse<UploadVideoResponse>>(
          `/videos/upload`,
          formData,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
          },
        )

        if (!response.success) {
          const msg = extractServerMessage(response.errors)
          throw new Error(msg ?? 'Something went wrong during upload')
        }

        return response
      } catch (error: unknown) {
        serverLog('UPLOAD_VIDEO_FILE_ERROR', error, false)

        throw error
      }
    },
    onError: error => {
      addToast(error.message, 'error')
    },
  })
}

export const useCreateVideoMetadata = () => {
  const addToast = useToastStore(state => state.addToast)
  // const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (formData: FormData) => await createVideoMetadataAction(formData),
    onSuccess: response => {
      if (response.success) {
        addToast('Video successfully published!', 'success')

        // queryClient.invalidateQueries({ queryKey: ['videos'] })
      } else {
        const msg = extractServerMessage(response.errors)
        addToast(msg || 'Error publishing video', 'error')
      }
    },
    onError: () => {
      addToast('Critical error while publishing', 'error')
    },
  })
}
