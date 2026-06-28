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
  onEditMessage: (msg: { id: string; content: string }) => void
}

export default function ChatMessagesList({ chatId, currentUserId, onEditMessage }: Props) {
  const { data, isLoading, hasNextPage, fetchNextPage, isFetchingNextPage } = useGetChatMessagesQuery(chatId)
  const { mutate: deleteMessage, isPending: isDeletingMessage, variables: deleteVars } = useDeleteMessageMutation()

  const rawMessages = data?.pages.flatMap(page => page.data?.items || []) || []
  const messages = Array.from(new Map(rawMessages.map(msg => [msg.id, msg])).values())

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

  const handleEdit = () => {
    if (contextMenu.messageId && contextMenu.content) {
      onEditMessage({ id: contextMenu.messageId, content: contextMenu.content })
    }
  }

  return (
    <>
      <div className="flex-1 overflow-y-auto p-6 flex flex-col-reverse gap-2 relative">
        <div className="shrink-0 h-px w-full" style={{ overflowAnchor: 'auto' }} />
        {isLoading && messages.length === 0 ? (
          <ChatMessagesListSkeleton />
        ) : (
          messages.map((msg, index) => {
            const isMe = msg?.sender?.id === currentUserId
            const isLast = index === messages.length - 1

            const isDeleting = isDeletingMessage && deleteVars?.messageId === msg.id

            const messageBlock = (
              <div
                className={`flex flex-col max-w-[70%] transition-all duration-300 ease-out ${
                  isMe ? 'self-end items-end' : 'self-start items-start'
                }
                  ${isDeleting ? 'opacity-40 scale-[0.98] pointer-events-none blur-[0.5px]' : 'opacity-100'}`}
                onContextMenu={e => handleContextMenu(e, msg, isMe)}
              >
                <div
                  className={`px-4 py-2.5 flex items-end gap-3 shadow-sm cursor-context-menu ${
                    isMe
                      ? `bg-emerald-500 dark:text-foreground-inverse text-neutral-900 rounded-2xl rounded-br-sm border
                        border-border/50`
                      : 'bg-background text-foreground-subtle rounded-2xl rounded-bl-sm font-medium'
                    }`}
                >
                  <p className="text-[15px] leading-relaxed break-all whitespace-pre-wrap">{msg.content}</p>

                  <div
                    className={`text-[10px] shrink-0 translate-y-0.5 flex gap-1.5 items-center ${
                      isMe ? 'text-emerald-900/75' : 'text-foreground-disabled'
                    }`}
                  >
                    {msg.isEdited && <span>edited</span>}
                    <span>{formatTime(msg.createdAt)}</span>
                  </div>
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
          <div className="text-center text-xs text-foreground-faint py-2">Loading older messages...</div>
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
        onEdit={handleEdit}
      />
    </>
  )
}
