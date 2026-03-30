'use client'

import { useEffect, useState } from 'react'

import { AnimatePresence } from 'framer-motion'
import { useSearchParams } from 'next/navigation'
import { useDebouncedCallback } from 'use-debounce'

import { useSearchPlaylistVideosQuery } from '@/lib/hooks/api/playlist/useSearchPlaylistVideos'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
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

  const [query, setQuery] = useState('')

  const debouncedSearch = useDebouncedCallback((value: string) => {
    setQuery(value)
  }, 400)

  const {
    data: searchData,
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  } = useSearchPlaylistVideosQuery(playlistId, query)

  const uiVideos = searchData?.pages.flatMap(page => page?.data?.items || []) || []

  const {
    data: queueData,
    hasNextPage: hasNextQueuePage,
    fetchNextPage: fetchNextQueuePage,
    isFetchingNextPage: isFetchingQueue,
  } = useSearchPlaylistVideosQuery(playlistId, '')

  const queueVideos = queueData?.pages.flatMap(page => page?.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  usePlaylistAutoPlay({
    videos: queueVideos,
    currentV,
    playlistId,
    setOptimisticId,
    hasNextQueuePage,
    fetchNextQueuePage,
    isFetchingQueue,
  })

  useEffect(() => {
    setOptimisticId(null)
  }, [currentV])

  return (
    <div className="w-[300px] xl:w-[320px] 2xl:w-[368px] flex flex-col h-[90vh] shrink-0">
      <PlaylistContainerHeader playlistName={playlistName} hiddenVideosCount={hiddenVideosCount} />
      <Search loading={isLoading} handleSearchChange={debouncedSearch} />
      <QueueContainer>
        {isLoading && uiVideos.length === 0 ? (
          <p className="text-center text-neutral-500 py-10">Loading queue...</p>
        ) : (
          <>
            <AnimatePresence mode="popLayout">
              {uiVideos.map((video, index) => {
                const isLast = uiVideos.length === index + 1
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
              fetchingMore={isFetchingNextPage}
              hasMore={hasNextPage}
              videosLength={uiVideos.length}
              loading={isLoading}
              appliedQuery={query}
            />
          </>
        )}
      </QueueContainer>
    </div>
  )
}
