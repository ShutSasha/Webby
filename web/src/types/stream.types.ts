import { PaginatedData } from './general.types'

export type StreamUser = {
  userId: string
  username: string
  avatarUrl: string
  isFollowed: boolean
}

export type Stream = {
  streamId: string
  name: string
  source: 'Twitch' | 'Webby' | string
  previewUrl: string
  startedAt: string
  viewers: number
  streamUrl: string
  user: StreamUser
}

export type SearchStreamsResponse = PaginatedData<Stream>