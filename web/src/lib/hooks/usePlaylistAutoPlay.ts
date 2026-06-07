import { useEffect, useRef } from 'react'

import { usePlayerPlayStore } from '@/stores/player.store'
import { usePlaylistStore } from '@/stores/playlist.store'

import { PlaylistVideo } from '../actions/playlist.actions'

type UsePlaylistAutoPlayProps = {
  videos: PlaylistVideo[]
  hasNextQueuePage?: boolean
  fetchNextQueuePage?: () => void
  isFetchingQueue?: boolean
}

export function usePlaylistAutoPlay({
  videos,
  hasNextQueuePage,
  fetchNextQueuePage,
  isFetchingQueue,
}: UsePlaylistAutoPlayProps) {
  const endedSignal = usePlayerPlayStore(state => state.endedSignal)

  const activeVideoId = usePlaylistStore(state => state.activeVideoId)

  const latestVideos = useRef(videos)

  useEffect(() => {
    latestVideos.current = videos
  }, [videos])

  useEffect(() => {
    if (!activeVideoId || videos.length === 0) return
    if (!hasNextQueuePage || !fetchNextQueuePage || isFetchingQueue) return

    const currentIndex = videos.findIndex(v => v.videoId === activeVideoId)

    if (currentIndex !== -1 && currentIndex >= videos.length - 3) {
      fetchNextQueuePage()
    }
  }, [activeVideoId, videos.length, hasNextQueuePage, fetchNextQueuePage, isFetchingQueue, videos])

  useEffect(() => {
    if (endedSignal === 0) return

    const currentVideos = latestVideos.current

    const currentActiveId = usePlaylistStore.getState().activeVideoId
    const currentIndex = currentVideos.findIndex(v => v.videoId === currentActiveId)

    if (currentIndex !== -1 && currentIndex < currentVideos.length - 1) {
      const nextVideoId = currentVideos[currentIndex + 1].videoId

      usePlaylistStore.getState().setActiveVideoId(nextVideoId)
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
