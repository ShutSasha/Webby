import { InfiniteData, useMutation, useQueryClient } from '@tanstack/react-query'

import { deletePlaylist } from '@/lib/actions/playlist.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { removePaginatedCacheItem } from '@/lib/utils/query.utils'
import { useToastStore } from '@/stores/toast-store'
import { PaginatedData } from '@/types/general.types'
import { CachedPlaylist } from '@/types/playlist.types'

type PlaylistPage = PaginatedData<CachedPlaylist>

export function useDeletePlaylist() {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ playlistId }: { playlistId: string }) => {
      return unwrapServerAction(await deletePlaylist(playlistId))
    },

    onSuccess: (_, variables) => {
      addToast('User playlist has been deleted successfully', 'success')

      const { playlistId } = variables

      const updatePaginationCache = (oldData: InfiniteData<PlaylistPage> | undefined) => {
        return removePaginatedCacheItem(oldData, playlistId, 'playlistId')
      }

      queryClient.setQueriesData({ queryKey: ['search-user-playlists'] }, updatePaginationCache)
      queryClient.setQueriesData({ queryKey: ['search-playlists'] }, updatePaginationCache)
    },
    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message
      addToast(errorMessage, 'error')
    },
  })
}
