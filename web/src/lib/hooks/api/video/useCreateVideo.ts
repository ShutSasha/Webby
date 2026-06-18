import { useMutation, useQueryClient } from '@tanstack/react-query'
import { AxiosProgressEvent } from 'axios'
import { useRouter } from 'next/navigation' 

import { createVideoMetadataAction } from '@/lib/actions/video.actions'
import $api from '@/lib/config/api.config'
import { extractServerMessage, serverLog } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import { BaseServerResponse } from '@/types/general.types'
import { UploadVideoResponse } from '@/types/video.types'

type UploadVideoPayload = {
  formData: FormData
  onUploadProgress?: (progressEvent: AxiosProgressEvent) => void
}

export const useUploadVideoFile = () => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ formData, onUploadProgress }: UploadVideoPayload) => {
      try {
        const { data: response } = await $api.post<BaseServerResponse<UploadVideoResponse>>(
          `/videos/upload`,
          formData,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
            onUploadProgress,
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
    onSuccess: () => {
      addToast('User video has been uploaded successfully', 'success')

      queryClient.invalidateQueries({
        predicate: query => query.queryKey[0] === 'user-videos' && query.queryKey[2] === true,
      })
    },
    onError: error => {
      addToast(error.message, 'error')
    },
  })
}


export const useCreateVideoMetadata = (options?: { onSuccess?: () => void }) => {
  const addToast = useToastStore(state => state.addToast)
  const queryClient = useQueryClient()
  const router = useRouter()

  return useMutation({
    mutationFn: async (formData: FormData) => await createVideoMetadataAction(formData),
    onSuccess: response => {
      if (response.success) {
        addToast('Video successfully published!', 'success')

        
        router.push('/studio')

        queryClient.invalidateQueries({ queryKey: ['user-videos'] })
        queryClient.invalidateQueries({ queryKey: ['search-videos'] })

        if (options?.onSuccess) {
          options.onSuccess()
        }
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
