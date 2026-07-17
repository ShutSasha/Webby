'use client'
import { useSession } from 'next-auth/react'

import { useCheckPlaylistVideoQuery, useVideoInfoQuery } from '@/lib/hooks/api/video/useVideoInfoQuery'

import MediaDetails from './MediaDetails'
import VideoDetailsSkeleton from './VideoDetailsSkeleton'
import VideoNotFound from '../Videos/VideoNotFound'

type Props = {
  videoId: string | undefined
  playlistId?: string
}

export default function PlaylistMediaVideoContainer({ videoId, playlistId }: Props) {
  const { data: video, isLoading: isVideoLoading } = useVideoInfoQuery(videoId)
  const { data: isVideoInPlaylist, isLoading: isCheckLoading } = useCheckPlaylistVideoQuery(playlistId, videoId)
  const { data: session } = useSession()

  if (isVideoLoading || isCheckLoading) return <VideoDetailsSkeleton />
  if (!video || (playlistId && !isVideoInPlaylist)) return <VideoNotFound />

  const {
    videoId: fetchedVideoId,
    description,
    createdAt,
    views,
    name,
    user,
    videoUrl,
    isPrivate,
    videoTags,
    source,
  } = video

  if (!user) return <VideoNotFound />

  return (
    <MediaDetails
      playlistId={playlistId}
      resourceId={fetchedVideoId}
      description={description ?? ''}
      createdAt={createdAt}
      views={views}
      title={name}
      user={{ userId: user.userId, avatarUrl: user.avatarUrl, username: user.username, isFollowed: user.isFollowed }}
      currentUserId={session?.user?.id}
      source={source}
      isOwner={session?.user?.id === user.userId}
      mediaUrl={videoUrl}
      isMediaPrivate={isPrivate}
      mediaTags={videoTags}
    />
  )
}
