import { useEffect, useRef } from 'react'

import { useQueryClient } from '@tanstack/react-query'
import io from 'socket.io-client'

import { getWsTokenAction } from '@/lib/actions/room.actions'
import { useRoomStore } from '@/stores/room.store'

import { ChatMessage } from '../actions/chat.actions'

export const useRoomWebSocket = (
  roomId: string | undefined,
  chatId: string | undefined,
  userId: string | undefined,
) => {
  const socketRef = useRef<SocketIOClient.Socket | null>(null)
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!chatId || !roomId || !userId) return

    let isMounted = true

    const connectSocket = async () => {
      const tokenRes = await getWsTokenAction()

      if (!tokenRes.success || !tokenRes.data || !isMounted) return

      socketRef.current = io(process.env.NEXT_PUBLIC_WS_URL || 'http://localhost:5000', {
        query: {
          token: tokenRes.data,
          chat_id: chatId,
        },
        transports: ['websocket'],
        upgrade: false,
        reconnection: false,
      })

      const socket = socketRef.current

      socket.on('connect_error', (error: any) => {
        console.error('[WS ROOM] Connection Error Detailed:', error)
      })

      socket.on('disconnect', (reason: string) => {
        console.warn('[WS ROOM] Disconnected. Reason:', reason)
      })

      socket.on('REPORT_TIMECODE', (payload: { synchronizeId: string }) => {
        if (payload?.synchronizeId) {
          useRoomStore.getState().setSyncTriggerId(payload.synchronizeId)
        }
      })

      socket.on('SYNCHRONIZE', (payload: { timecode: number }) => {
        if (payload?.timecode !== undefined) {
          useRoomStore.getState().setSyncTargetTimecode(payload.timecode)
        }
      })

      socket.on('QUEUE_UPDATED', (payload: { position: number[] }) => {
        queryClient.invalidateQueries({ queryKey: ['room-queue', roomId] })
      })

      socket.on('NEW_MESSAGE', (payload: ChatMessage) => {
        if (!payload) return

        queryClient.setQueryData(['chat-messages', chatId], (oldData: any) => {
          if (!oldData || !oldData.pages || oldData.pages.length === 0) {
            return oldData
          }

          const newPages = [...oldData.pages]
          const firstPage = { ...newPages[0] }

          if (firstPage.data && firstPage.data.items) {
            firstPage.data = {
              ...firstPage.data,
              items: [payload, ...firstPage.data.items],
            }
            newPages[0] = firstPage
          }

          const updatedData = { ...oldData, pages: newPages }

          return updatedData
        })
      })

      socket.on('ROOM_POINTS_UPDATED', (payload: { added_points: number; totals: Record<string, number> }) => {
        if (!payload || !payload.totals) return

        const myNewPoints = payload.totals[userId]

        if (myNewPoints !== undefined) {
          queryClient.setQueryData(['room-member-points', roomId], (oldData: any) => {
            if (!oldData) return oldData

            return {
              ...oldData,
              data: {
                ...oldData.data,
                points: myNewPoints,
              },
            }
          })
        }
      })
    }

    connectSocket()

    return () => {
      isMounted = false
      if (socketRef.current) {
        socketRef.current.disconnect()
      }
    }
  }, [chatId, roomId, userId, queryClient])

  return socketRef.current
}
