'use client'

import { useSession } from 'next-auth/react'

import { useStreamQuery } from '@/lib/hooks/api/stream/useStreamQuery'
import { useCheckPlaylistVideoQuery } from '@/lib/hooks/api/video/useVideoInfoQuery'

import MediaDetails from './MediaDetails'
import VideoDetailsSkeleton from './VideoDetailsSkeleton'
import VideoNotFound from '../Videos/VideoNotFound'

type Props = {
  streamId: string | undefined
  playlistId?: string
}

export default function PlaylistMediaStreamContainer({ streamId, playlistId }: Props) {
  const { data: stream, isLoading: isStreamLoading } = useStreamQuery(streamId)
  const { data: isStreamInPlaylist, isLoading: isCheckLoading } = useCheckPlaylistVideoQuery(playlistId, streamId)
  const { data: session } = useSession()

  if (isStreamLoading || isCheckLoading) return <VideoDetailsSkeleton />
  if (!stream || (playlistId && !isStreamInPlaylist)) return <VideoNotFound />

  const { streamerId, description, startedAt, viewers, name, streamerInformation, streamUrl, streamTags, source } =
    stream

  if (!streamerInformation) return <VideoNotFound />

  return (
    <MediaDetails
      playlistId={playlistId}
      resourceId={streamerId}
      description={description ?? ''}
      createdAt={startedAt}
      views={viewers}
      title={name}
      user={{
        userId: streamerInformation.userId,
        avatarUrl: streamerInformation.avatarUrl,
        username: streamerInformation.username,
        isFollowed: streamerInformation.isFollowed,
      }}
      currentUserId={session?.user?.id}
      source={source}
      isOwner={session?.user?.id === streamerInformation.userId}
      mediaUrl={streamUrl}
      isMediaPrivate={false}
      mediaTags={streamTags}
    />
  )
}
