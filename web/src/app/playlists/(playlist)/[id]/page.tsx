import { Suspense } from 'react'

import { notFound, redirect } from 'next/navigation'

import { getPlaylistInfo } from '@/app/api/playlists'
import PlayerSkeleton from '@/ui/components/modules/Playlists/PlayerSkeleton'
import PlaylistQueueContainer from '@/ui/components/modules/Playlists/PlaylistQueueContainer/PlaylistQueueContainer'
import VideoDetails from '@/ui/components/modules/Playlists/VideoDetails'
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

  if (!v && playlistInfoResponse.data.firstVideo) {
    redirect(`/playlists/${playlistId}?v=${playlistInfoResponse.data.firstVideo.videoId}`)
  }

  return (
    <div className="flex gap-5">
      <Suspense key={v} fallback={<PlayerSkeleton />}>
        <VideoDetails userId={session?.user.id} v={v} />
      </Suspense>
      <PlaylistQueueContainer playlistId={playlistId} />
    </div>
  )
}
