import Image from 'next/image'
import { notFound } from 'next/navigation'

import { getVideoInfo } from '@/app/api/videos'
import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { BLUR_DATA_URLS } from '@/ui/images'

import PlayerContainer from './PlayerContainer'
import ComplaintButton from '../Profile/ComplaintButton'
import FollowButton from '../Profile/FollowButton'
import VideoDescription from '../Videos/VideoDescription'

type Props = {
  v: string
  userId: string | undefined
}

export default async function VideoDetails({ v, userId }: Props) {
  const video = await getVideoInfo(v)

  if (!video.success || !video.data) {
    notFound()
  }

  const { videoId, description, createdAt, views, name, user } = video.data

  return (
    <div className="h-fit min-w-0 w-full">
      <PlayerContainer videoId={v} />

      <div className="flex items-center justify-between mt-3 mb-2">
        <p className="text-neutral-300 text-[20px] font-bold">{name}</p>
        <div className="flex gap-3 items-center">
          <button
            className="flex items-center gap-2 px-4 py-2 bg-neutral-800 rounded-full transition-all duration-300
              hover:bg-neutral-700/40 cursor-pointer active:scale-90 group border border-transparent"
          >
            <PlusIcon className="size-4 text-emerald-500" />
            <p className="text-sm leading-3.5">Add to playlist</p>
          </button>

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
