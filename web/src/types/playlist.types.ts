export interface CachedPlaylist {
  playlistId: string
  name: string
  isPrivate: boolean
  [key: string]: unknown
}
