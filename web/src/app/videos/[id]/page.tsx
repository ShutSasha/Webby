import Image from 'next/image'
import { notFound } from 'next/navigation'

import { getVideoInfo } from '@/app/api/videos'
import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { clog } from '@/lib/utils/utils'
import MainLayout from '@/ui/components/MainLayout'
import CustomPlayer from '@/ui/components/modules/Player/CustomPlayer'
import ComplaintButton from '@/ui/components/modules/Profile/ComplaintButton'
import FollowButton from '@/ui/components/modules/Profile/FollowButton'
import AsideVideoCard from '@/ui/components/modules/Videos/AsideVideoCard'
import VideoDescription from '@/ui/components/modules/Videos/VideoDescription'
import { BLUR_DATA_URLS } from '@/ui/images'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
}

export default async function VideoPage({ params }: Props) {
  const { id } = await params
  const session = await auth()
  const videoInfoResponse = await getVideoInfo(id)

  clog('video data', videoInfoResponse)

  if (!videoInfoResponse.success || !videoInfoResponse.data) {
    notFound()
  }

  const { videoId, videoUrl, name, user, description, views, createdAt } = videoInfoResponse.data

  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <div className="flex gap-5">
          <div className="flex-1 min-w-0">
            <CustomPlayer videoUrl={videoUrl} />

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
                {/* TODO: change it to right API for video complaint */}
                <ComplaintButton authorId={session?.user.id} targetId={videoId} />
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

            <VideoDescription text={description} views={views} date={createdAt} />
          </div>

          <div className="w-[418px] flex flex-col gap-3">
            {[...new Array(20)].map((_, index) => (
              <AsideVideoCard
                key={index}
                video={{
                  videoId: 'e68b32e3-473f-40d8-ad40-b5a64a9436d8',
                  previewUrl:
                    'https://webby-watch-platform-bucket.s3.eu-north-1.amazonaws.com/previews/e68b32e3-473f-40d8-ad40-b5a64a9436d8/14:26:38videoframe_3294840.png',
                }}
              />
            ))}
          </div>
        </div>
      </div>
    </MainLayout>
  )
}
