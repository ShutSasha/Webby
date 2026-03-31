import { useMutation, useQueryClient } from '@tanstack/react-query'

import { deletePlaylist } from '@/lib/actions/playlist.actions'
import { extractServerMessage, serverLog } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

export function useDeletePlaylist() {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  return useMutation({
    mutationFn: async ({ playlistId }: { playlistId: string }) => {
      const response = await deletePlaylist(playlistId)

      if (!response.success) {
        const msg = extractServerMessage(response.errors)
        throw new Error(msg ?? `Something went wrong while deleting playlist`)
      }

      return response
    },
    onSuccess: () => {
      addToast('User playlist has been deleted successfully', 'success')

      queryClient.invalidateQueries({
        queryKey: ['search-user-playlists'],
      })

      queryClient.invalidateQueries({
        queryKey: ['search-playlists'],
      })
    },
    onError: error => {
      serverLog('Handle submit for delete user playlist', error, true)
      addToast(error.message, 'error')
    },
  })
}
