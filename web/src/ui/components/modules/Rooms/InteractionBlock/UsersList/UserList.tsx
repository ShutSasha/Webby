'use client'

import { useState, useMemo } from 'react'

import { useDebounce } from 'use-debounce'

import SearchIcon from '@/assets/icons/ic_search.svg'
import { useRoomMembersQuery } from '@/lib/hooks/api/room/useRoomMembersQuery'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import { RoomMember } from '@/types/room.types'

import UserItem from './UserItem'

type Props = {
  roomId: string
}

export default function UserList({ roomId }: Props) {
  const [search, setSearch] = useState('')

  const [debouncedSearch] = useDebounce(search, 500)

  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useRoomMembersQuery(
    roomId,
    debouncedSearch,
  )

  const members: RoomMember[] = useMemo(() => {
    return data?.pages.flatMap(page => page?.data?.items || []) || []
  }, [data])

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  return (
    <>
      <div className="relative group mb-2 shrink-0">
        <SearchIcon
          className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-foreground-faint
            group-focus-within:text-emerald-500 transition-colors stroke-[1.5px]"
        />
        <input
          type="text"
          placeholder="Search a user"
          value={search}
          onChange={e => setSearch(e.target.value)}
          className="w-full bg-neutral-900 rounded-md py-2 pl-10 pr-4 text-sm outline-none border border-transparent
            focus:border-emerald-500/70 transition-all placeholder:text-foreground-faint"
        />
      </div>

      <div
        className="flex-1 overflow-y-auto min-h-0 flex flex-col gap-1 pr-1 [&::-webkit-scrollbar]:w-1.5
          [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-neutral-900
          [&::-webkit-scrollbar-thumb]:border-0 [&::-webkit-scrollbar-thumb]:rounded-full
          hover:[&::-webkit-scrollbar-thumb]:bg-neutral-800"
      >
        {isLoading && members.length === 0 ? (
          <p className="text-center text-foreground-faint py-4 text-sm">Loading users...</p>
        ) : members.length === 0 ? (
          <p className="text-center text-foreground-faint py-4 text-sm">No users found</p>
        ) : (
          members.map((user, index) => {
            const isLast = members.length === index + 1
            const item = <UserItem key={user.userId} user={user} roomId={roomId} />

            if (isLast) {
              return (
                <div ref={lastElementRef} key={`last-${user.userId}`}>
                  {item}
                </div>
              )
            }

            return item
          })
        )}

        {isFetchingNextPage && (
          <div className="flex justify-center py-4 shrink-0">
            <div className="size-5 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
          </div>
        )}
      </div>
    </>
  )
}
