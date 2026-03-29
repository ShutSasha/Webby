'use server'

import $api from '@/lib/api/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/utils'
import { BaseServerResponse } from '@/types/general'

const endpoint = '/achievements'

type BaseAchievement = {
  achievementId: string
  title: string
  description: string
  iconUrl: string
}

type AdminAchievement = BaseAchievement & {
  code: string
}

type UserAchievement = BaseAchievement & {
  isUnlocked: boolean
}

type GetUserAchievementsResponse = {
  pinnedAchievements: UserAchievement[]
  achievements: UserAchievement[]
}

export async function getAllAchievements(): Promise<BaseServerResponse<AdminAchievement[]>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<AdminAchievement[]>>(`${endpoint}`)
    return response
  } catch (error: unknown) {
    serverLog('FETCH_ALL_ACHIEVEMENTS_ERROR', error, true)
    return {
      success: false,
      data: null,
      message: 'Get all achievements error',
      errors: parseAxiosError(error),
    }
  }
}

export async function getUserAchievements(id: string): Promise<BaseServerResponse<GetUserAchievementsResponse>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<GetUserAchievementsResponse>>(`${endpoint}/${id}`)
    return response
  } catch (error: unknown) {
    serverLog('FETCH_USER_ACHIEVEMENTS_ERROR', error, true)
    return {
      success: false,
      data: null,
      message: 'Get user achievements error',
      errors: parseAxiosError(error),
    }
  }
}

export async function pinAchievement(achievementId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/users/achievements/${achievementId}`)
    return response
  } catch (error: unknown) {
    serverLog('PIN_ACHIEVEMENT_ERROR', error, true)
    return {
      success: false,
      data: null,
      message: 'Pin achievement error',
      errors: parseAxiosError(error),
    }
  }
}

export async function unpinAchievement(achievementId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.delete<BaseServerResponse<null>>(`/users/achievements/${achievementId}`)
    return response
  } catch (error: unknown) {
    serverLog('UNPIN_ACHIEVEMENT_ERROR', error, true)
    return {
      success: false,
      data: null,
      message: 'Unpin achievement error',
      errors: parseAxiosError(error),
    }
  }
}
