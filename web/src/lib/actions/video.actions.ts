'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'
import { SearchVideosResponse, Video } from '@/types/video.types'

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

export async function searchVideos(
  query: string,
  page: number,
  pageSize: number,
): Promise<BaseServerResponse<SearchVideosResponse>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<SearchVideosResponse>>(
      `${endpoint}/search?searchText=${encodeURIComponent(query)}&page=${page}&pageSize=${pageSize}`,
    )

    return response
  } catch (error: unknown) {
    serverLog('GET_PUBLIC_VIDEOS_WHILE_SEARCH_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: `Failed to retrieve public videos from search`,
      errors: parseAxiosError(error),
    }
  }
}
