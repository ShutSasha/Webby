import { playlists } from '@/lib/placeholder-data/profile'

import PlaylistItem from '../Playlist/PlaylistItem'

export default async function UserPlaylists() {
  // await, sync Public playlists
  await new Promise(r => setTimeout(r, 1500))

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
      {playlists.map(playlist => (
        <PlaylistItem key={playlist.id} {...playlist} />
      ))}
    </div>
  )
}

export function PlaylistItemSkeleton() {
  return (
    <div className="relative">
      {/* Thumbnail */}
      <div className="w-full aspect-video rounded-2xl mb-1 bg-neutral-800 animate-pulse" />

      {/* Playlist name */}
      <div className="h-4 w-3/4 bg-neutral-800 rounded-md mb-1 animate-pulse" />

      {/* Creator */}
      <div className="h-3 w-1/3 bg-neutral-800 rounded-md animate-pulse" />

      <div className="absolute top-2 right-2 h-[26px] w-[58px] rounded-lg bg-neutral-800 animate-pulse" />
    </div>
  )
}

export function UserPlaylistsSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
      {Array.from({ length: 12 }).map((_, index) => (
        <PlaylistItemSkeleton key={index} />
      ))}
    </div>
  )
}
