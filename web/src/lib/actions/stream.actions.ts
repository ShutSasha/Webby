'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'
import { SearchStreamsResponse } from '@/types/stream.types'

const endpoint = '/streams'

export async function searchStreams(
  query: string,
  page: number = 1,
  pageSize: number = 20,
): Promise<BaseServerResponse<SearchStreamsResponse>> {
  try {
    const params = new URLSearchParams()
    if (query) params.append('searchText', query)
    params.append('page', page.toString())
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
