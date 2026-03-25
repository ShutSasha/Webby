'use client'

import { useCallback, useEffect, useRef, useState } from 'react'

import { AnimatePresence } from 'framer-motion'
import { useSearchParams } from 'next/navigation'
import { useDebouncedCallback } from 'use-debounce'

import { getPlaylistVideos, PlaylistVideo } from '@/app/api/playlists'
import { usePlaylistAutoPlay } from '@/lib/hooks/usePlaylistAutoPlay'
import { serverLog } from '@/lib/utils/utils'
import VideoItem from '@/ui/components/modules/Playlists/VideoItem'

import PlaylistQueueFooter from './PlaylistQueueFooter'
import QueueContainer from './QueueContainer'
import Search from './Search'

type Props = {
  playlistId: string
}

const PAGE_SIZE = 20

export default function PlaylistQueueContainer({ playlistId }: Props) {
  const [videos, setVideos] = useState<PlaylistVideo[]>([])

  const [loading, setLoading] = useState<boolean>(true)
  const [appliedQuery, setAppliedQuery] = useState<string>('')

  const [page, setPage] = useState<number>(1)
  const [fetchingMore, setFetchingMore] = useState<boolean>(false)
  const [hasMore, setHasMore] = useState<boolean>(true)

  const observer = useRef<IntersectionObserver | null>(null)

  const searchParams = useSearchParams()
  const currentV = searchParams.get('v')

  const [optimisticId, setOptimisticId] = useState<string | null>(null)

  usePlaylistAutoPlay({
    videos,
    currentV,
    playlistId,
    setOptimisticId,
  })

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
        serverLog('FETCH_ERROR', error)
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
    setAppliedQuery(value)
    fetchVideos(value, 1, true)
  }, 400)

  useEffect(() => {
    setOptimisticId(null)
  }, [currentV])

  useEffect(() => {
    if (page > 1) {
      fetchVideos(appliedQuery, page, false)
    }
  }, [page, appliedQuery, fetchVideos])

  useEffect(() => {
    fetchVideos('', 1, true)
  }, [fetchVideos])

  const handleSearchChange = (value: string) => {
    debouncedSearch(value)
  }

  return (
    <div className="w-[300px] xl:w-[320px] 2xl:w-[368px] flex flex-col gap-3 h-[90vh] shrink-0">
      <Search loading={loading} handleSearchChange={handleSearchChange} />
      <QueueContainer>
        {loading && videos.length === 0 ? (
          <p className="text-center text-neutral-500 py-10">Loading queue...</p>
        ) : (
          <>
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
                        optimisticId={optimisticId}
                        onOptimisticClick={() => setOptimisticId(video.videoId)}
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
                    optimisticId={optimisticId}
                    onOptimisticClick={() => setOptimisticId(video.videoId)}
                  />
                )
              })}
            </AnimatePresence>

            <PlaylistQueueFooter
              fetchingMore={fetchingMore}
              hasMore={hasMore}
              videosLength={videos.length}
              loading={loading}
              appliedQuery={appliedQuery}
            />
          </>
        )}
      </QueueContainer>
    </div>
  )
}
