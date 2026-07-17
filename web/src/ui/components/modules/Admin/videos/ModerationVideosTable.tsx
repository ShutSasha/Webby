'use client'

import { useMemo } from 'react'

import { useSearchModerationVideosQuery } from '@/lib/hooks/api/admin/useSearchModerationVideos'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'

import ModerationVideosTableSkeleton from './ModerationVideosTableSkeleton'
import VideoTableRow from './VideoTableRow'

type Props = {
  searchQuery: string
}

export default function ModerationVideosTable({ searchQuery }: Props) {
  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useSearchModerationVideosQuery(
    searchQuery,
    20,
  )

  const videos = useMemo(() => {
    return data?.pages.flatMap(page => page?.items || []) || []
  }, [data])

  const lastElementRef = useInfiniteScroll({
    isLoading,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  })

  if (isLoading && videos.length === 0) {
    return <ModerationVideosTableSkeleton />
  }

  if (!isLoading && videos.length === 0) {
    return (
      <div
        className="flex flex-col items-center justify-center py-20 bg-surface border border-border rounded-2xl
          text-foreground-faint"
      >
        <p>No videos found matching &quot;{searchQuery}&quot;</p>
      </div>
    )
  }

  return (
    <div className="bg-surface border border-border rounded-2xl overflow-hidden flex flex-col">
      <div className="overflow-x-auto custom-scrollbar">
        <table className="w-full text-left border-collapse min-w-[800px]">
          <thead>
            <tr className="border-b border-border bg-surface/50 text-foreground-faint text-sm">
              <th className="py-4 px-6 font-medium">Video</th>
              <th className="py-4 px-6 font-medium">Status</th>
              <th className="py-4 px-6 font-medium">Created At</th>
              <th className="py-4 px-6 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-background/60">
            {videos.map((video, index) => {
              const isLast = index === videos.length - 1
              return <VideoTableRow key={video.videoId} video={video} lastElementRef={isLast ? lastElementRef : null} />
            })}
          </tbody>
        </table>
      </div>

      {isFetchingNextPage && (
        <div className="flex justify-center py-4 border-t border-background/60">
          <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}
    </div>
  )
}
