import PlaylistItem from '../Playlist/PlaylistItem'

export default async function UserPlaylists() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {[...new Array(10)].map((item, index) => (
        <PlaylistItem key={index} />
      ))}
    </div>
  )
}
