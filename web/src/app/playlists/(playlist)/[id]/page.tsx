import Image from 'next/image'
import { notFound, redirect } from 'next/navigation'

import { getPlaylistInfo } from '@/app/api/playlists'
import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import PlayerContainer from '@/ui/components/modules/Playlists/PlayerContainer'
import PlaylistQueueContainer from '@/ui/components/modules/Playlists/PlaylistQueueContainer'
import ComplaintButton from '@/ui/components/modules/Profile/ComplaintButton'
import FollowButton from '@/ui/components/modules/Profile/FollowButton'
import VideoDescription from '@/ui/components/modules/Videos/VideoDescription'
import { BLUR_DATA_URLS } from '@/ui/images'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
  searchParams: Promise<{ v?: string }>
}

export default async function PlaylistPage({ params, searchParams }: Props) {
  const { id: playlistId } = await params
  const session = await auth()
  const playlistInfoResponse = await getPlaylistInfo(playlistId)
  const { v } = await searchParams

  if (!playlistInfoResponse.success || !playlistInfoResponse.data) {
    notFound()
  }

  const { videoId, description, createdAt, views, name, user } = playlistInfoResponse.data.firstVideo

  if (!v) {
    redirect(`/playlists/${playlistId}?v=${videoId}`)
  }

  return (
    <div className="flex gap-5">
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

            <ComplaintButton authorId={session?.user.id} targetId={videoId} targetType="Video" />
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

      <PlaylistQueueContainer playlistId={playlistId} />
    </div>
  )
}
