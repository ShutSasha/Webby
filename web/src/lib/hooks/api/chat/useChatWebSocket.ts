import { useEffect, useRef } from 'react'

import { useQueryClient } from '@tanstack/react-query'
import io from 'socket.io-client'

import { ChatMessage } from '@/lib/actions/chat.actions'
import { getWsTokenAction } from '@/lib/actions/room.actions'

export const useChatWebSocket = (chatId: string) => {
  const socketRef = useRef<SocketIOClient.Socket | null>(null)
  const queryClient = useQueryClient()

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

          return { ...oldData, pages: newPages }
        })
      })
    }

    connectSocket()

    return () => {
      isMounted = false
      if (socketRef.current) {
        socketRef.current.disconnect()
      }
    }
  }, [chatId, queryClient])

  return socketRef.current
}
