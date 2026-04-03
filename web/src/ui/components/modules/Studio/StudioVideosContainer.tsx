'use client'

import { notFound } from 'next/navigation'
import { useSession } from 'next-auth/react'

import { useGetUserVideosQuery } from '@/lib/hooks/api/video/useGetUserVideosQuery'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import { StudioVideoRow } from '@/ui/components/modules/Studio/StudioVideoRow'

export default function StudioVideosContainer() {
  const { data: session } = useSession()
  const userId = session?.user?.id

  if (!userId) {
    notFound()
  }

  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useGetUserVideosQuery(userId)

  const videos = data?.pages.flatMap(page => page?.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  if (isLoading && videos.length === 0) {
    return (
      <div className="flex flex-col w-full animate-pulse">
        {[...new Array(5)].map((_, i) => (
          <div key={i} className="h-24 w-full border-b border-neutral-800 bg-neutral-800/20" />
        ))}
      </div>
    )
  }

  if (!isLoading && videos.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center py-20">
        <p className="text-neutral-500 text-center">No videos found. Upload your first video!</p>
      </div>
    )
  }

  return (
    <div className="flex flex-col pb-10">
      {videos.map((video, index) => {
        const isLast = videos.length === index + 1
        const item = <StudioVideoRow key={video.videoId} video={video} />

        if (isLast) {
          return (
            <div ref={lastElementRef} key={`last-${video.videoId}`}>
              {item}
            </div>
          )
        }

        return item
      })}

      {isFetchingNextPage && (
        <div className="w-full flex justify-center py-6">
          <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </div>
  )
}
