'use client'

import { useSession } from 'next-auth/react'

import { useSearchVideosQuery } from '@/lib/hooks/api/video/useSearchVideos'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'

import VideoCard, { VideoCardSkeleton } from './VideoCard'
import GridCardsContainer from '../../shared/GridCardsContainer'

type Props = {
  query: string
}

export default function VideosPublicSearchContainer({ query }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useSearchVideosQuery(query)
  const { data: session } = useSession()

  const videos = data?.pages.flatMap(page => page?.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  if (isLoading && videos.length === 0) {
    return (
      <GridCardsContainer>
        {[...new Array(20)].map((_, idx) => (
          <VideoCardSkeleton key={idx} />
        ))}
      </GridCardsContainer>
    )
  }

  if (!isLoading && videos.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center py-20">
        {query.length > 0 && (
          <p className="text-neutral-500 text-center">
            Videos by query <span className="text-neutral-300">{`'${query}'`}</span> not found
          </p>
        )}
        {query.length === 0 && <p className="text-neutral-500 text-center">No videos found</p>}
      </div>
    )
  }

  return (
    <>
      <GridCardsContainer>
        {videos.map((video, index) => {
          const isLast = videos.length === index + 1

          const item = (
            <VideoCard
              key={video.videoId}
              videoId={video.videoId}
              title={video.name}
              previewUrl={video.previewUrl}
              duration={video.duration}
              creator={video.user.username}
              views={video.views}
              createAt={video.createdAt}
              userAvatar={video.user.avatarUrl}
              currentUserId={session?.user.id}
            />
          )

          if (isLast) {
            return (
              <div ref={lastElementRef} key={`last-${video.videoId}`}>
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
