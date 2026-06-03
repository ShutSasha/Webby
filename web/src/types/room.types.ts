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

export type RoomQueueItem = {
  id: string
  videoId: string
  title: string
  thumbnail: string
  videoUrl: string
  isActive: boolean
  position: number
  isFolder?: boolean
  children?: RoomQueueItem[]
}

export type GetRoomQueueResponse = PaginatedData<RoomQueueItem>

export type RoomMember = {
  userId: string
  username: string
  avatarUrl: string
  roomPoints: number
}

export type GetRoomMembersResponse = PaginatedData<RoomMember>
