'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'
import { GetPublicRooms, GetUserRoomsResponse, Room } from '@/types/room.types'

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

// TODO-ROOMS: check whether search params set correct or not
export async function getUserRooms(
  query: string,
  page: number,
  pageSize: number,
): Promise<BaseServerResponse<GetUserRoomsResponse>> {
  try {
    const params = new URLSearchParams()

    if (query) params.append('search', query)
    if (page) params.append('page', page.toString())
    if (pageSize) params.append('limit', pageSize.toString())

    const queryString = params.toString()
    // TODO-ROOMS: only url with queryString must be exist, replace it later when server is ready
    const url = queryString ? `${endpoint}/my?${queryString}` : `${endpoint}/my`

    const { data: response } = await $api.get<BaseServerResponse<GetUserRoomsResponse>>(url)

    return response
  } catch (error: unknown) {
    serverLog('GET_USER_ROOMS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: `Failed to retrieve user rooms`,
      errors: parseAxiosError(error),
    }
  }
}

export async function createRoomAction(formData: FormData): Promise<BaseServerResponse<Room>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<Room>>(`${endpoint}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })

    return response
  } catch (error: unknown) {
    serverLog('CREATE_ROOM_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: `Failed to create room. Please try again later.`,
      errors: parseAxiosError(error),
    }
  }
}

export async function deleteRoomAction(roomId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.delete<BaseServerResponse<null>>(`${endpoint}/${roomId}`)
    return response
  } catch (error: unknown) {
    serverLog('DELETE_ROOM_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: `Failed to delete room`,
      errors: parseAxiosError(error),
    }
  }
}

export async function updateRoomAction(roomId: string, formData: FormData): Promise<BaseServerResponse<Room>> {
  try {
    const { data: response } = await $api.put<BaseServerResponse<Room>>(`${endpoint}/${roomId}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response
  } catch (error: unknown) {
    serverLog('UPDATE_ROOM_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: `Failed to update room`,
      errors: parseAxiosError(error),
    }
  }
}
