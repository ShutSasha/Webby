import { useMutation, useQueryClient } from '@tanstack/react-query'

import { togglePlaylistVideo } from '@/lib/actions/playlist.actions'
import { VideoSource } from '@/types/video.types'

export const useTogglePlaylistVideoMutation = (userId: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({
      playlistId,
      videoId,
      videoSource = 'Webby',
    }: {
      playlistId: string
      videoId: string
      videoSource?: VideoSource
    }) => {
      const response = await togglePlaylistVideo(playlistId, videoId, videoSource)

      if (!response.success) {
        throw new Error(response.message || 'Failed to toggle video')
      }

      return response
    },

    onSettled: (_, __, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['search-user-playlists', userId],
      })

      queryClient.invalidateQueries({
        queryKey: ['search-playlists'],
      })

      queryClient.invalidateQueries({
        queryKey: ['search-playlist-videos', variables.playlistId],
      })
    },
  })
}
