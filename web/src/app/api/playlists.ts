'use server'

import $api from '@/app/api'
import { parseAxiosError, serverLog } from '@/lib/utils/utils'
import { BaseServerResponse, Optional } from '@/types/general'

import { Video } from './videos'

const endpoint = '/playlists'

export type BasePlaylist = {
  playlistId: string
  name: string
  countOfVideos: number
  playlistCover: string
  isPrivate: boolean
}

export type SearchPlaylsit = BasePlaylist & {
  userId: string
  username: string
}

export type CreatePlaylist = BasePlaylist & {
  userId: string
}

export type PlaylistInfoDetails = BasePlaylist & {
  userId: string
}

export type PlaylistInfo = {
  playlist: PlaylistInfoDetails
  firstVideo?: Omit<Video, 'videoUrl'>
}

export async function getPlaylistInfo(playlistId: string): Promise<BaseServerResponse<PlaylistInfo>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<PlaylistInfo>>(`${endpoint}/${playlistId}/details`)

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
): Promise<BaseServerResponse<CreatePlaylist>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<CreatePlaylist>>(`${endpoint}`, {
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

type PlaylistsSearchData = {
  items: SearchPlaylsit[]
  page: number
  pageSize: number
  totalCount: number
}

export async function searchPlaylists(
  query: string,
  page: number,
  pageSize: number,
): Promise<BaseServerResponse<PlaylistsSearchData>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<PlaylistsSearchData>>(
      `${endpoint}/search?SearchText=${encodeURIComponent(query)}&Page=${page}&PageSize=${pageSize}`,
    )

    return response
  } catch (error: unknown) {
    serverLog('GET_PLAYLISTS_WHILE_SEARCH_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: `Failed to retrieve playlists from playlist search`,
      errors: parseAxiosError(error),
    }
  }
}

export type UserPlaylistDetails = BasePlaylist & {
  isVideoAdded: boolean
}

export type UserPlaylistsSearchData = {
  items: UserPlaylistDetails[]
  page: number
  pageSize: number
  totalCount: number
}

export async function searchUserPlaylists(
  userId: string,
  query: string,
  page: number,
  pageSize: number,
  targetVideoId?: string | undefined,
): Promise<BaseServerResponse<UserPlaylistsSearchData>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<UserPlaylistsSearchData>>(
      `${endpoint}/${userId}?SearchText=${encodeURIComponent(query)}&Page=${page}&PageSize=${pageSize}&videoId=${targetVideoId ?? ''}`,
    )

    return response
  } catch (error: unknown) {
    serverLog('GET_USER_PLAYLISTS_WHILE_SEARCH_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: `Failed to retrieve user playlists from playlist search`,
      errors: parseAxiosError(error),
    }
  }
}
