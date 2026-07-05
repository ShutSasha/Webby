import { InfiniteData, useMutation, useQueryClient } from '@tanstack/react-query'

import { AdminVideoRecord, banVideoAction } from '@/lib/actions/admin.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import { PaginatedData } from '@/types/general.types'

export const useBanVideoMutation = () => {
  const addToast = useToastStore(state => state.addToast)
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (videoId: string) => {
      return unwrapServerAction(await banVideoAction(videoId))
    },
    onSuccess: (_, videoId) => {
      queryClient.setQueriesData(
        { queryKey: ['moderation-videos-search'] },
        (oldData: InfiniteData<PaginatedData<AdminVideoRecord>> | undefined) => {
          if (!oldData) return oldData

          return {
            ...oldData,
            pages: oldData.pages.map(page => ({
              ...page,
              items: page.items.map(video => (video.videoId === videoId ? { ...video, isBanned: true } : video)),
            })),
          }
        },
      )
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message
      addToast(errorMessage, 'error')
    },
  })
}
