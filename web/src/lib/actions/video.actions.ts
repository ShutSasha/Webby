'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'
import { GetUserVideosResponse, SearchVideosResponse, Video, VideoSource } from '@/types/video.types'

const endpoint = '/videos'

export async function getVideoInfo(
  videoId: string,
  platform: VideoSource = 'Webby',
): Promise<BaseServerResponse<Video>> {
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
  searchPlatform?: 'Webby' | 'YouTube',
  nextPageToken?: string | null,
  contentSeed?: number | null,
): Promise<BaseServerResponse<SearchVideosResponse>> {
  try {
    const params = new URLSearchParams()
    if (query) params.append('searchText', query)

    params.append('page', page.toString())
    params.append('pageSize', pageSize.toString())
    if (searchPlatform) params.append('searchPlatform', searchPlatform)

    if (nextPageToken) params.append('nextPageToken', nextPageToken)
    if (contentSeed) params.append('contentSeed', contentSeed.toString())

    const { data: response } = await $api.get<BaseServerResponse<SearchVideosResponse>>(
      `${endpoint}/search?${params.toString()}`,
    )
    return response
  } catch (error: unknown) {
    serverLog('GET_PUBLIC_VIDEOS_WHILE_SEARCH_ERROR', error, true)
    return { data: null, success: false, message: `Failed to retrieve videos`, errors: parseAxiosError(error) }
  }
}

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

export async function getUserVideos(
  userId: string,
  page: number = 1,
  pageSize: number = 10,
  shouldShowDrafts: boolean = false,
): Promise<BaseServerResponse<GetUserVideosResponse>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<GetUserVideosResponse>>(
      `${endpoint}/users/${userId}?page=${page}&pageSize=${pageSize}&shouldShowDrafts=${shouldShowDrafts}`,
    )

    return response
  } catch (error: unknown) {
    serverLog('GET_USER_VIDEOS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve user videos',
      errors: parseAxiosError(error),
    }
  }
}

export async function deleteVideo(videoId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.delete<BaseServerResponse<null>>(`${endpoint}/${videoId}`)

    return response
  } catch (error: unknown) {
    serverLog('DELETE_VIDEO_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: `Failed to delete the video`,
      errors: parseAxiosError(error),
    }
  }
}

export async function checkVideoUploadStatusAction(videoId: string): Promise<BaseServerResponse<boolean>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<boolean>>(`${endpoint}/${videoId}/check-upload-status`)

    return response
  } catch (error: unknown) {
    serverLog('CHECK_VIDEO_UPLOAD_STATUS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to check video upload status',
      errors: parseAxiosError(error),
    }
  }
}

export async function cancelVideoUploadAction(videoId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.delete<BaseServerResponse<null>>(`${endpoint}/${videoId}/cancel`)

    return response
  } catch (error: unknown) {
    serverLog('CANCEL_VIDEO_UPLOAD_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to cancel video upload',
      errors: parseAxiosError(error),
    }
  }
}

export async function updateVideoMetadataAction(formData: FormData): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.patch<BaseServerResponse<null>>(`${endpoint}`, formData)

    return response
  } catch (error: unknown) {
    serverLog('UPDATE_VIDEO_METADATA_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to update video metadata',
      errors: parseAxiosError(error),
    }
  }
}

export async function incrementVideoView(videoId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`${endpoint}/increment-view`, {
      videoId,
    })

    return response
  } catch (error: unknown) {
    serverLog('INCREMENT_VIDEO_VIEW_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to increment view count',
      errors: parseAxiosError(error),
    }
  }
}
