'use client'

import { useEffect } from 'react'

import { usePlaylistStore } from '@/stores/playlist.store'

import VideoDetails from './VideoDetails'

type Props = {
  playlistId: string
  initialVideoId?: string
  currentUserId: string | undefined
}

export default function PlaylistVideoContainer({ playlistId, initialVideoId, currentUserId }: Props) {
  const initPlaylist = usePlaylistStore(state => state.initPlaylist)

  useEffect(() => {
    initPlaylist(initialVideoId ?? '')
  }, [playlistId, initialVideoId, initPlaylist])

  const activeVideoId = usePlaylistStore(state => state.activeVideoId)

  return (
    <VideoDetails
      videoId={activeVideoId || initialVideoId || ''}
      playlistId={playlistId}
      currentUserId={currentUserId}
    />
  )
}
