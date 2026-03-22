'use server'

import $api from '@/app/api'
import { parseAxiosError, serverLog } from '@/lib/utils/utils'
import { BaseServerResponse } from '@/types/general'

import { Video } from './videos'

const endpoint = '/playlists'

export type PlaylistDetails = {
  playlistId: string
  userId: string
  name: string
  description: string
  countOfVideos: number
  playlistCover: string
}

export type PlaylistData = {
  playlist: PlaylistDetails
  video: Video
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
