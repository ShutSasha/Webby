'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'
import { GetNotificationsResponse } from '@/types/notification.types'

const endpoint = '/notifications'

export async function getUserNotifications(
  page: number = 1,
  pageSize: number = 10,
): Promise<BaseServerResponse<GetNotificationsResponse>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<GetNotificationsResponse>>(
      `${endpoint}?page=${page}&pageSize=${pageSize}`,
    )

    return response
  } catch (error: unknown) {
    serverLog('GET_USER_NOTIFICATIONS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve user notifications',
      errors: parseAxiosError(error),
    }
  }
}

export async function getUnreadNotificationsCount(): Promise<BaseServerResponse<number>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<number>>(`${endpoint}/count-unread-messages`)

    return response
  } catch (error: unknown) {
    serverLog('GET_UNREAD_NOTIFICATIONS_COUNT_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve unread notifications count',
      errors: parseAxiosError(error),
    }
  }
}

export async function getUnreadNotifications(
  page: number = 1,
  pageSize: number = 10,
): Promise<BaseServerResponse<GetNotificationsResponse>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<GetNotificationsResponse>>(
      `${endpoint}/unread?page=${page}&pageSize=${pageSize}`,
    )
    return response
  } catch (error: unknown) {
    serverLog('GET_UNREAD_NOTIFICATIONS_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to retrieve unread notifications',
      errors: parseAxiosError(error),
    }
  }
}

export async function deleteNotificationAction(notificationId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.delete<BaseServerResponse<null>>(`${endpoint}/${notificationId}`)
    return response
  } catch (error: unknown) {
    serverLog('DELETE_NOTIFICATION_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to delete notification',
      errors: parseAxiosError(error),
    }
  }
}

export async function setNotificationsReadStatusAction(notificationIds: string[]): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`${endpoint}/set-read-status`, notificationIds)
    return response
  } catch (error: unknown) {
    serverLog('SET_NOTIFICATIONS_READ_STATUS_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Failed to set notifications read status',
      errors: parseAxiosError(error),
    }
  }
}
