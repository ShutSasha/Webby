import { useEffect, useRef } from 'react'

import { InfiniteData, useQueryClient } from '@tanstack/react-query'
import io from 'socket.io-client'

import { ChatMessage } from '@/lib/actions/chat.actions'
import { getWsTokenAction } from '@/lib/actions/room.actions'
import { BaseServerResponse, PaginatedData } from '@/types/general.types'

export const useChatWebSocket = (chatId: string) => {
  const socketRef = useRef<SocketIOClient.Socket | null>(null)
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!chatId) return

    let isMounted = true

    const connectSocket = async () => {
      const tokenRes = await getWsTokenAction()

      if (!tokenRes.success || !tokenRes.data || !isMounted) return

      socketRef.current = io(process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000', {
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

        type ChatQueryData = InfiniteData<BaseServerResponse<PaginatedData<ChatMessage>>>

        queryClient.setQueryData<ChatQueryData>(['chat-messages', chatId], oldData => {
          if (!oldData || !oldData.pages || oldData.pages.length === 0) {
            return oldData
          }

          const allItems = oldData.pages.flatMap(p => p.data?.items || [])
          if (allItems.some(msg => msg.id === payload.id)) {
            return oldData
          }

          const newPages = [...oldData.pages]
          const firstPage = { ...newPages[0] }

          if (firstPage.data && firstPage.data.items) {
            const tempMsgIndex = firstPage.data.items.findIndex(
              msg =>
                msg.id.startsWith('temp-') && msg.content === payload.content && msg.sender?.id === payload.sender?.id,
            )

            if (tempMsgIndex !== -1) {
              const newItems = [...firstPage.data.items]
              newItems[tempMsgIndex] = payload
              firstPage.data = {
                ...firstPage.data,
                items: newItems,
              }
            } else {
              firstPage.data = {
                ...firstPage.data,
                items: [payload, ...firstPage.data.items],
              }
            }
            newPages[0] = firstPage
          }

          return { ...oldData, pages: newPages }
        })

        queryClient.setQueriesData({ queryKey: ['chat-history'] }, (oldData: any) => {
          if (!oldData || !oldData.pages) return oldData

          let updatedChat: any = null

          const newPages = oldData.pages.map((page: any) => {
            if (!page.data || !page.data.items) return page

            const filteredItems = page.data.items.filter((chat: any) => {
              if (chat.chatId === chatId) {
                updatedChat = { ...chat }
                return false
              }
              return true
            })

            return {
              ...page,
              data: { ...page.data, items: filteredItems },
            }
          })

          if (updatedChat) {
            updatedChat.lastMessage = {
              content: payload.content,
              createdAt: payload.createdAt,
            }

            if (newPages[0]?.data?.items) {
              newPages[0].data.items.unshift(updatedChat)
            }
          }

          return { ...oldData, pages: newPages }
        })
      })

      socket.on('MESSAGE_DELETED', (payload: { id: string }) => {
        if (!payload || !payload.id) return

        queryClient.setQueryData(['chat-messages', chatId], (oldData: any) => {
          if (!oldData || !oldData.pages) return oldData

          const newPages = oldData.pages.map((page: any) => {
            if (!page.data || !page.data.items) return page

            return {
              ...page,
              data: {
                ...page.data,
                items: page.data.items.filter((msg: any) => msg.id !== payload.id),
              },
            }
          })

          return { ...oldData, pages: newPages }
        })
      })

      socket.on('MESSAGE_UPDATED', (payload: ChatMessage) => {
        if (!payload || !payload.id) return

        queryClient.setQueryData(['chat-messages', chatId], (oldData: any) => {
          if (!oldData || !oldData.pages) return oldData

          const newPages = oldData.pages.map((page: any) => {
            if (!page.data || !page.data.items) return page

            return {
              ...page,
              data: {
                ...page.data,
                items: page.data.items.map((msg: any) =>
                  msg.id === payload.id
                    ? { ...msg, content: payload.content, isEdited: payload.isEdited ?? true }
                    : msg,
                ),
              },
            }
          })

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
