'use client'

import { useCallback, useEffect } from 'react'

import { SearchPlaylist, searchPlaylists } from '@/app/api/playlists'
import { useInfiniteSearch } from '@/lib/hooks/useInfiniteSearch'

import PlaylistItem, { PlaylistItemSkeleton } from './PlaylistItem'
import GridCardsContainer from '../../shared/GridCardsContainer'

type Props = {
  query: string
}

const PAGE_SIZE = 20

export default function PlaylistPublicSearchContainer({ query }: Props) {
  const fetchPlaylistsFn = useCallback((searchQuery: string, page: number, pageSize: number) => {
    return searchPlaylists(searchQuery, page, pageSize)
  }, [])

  const {
    items: playlists,
    loading,
    fetchingMore,
    lastElementRef,
    searchImmediate,
  } = useInfiniteSearch<SearchPlaylist>({
    fetchFn: fetchPlaylistsFn,
    pageSize: PAGE_SIZE,
  })

  useEffect(() => {
    searchImmediate(query)
  }, [query, searchImmediate])

  if (loading && playlists.length === 0) {
    return (
      <GridCardsContainer>
        {[...new Array(20)].map((_, idx) => (
          <PlaylistItemSkeleton key={idx} />
        ))}
      </GridCardsContainer>
    )
  }

  if (!loading && playlists.length === 0) {
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

      {fetchingMore && (
        <div className="w-full flex justify-center py-8">
          <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </>
  )
}
