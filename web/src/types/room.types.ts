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
  chatId: string
}

export type GetPublicRooms = PaginatedData<Room>

// TODO-ROOMS: change UserRoom type when server updates its response
export type UserRoom = {
  id: string
  name: string
  categoryName: string
  isPrivate: boolean
  hostId: string
  thumbnail: string
}

export type GetUserRoomsResponse = PaginatedData<UserRoom>
