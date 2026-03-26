import Image from 'next/image'
import Link from 'next/link'

import { getVideoInfo } from '@/app/api/videos'
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
}

export default async function VideoDetails({ v, userId }: Props) {
  if (!v) {
    return <VideoNotFound />
  }

  const video = await getVideoInfo(v)

  if (!video.success || !video.data) {
    return <VideoNotFound />
  }

  const { videoId, description, createdAt, views, name, user, videoUrl } = video.data
  const isOwner = userId === user.userId

  return (
    <div className="h-fit min-w-0 w-full">
      <CustomPlayer videoUrl={videoUrl} />

      <div className="flex items-center justify-between mt-3 mb-2">
        <p className="text-neutral-300 text-[20px] font-bold">{name}</p>
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
