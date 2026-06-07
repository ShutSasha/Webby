import { useMutation, useQueryClient } from '@tanstack/react-query'

import { togglePlaylistVideo, togglePlaylistStream } from '@/lib/actions/playlist.actions'
import { MediaType } from '@/types/general.types'

export const useTogglePlaylistMediaMutation = (userId: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({
      playlistId,
      mediaId,
      mediaType,
    }: {
      playlistId: string
      mediaId: string
      mediaType: MediaType
    }) => {
      const response =
        mediaType === 'LiveStream'
          ? await togglePlaylistStream(playlistId, mediaId)
          : await togglePlaylistVideo(playlistId, mediaId)

      if (!response.success) {
        throw new Error(response.message || 'Failed to toggle media')
      }

      return response
    },

    onSettled: async (_, __, variables) => {
      // TODO: remove promise
      await new Promise(resolve => setTimeout(resolve, 300))

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
