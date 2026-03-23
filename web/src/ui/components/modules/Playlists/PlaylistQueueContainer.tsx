'use client'

import { useCallback, useEffect, useState } from 'react'

import { useDebouncedCallback } from 'use-debounce'

import { getPlaylistVideos, PlaylistVideo } from '@/app/api/playlists'
import SearchIcon from '@/assets/icons/ic_search.svg'
import { clog } from '@/lib/utils/utils'
import VideoItem from '@/ui/components/modules/Playlists/VideoItem'

type Props = {
  playlistId: string
}

export default function PlaylistQueueContainer({ playlistId }: Props) {
  const [videos, setVideos] = useState<PlaylistVideo[]>([])
  const [loading, setLoading] = useState(true)
  const [query, setQuery] = useState('')

  const fetchVideos = useCallback(
    async (searchQuery: string) => {
      setLoading(true)
      try {
        const response = await getPlaylistVideos(playlistId, searchQuery, 1, 50)

        if (response.success && response.data) {
          setVideos(response.data.items)
        }
      } catch (error) {
        clog('SEARCH_ERROR', error)
      } finally {
        setLoading(false)
      }
    },
    [playlistId],
  )

  const debouncedSearch = useDebouncedCallback((value: string) => {
    fetchVideos(value)
  }, 300)

  useEffect(() => {
    fetchVideos('')
  }, [fetchVideos])

  const handleSearchChange = (value: string) => {
    setQuery(value)
    debouncedSearch(value)
  }

  return (
    <div className="w-[368px] flex flex-col gap-3 min-h-0">
      {/* Search */}
      <div className="relative">
        <SearchIcon className="absolute top-1/2 -translate-y-1/2 left-4 h-5 w-5 text-neutral-600" aria-hidden="true" />
        <input
          id="search"
          type="text"
          autoComplete="off"
          className="focus:border-emerald-500 focus:ring-emerald-500 ring-[0.3px] ring-transparent block
            placeholder:text-neutral-600 focus:outline-none pl-12 py-3 rounded-2xl font-medium text-sm border
            border-border w-full bg-neutral-900/50 text-neutral-200 transition-all"
          placeholder="Search a video in playlist"
          value={query}
          onChange={e => handleSearchChange(e.target.value)}
        />
        {loading && query && (
          <div className="absolute right-4 top-1/2 -translate-y-1/2">
            <div className="size-4 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
          </div>
        )}
      </div>

      {/* videos list */}
      <div className="flex-1 flex flex-col gap-2 overflow-y-auto pr-1">
        {loading && videos.length === 0 ? (
          <p className="text-center text-neutral-500 py-10">Loading queue...</p>
        ) : videos.length > 0 ? (
          videos.map(video => (
            <VideoItem
              key={video.videoId}
              id={video.videoId}
              title={video.name}
              thumbnail={video.previewUrl}
              isActive={false}
            />
          ))
        ) : (
          <p className="text-center text-neutral-500 py-10">No videos found</p>
        )}
      </div>
    </div>
  )
}
