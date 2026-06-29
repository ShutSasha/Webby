'use client'

import { useEffect } from 'react'

import { useQuery, useQueryClient } from '@tanstack/react-query'

import { checkVideoUploadStatusAction } from '@/lib/actions/video.actions'
import { unwrapServerAction } from '@/lib/utils/general.utils'

export const useCheckVideoUploadStatus = (videoId: string, isUploading: boolean) => {
  const queryClient = useQueryClient()

  const query = useQuery({
    queryKey: ['check-video-status', videoId],
    queryFn: async () => {
      return unwrapServerAction(await checkVideoUploadStatusAction(videoId))
    },
    enabled: isUploading,
    refetchInterval: isUploading ? 5000 : false,
  })

  useEffect(() => {
    if (query.data === true) {
      queryClient.invalidateQueries({ queryKey: ['user-videos'] })
    }
  }, [query.data, queryClient])

  return query
}
