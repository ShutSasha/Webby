'use client'

import Link from 'next/link'

import { DEFAULT_VIDEO_THUMBNAIL } from '@/lib/constants/url.constasts'
import { useGetUserVideosQuery } from '@/lib/hooks/api/video/useGetUserVideosQuery'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import { formatTimeAgo } from '@/lib/utils/date.utils'
import { formatViews } from '@/lib/utils/video.utils'

import ImageBackground from './ImageBackground'

type Props = {
  userId: string
}

export default function UserVideos({ userId }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useGetUserVideosQuery(userId)

  const videos = data?.pages.flatMap(page => page?.data?.items || []) || []

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  if (isLoading && videos.length === 0) {
    return <UserVideosSkeleton />
  }

  if (!isLoading && videos.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center py-20">
        <p className="text-neutral-500 text-center">No videos found</p>
      </div>
    )
  }

  return (
    <>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
        {videos.map((video, index) => {
          const isLast = videos.length === index + 1

          return (
            <Link
              href={`/videos/${video.videoId}`}
              key={video.videoId}
              ref={isLast ? lastElementRef : null}
              className="group cursor-pointer"
            >
              <ImageBackground src={video.previewUrl || DEFAULT_VIDEO_THUMBNAIL} />
              <p className="text-sm font-medium wrap-break-word line-clamp-1">{video.name}</p>
              <div className="flex gap-2 text-[12px] text-neutral-500">
                {formatViews(video.views)} • {formatTimeAgo(video.createdAt)}
              </div>
            </Link>
          )
        })}
      </div>

      {isFetchingNextPage && (
        <div className="w-full flex justify-center py-8">
          <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </>
  )
}

export function UserVideosSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
      {Array.from({ length: 12 }).map((_, index) => (
        <div key={index}>
          <div className="w-full aspect-video rounded-2xl mb-1 bg-neutral-800 animate-pulse" />

          {/* Title */}
          <div className="h-4 w-3/4 bg-neutral-800 rounded-md mb-1.5 animate-pulse mt-1" />

          {/* Meta info (Views and Date) */}
          <div className="flex gap-2 items-center">
            <div className="h-3 w-16 bg-neutral-800 rounded-md animate-pulse" />
            <div className="h-3 w-20 bg-neutral-800 rounded-md animate-pulse" />
          </div>
        </div>
      ))}
    </div>
  )
}
