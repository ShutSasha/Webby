'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse, PaginatedData } from '@/types/general.types'
import { Role } from '@/types/user.types'

export type AbsoluteStats = {
  activeReportsCount: number
  monthRevenue: number
  totalRegistrations: number
  totalRooms: number
}

export type MonthlyStats = Record<string, number>

export async function getAbsoluteStatsAction(): Promise<BaseServerResponse<AbsoluteStats>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<AbsoluteStats>>('/stats')

    return response
  } catch (error: unknown) {
    serverLog('GET_ABSOLUTE_STATS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve absolute statistics',
      errors: parseAxiosError(error),
    }
  }
}

export async function getRegistrationsStatsAction(): Promise<BaseServerResponse<MonthlyStats>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<MonthlyStats>>('/stats/registrations')

    return response
  } catch (error: unknown) {
    serverLog('GET_REGISTRATIONS_STATS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve registrations statistics',
      errors: parseAxiosError(error),
    }
  }
}

export async function getSubscriptionsStatsAction(): Promise<BaseServerResponse<MonthlyStats>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<MonthlyStats>>('/stats/subscriptions')

    return response
  } catch (error: unknown) {
    serverLog('GET_SUBSCRIPTIONS_STATS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve subscriptions statistics',
      errors: parseAxiosError(error),
    }
  }
}

export type Complaint = {
  id: string
  complainer: {
    id: string
    username: string
  }
  target: {
    type: 'Video' | 'User'
    id: string
    name: string
  }
  reasonType: string
  additionalInfo: string
  createdAt: string
}

export async function getComplaintsAction(
  page: number = 1,
  limit: number = 10,
): Promise<BaseServerResponse<PaginatedData<Complaint>>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<PaginatedData<Complaint>>>('/complaints', {
      params: { page, limit },
    })

    return response
  } catch (error: unknown) {
    serverLog('GET_COMPLAINTS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve complaints',
      errors: parseAxiosError(error),
    }
  }
}

export async function acceptComplaintAction(complaintId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/complaints/${complaintId}/accept`)

    return response
  } catch (error: unknown) {
    serverLog('ACCEPT_COMPLAINT_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to accept complaint',
      errors: parseAxiosError(error),
    }
  }
}

export async function denyComplaintAction(complaintId: string, reason: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/complaints/${complaintId}/deny`, {
      reason,
    })

    return response
  } catch (error: unknown) {
    serverLog('DENY_COMPLAINT_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to deny complaint',
      errors: parseAxiosError(error),
    }
  }
}

export type AdminUserRecord = {
  userId: string
  email: string
  username: string
  about: string
  avatarUrl: string
  isBanned: boolean
  createdAt: string
  role: Role
}

export async function searchAdminUsersAction(
  searchText: string = '',
  page: number = 1,
  limit: number = 10,
): Promise<BaseServerResponse<PaginatedData<AdminUserRecord>>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<PaginatedData<AdminUserRecord>>>('/users/search', {
      params: {
        searchText,
        page,
        pageSize: limit,
      },
    })

    return response
  } catch (error: unknown) {
    serverLog('SEARCH_ADMIN_USERS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to search users',
      errors: parseAxiosError(error),
    }
  }
}

export async function banUserAction(userId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/users/blocks/${userId}/ban`)

    return response
  } catch (error: unknown) {
    serverLog('BAN_USER_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to ban user',
      errors: parseAxiosError(error),
    }
  }
}

export async function unbanUserAction(userId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/users/blocks/${userId}/unban`)

    return response
  } catch (error: unknown) {
    serverLog('UNBAN_USER_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to unban user',
      errors: parseAxiosError(error),
    }
  }
}

export async function changeUserRoleAction(userId: string, userRole: Role): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/users/roles/${userId}/change-role`, null, {
      params: { userRole },
    })

    return response
  } catch (error: unknown) {
    serverLog('CHANGE_USER_ROLE_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to change user role',
      errors: parseAxiosError(error),
    }
  }
}

export type AdminVideoUser = {
  avatarUrl: string
  isFollowed: boolean
  userId: string
  username: string
}

export type AdminVideoRecord = {
  videoId: string
  previewUrl: string
  name: string
  isBanned: boolean
  createdAt: string
  user: AdminVideoUser
}

export async function searchModerationVideosAction(
  searchText: string = '',
  page: number = 1,
  limit: number = 20,
): Promise<BaseServerResponse<PaginatedData<AdminVideoRecord>>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<PaginatedData<AdminVideoRecord>>>(
      '/videos/moderation/search',
      {
        params: {
          searchText,
          page,
          pageSize: limit,
        },
      },
    )

    return response
  } catch (error: unknown) {
    serverLog('SEARCH_MODERATION_VIDEOS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to search moderation videos',
      errors: parseAxiosError(error),
    }
  }
}

export async function banVideoAction(videoId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/videos/blocks/${videoId}/ban`)

    return response
  } catch (error: unknown) {
    serverLog('BAN_VIDEO_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to ban video',
      errors: parseAxiosError(error),
    }
  }
}

export async function unbanVideoAction(videoId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/videos/blocks/${videoId}/unban`)

    return response
  } catch (error: unknown) {
    serverLog('UNBAN_VIDEO_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to unban video',
      errors: parseAxiosError(error),
    }
  }
}
