'use client'

import { useCallback, useEffect, useRef, useState } from 'react'

import { AnimatePresence } from 'framer-motion'
import { useDebouncedCallback } from 'use-debounce'

import { getPlaylistVideos, PlaylistVideo } from '@/app/api/playlists'
import SearchIcon from '@/assets/icons/ic_search.svg'
import { clog } from '@/lib/utils/utils'
import VideoItem from '@/ui/components/modules/Playlists/VideoItem'

type Props = {
  playlistId: string
}

const PAGE_SIZE = 20

export default function PlaylistQueueContainer({ playlistId }: Props) {
  const [videos, setVideos] = useState<PlaylistVideo[]>([])
  const [loading, setLoading] = useState(true)
  const [query, setQuery] = useState('')
  const [page, setPage] = useState(1)
  const [fetchingMore, setFetchingMore] = useState(false)
  const [hasMore, setHasMore] = useState(true)
  const observer = useRef<IntersectionObserver | null>(null)

  const lastVideoElementRef = useCallback(
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

  const fetchVideos = useCallback(
    async (searchQuery: string, targetPage: number, isInitial: boolean) => {
      if (isInitial) setLoading(true)
      else setFetchingMore(true)

      try {
        const response = await getPlaylistVideos(playlistId, searchQuery, targetPage, PAGE_SIZE)

        if (response.success && response.data) {
          const newItems = response.data.items

          setVideos(prev => (isInitial ? newItems : [...prev, ...newItems]))

          setHasMore(newItems.length === PAGE_SIZE)
        }
      } catch (error) {
        clog('FETCH_ERROR', error)
      } finally {
        setLoading(false)
        setFetchingMore(false)
      }
    },
    [playlistId],
  )

  const debouncedSearch = useDebouncedCallback((value: string) => {
    setPage(1)
    setVideos([])
    setHasMore(true)
    fetchVideos(value, 1, true)
  }, 400)

  useEffect(() => {
    if (page > 1) {
      fetchVideos(query, page, false)
    }
  }, [page])

  useEffect(() => {
    fetchVideos('', 1, true)
  }, [])

  const handleSearchChange = (value: string) => {
    setQuery(value)
    debouncedSearch(value)
  }

  return (
    <div className="w-[368px] flex flex-col gap-3 h-0 min-h-full shrink-0">
      {/* Search */}
      <div className="relative">
        <SearchIcon className="absolute top-1/2 -translate-y-1/2 left-4 h-5 w-5 text-neutral-600" aria-hidden="true" />
        <input
          id="search"
          type="text"
          autoComplete="off"
          className="focus:border-emerald-500 focus:ring-emerald-500 ring-[0.3px] ring-transparent block
            placeholder:text-neutral-600 focus:outline-none pl-12 py-3 rounded-2xl font-medium text-sm border
            border-border w-full bg-neutral-900/50 text-neutral-200 transition-all"
          placeholder="Search a video in playlist"
          value={query}
          onChange={e => handleSearchChange(e.target.value)}
        />
        {loading && query && (
          <div className="absolute right-4 top-1/2 -translate-y-1/2">
            <div className="size-4 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
          </div>
        )}
      </div>

      {/* videos list */}
      <div
        className="flex-1 flex flex-col gap-2 overflow-y-auto pr-1 [&::-webkit-scrollbar]:w-1.5
          [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-neutral-800
          [&::-webkit-scrollbar-thumb]:border-0 [&::-webkit-scrollbar-thumb]:rounded-full
          hover:[&::-webkit-scrollbar-thumb]:bg-neutral-700/50"
      >
        {loading && videos.length === 0 ? (
          <p className="text-center text-neutral-500 py-10">Loading queue...</p>
        ) : (
          <AnimatePresence mode="popLayout">
            {videos.map((video, index) => {
              if (videos.length === index + 1) {
                return (
                  <div ref={lastVideoElementRef} key={video.videoId}>
                    <VideoItem
                      id={video.videoId}
                      title={video.name}
                      thumbnail={video.previewUrl}
                      playlistId={playlistId}
                    />
                  </div>
                )
              }
              return (
                <VideoItem
                  key={video.videoId}
                  id={video.videoId}
                  title={video.name}
                  thumbnail={video.previewUrl}
                  playlistId={playlistId}
                />
              )
            })}

            {fetchingMore && (
              <div className="flex justify-center py-4">
                <div className="size-5 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
              </div>
            )}

            {!hasMore && videos.length > 0 && (
              <p className="text-center text-xs text-neutral-600 py-4 italic">End of playlist</p>
            )}

            {!loading && videos.length === 0 && <p className="text-center text-neutral-500 py-10">No videos found</p>}
          </AnimatePresence>
        )}
      </div>
    </div>
  )
}
