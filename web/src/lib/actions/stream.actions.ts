'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'
import { SearchStreamsResponse, Stream } from '@/types/stream.types'

const endpoint = '/streams'

export async function searchStreams(
  query: string,
  nextPageToken?: string | null,
  pageSize: number = 20,
): Promise<BaseServerResponse<SearchStreamsResponse>> {
  try {
    const params = new URLSearchParams()
    if (query) params.append('searchText', query)

    if (nextPageToken) params.append('nextPageToken', nextPageToken)

    params.append('pageSize', pageSize.toString())

    const { data: response } = await $api.get<BaseServerResponse<SearchStreamsResponse>>(
      `${endpoint}/search?${params.toString()}`,
    )
    return response
  } catch (error: unknown) {
    serverLog('GET_STREAMS_SEARCH_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to retrieve streams',
      errors: parseAxiosError(error),
    }
  }
}

export async function getStreamInfo(streamerId: string): Promise<BaseServerResponse<Stream>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<Stream>>(`${endpoint}/${streamerId}`)
    return response
  } catch (error: unknown) {
    serverLog('GET_STREAM_INFO_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to retrieve stream information',
      errors: parseAxiosError(error),
    }
  }
}
