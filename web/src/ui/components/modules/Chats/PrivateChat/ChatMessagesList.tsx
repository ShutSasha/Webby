'use client'

import { useState } from 'react'

import { ChatMessage } from '@/lib/actions/chat.actions'
import { useDeleteMessageMutation } from '@/lib/hooks/api/chat/useDeleteMessage'
import { useGetChatMessagesQuery } from '@/lib/hooks/api/chat/useGetChatMessages'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'

import ChatMessagesListSkeleton from './ChatMessagesListSkeleton'
import MessageContextMenu from './MessageContextMenu'

type Props = {
  chatId: string
  currentUserId: string | undefined
}

export default function ChatMessagesList({ chatId, currentUserId }: Props) {
  const { data, isLoading, hasNextPage, fetchNextPage, isFetchingNextPage } = useGetChatMessagesQuery(chatId)
  const { mutate: deleteMessage } = useDeleteMessageMutation()

  const messages = data?.pages.flatMap(page => page.data?.items || []) || []

  const [contextMenu, setContextMenu] = useState<{
    isOpen: boolean
    x: number
    y: number
    messageId: string
    isOwnMessage: boolean
    content: string
  }>({ isOpen: false, x: 0, y: 0, messageId: '', isOwnMessage: false, content: '' })

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  const handleContextMenu = (e: React.MouseEvent, msg: ChatMessage, isMe: boolean) => {
    e.preventDefault()
    setContextMenu({
      isOpen: true,
      x: e.clientX,
      y: e.clientY,
      messageId: msg.id,
      isOwnMessage: isMe,
      content: msg.content,
    })
  }

  const handleCopy = () => {
    if (contextMenu.content) {
      navigator.clipboard.writeText(contextMenu.content)
    }
  }

  const handleDelete = () => {
    if (contextMenu.messageId) {
      deleteMessage({ chatId, messageId: contextMenu.messageId })
    }
  }

  return (
    <>
      <div className="flex-1 overflow-y-auto custom-scrollbar p-6 flex flex-col-reverse gap-2 relative">
        {isLoading && messages.length === 0 ? (
          <ChatMessagesListSkeleton />
        ) : (
          messages.map((msg, index) => {
            const isMe = msg.sender.id === currentUserId
            const isLast = index === messages.length - 1

            const messageBlock = (
              <div
                className={`flex flex-col max-w-[70%] ${isMe ? 'self-end items-end' : 'self-start items-start'}`}
                onContextMenu={e => handleContextMenu(e, msg, isMe)}
              >
                <div
                  className={`px-4 py-2.5 flex items-end gap-3 shadow-sm cursor-context-menu ${
                    isMe
                      ? 'bg-emerald-500 text-neutral-950 rounded-2xl rounded-br-sm border border-neutral-800/50'
                      : 'bg-neutral-800 text-neutral-300 rounded-2xl rounded-bl-sm font-medium'
                    }`}
                >
                  <p className="text-[15px] leading-relaxed break-all whitespace-pre-wrap">{msg.content}</p>
                  <span
                    className={`text-[10px] shrink-0 translate-y-0.5 ${
                      isMe ? 'text-emerald-900/75' : 'text-neutral-600'
                    }`}
                  >
                    {formatTime(msg.createdAt)}
                  </span>
                </div>
              </div>
            )

            if (isLast) {
              return (
                <div key={msg.id} ref={lastElementRef} className="flex flex-col">
                  {messageBlock}
                </div>
              )
            }

            return (
              <div key={msg.id} className="flex flex-col">
                {messageBlock}
              </div>
            )
          })
        )}

        {isFetchingNextPage && (
          <div className="text-center text-xs text-neutral-500 py-2">Loading older messages...</div>
        )}
      </div>

      <MessageContextMenu
        isOpen={contextMenu.isOpen}
        x={contextMenu.x}
        y={contextMenu.y}
        isOwnMessage={contextMenu.isOwnMessage}
        onClose={() => setContextMenu(prev => ({ ...prev, isOpen: false }))}
        onCopy={handleCopy}
        onDelete={handleDelete}
      />
    </>
  )
}
