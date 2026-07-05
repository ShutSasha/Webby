import { useEffect } from 'react'

import { useQueryClient } from '@tanstack/react-query'

import { useRoomStore } from '@/stores/room.store'

export const usePlayerSocketListeners = (roomId: string | undefined) => {
  const socket = useRoomStore(state => state.socket)
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!socket || !roomId) return

    const handleReportTimecode = (payload: { synchronizeId: string }) => {
      if (payload?.synchronizeId) {
        useRoomStore.getState().setSyncTriggerId(payload.synchronizeId)
      }
    }

    const handleSynchronize = (payload: { timecode: number }) => {
      if (payload?.timecode !== undefined) {
        useRoomStore.getState().setSyncTargetTimecode(payload.timecode)
      }
    }

    const handleQueueUpdated = () => {
      queryClient.invalidateQueries({ queryKey: ['room-queue', roomId] })
    }

    socket.on('REPORT_TIMECODE', handleReportTimecode)
    socket.on('SYNCHRONIZE', handleSynchronize)
    socket.on('QUEUE_UPDATED', handleQueueUpdated)

    return () => {
      socket.off('REPORT_TIMECODE', handleReportTimecode)
      socket.off('SYNCHRONIZE', handleSynchronize)
      socket.off('QUEUE_UPDATED', handleQueueUpdated)
    }
  }, [socket, roomId, queryClient])
}
