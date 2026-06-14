import { useQuery } from '@tanstack/react-query'

import { checkVideoInPlaylist } from '@/lib/actions/playlist.actions'
import { getVideoInfo } from '@/lib/actions/video.actions'

export const useVideoInfoQuery = (videoId: string | undefined) => {
  return useQuery({
    queryKey: ['video-info', videoId],
    queryFn: async () => {
      if (!videoId) throw new Error('Video ID is required')
      const res = await getVideoInfo(videoId)
      if (!res.success || !res.data) throw new Error(res.message)
      return res.data
    },
    enabled: !!videoId,
    staleTime: 5 * 60 * 1000,
  })
}

export const useCheckPlaylistVideoQuery = (playlistId: string | undefined, videoId: string | undefined) => {
  return useQuery({
    queryKey: ['check-playlist-video', playlistId, videoId],
    queryFn: async () => {
      if (!playlistId || !videoId) return false
      const res = await checkVideoInPlaylist(playlistId, videoId)
      return res.data === true
    },
    enabled: !!playlistId && !!videoId,
    staleTime: 5 * 60 * 1000,
  })
}
