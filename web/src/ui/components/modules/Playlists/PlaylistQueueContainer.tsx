'use client'

import { useMemo } from 'react'

import Link from 'next/link'
import { useSearchParams } from 'next/navigation'

import VideoItem from '@/ui/components/modules/Playlists/VideoItem'

import Search from '../../Search'

type Props = {
  playlistId: string
}

export default function PlaylistQueueContainer({ playlistId }: Props) {
  const searchParams = useSearchParams()
  // const query = searchParams.get('query')?.toLowerCase() || ''

  return (
    <div className="w-[368px] flex flex-col gap-3 min-h-0">
      <Search
        placeholder="Search video in queue"
        inputClassName="pl-10 py-2 rounded-xl bg-neutral-900 border-none text-sm"
        iconClassName="left-3 h-4 w-4 text-neutral-500"
      />

      <div className="flex flex-col gap-2 overflow-y-auto pr-1 custom-scrollbar">
        {/* {filteredVideos.map(video => (
          <Link key={video.videoId} href={`/playlists/${playlistId}?v=${video.videoId}`} scroll={false}>
            <VideoItem
              id={video.videoId}
              title={video.name}
              thumbnail={video.previewUrl}
              isActive={video.videoId === activeVideoId}
            />
          </Link>
        ))} */}
      </div>
    </div>
  )
}
