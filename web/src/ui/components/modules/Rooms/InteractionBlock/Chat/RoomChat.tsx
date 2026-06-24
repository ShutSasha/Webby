'use client'

import { useGetChatMessagesQuery } from '@/lib/hooks/api/chat/useGetChatMessages'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'

import RoomChatMessage, { RoomChatMessageSkeleton } from './RoomChatMessage'

type Props = {
  chatId: string
}

export default function RoomChat({ chatId }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useGetChatMessagesQuery(chatId)

  const messages = data?.pages.flatMap(page => page.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  if (isLoading && messages.length === 0) {
    return (
      <div className="flex-1 overflow-hidden flex flex-col-reverse gap-2 pr-1">
        {Array.from({ length: 30 }).map((_, i) => (
          <RoomChatMessageSkeleton key={`skeleton-${i}`} index={i} />
        ))}
      </div>
    )
  }

  return (
    <div
      className="flex-1 overflow-y-auto min-h-0 flex flex-col-reverse gap-1 pr-1 [&::-webkit-scrollbar]:w-1.5
        [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-surface
        [&::-webkit-scrollbar-thumb]:border-0 [&::-webkit-scrollbar-thumb]:rounded-full
        hover:[&::-webkit-scrollbar-thumb]:bg-background"
    >
      {messages.length === 0 ? (
        <div
          className="h-full flex items-center justify-center text-foreground-disabled text-sm italic rotate-180
            transform"
        >
          <span className="rotate-180">No messages yet. Be the first to say hello!</span>
        </div>
      ) : (
        messages.map((msg, index) => {
          const isLast = messages.length === index + 1

          const messageNode = (
            <RoomChatMessage
              key={msg.id}
              senderId={msg.sender.id}
              username={msg.sender.username}
              message={msg.content}
            />
          )

          if (isLast) {
            return (
              <div key={`last-${msg.id}`} ref={lastElementRef}>
                {isFetchingNextPage && (
                  <div className="flex justify-center py-2">
                    <div
                      className="size-4 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin"
                    />
                  </div>
                )}
                {messageNode}
              </div>
            )
          }

          return messageNode
        })
      )}
    </div>
  )
}
