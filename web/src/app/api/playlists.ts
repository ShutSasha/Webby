'use server'

import $api from '@/app/api'
import { parseAxiosError, serverLog } from '@/lib/utils/utils'
import { BaseServerResponse } from '@/types/general'

const endpoint = '/playlists'

type Playlist = {
  pip: string
}

export async function getPlaylistInfo(playlistId: string): Promise<Playlist | BaseServerResponse | null> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<Playlist>>(`${endpoint}/${playlistId}/details`)

    return response.data
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
