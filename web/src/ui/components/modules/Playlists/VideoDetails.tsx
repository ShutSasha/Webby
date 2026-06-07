'use client'

import { useEffect } from 'react'

import Image from 'next/image'
import Link from 'next/link'

import LockIcon from '@/assets/icons/shared/lock.svg'
import { useCheckPlaylistVideoQuery, useVideoInfoQuery } from '@/lib/hooks/api/video/useVideoInfoQuery'
import { usePlaylistStore } from '@/stores/playlist.store'
import { BLUR_DATA_URLS } from '@/ui/images'

import SaveToPlaylistButton from './SaveToPlaylistBtn/SaveToPlaylistButton'
import VideoDetailsSkeleton from './VideoDetailsSkeleton'
import ComplaintButton from '../../shared/ComplaintButton'
import CustomPlayer from '../Player/CustomPlayer'
import FollowButton from '../Profile/FollowButton'
import VideoDescription from '../Videos/VideoDescription'
import VideoNotFound from '../Videos/VideoNotFound'

type Props = {
  videoId: string
  currentUserId: string | undefined
  playlistId?: string
}

export default function VideoDetails({ videoId, currentUserId, playlistId }: Props) {
  const setActiveMediaType = usePlaylistStore(state => state.setActiveMediaType)

  const { data: video, isLoading: isVideoLoading } = useVideoInfoQuery(videoId)
  const { data: isVideoInPlaylist, isLoading: isCheckLoading } = useCheckPlaylistVideoQuery(playlistId, videoId)

  useEffect(() => {
    if (video?.mediaType) {
      setActiveMediaType(video.mediaType)
    }
  }, [video?.mediaType, setActiveMediaType])

  if (!videoId) return <VideoNotFound />
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
  const isOwner = currentUserId === user.userId
  if (isPrivate && !isOwner) return <VideoNotFound />

  return (
    <div className="h-fit min-w-0 w-full">
      <CustomPlayer videoUrl={videoUrl} videoId={fetchedVideoId} />

      <div className="flex items-start justify-between mt-3 mb-2">
        <div className="flex flex-col gap-0.5">
          <p className="text-neutral-300 text-[20px] font-bold">{name}</p>
          {isPrivate && (
            <div className="bg-black/40 py-0.5 px-2 rounded-sm flex items-center gap-1 w-fit">
              <LockIcon className="size-3 text-neutral-500" />
              <p className="text-neutral-500 text-[12px]">Private</p>
            </div>
          )}
        </div>
      </div>
      {source === 'Webby' && (
        <div className="flex items-center justify-between flex-wrap gap-2">
          <div className="flex items-center gap-3">
            <Link href={`/profile/${user.userId}`} className="flex items-center gap-3">
              <Image
                src={user.avatarUrl}
                className="size-9 object-cover rounded-full"
                alt=""
                width={50}
                height={50}
                loading="lazy"
                placeholder="blur"
                blurDataURL={BLUR_DATA_URLS['neutral900']}
              />
              <p className="text-[16px] font-medium">{user.username}</p>
            </Link>
            {!isOwner && (
              <FollowButton
                targetUserId={user.userId}
                initialIsFollowing={user.isFollowed}
                currentUserId={currentUserId}
              />
            )}
          </div>
          <div className="flex gap-3 items-center">
            {currentUserId && <SaveToPlaylistButton userId={currentUserId} videoId={fetchedVideoId} />}
            <ComplaintButton authorId={currentUserId} targetId={fetchedVideoId.slice(3)} targetType="Video" />
          </div>
        </div>
      )}
      {currentUserId && source !== 'Webby' && (
        <div className="flex flex-row-reverse">
          <SaveToPlaylistButton userId={currentUserId} videoId={fetchedVideoId} />
        </div>
      )}
      <VideoDescription text={description ?? ''} views={views} date={createdAt} videoTags={videoTags} />
    </div>
  )
}
