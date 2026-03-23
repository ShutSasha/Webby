'use server'

import $api from '@/app/api'
import { clog, parseAxiosError, serverLog } from '@/lib/utils/utils'
import { BaseServerResponse, Optional } from '@/types/general'

import { Video } from './videos'

const endpoint = '/playlists'

export type PlaylistDetails = {
  playlistId: string
  userId: string
  name: string
  description: string
  countOfVideos: number
  playlistCover: string
  isPrivate: boolean
}

export type PlaylistData = {
  playlist: PlaylistDetails
  firstVideo: Omit<Video, 'videoUrl'>
}

export async function getPlaylistInfo(playlistId: string): Promise<BaseServerResponse<PlaylistData>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<PlaylistData>>(`${endpoint}/${playlistId}/details`)

    return response
  } catch (error: unknown) {
    serverLog('GET_PLAYLIST_INFO_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve playlist information',
      errors: parseAxiosError(error),
    }
  }
}

export type PlaylistVideo = Optional<Video, 'user'>

type PlaylistVideosData = {
  items: PlaylistVideo[]
  page: number
  pageSize: number
  totalCount: number
}

export async function getPlaylistVideos(
  playlistId: string,
  query: string,
  page: number,
  pageSize: number,
): Promise<BaseServerResponse<PlaylistVideosData>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<PlaylistVideosData>>(
      `/videos/${playlistId}/search?SearchText=${encodeURIComponent(query)}&Page=${page}&PageSize=${pageSize}`,
    )

    clog('videos', response)

    return response
  } catch (error: unknown) {
    serverLog('GET_PLAYLIST_VIDEOS_WHILE_SEARCH_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: `Failed to retrieve videos from playlist search`,
      errors: parseAxiosError(error),
    }
  }
}
