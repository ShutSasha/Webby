'use client'

import Image from 'next/image'
import Link from 'next/link'

import SearchIcon from '@/assets/icons/ic_search.svg'

const MOCK_CHATS = Array.from({ length: 10 }).map((_, i) => ({
  id: `chat-${i}`,
  username: '@guntersteam',
  lastMessage: 'Привіт, нещодавно дуже сильно хотів тебе запитати, щодо такої цікавої справи, як...',
  avatar: 'https://i.pravatar.cc/150?u=' + i,
}))

export default function ChatSidebar() {
  return (
    <div className="w-[340px] shrink-0 rounded-2xl border border-neutral-800/60 flex flex-col overflow-hidden">
      <div className="p-4 border-b border-neutral-800/60">
        <div className="relative">
          <input
            type="text"
            placeholder="Search for username..."
            className="w-full bg-neutral-900 text-neutral-200 placeholder:text-neutral-600 text-sm rounded-xl py-2.5
              pl-4 pr-10 outline-none ring-0 border border-transparent focus:border-neutral-700 transition-colors"
          />
          <SearchIcon className="absolute right-3 top-1/2 -translate-y-1/2 size-4 text-neutral-500" />
        </div>
      </div>
      {/* Chat List */}
      <div className="flex-1 overflow-y-auto custom-scrollbar p-2 flex flex-col gap-1">
        {MOCK_CHATS.map(chat => (
          <Link
            key={chat.id}
            href={`/chats/${chat.id}`}
            className="flex items-center gap-3 p-3 rounded-xl hover:bg-neutral-800/50 transition-colors cursor-pointer"
          >
            <Image
              src={chat.avatar}
              alt={chat.username}
              width={48}
              height={48}
              className="rounded-full object-cover shrink-0 size-12"
            />
            <div className="flex flex-col overflow-hidden">
              <span className="text-sm font-semibold text-neutral-200">{chat.username}</span>
              <span className="text-xs text-neutral-500 truncate mt-0.5">{chat.lastMessage}</span>
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
