import { useEffect, useRef } from 'react'

import io from 'socket.io-client'

import { getWsTokenAction } from '@/lib/actions/room.actions'
import { useRoomStore } from '@/stores/room.store'

export const useRoomWebSocket = (chatId: string | undefined) => {
  const socketRef = useRef<SocketIOClient.Socket | null>(null)

  useEffect(() => {
    if (!chatId) return

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
    }

    connectSocket()

    return () => {
      isMounted = false
      if (socketRef.current) {
        socketRef.current.disconnect()
      }
    }
  }, [chatId])

  return socketRef.current
}
