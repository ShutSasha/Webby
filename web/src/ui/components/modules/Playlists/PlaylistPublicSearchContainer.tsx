'use client'

import { useSearchPlaylistsQuery } from '@/lib/hooks/api/playlist/useSearchPlaylists'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'

import PlaylistItem, { PlaylistItemSkeleton } from './PlaylistItem'
import GridCardsContainer from '../../shared/GridCardsContainer'

type Props = {
  query: string
}

export default function PlaylistPublicSearchContainer({ query }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useSearchPlaylistsQuery(query)

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
        {query.length > 0 && (
          <p className="text-neutral-500 text-center">
            Playlists by query <span className="text-neutral-300">{`'${query}'`}</span> not found
          </p>
        )}
        {query.length === 0 && <p className="text-neutral-500 text-center">No playlists found</p>}
      </div>
    )
  }

  return (
    <>
      <GridCardsContainer>
        {playlists.map((playlist, index) => {
          const isLast = playlists.length === index + 1

          const item = (
            <PlaylistItem
              key={playlist.playlistId}
              id={playlist.playlistId}
              src={playlist.playlistCover}
              name={playlist.name}
              creator={playlist.username}
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
