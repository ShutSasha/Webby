import { getPlaylistInfo } from '@/lib/actions/playlist.actions'
import PlaylistMediaContainer from '@/ui/components/modules/Playlists/PlaylistMediaContainer'
import PlaylistQueueContainer from '@/ui/components/modules/Playlists/PlaylistQueueContainer/PlaylistQueueContainer'
import EmptyState from '@/ui/components/shared/EmptyState'
import { auth } from '@/workspace/auth'

type Props = { params: Promise<{ id: string }> }

export default async function PlaylistPage({ params }: Props) {
  const { id: playlistId } = await params
  const sessionPromise = auth()
  const playlistInfoPromise = getPlaylistInfo(playlistId)
  const [session, playlistInfoResponse] = await Promise.all([sessionPromise, playlistInfoPromise])

  if (!playlistInfoResponse.success || !playlistInfoResponse.data) {
    return (
      <div className="flex-1 flex items-center justify-center bg-surface/20 rounded-[20px]">
        <EmptyState
          title="Playlist not found"
          description="This playlist doesn't exist, is private, or has been deleted."
        />
      </div>
    )
  }

  const { firstVideo, hiddenVideosCount, playlist } = playlistInfoResponse.data

  return (
    <div className="flex flex-col xl:flex-row gap-5">
      <PlaylistMediaContainer
        playlistId={playlistId}
        mediaType={firstVideo?.mediaType}
        initialMediaId={firstVideo?.videoId}
      />
      <PlaylistQueueContainer
        playlistId={playlistId}
        hiddenVideosCount={hiddenVideosCount}
        playlistName={playlist.name}
        authorId={playlist.userId}
        guestUserId={session?.user?.id}
      />
    </div>
  )
}
