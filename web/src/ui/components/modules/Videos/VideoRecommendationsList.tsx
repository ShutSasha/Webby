'use client'

import { useVideoRecommendationsQuery } from '@/lib/hooks/api/video/useVideoRecommendations'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'

import AsideVideoCard, { AsideVideoCardSkeleton } from './AsideVideoCard'

type Props = {
  videoId: string
}

export default function VideoRecommendationsList({ videoId }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useVideoRecommendationsQuery(videoId)

  const rawVideos = data?.pages.flatMap(p => p.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({ isLoading, isFetchingNextPage, hasNextPage, fetchNextPage })

  if (isLoading && rawVideos.length === 0) {
    return (
      <>
        {[...Array(15)].map((_, i) => (
          <AsideVideoCardSkeleton key={i} />
        ))}
      </>
    )
  }

  if (rawVideos.length === 0) return null

  return (
    <>
      {rawVideos.map((video, index) => {
        const isLast = index === rawVideos.length - 1
        const card = <AsideVideoCard key={video.videoId} video={video} />

        if (isLast) {
          return (
            <div key={`last-${video.videoId}`} ref={lastElementRef}>
              {card}
            </div>
          )
        }
        return card
      })}

      {isFetchingNextPage && (
        <div className="flex justify-center py-4 w-full xl:col-span-1 sm:col-span-2 lg:col-span-3">
          <div className="size-5 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </>
  )
}
