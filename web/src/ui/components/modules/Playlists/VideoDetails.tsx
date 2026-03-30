import Image from 'next/image'
import Link from 'next/link'

import LockIcon from '@/assets/icons/shared/lock.svg'
import { checkVideoInPlaylist } from '@/lib/actions/playlist.actions'
import { getVideoInfo } from '@/lib/actions/video.actions'
import { BLUR_DATA_URLS } from '@/ui/images'

import SaveToPlaylistButton from './SaveToPlaylistBtn/SaveToPlaylistButton'
import CustomPlayer from '../Player/CustomPlayer'
import ComplaintButton from '../Profile/ComplaintButton'
import FollowButton from '../Profile/FollowButton'
import VideoDescription from '../Videos/VideoDescription'
import VideoNotFound from '../Videos/VideoNotFound'

type Props = {
  v: string | undefined
  userId: string | undefined
  playlistId?: string
}

export default async function VideoDetails({ v, userId, playlistId }: Props) {
  if (!v) {
    return <VideoNotFound />
  }

  const [video, playlistCheck] = await Promise.all([
    getVideoInfo(v),
    playlistId ? checkVideoInPlaylist(playlistId, v) : Promise.resolve({ success: true, data: true }),
  ])

  if (!video.success || !video.data) {
    return <VideoNotFound />
  }

  if (playlistId && !playlistCheck.data) {
    return <VideoNotFound />
  }

  const { videoId, description, createdAt, views, name, user, videoUrl, isPrivate } = video.data
  const isOwner = userId === user.userId

  if (isPrivate && !isOwner) {
    return <VideoNotFound />
  }

  return (
    <div className="h-fit min-w-0 w-full">
      <CustomPlayer videoUrl={videoUrl} />

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
        <div className="flex gap-3 items-center">
          {userId && <SaveToPlaylistButton userId={userId} videoId={videoId} />}

          <ComplaintButton authorId={userId} targetId={videoId} targetType="Video" />
        </div>
      </div>

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
        {!isOwner && <FollowButton targetUserId={user.userId} initialIsFollowing={user.isFollowed} />}
      </div>

      <VideoDescription text={description ?? ''} views={views} date={createdAt} />
    </div>
  )
}
