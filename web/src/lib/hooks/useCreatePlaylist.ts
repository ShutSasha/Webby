import { useState, useTransition } from 'react'

import { useRouter } from 'next/navigation'

import { createUserPlaylist } from '@/app/api/playlists'
import { extractServerMessage, serverLog } from '@/lib/utils/utils'
import { useToastStore } from '@/stores/toast-store'

type UseCreatePlaylistProps = {
  onSuccess: () => void
}

export function useCreatePlaylist({ onSuccess }: UseCreatePlaylistProps) {
  const router = useRouter()
  const addToast = useToastStore(state => state.addToast)
  const [isPending, startTransition] = useTransition()
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleCreate = async (name: string, isPrivate: boolean) => {
    const trimmedName = name.trim()
    if (!trimmedName) {
      addToast(`Playlist's name must be not empty`, 'info')
      return
    }

    try {
      setIsSubmitting(true)
      const response = await createUserPlaylist(trimmedName, isPrivate)

      if (response.success) {
        addToast('User playlist has been created successfully', 'success')
        onSuccess()

        startTransition(() => {
          router.refresh()
        })
      } else {
        const msg = extractServerMessage(response.errors)
        addToast(msg ?? `Something went wrong while creating ${name} playlist`, 'error')
      }
    } catch (error: unknown) {
      serverLog('Handle submit for create user playlist', error, true)
    } finally {
      setIsSubmitting(false)
    }
  }

  return {
    handleCreate,
    isLoading: isSubmitting || isPending,
  }
}
