import { useEffect } from 'react'

import { useQueryClient } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'

import { useRoomStore } from '@/stores/room.store'
import { useToastStore } from '@/stores/toast-store'

type RemoveMemberPayload = {
  memberId: string
}

export const useRoomMemberSocketListeners = (roomId: string | undefined, currentUserId: string | undefined) => {
  const socket = useRoomStore(state => state.socket)
  const addToast = useToastStore(state => state.addToast)
  const router = useRouter()
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!socket || !roomId) return

    const handleRemoveMember = (payload: RemoveMemberPayload) => {
      if (!payload || !payload.memberId) return

      if (payload.memberId === currentUserId) {
        queryClient.invalidateQueries({ queryKey: ['room', roomId] })
        addToast('You have been removed from the room by the host.', 'error')
        router.replace('/')
        router.refresh()
        return
      }

      queryClient.invalidateQueries({ queryKey: ['room-members', roomId] })
    }

    socket.on('REMOVE_MEMBER', handleRemoveMember)

    return () => {
      socket.off('REMOVE_MEMBER', handleRemoveMember)
    }
  }, [socket, roomId, currentUserId, router, addToast, queryClient])
}
