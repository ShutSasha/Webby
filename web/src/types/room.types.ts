import { PaginatedData } from './general.types'

export type Room = {
  id: string
  name: string
  categoryName: string
  isPrivate: boolean
  hostId: string
  hostUsername: string
  hostAvatarUrl: string
  thumbnail: string
  token: string
}

export type GetPublicRooms = PaginatedData<Room>
