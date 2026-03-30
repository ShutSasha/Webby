import { useMutation, useQueryClient } from '@tanstack/react-query'

import { createUserPlaylist } from '@/lib/actions/playlist.actions'
import { extractServerMessage, serverLog } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'

type UseCreatePlaylistProps = {
  onSuccess: () => void
}

export function useCreatePlaylist({ onSuccess }: UseCreatePlaylistProps) {
  const queryClient = useQueryClient()
  const addToast = useToastStore(state => state.addToast)

  const { mutate, isPending } = useMutation({
    mutationFn: async ({ name, isPrivate }: { name: string; isPrivate: boolean }) => {
      const response = await createUserPlaylist(name, isPrivate)

      if (!response.success) {
        const msg = extractServerMessage(response.errors)
        throw new Error(msg ?? `Something went wrong while creating ${name} playlist`)
      }

      return response
    },
    onSuccess: () => {
      addToast('User playlist has been created successfully', 'success')
      onSuccess()

      queryClient.invalidateQueries({
        queryKey: ['search-user-playlists'],
      })

      queryClient.invalidateQueries({
        queryKey: ['search-playlists'],
      })
    },
    onError: error => {
      serverLog('Handle submit for create user playlist', error, true)
      addToast(error.message, 'error')
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
