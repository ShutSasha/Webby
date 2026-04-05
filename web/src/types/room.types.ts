import { PaginatedData } from './general.types'

// TODO change Room type
export type Room = {
  id: string
  categoryId: string
  name: string
  isPrivate: boolean
  hostId: string
  thumbnail: string
  token: string
}

export type GetPublicRooms = PaginatedData<Room>
