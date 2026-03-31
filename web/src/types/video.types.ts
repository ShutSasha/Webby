import { PaginatedData } from './general.types'

export type VideoUser = {
  userId: string
  username: string
  avatarUrl: string
  isFollowed: boolean
}

export type UploadStatus = 'Ready' | 'Uploading' | 'Failed' | 'Canceled'

export type Video = {
  videoId: string
  videoUrl: string
  name: string
  views: number
  previewUrl: string
  isPrivate: boolean
  createdAt: string
  videoTags: string[] | null
  description: string | null
  duration: number
  videotags: string[]
  videoUploadStatus: UploadStatus
  user: VideoUser
}

export type SearchVideosResponse = PaginatedData<Video>
