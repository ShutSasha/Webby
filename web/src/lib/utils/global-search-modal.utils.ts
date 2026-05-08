import { Route } from 'next'

import { FallbackType } from '@/ui/components/shared/SafeImage'

export type EntityType = 'Video' | 'Room' | 'Playlist' | 'Stream' | 'User' | 'YouTube'

export const getRoute = (type: EntityType, id: string): Route => {
  const routes: Record<EntityType, string> = {
    Video: `/videos/${id}`,
    Room: `/rooms/${id}`,
    Playlist: `/playlists/${id}`,
    Stream: `/streams/${id}`,
    User: `/profile/${id}`,
    YouTube: `/videos/${id}`,
  }
  return routes[type] as Route
}

export const FALLBACK_MAP: Record<EntityType, FallbackType> = {
  User: 'user',
  Video: 'video',
  Room: 'room',
  Playlist: 'playlist',
  Stream: 'video',
  YouTube: 'video',
}
