'use client'

import { useState, useMemo } from 'react'

import { useParams } from 'next/navigation'

import SearchIcon from '@/assets/icons/ic_search.svg'
import { useRoomQueueQuery } from '@/lib/hooks/api/room/useRoomQueueQuery'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import { RoomQueueItem } from '@/types/room.types'

import QueueItem from './QueueItem'
import QueueItemSkeleton from './QueueItemSkeleton'

export default function RoomQueue() {
  const params = useParams()
  const roomId = params?.id as string

  const [search, setSearch] = useState('')

  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useRoomQueueQuery(roomId)

  const queueItems: RoomQueueItem[] = useMemo(() => {
    return data?.pages.flatMap(page => page?.data?.items || []) || []
  }, [data])

  const filteredPlaylist = useMemo(() => {
    if (!search.trim()) return queueItems
    return queueItems.filter(v => (v?.title || '').toLowerCase().includes(search.toLowerCase()))
  }, [queueItems, search])

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  return (
    <div className="flex flex-col h-full min-h-0">
      <div className="relative group mb-2 shrink-0">
        <SearchIcon
          className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-foreground-faint
            group-focus-within:text-emerald-500 transition-colors stroke-[1.5px]"
        />
        <input
          type="text"
          placeholder="Search video in queue"
          value={search}
          onChange={e => setSearch(e.target.value)}
          className="w-full bg-neutral-900 rounded-md py-2 pl-10 pr-4 text-sm outline-none border border-transparent
            focus:border-emerald-500/70 transition-all placeholder:text-foreground-faint"
        />
      </div>

      <div className="flex-1 overflow-y-auto min-h-0 flex flex-col gap-2 scrollbar-hide">
        {isLoading && queueItems.length === 0 ? (
          <div className="flex flex-col gap-2">
            {Array.from({ length: 6 }).map((_, i) => (
              <QueueItemSkeleton key={i} />
            ))}
          </div>
        ) : filteredPlaylist.length === 0 ? (
          <p className="text-foreground-faint text-center py-10 text-sm">
            {search ? 'No videos match your search' : 'Queue is empty'}
          </p>
        ) : (
          filteredPlaylist.map((video, index) => {
            if (!video || !video.id) return null

            const isLast = filteredPlaylist.length === index + 1
            const item = <QueueItem key={video.id} roomId={roomId} video={video} />

            if (isLast) {
              return (
                <div ref={lastElementRef} key={`last-${video.id}`}>
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
    </div>
  )
}
