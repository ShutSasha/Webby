'use client'

import { useCallback, useEffect, useState } from 'react'

import { AnimatePresence } from 'framer-motion'
import { useSearchParams } from 'next/navigation'

import { getPlaylistVideos, PlaylistVideo } from '@/app/api/playlists'
import { useInfiniteSearch } from '@/lib/hooks/useInfiniteSearch'
import { usePlaylistAutoPlay } from '@/lib/hooks/usePlaylistAutoPlay'
import VideoItem from '@/ui/components/modules/Playlists/VideoItem'

import PlaylistContainerHeader from './PlaylistContainerHeader'
import PlaylistQueueFooter from './PlaylistQueueFooter'
import QueueContainer from './QueueContainer'
import Search from './Search'

type Props = {
  playlistId: string
  hiddenVideosCount: number
  playlistName: string
}

export default function PlaylistQueueContainer({ playlistId, hiddenVideosCount, playlistName }: Props) {
  const searchParams = useSearchParams()
  const currentV = searchParams.get('v')
  const [optimisticId, setOptimisticId] = useState<string | null>(null)

  const fetchVideosFn = useCallback(
    (query: string, page: number, pageSize: number) => {
      return getPlaylistVideos(playlistId, query, page, pageSize)
    },
    [playlistId],
  )

  const {
    items: videos,
    loading,
    fetchingMore,
    hasMore,
    appliedQuery,
    lastElementRef,
    handleSearchChange,
  } = useInfiniteSearch<PlaylistVideo>({
    fetchFn: fetchVideosFn,
    pageSize: 20,
  })

  usePlaylistAutoPlay({
    videos,
    currentV,
    playlistId,
    setOptimisticId,
  })

  useEffect(() => {
    setOptimisticId(null)
  }, [currentV])

  return (
    <div className="w-[300px] xl:w-[320px] 2xl:w-[368px] flex flex-col h-[90vh] shrink-0">
      <PlaylistContainerHeader playlistName={playlistName} hiddenVideosCount={hiddenVideosCount} />
      <Search loading={loading} handleSearchChange={handleSearchChange} />
      <QueueContainer>
        {loading && videos.length === 0 ? (
          <p className="text-center text-neutral-500 py-10">Loading queue...</p>
        ) : (
          <>
            <AnimatePresence mode="popLayout">
              {videos.map((video, index) => {
                const isLast = videos.length === index + 1
                const item = (
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

                if (isLast) {
                  return (
                    <div ref={lastElementRef} key={video.videoId}>
                      {item}
                    </div>
                  )
                }
                return item
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
