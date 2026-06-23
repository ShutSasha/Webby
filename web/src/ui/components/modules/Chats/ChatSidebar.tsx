'use client'

import { useState } from 'react'

import Link from 'next/link'
import { useParams } from 'next/navigation'
import { useDebounce } from 'use-debounce'

import SearchIcon from '@/assets/icons/ic_search.svg'
import { useGetChatHistoryQuery } from '@/lib/hooks/api/chat/useGetChatHistory'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import { cn } from '@/lib/utils/general.utils'
import SafeImage from '@/ui/components/shared/SafeImage'

import ChatItemSkeleton from './ChatItemSkeleton'

export default function ChatSidebar() {
  const [searchValue, setSearchValue] = useState('')
  const [debouncedSearch] = useDebounce(searchValue, 300)
  const params = useParams()
  const currentChatId = params?.chatId as string | undefined

  const { data, isLoading, hasNextPage, fetchNextPage, isFetchingNextPage } = useGetChatHistoryQuery(debouncedSearch)

  const chats = data?.pages.flatMap(page => page.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  return (
    <div className="w-[300px] lg:w-[340px] shrink-0 rounded-2xl border border-neutral-800/60 flex flex-col
      overflow-hidden">
      <div className="p-4 border-b border-neutral-800/60">
        <div className="relative">
          <input
            type="text"
            value={searchValue}
            onChange={e => setSearchValue(e.target.value)}
            placeholder="Search for username..."
            className="w-full bg-neutral-900 text-foreground-tertiary placeholder:text-foreground-disabled text-sm
              rounded-xl py-2.5 pl-4 pr-10 outline-none ring-0 border border-transparent focus:border-neutral-700
              transition-colors"
          />
          <SearchIcon className="absolute right-3 top-1/2 -translate-y-1/2 size-4 text-foreground-faint" />
        </div>
      </div>

      {/* Chat List */}
      <div className="flex-1 overflow-y-auto p-2 flex flex-col gap-1">
        {isLoading && chats.length === 0 ? (
          Array.from({ length: 9 }).map((_, index) => <ChatItemSkeleton key={index} />)
        ) : chats.length > 0 ? (
          chats.map((chat, index) => {
            const isLast = chats.length === index + 1

            const chatItem = (
              <Link
                href={`/chats/${chat.chatId}`}
                className={cn(
                  'flex items-center gap-3 p-3 rounded-xl transition-colors cursor-pointer border',
                  currentChatId === chat.chatId
                    ? 'bg-emerald-500/10 border-emerald-500/20'
                    : 'border-transparent hover:bg-neutral-800/50',
                )}
              >
                <SafeImage
                  src={chat.user.avatarUrl}
                  fallbackType="user"
                  alt={chat.user.username}
                  width={48}
                  height={48}
                  className="rounded-full object-cover shrink-0 size-12"
                />
                <div className="flex flex-col overflow-hidden">
                  <span className="text-sm font-semibold text-foreground-tertiary">
                    {chat.user.username.startsWith('@') ? chat.user.username : `@${chat.user.username}`}
                  </span>
                  <span className="text-xs text-foreground-faint truncate mt-0.5">
                    {chat.lastMessage?.content || 'No messages yet'}
                  </span>
                </div>
              </Link>
            )

            if (isLast) {
              return (
                <div ref={lastElementRef} key={`last-${chat.chatId}`}>
                  {chatItem}
                </div>
              )
            }

            return <div key={chat.chatId}>{chatItem}</div>
          })
        ) : (
          <div className="p-4 text-center text-sm text-foreground-faint">
            {debouncedSearch ? 'No chats found for this search' : 'No chats yet'}
          </div>
        )}

        {isFetchingNextPage && <div className="p-2 text-center text-xs text-foreground-faint">Loading more...</div>}
      </div>
    </div>
  )
}
