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
  const [{ id: playlistId }, { v }] = await Promise.all([params, searchParams])
  const sessionPromise = auth()
  const playlistInfoPromise = getPlaylistInfo(playlistId)
  const [session, playlistInfoResponse] = await Promise.all([sessionPromise, playlistInfoPromise])

  if (!playlistInfoResponse.success || !playlistInfoResponse.data) {
    notFound()
  }

  const { firstVideo, hiddenVideosCount, playlist } = playlistInfoResponse.data

  if (!v && firstVideo) {
    redirect(`/playlists/${playlistId}?v=${firstVideo.videoId}`)
  }

  return (
    <div className="flex gap-5">
      <Suspense key={v} fallback={<PlayerSkeleton />}>
        <VideoDetails userId={session?.user.id} v={v} playlistId={playlistId} />
      </Suspense>
      <PlaylistQueueContainer
        playlistId={playlistId}
        hiddenVideosCount={hiddenVideosCount}
        playlistName={playlist.name}
      />
    </div>
  )
}
