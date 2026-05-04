import { useEffect, useRef } from 'react'

import { useRouter } from 'next/navigation'

import { usePlayerPlayStore } from '@/stores/player.store'

import { PlaylistVideo } from '../actions/playlist.actions'

type UsePlaylistAutoPlayProps = {
  videos: PlaylistVideo[]
  currentV: string | null
  playlistId: string
  setOptimisticId: (id: string | null) => void
  hasNextQueuePage?: boolean
  fetchNextQueuePage?: () => void
  isFetchingQueue?: boolean
}

export function usePlaylistAutoPlay({
  videos,
  currentV,
  playlistId,
  setOptimisticId,
  hasNextQueuePage,
  fetchNextQueuePage,
  isFetchingQueue,
}: UsePlaylistAutoPlayProps) {
  const router = useRouter()
  const endedSignal = usePlayerPlayStore(state => state.endedSignal)

  const latestVideos = useRef(videos)

  useEffect(() => {
    latestVideos.current = videos
  }, [videos])

  useEffect(() => {
    if (!currentV || videos.length === 0) return
    if (!hasNextQueuePage || !fetchNextQueuePage || isFetchingQueue) return

    const currentIndex = videos.findIndex(v => v.videoId === currentV)

    if (currentIndex !== -1 && currentIndex >= videos.length - 3) {
      fetchNextQueuePage()
    }
  }, [currentV, videos.length, hasNextQueuePage, fetchNextQueuePage, isFetchingQueue, videos])

  useEffect(() => {
    if (endedSignal === 0) return

    const currentVideos = latestVideos.current
    const currentIndex = currentVideos.findIndex(v => v.videoId === currentV)

    if (currentIndex !== -1 && currentIndex < currentVideos.length - 1) {
      const nextVideoId = currentVideos[currentIndex + 1].videoId
      const nextVideoSource = currentVideos[currentIndex + 1].source

      setOptimisticId(nextVideoId)

      router.push(`/playlists/${playlistId}?v=${nextVideoId}&source=${nextVideoSource}`, { scroll: false })
    } else {
      const isPlaying = usePlayerPlayStore.getState().playing
      if (isPlaying) {
        setTimeout(() => {
          usePlayerPlayStore.getState().togglePlay()
        }, 50)
      }
    }
  }, [endedSignal])
}
