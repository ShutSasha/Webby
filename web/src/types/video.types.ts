import { PaginatedData } from './general.types'

export type VideoUser = {
  userId: string
  username: string
  avatarUrl: string
  isFollowed: boolean
}

export type UploadStatus = 'Ready' | 'Uploading' | 'Failed' | 'Canceled'

export type VideoSource = 'YouTube' | 'Webby'

export type Video = {
  videoId: string
  videoUrl: string
  name: string
  views: number
  previewUrl: string
  isPrivate: boolean
  isPublished: boolean
  createdAt: string
  videoTags: string[] | null
  description: string | null
  duration: number
  videotags: string[]
  source: VideoSource
  videoUploadStatus: UploadStatus
  user: VideoUser
}

export type RecommendedVideo = {
  videoId: string
  name: string
  previewUrl: string
  views: number
  isPrivate: boolean
  createdAt: string
  duration: number
  user: VideoUser
}

export type GetRecommendationsResponse = PaginatedData<RecommendedVideo>

export type UploadVideoResponse = {
  videoId: string
  userId: string
  name: string
  videoUploadStatus: UploadStatus
  isPrivate: boolean
}

export type SearchVideosResponse = PaginatedData<Video>

export type GetUserVideosResponse = PaginatedData<Video>
