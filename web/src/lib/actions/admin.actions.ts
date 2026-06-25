'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'

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
