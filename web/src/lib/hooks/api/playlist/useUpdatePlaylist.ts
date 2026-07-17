import { InfiniteData, useMutation, useQueryClient } from '@tanstack/react-query'

import { updatePlaylistAction } from '@/lib/actions/playlist.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { updatePaginatedCacheItem } from '@/lib/utils/query.utils'
import { useToastStore } from '@/stores/toast-store'
import { PaginatedData } from '@/types/general.types'
import { CachedPlaylist } from '@/types/playlist.types'

type UpdatePlaylistArgs = {
  playlistId: string
  name: string
  isPrivate: boolean
}

type PlaylistPage = PaginatedData<CachedPlaylist>

export const useUpdatePlaylist = (options?: { onSuccess?: () => void }) => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ playlistId, name, isPrivate }: UpdatePlaylistArgs) => {
      return unwrapServerAction(await updatePlaylistAction(playlistId, name, isPrivate))
    },

    onSuccess: updatedPlaylist => {
      addToast('Playlist updated successfully!', 'success')

      const updateCache = (oldData: InfiniteData<PlaylistPage> | undefined) => {
        return updatePaginatedCacheItem(oldData, updatedPlaylist, 'playlistId')
      }

      queryClient.setQueriesData({ queryKey: ['search-user-playlists'] }, updateCache)
      queryClient.setQueriesData({ queryKey: ['search-playlists'] }, updateCache)

      if (options?.onSuccess) {
        options.onSuccess()
      }
    },

    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message
      addToast(errorMessage, 'error')
    },
  })
}
