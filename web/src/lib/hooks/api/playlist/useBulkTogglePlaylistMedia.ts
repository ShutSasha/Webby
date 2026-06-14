import { useMutation, useQueryClient } from '@tanstack/react-query'

import { bulkTogglePlaylistsStream, bulkTogglePlaylistsVideo } from '@/lib/actions/playlist.actions'
import { MediaType } from '@/types/general.types'

export const useBulkTogglePlaylistMediaMutation = (userId: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({
      playlistIds,
      mediaId,
      mediaType,
    }: {
      playlistIds: string[]
      mediaId: string
      mediaType: MediaType
    }) => {
      const response =
        mediaType === 'LiveStream'
          ? await bulkTogglePlaylistsStream(playlistIds, mediaId)
          : await bulkTogglePlaylistsVideo(playlistIds, mediaId)

      if (!response.success) {
        throw new Error(response.message || 'Failed to update playlists')
      }

      return response
    },

    onSuccess: async (_, variables) => {
      const invalidationPromises = [
        queryClient.invalidateQueries({ queryKey: ['search-user-playlists', userId] }),
        queryClient.invalidateQueries({ queryKey: ['search-playlists'] }),
      ]

      variables.playlistIds.forEach(id => {
        invalidationPromises.push(queryClient.invalidateQueries({ queryKey: ['search-playlist-videos', id] }))
      })

      await Promise.all(invalidationPromises)
    },
  })
}
