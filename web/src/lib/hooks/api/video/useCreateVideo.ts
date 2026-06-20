import { useMutation, useQueryClient } from '@tanstack/react-query'
import { AxiosProgressEvent } from 'axios'
import { useRouter } from 'next/navigation'

import { createVideoMetadataAction } from '@/lib/actions/video.actions'
import $api from '@/lib/config/api.config'
import { extractServerMessage, ServerActionError, serverLog, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import { useVideoDraftStore } from '@/stores/video-draft.store'
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
      queryClient.invalidateQueries({
        predicate: query => query.queryKey[0] === 'user-videos' && query.queryKey[2] === true,
      })
    },
    onError: error => {
      addToast(error.message, 'error')
    },
  })
}

export const useCreateVideoMetadata = () => {
  const addToast = useToastStore(state => state.addToast)
  const clearDraft = useVideoDraftStore(state => state.clearDraft)
  const queryClient = useQueryClient()
  const router = useRouter()

  return useMutation({
    mutationFn: async (formData: FormData) => {
      return unwrapServerAction(await createVideoMetadataAction(formData))
    },
    onSuccess: () => {
      router.push('/studio')

      queryClient.invalidateQueries({ queryKey: ['user-videos'] })
      queryClient.invalidateQueries({ queryKey: ['search-videos'] })

      // setTimeout here for clear draft after redirect to studio page.
      setTimeout(() => {
        clearDraft()
      }, 3000)
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message
      addToast(errorMessage, 'error')
    },
  })
}
