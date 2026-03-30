'use client'

import { useSearchUserPlaylistsQuery } from '@/lib/hooks/api/playlist/useSearchUserPlaylists'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'

import { PlaylistItemSkeleton } from './PlaylistItem'
import UserPlaylistItem from './UserPlaylistItem'
import GridCardsContainer from '../../shared/GridCardsContainer'

type Props = {
  userId: string
  query: string
}

export default function PlaylistUserSearchContainer({ userId, query }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useSearchUserPlaylistsQuery(userId, query)

  const playlists = data?.pages.flatMap(page => page?.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  if (isLoading && playlists.length === 0) {
    return (
      <GridCardsContainer>
        {[...new Array(20)].map((_, idx) => (
          <PlaylistItemSkeleton key={idx} />
        ))}
      </GridCardsContainer>
    )
  }

  if (!isLoading && playlists.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center py-20">
        <p className="text-neutral-500 text-center">
          Playlists by query <span className="text-neutral-300">{`'${query}'`}</span> not found
        </p>
      </div>
    )
  }

  return (
    <>
      <GridCardsContainer>
        {playlists.map((playlist, index) => {
          const isLast = playlists.length === index + 1

          const item = (
            <UserPlaylistItem
              key={playlist.playlistId}
              id={playlist.playlistId}
              src={playlist.playlistCover}
              name={playlist.name}
              isPrivate={playlist.isPrivate}
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
      </GridCardsContainer>

      {isFetchingNextPage && (
        <div className="w-full flex justify-center py-8">
          <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </>
  )
}
