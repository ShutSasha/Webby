import { useMutation, useQueryClient } from '@tanstack/react-query'

import { updatePlaylistAction } from '@/lib/actions/playlist.actions'
import { extractServerMessage } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

type UpdatePlaylistArgs = {
  playlistId: string
  name: string
  isPrivate: boolean
}

export const useUpdatePlaylist = (options?: { onSuccess?: () => void }) => {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ playlistId, name, isPrivate }: UpdatePlaylistArgs) =>
      await updatePlaylistAction(playlistId, name, isPrivate),
    onSuccess: response => {
      if (response.success) {
        addToast('Playlist updated successfully!', 'success')

        queryClient.invalidateQueries({
          queryKey: ['search-user-playlists'],
        })
        queryClient.invalidateQueries({
          queryKey: ['search-playlists'],
        })

        if (options?.onSuccess) {
          options.onSuccess()
        }
      } else {
        const msg = extractServerMessage(response.errors)
        addToast(msg || 'Error updating playlist', 'error')
      }
    },
    onError: () => {
      addToast('Critical error while updating playlist', 'error')
    },
  })
}
