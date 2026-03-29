import { useEffect, useRef } from 'react'

import { useRouter } from 'next/navigation'

import { usePlayerPlayStore } from '@/stores/player.store'

import { PlaylistVideo } from '../actions/playlist.actions'

type UsePlaylistAutoPlayProps = {
  videos: PlaylistVideo[]
  currentV: string | null
  playlistId: string
  setOptimisticId: (id: string | null) => void
}

export function usePlaylistAutoPlay({ videos, currentV, playlistId, setOptimisticId }: UsePlaylistAutoPlayProps) {
  const router = useRouter()
  const endedSignal = usePlayerPlayStore(state => state.endedSignal)

  const latestVideos = useRef(videos)

  useEffect(() => {
    latestVideos.current = videos
  }, [videos])

  useEffect(() => {
    if (endedSignal === 0) return

    const currentVideos = latestVideos.current
    const currentIndex = currentVideos.findIndex(v => v.videoId === currentV)

    if (currentIndex !== -1 && currentIndex < currentVideos.length - 1) {
      const nextVideoId = currentVideos[currentIndex + 1].videoId

      setOptimisticId(nextVideoId)

      router.push(`/playlists/${playlistId}?v=${nextVideoId}`, { scroll: false })
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
