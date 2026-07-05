import { useEffect } from 'react'

import io from 'socket.io-client'

import { getWsTokenAction } from '@/lib/actions/room.actions'
import { useRoomStore } from '@/stores/room.store'

export const useRoomSocketConnection = (roomId?: string, chatId?: string) => {
  const setSocket = useRoomStore(state => state.setSocket)

  useEffect(() => {
    if (!chatId || !roomId) return

    let isMounted = true
    let socket: SocketIOClient.Socket | null = null

    const connect = async () => {
      const tokenRes = await getWsTokenAction()

      if (!tokenRes.success || !tokenRes.data || !isMounted) return

      socket = io(process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000', {
        query: {
          token: tokenRes.data,
          chat_id: chatId,
        },
        transports: ['websocket'],
        upgrade: false,
        reconnection: false,
      })

      socket.on('connect_error', (error: unknown) => {
        console.error('[WS ROOM] Connection Error Detailed:', error)
      })

      socket.on('disconnect', (reason: string) => {
        console.warn('[WS ROOM] Disconnected. Reason:', reason)
      })

      setSocket(socket)
    }

    connect()

    return () => {
      isMounted = false
      if (socket) {
        socket.disconnect()
        setSocket(null)
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [chatId, roomId])
}
