import { PaginatedData } from './general.types'
import { VideoSource } from './video.types'

export type StreamerInformation = {
  userId: string
  username: string
  avatarUrl: string
  isFollowed: boolean
}

export type Stream = {
  streamerId: string
  name: string
  source: VideoSource
  previewUrl: string
  description: string | null
  startedAt: string
  viewers: number
  streamUrl: string
  streamTags: string[]
  streamerInformation: StreamerInformation
}

export type SearchStreamsResponse = PaginatedData<Stream>
