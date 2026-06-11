'use client'

import Image from 'next/image'
import Link from 'next/link'

import LockIcon from '@/assets/icons/shared/lock.svg'
import { MediaType } from '@/types/general.types'
import { VideoSource } from '@/types/video.types'
import { BLUR_DATA_URLS } from '@/ui/images'

import SaveToPlaylistButton from './SaveToPlaylistBtn/SaveToPlaylistButton'
import ComplaintButton from '../../shared/ComplaintButton'
import CustomPlayer from '../Player/CustomPlayer'
import FollowButton from '../Profile/FollowButton'
import VideoDescription from '../Videos/VideoDescription'
import VideoNotFound from '../Videos/VideoNotFound'

type Props = {
  resourceId: string
  currentUserId: string | undefined
  playlistId?: string
  user: { userId: string; avatarUrl: string; username: string; isFollowed: boolean }
  source: VideoSource
  title: string
  mediaUrl: string
  isOwner: boolean
  isMediaPrivate: boolean
  description: string
  views: number
  createdAt: string
  mediaTags?: string[] | null
  mediaType?: MediaType
}

export default function MediaDetails({
  resourceId,
  currentUserId,
  user,
  source,
  title,
  mediaUrl,
  isOwner,
  isMediaPrivate,
  description,
  views,
  createdAt,
  mediaTags,
  mediaType,
}: Props) {
  if (isMediaPrivate && !isOwner) return <VideoNotFound />

  return (
    <div className="h-fit min-w-0 w-full">
      <CustomPlayer videoUrl={mediaUrl} videoId={resourceId} />

      <div className="flex items-start justify-between mt-3 mb-2">
        <div className="flex flex-col gap-0.5">
          <p className="text-neutral-300 text-[20px] font-bold">{title}</p>
          {isMediaPrivate && (
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
            {currentUserId && (
              <SaveToPlaylistButton userId={currentUserId} videoId={resourceId} mediaType={mediaType} />
            )}
            <ComplaintButton authorId={currentUserId} targetId={resourceId.slice(3)} targetType="Video" />
          </div>
        </div>
      )}
      {currentUserId && source !== 'Webby' && (
        <div className="flex flex-row-reverse">
          <SaveToPlaylistButton userId={currentUserId} videoId={resourceId} mediaType={mediaType} />
        </div>
      )}
      <VideoDescription text={description ?? ''} views={views} date={createdAt} videoTags={mediaTags} />
    </div>
  )
}
