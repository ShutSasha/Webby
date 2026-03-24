import Image from 'next/image'

import { getVideoInfo } from '@/app/api/videos'
import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { BLUR_DATA_URLS } from '@/ui/images'

import CustomPlayer from '../Player/CustomPlayer'
import ComplaintButton from '../Profile/ComplaintButton'
import FollowButton from '../Profile/FollowButton'
import ProfileActionButton from '../Profile/ProfileActionButton'
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

  return (
    <div className="h-fit min-w-0 w-full">
      <CustomPlayer videoUrl={videoUrl} />

      <div className="flex items-center justify-between mt-3 mb-2">
        <p className="text-neutral-300 text-[20px] font-bold">{name}</p>
        <div className="flex gap-3 items-center">
          <ProfileActionButton label="Add to playlist" btnClassName="self-end">
            <PlusIcon className="size-4 text-emerald-500" />
          </ProfileActionButton>

          <ComplaintButton authorId={userId} targetId={videoId} targetType="Video" />
        </div>
      </div>

      <div className="flex items-center gap-3">
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
        <FollowButton targetUserId={user.userId} initialIsFollowing={user.isFollowed} />
      </div>

      <VideoDescription text={description ?? ''} views={views} date={createdAt} />
    </div>
  )
}
