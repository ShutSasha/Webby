'use client'

import { useState } from 'react'

import { useSession } from 'next-auth/react'

import SendIcon from '@/assets/icons/shared/send_message.svg'
import { useGetChatDetailsQuery } from '@/lib/hooks/api/chat/useGetChatDetails'
import { useGetChatMessagesQuery } from '@/lib/hooks/api/chat/useGetChatMessages'
import { useSendMessageMutation } from '@/lib/hooks/api/chat/useSendMessage'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import SafeImage from '@/ui/components/shared/SafeImage'

type Props = {
  chatId: string
}

export default function ChatArea({ chatId }: Props) {
  const { data: session } = useSession()
  const currentUserId = session?.user?.id

  const [messageText, setMessageText] = useState('')

  const { data, isLoading, hasNextPage, fetchNextPage, isFetchingNextPage } = useGetChatMessagesQuery(chatId)
  const { mutate: sendMessage, isPending: isSending } = useSendMessageMutation()

  const messages = data?.pages.flatMap(page => page.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  const { data: chatDetailsResponse, isLoading: isChatDetailsLoading } = useGetChatDetailsQuery(chatId)
  const chatDetails = chatDetailsResponse?.data
  const otherUser = chatDetails?.user

  const handleSendMessage = () => {
    if (!messageText.trim() || isSending) return

    sendMessage(
      { chatId, content: messageText.trim() },
      {
        onSuccess: () => {
          setMessageText('')
        },
      },
    )
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      handleSendMessage()
    }
  }

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  return (
    <div className="flex flex-col h-full w-full">
      {/* Header */}
      <div
        className="h-[72px] shrink-0 border-b border-neutral-800/60 flex items-center justify-between px-6
          bg-neutral-900/20"
      >
        <div className="flex items-center gap-3">
          <SafeImage
            src={otherUser?.avatarUrl || ''}
            fallbackType="user"
            alt="Avatar"
            width={40}
            height={40}
            className="rounded-full object-cover size-10"
          />
          <span className="font-medium text-neutral-200">
            {isChatDetailsLoading
              ? 'Loading...'
              : otherUser?.username?.startsWith('@')
                ? otherUser.username
                : `@${otherUser?.username}`}
          </span>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto custom-scrollbar p-6 flex flex-col-reverse gap-2">
        {isLoading && messages.length === 0 ? (
          <div className="text-center text-neutral-500 my-auto">Loading messages...</div>
        ) : (
          messages.map((msg, index) => {
            const isMe = msg.sender.id === currentUserId
            const isLast = index === messages.length - 1

            const messageBlock = (
              <div className={`flex flex-col max-w-[70%] ${isMe ? 'self-end items-end' : 'self-start items-start'}`}>
                <div
                  className={`px-4 py-2.5 flex items-end gap-3 shadow-sm ${
                    isMe
                      ? 'bg-emerald-500 text-neutral-950 rounded-2xl rounded-br-sm border border-neutral-800/50'
                      : 'bg-neutral-800 text-neutral-300 rounded-2xl rounded-bl-sm font-medium'
                    }`}
                >
                  <p className="text-[15px] leading-relaxed break-all whitespace-pre-wrap">{msg.content}</p>
                  <span
                    className={`text-[10px] shrink-0 translate-y-0.5 ${
                      isMe ? 'text-emerald-900/60' : ' text-neutral-600'
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

      {/* Input Area */}
      <div className="p-4 bg-neutral-900/20 border-t border-neutral-800/60">
        <div className="relative flex items-center">
          <input
            type="text"
            value={messageText}
            onChange={e => setMessageText(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Write a message..."
            disabled={isSending}
            className="w-full bg-[#141414] border border-neutral-800 text-neutral-200 placeholder:text-neutral-600
              rounded-xl py-3.5 pl-5 pr-12 outline-none focus:border-neutral-600 transition-colors disabled:opacity-50"
          />
          <button
            onClick={handleSendMessage}
            disabled={!messageText.trim() || isSending}
            className="absolute right-3 p-1.5 text-neutral-500 hover:text-emerald-500 transition-colors
              disabled:hover:text-neutral-500 disabled:opacity-50"
          >
            <SendIcon className="size-5" />
          </button>
        </div>
      </div>
    </div>
  )
}
