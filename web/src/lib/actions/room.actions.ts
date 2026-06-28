'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'
import {
  GetPublicRooms,
  GetRoomMembersResponse,
  GetRoomQueueResponse,
  GetUserRoomsResponse,
  Room,
} from '@/types/room.types'

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

    const url = `${endpoint}/public?${queryString}`

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

    const url = `${endpoint}/my?${queryString}`

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

export async function initiateSyncAction(roomId: string | undefined): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`${endpoint}/${roomId}/sync`)
    return response
  } catch (error: unknown) {
    serverLog('INITIATE_SYNC_ERROR', error, true)
    return { data: null, success: false, message: 'Failed to initiate sync', errors: parseAxiosError(error) }
  }
}

export async function reportTimecodeAction(
  roomId: string,
  syncId: string,
  timecode: number,
): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`${endpoint}/${roomId}/sync/report`, {
      syncId,
      timecode,
    })
    return response
  } catch (error: unknown) {
    serverLog('REPORT_TIMECODE_ERROR', error, true)
    return { data: null, success: false, message: 'Failed to report timecode', errors: parseAxiosError(error) }
  }
}

export async function getWsTokenAction(): Promise<BaseServerResponse<string>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<string>>(`/ws-token`)
    return response
  } catch (error: unknown) {
    serverLog('GET_WS_TOKEN_ERROR', error, true)
    return { data: null, success: false, message: 'Failed to get WS token', errors: parseAxiosError(error) }
  }
}

export async function getRoomById(roomId: string): Promise<BaseServerResponse<Room>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<Room>>(`${endpoint}/${roomId}`)
    return response
  } catch (error: unknown) {
    serverLog('GET_ROOM_BY_ID_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: `Failed to retrieve room details`,
      errors: parseAxiosError(error),
    }
  }
}

export async function getRoomQueueAction(
  roomId: string,
  page: number,
  limit: number = 20,
): Promise<BaseServerResponse<GetRoomQueueResponse>> {
  try {
    const params = new URLSearchParams()
    params.append('page', page.toString())
    params.append('limit', limit.toString())

    const url = `${endpoint}/${roomId}/queue?${params.toString()}`

    const { data: response } = await $api.get<BaseServerResponse<GetRoomQueueResponse>>(url)
    return response
  } catch (error: unknown) {
    serverLog('GET_ROOM_QUEUE_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to retrieve room queue',
      errors: parseAxiosError(error),
    }
  }
}

export async function activateQueueItemAction(roomId: string, itemId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.patch<BaseServerResponse<null>>(
      `${endpoint}/${roomId}/queue/${itemId}/activate`,
    )
    return response
  } catch (error: unknown) {
    serverLog('ACTIVATE_QUEUE_ITEM_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to activate queue item',
      errors: parseAxiosError(error),
    }
  }
}

export async function removeQueueItemAction(roomId: string, itemId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.delete<BaseServerResponse<null>>(`${endpoint}/${roomId}/queue/${itemId}`)
    return response
  } catch (error: unknown) {
    serverLog('REMOVE_QUEUE_ITEM_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to remove item from queue',
      errors: parseAxiosError(error),
    }
  }
}

export async function addQueueItemAction(roomId: string, videoId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`${endpoint}/${roomId}/queue`, { videoId })
    return response
  } catch (error: unknown) {
    serverLog('ADD_QUEUE_ITEM_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to add item to queue',
      errors: parseAxiosError(error),
    }
  }
}

export async function getRoomMembersAction(
  roomId: string,
  page: number,
  limit: number = 10,
  search: string = '',
): Promise<BaseServerResponse<GetRoomMembersResponse>> {
  try {
    const params = new URLSearchParams()
    params.append('page', page.toString())
    params.append('limit', limit.toString())

    if (search) {
      params.append('search', search)
    }

    const { data: response } = await $api.get<BaseServerResponse<GetRoomMembersResponse>>(
      `${endpoint}/${roomId}/members?${params.toString()}`,
    )
    return response
  } catch (error: unknown) {
    serverLog('GET_ROOM_MEMBERS_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to retrieve room members',
      errors: parseAxiosError(error),
    }
  }
}

export async function removeRoomMemberAction(roomId: string, memberId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.delete<BaseServerResponse<null>>(`${endpoint}/${roomId}/members/${memberId}`)
    return response
  } catch (error: unknown) {
    serverLog('REMOVE_ROOM_MEMBER_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to remove member from room',
      errors: parseAxiosError(error),
    }
  }
}

export async function addRoomMembersAction(roomId: string, userIds: string[]): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`${endpoint}/${roomId}/members`, {
      userIds,
    })
    return response
  } catch (error: unknown) {
    serverLog('ADD_ROOM_MEMBERS_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to add members to the room',
      errors: parseAxiosError(error),
    }
  }
}

export type MemberPoints = {
  points: number
}

export async function getRoomMemberPointsAction(roomId: string): Promise<BaseServerResponse<MemberPoints>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<MemberPoints>>(`/rooms/${roomId}/members/points`)
    return response
  } catch (error: unknown) {
    serverLog('GET_ROOM_MEMBER_POINTS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve member points',
      errors: parseAxiosError(error),
    }
  }
}
