import Image from 'next/image'
import { notFound } from 'next/navigation'

import { getPlaylistInfo } from '@/app/api/playlists'
import { getVideoInfo } from '@/app/api/videos'
import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { clog } from '@/lib/utils/utils'
import MainLayout from '@/ui/components/MainLayout'
import CustomPlayer from '@/ui/components/modules/Player/CustomPlayer'
import PlaylistQueueContainer from '@/ui/components/modules/Playlists/PlaylistQueueContainer'
import VideoItem from '@/ui/components/modules/Playlists/VideoItem'
import ComplaintButton from '@/ui/components/modules/Profile/ComplaintButton'
import FollowButton from '@/ui/components/modules/Profile/FollowButton'
import PlaylistItem from '@/ui/components/modules/Rooms/InteractionBlock/Playlists/PlaylistItem'
import AsideVideoCard from '@/ui/components/modules/Videos/AsideVideoCard'
import VideoDescription from '@/ui/components/modules/Videos/VideoDescription'
import { BLUR_DATA_URLS } from '@/ui/images'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
  searchParams: Promise<{ v?: string; query?: string }>
}

export default async function PlaylistPage({ params, searchParams }: Props) {
  const { id: playlistId } = await params
  const { v: videoIdFromUrl } = await searchParams
  const session = await auth()
  const playlistInfoResponse = await getPlaylistInfo(playlistId)
  const video = await getVideoInfo('22f510be-a044-47b2-9710-7660caba5192')

  if (!playlistInfoResponse.success || !playlistInfoResponse.data || !video.success || !video.data) {
    notFound()
  }

  clog('playlistData', playlistInfoResponse)

  return (
    <div className="flex gap-5">
      <div className="flex-1 min-w-0">
        <CustomPlayer videoUrl={video.data.videoUrl} />

        <div className="flex items-center justify-between mt-3 mb-2">
          <p className="text-neutral-300 text-[20px] font-bold">123123123</p>
          <div className="flex gap-3 items-center">
            <button
              className="flex items-center gap-2 px-4 py-2 bg-neutral-800 rounded-full transition-all duration-300
                hover:bg-neutral-700/40 cursor-pointer active:scale-90 group border border-transparent"
            >
              <PlusIcon className="size-4 text-emerald-500" />
              <p className="text-sm leading-3.5">Add to playlist</p>
            </button>

            <ComplaintButton authorId={session?.user.id} targetId={video.data.videoId} targetType="Video" />
          </div>
        </div>

        <div className="flex items-center gap-3">
          <Image
            src={'https://i.ibb.co/60Ns8j8r/cd4af4dd04fcfba0a358cfdee5c039f7.jpg'}
            className="size-9 object-cover rounded-full"
            alt=""
            width={50}
            height={50}
            loading="lazy"
            placeholder="blur"
            blurDataURL={BLUR_DATA_URLS['neutral900']}
          />
          <p className="text-[16px] font-medium">{'123'}</p>
          <FollowButton targetUserId={'123'} initialIsFollowing={false} />
        </div>

        <VideoDescription text={'videoData.description'} views={123} date={'videoData.createdAt'} />
      </div>

      <div className="w-[368px] flex flex-col gap-3">
        {[...new Array(20)].map((_, index) => (
          <VideoItem
            key={index}
            id={'1'}
            title="title123123"
            thumbnail="https://i.ibb.co/PGL4ymBS/thumb-1920-415519.jpg"
            isActive={index === 2 ? true : false}
          />
        ))}
      </div>
      {/* <PlaylistQueueContainer playlistId={playlistId} /> */}
    </div>
  )
}
