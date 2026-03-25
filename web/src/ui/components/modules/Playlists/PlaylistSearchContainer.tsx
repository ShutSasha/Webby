'use client'

import { useCallback, useEffect, useRef, useState } from 'react'

import { searchPlaylists, SearchPlaylsit } from '@/app/api/playlists'

import PlaylistItem from './PlaylistItem'
import GridCardsContainer from '../../shared/GridCardsContainer'

type Props = {
  query: string
}

const PAGE_SIZE = 20

export default function PlaylistSearchContainer({ query }: Props) {
  const [playlists, setPlaylists] = useState<SearchPlaylsit[]>([])

  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [fetchingMore, setFetchingMore] = useState(false)
  const [hasMore, setHasMore] = useState(true)

  const observer = useRef<IntersectionObserver | null>(null)

  const lastElementRef = useCallback(
    (node: HTMLDivElement) => {
      if (loading || fetchingMore) return
      if (observer.current) observer.current.disconnect()

      observer.current = new IntersectionObserver(entries => {
        if (entries[0].isIntersecting && hasMore) {
          setPage(prevPage => prevPage + 1)
        }
      })

      if (node) observer.current.observe(node)
    },
    [loading, fetchingMore, hasMore],
  )

  const fetchPlaylists = useCallback(async (searchQuery: string, targetPage: number, isInitial: boolean) => {
    if (isInitial) setLoading(true)
    else setFetchingMore(true)

    try {
      const response = await searchPlaylists(searchQuery, targetPage, PAGE_SIZE)

      if (response.success && response.data) {
        const newItems = response.data.items

        setPlaylists(prev => (isInitial ? newItems : [...prev, ...newItems]))
        setHasMore(newItems.length === PAGE_SIZE)
      }
    } catch (error) {
      console.error('FETCH_PLAYLISTS_ERROR', error)
    } finally {
      setLoading(false)
      setFetchingMore(false)
    }
  }, [])

  useEffect(() => {
    setPage(1)
    setPlaylists([])
    setHasMore(true)
    fetchPlaylists(query, 1, true)
  }, [query, fetchPlaylists])

  useEffect(() => {
    if (page > 1) {
      fetchPlaylists(query, page, false)
    }
  }, [page, query, fetchPlaylists])

  if (loading && playlists.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center py-20">
        <div className="size-8 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
      </div>
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
