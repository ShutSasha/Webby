import { useEffect } from 'react'

import { useQueryClient } from '@tanstack/react-query'

import { MemberPoints } from '@/lib/actions/room.actions'
import { useRoomStore } from '@/stores/room.store'

export type ReactionSentPayload = {
  id: string
  name: string
  cost: number
  stickerUrl: string
}

export const useEntertainmentRoomListeners = (userId: string | undefined, roomId: string | undefined) => {
  const socket = useRoomStore(state => state.socket)
  const addReaction = useRoomStore(state => state.addReaction)
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!socket || !userId || !roomId) return

    const handleRoomPointsUpdated = (payload: { added_points: number; totals: Record<string, number> }) => {
      if (!payload || !payload.totals) return

      const myNewPoints = payload.totals[userId]

      if (myNewPoints !== undefined) {
        queryClient.setQueryData(['room-member-points', roomId], (oldData: MemberPoints | undefined) => {
          if (!oldData) return { points: myNewPoints }

          return {
            ...oldData,
            points: myNewPoints,
          }
        })
      }
    }

    const handleReactionSent = (payload: ReactionSentPayload) => {
      if (!payload) return

      addReaction(payload)
    }

    socket.on('ROOM_POINTS_UPDATED', handleRoomPointsUpdated)
    socket.on('REACTION_SENT', handleReactionSent)

    return () => {
      socket.off('ROOM_POINTS_UPDATED', handleRoomPointsUpdated)
      socket.off('REACTION_SENT', handleReactionSent)
    }
  }, [socket, roomId, queryClient, userId, addReaction])
}
