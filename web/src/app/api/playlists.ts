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

export async function createUserPlaylist(
  name: string,
  isPrivate: boolean,
): Promise<BaseServerResponse<PlaylistDetails>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<PlaylistDetails>>(`${endpoint}`, {
      name,
      isPrivate,
    })

    return response
  } catch (error: unknown) {
    serverLog('CREATE_USER_PLAYLIST_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: `Failed to create playlist "${name}". Please try again later.`,
      errors: parseAxiosError(error),
    }
  }
}
