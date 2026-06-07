'use client'

import { useEffect } from 'react'

import { usePlaylistStore } from '@/stores/playlist.store'

import VideoDetails from '../Playlists/VideoDetails'

type Props = {
  videoId: string
  currentUserId: string | undefined
}

export default function StandaloneVideoContainer({ videoId, currentUserId }: Props) {
  const setActiveVideoId = usePlaylistStore(state => state.setActiveVideoId)
  const setActiveMediaType = usePlaylistStore(state => state.setActiveMediaType)

  useEffect(() => {
    setActiveVideoId(null)
    setActiveMediaType(null)
  }, [setActiveVideoId, setActiveMediaType])

  return <VideoDetails videoId={videoId} currentUserId={currentUserId} />
}
