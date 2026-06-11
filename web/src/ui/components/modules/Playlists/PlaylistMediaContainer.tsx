'use client'

import { useEffect } from 'react'

import { usePlaylistStore } from '@/stores/playlist.store'
import { MediaType } from '@/types/general.types'

import PlaylistMediaStreamContainer from './PlaylistMediaStreamContainer'
import PlaylistMediaVideoContainer from './PlaylistMediaVideoContainer'
import VideoNotFound from '../Videos/VideoNotFound'

type Props = {
  playlistId?: string
  mediaType: MediaType | undefined
  initialMediaId: string | undefined
}

export default function PlaylistMediaContainer({ playlistId, initialMediaId, mediaType }: Props) {
  const setActiveResourceId = usePlaylistStore(state => state.setActiveResourceId)
  const setActiveMediaType = usePlaylistStore(state => state.setActiveMediaType)
  const resetStore = usePlaylistStore(state => state.reset)

  const activeVideoId = usePlaylistStore(state => state.activeResourceId)
  const activeMediaType = usePlaylistStore(state => state.activeMediaType)

  const currentMediaId = activeVideoId || initialMediaId
  const currentMediaType = activeMediaType || mediaType

  useEffect(() => {
    return () => resetStore()
  }, [resetStore])

  useEffect(() => {
    const currentStoreId = usePlaylistStore.getState().activeResourceId
    if (!currentStoreId && initialMediaId && mediaType) {
      setActiveResourceId(initialMediaId)
      setActiveMediaType(mediaType)
    }
  }, [initialMediaId, mediaType, setActiveResourceId, setActiveMediaType])

  if (!currentMediaType || !currentMediaId) return <VideoNotFound />

  const isStream = typeof currentMediaType === 'string' && currentMediaType === 'LiveStream'

  return isStream ? (
    <PlaylistMediaStreamContainer streamId={currentMediaId} playlistId={playlistId} />
  ) : (
    <PlaylistMediaVideoContainer videoId={currentMediaId} playlistId={playlistId} />
  )
}
