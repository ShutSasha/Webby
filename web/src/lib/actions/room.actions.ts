'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'
import { GetPublicRooms } from '@/types/room.types'

const endpoint = '/rooms'

export async function getPublicRooms(
  query: string,
  page: number,
  pageSize: number,
  category: string,
): Promise<BaseServerResponse<GetPublicRooms>> {
  try {
    const params = new URLSearchParams()

    if (query) params.append('search', query)
    if (page) params.append('page', page.toString())
    if (pageSize) params.append('limit', pageSize.toString())
    if (category) params.append('category', category.toString())

    const queryString = params.toString()

    const url = queryString ? `${endpoint}/public?${queryString}` : endpoint

    const { data: response } = await $api.get<BaseServerResponse<GetPublicRooms>>(url)

    return response
  } catch (error: unknown) {
    serverLog('GET_PUBLIC_ROOMS_WHILE_SEARCH_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: `Failed to retrieve rooms from search`,
      errors: parseAxiosError(error),
    }
  }
}
