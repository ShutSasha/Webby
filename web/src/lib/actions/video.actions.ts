'use server'

import $api from '@/lib/config/api.config'
import { clog, parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse, Platform } from '@/types/general.types'
import { SearchVideosResponse, UploadVideoResponse, Video } from '@/types/video.types'

const endpoint = '/videos'

export async function getVideoInfo(videoId: string, platform: Platform = 'Webby'): Promise<BaseServerResponse<Video>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<Video>>(`${endpoint}/${videoId}?platform=${platform}`)

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

// TODO: maybe do it in server action
// export async function uploadVideoFileAction(formData: FormData): Promise<BaseServerResponse<UploadVideoResponse>> {
//   try {
//     const file = formData.get('VideoFile') as File | null
//     if (!file) {
//       throw new Error('VideoFile is missing in the request')
//     }

//     const { data: response } = await $api.post<BaseServerResponse<UploadVideoResponse>>(`${endpoint}/upload`, file)

//     return response
//   } catch (error: unknown) {
//     serverLog('UPLOAD_VIDEO_FILE_ERROR', error, true)

//     return {
//       data: null,
//       success: false,
//       message: 'Failed to upload video file',
//       errors: parseAxiosError(error),
//     }
//   }
// }

export async function createVideoMetadataAction(formData: FormData): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`${endpoint}`, formData)

    return response
  } catch (error: unknown) {
    serverLog('CREATE_VIDEO_METADATA_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to save video metadata',
      errors: parseAxiosError(error),
    }
  }
}
