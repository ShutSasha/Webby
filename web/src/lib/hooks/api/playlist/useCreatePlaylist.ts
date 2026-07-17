import { InfiniteData, useMutation, useQueryClient } from '@tanstack/react-query'

import { createUserPlaylist } from '@/lib/actions/playlist.actions'
import { extractServerMessage, ServerActionError, unwrapServerAction } from '@/lib/utils/general.utils'
import { addPaginatedCacheItem } from '@/lib/utils/query.utils'
import { useToastStore } from '@/stores/toast-store'
import { PaginatedData } from '@/types/general.types'
import { CachedPlaylist } from '@/types/playlist.types'

type UseCreatePlaylistProps = {
  onSuccess: () => void
}

type PlaylistPage = PaginatedData<CachedPlaylist>

export function useCreatePlaylist({ onSuccess }: UseCreatePlaylistProps) {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  const { mutate, isPending } = useMutation({
    mutationFn: async ({ name, isPrivate }: { name: string; isPrivate: boolean }) => {
      return unwrapServerAction(await createUserPlaylist(name, isPrivate))
    },

    onSuccess: newPlaylist => {
      addToast('User playlist has been created successfully', 'success')
      onSuccess()

      const updateCache = (oldData: InfiniteData<PlaylistPage> | undefined) => {
        return addPaginatedCacheItem(oldData, newPlaylist)
      }

      queryClient.setQueriesData({ queryKey: ['search-user-playlists'] }, updateCache)
    },

    onError: (error: ServerActionError) => {
      const errorMessage = extractServerMessage(error.errors) || error.message
      addToast(errorMessage, 'error')
    },
  })

  const handleCreate = (name: string, isPrivate: boolean) => {
    const trimmedName = name.trim()

    if (!trimmedName) {
      addToast(`Playlist's name must be not empty`, 'info')
      return
    }

    mutate({ name: trimmedName, isPrivate })
  }

  return {
    handleCreate,
    isLoading: isPending,
  }
}
