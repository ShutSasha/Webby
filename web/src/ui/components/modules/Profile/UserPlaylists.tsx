import PlaylistItem from '../Playlist/PlaylistItem'

export default async function UserPlaylists() {
  // await, sync Public playlists
  await new Promise(resolve => {
    setTimeout(() => {
      resolve('')
    }, 1300)
  })

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {[...new Array(10)].map((item, index) => (
        <PlaylistItem key={index} />
      ))}
    </div>
  )
}

export function PlaylistItemSkeleton() {
  return (
    <div className="relative">
      {/* Thumbnail */}
      <div
        className="w-full h-[260px] md:h-[200px] xl:h-[180px] rounded-2xl mb-2 bg-neutral-800
          animate-[shimmer_1.5s_infinite]"
      />

      {/* Playlist name */}
      <div className="h-4 w-2/3 bg-neutral-800 rounded-md mb-2 animate-pulse" />

      {/* Creator */}
      <div className="h-3 w-1/3 bg-neutral-800 rounded-md animate-pulse" />

      {/* Videos badge */}
      <div className="absolute top-2 right-2 h-6 w-16 rounded-lg bg-neutral-700/20 animate-pulse" />
    </div>
  )
}

export function UserPlaylistsSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {Array.from({ length: 20 }).map((_, index) => (
        <PlaylistItemSkeleton key={index} />
      ))}
    </div>
  )
}
