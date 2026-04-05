'use client'

import { useSearchUserPlaylistsQuery } from '@/lib/hooks/api/playlist/useSearchUserPlaylists' // Шлях до твого хука
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'

import PlaylistItem, { PlaylistItemSkeleton } from '../Playlists/PlaylistItem'

type Props = {
  userId: string
}

export default function UserPlaylists({ userId }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useSearchUserPlaylistsQuery(userId, '')

  const playlists = data?.pages.flatMap(page => page?.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  if (isLoading && playlists.length === 0) {
    return <UserPlaylistsSkeleton />
  }

  if (!isLoading && playlists.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center py-20">
        <p className="text-neutral-500 text-center">No playlists found</p>
      </div>
    )
  }

  return (
    <>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
        {playlists.map((playlist, index) => {
          const isLast = playlists.length === index + 1

          const item = (
            <PlaylistItem
              key={playlist.playlistId}
              id={playlist.playlistId}
              src={playlist.playlistCover}
              name={playlist.name}
              creator={'creator'}
              videoCount={playlist.countOfVideos}
            />
          )

          if (isLast) {
            return (
              <div ref={lastElementRef} key={`last-${playlist.playlistId}`}>
                {item}
              </div>
            )
          }

          return item
        })}
      </div>

      {isFetchingNextPage && (
        <div className="w-full flex justify-center py-8">
          <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </>
  )
}

export function UserPlaylistsSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
      {Array.from({ length: 12 }).map((_, index) => (
        <PlaylistItemSkeleton key={index} />
      ))}
    </div>
  )
}
