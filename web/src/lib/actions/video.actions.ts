'use server'

import $api from '@/lib/api/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/utils'
import { BaseServerResponse } from '@/types/general'

export type VideoUser = {
  userId: string
  username: string
  avatarUrl: string
  isFollowed: boolean
}

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
  user: VideoUser
}

const endpoint = '/videos'

export async function getVideoInfo(videoId: string): Promise<BaseServerResponse<Video>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<Video>>(`${endpoint}/${videoId}`)

    return response
  } catch (error: unknown) {
    serverLog('GET_VIDEO_INFO_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve video information',
      errors: parseAxiosError(error),
    }
  }
}
