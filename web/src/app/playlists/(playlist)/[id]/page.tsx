import { Suspense } from 'react'

import { redirect } from 'next/navigation'

import { getPlaylistInfo } from '@/lib/actions/playlist.actions'
import PlaylistQueueContainer from '@/ui/components/modules/Playlists/PlaylistQueueContainer/PlaylistQueueContainer'
import VideoDetails from '@/ui/components/modules/Playlists/VideoDetails'
import VideoDetailsSkeleton from '@/ui/components/modules/Playlists/VideoDetailsSkeleton'
import EmptyState from '@/ui/components/shared/EmptyState'
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
    return (
      <div className="flex-1 flex items-center justify-center bg-neutral-900/20 rounded-[20px]">
        <EmptyState
          title="Playlist not found"
          description="This playlist doesn't exist, is private, or has been deleted."
        />
      </div>
    )
  }

  const { firstVideo, hiddenVideosCount, playlist } = playlistInfoResponse.data

  if (!v && firstVideo) {
    redirect(`/playlists/${playlistId}?v=${firstVideo.videoId}`)
  }

  return (
    <div className="flex flex-col xl:flex-row gap-5">
      <Suspense key={v} fallback={<VideoDetailsSkeleton />}>
        <VideoDetails currentUserId={session?.user.id} v={v} playlistId={playlistId} />
      </Suspense>
      <PlaylistQueueContainer
        playlistId={playlistId}
        hiddenVideosCount={hiddenVideosCount}
        playlistName={playlist.name}
        authorId={playlist.userId}
        guestUserId={session?.user.id}
      />
    </div>
  )
}
