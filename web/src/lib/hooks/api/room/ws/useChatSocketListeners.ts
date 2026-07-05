import { useEffect } from 'react'

import { InfiniteData, useQueryClient } from '@tanstack/react-query'

import { ChatMessage } from '@/lib/actions/chat.actions'
import { useRoomStore } from '@/stores/room.store'
import { BaseServerResponse, PaginatedData } from '@/types/general.types'

type ChatQueryData = InfiniteData<BaseServerResponse<PaginatedData<ChatMessage>>>

export const useChatSocketListeners = (chatId: string | undefined) => {
  const socket = useRoomStore(state => state.socket)
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!socket || !chatId) return

    const handleNewMessage = (payload: ChatMessage) => {
      if (!payload) return

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
            msg => msg.id.startsWith('temp-') && msg.content === payload.content && msg.sender.id === payload.sender.id,
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
    }

    socket.on('NEW_MESSAGE', handleNewMessage)

    return () => {
      socket.off('NEW_MESSAGE', handleNewMessage)
    }
  }, [socket, chatId, queryClient])
}
