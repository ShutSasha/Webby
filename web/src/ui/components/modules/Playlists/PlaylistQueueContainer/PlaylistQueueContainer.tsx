'use client'

import { useState } from 'react'

import { useDebouncedCallback } from 'use-debounce'

import { useSearchPlaylistVideosQuery } from '@/lib/hooks/api/playlist/useSearchPlaylistVideos'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import { usePlaylistAutoPlay } from '@/lib/hooks/usePlaylistAutoPlay'
import VideoItem, { VideoItemSkeleton } from '@/ui/components/modules/Playlists/VideoItem'

import PlaylistContainerHeader from './PlaylistContainerHeader'
import PlaylistQueueFooter from './PlaylistQueueFooter'
import QueueContainer from './QueueContainer'
import Search from './Search'

type Props = {
  playlistId: string
  hiddenVideosCount: number
  playlistName: string
  authorId: string
  guestUserId: string | undefined
}

export default function PlaylistQueueContainer({
  playlistId,
  hiddenVideosCount,
  playlistName,
  authorId,
  guestUserId,
}: Props) {
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

  const uiVideos = searchData?.pages.flatMap(page => page?.items || []) || []

  const {
    data: queueData,
    hasNextPage: hasNextQueuePage,
    fetchNextPage: fetchNextQueuePage,
    isFetchingNextPage: isFetchingQueue,
  } = useSearchPlaylistVideosQuery(playlistId, '')

  const queueVideos = queueData?.pages.flatMap(page => page?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  usePlaylistAutoPlay({
    videos: queueVideos,
    hasNextQueuePage,
    fetchNextQueuePage,
    isFetchingQueue,
  })

  return (
    <div className="w-full xl:w-[320px] 2xl:w-[388px] flex flex-col h-[90vh] shrink-0">
      <PlaylistContainerHeader playlistName={playlistName} hiddenVideosCount={hiddenVideosCount} />
      <Search loading={isLoading} handleSearchChange={debouncedSearch} />
      <QueueContainer>
        {isLoading && uiVideos.length === 0 ? (
          <div className="flex flex-col gap-2">
            {Array.from({ length: 8 }).map((_, index) => (
              <VideoItemSkeleton key={`skeleton-${index}`} />
            ))}
          </div>
        ) : (
          <>
            {uiVideos.map((video, index) => {
              const isLast = uiVideos.length === index + 1
              const item = (
                <VideoItem
                  key={video.videoId}
                  id={video.videoId}
                  title={video.name}
                  thumbnail={video.previewUrl}
                  playlistId={playlistId}
                  isOwner={guestUserId === authorId}
                  userId={guestUserId || ''}
                  mediaType={video.mediaType}
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
