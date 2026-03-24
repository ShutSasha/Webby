'use server'

import $api from '@/app/api'
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

export async function getAllAchievements(): Promise<BaseServerResponse<AdminAchievement[]> | null> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<AdminAchievement[]>>(`${endpoint}`)
    return response
  } catch (error: unknown) {
    serverLog('FETCH_ALL_ACHIEVEMENTS_ERROR', error, true)
    return null
  }
}

export async function getUserAchievements(
  id: string,
): Promise<BaseServerResponse<{ pinnedAchievements: UserAchievement[]; achievements: UserAchievement[] }> | null> {
  try {
    const { data: response } = await $api.get<
      BaseServerResponse<{ pinnedAchievements: UserAchievement[]; achievements: UserAchievement[] }>
    >(`${endpoint}/${id}`)
    return response
  } catch (error: unknown) {
    serverLog('FETCH_USER_ACHIEVEMENTS_ERROR', error, true)
    return null
  }
}

export async function pinAchievement(
  achievementId: string,
): Promise<BaseServerResponse<null> | { error: Record<string, string> }> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/users/achievements/${achievementId}`)
    return response
  } catch (error: unknown) {
    serverLog('PIN_ACHIEVEMENT_ERROR', error, true)
    return {
      error: parseAxiosError(error),
      success: false,
    }
  }
}

export async function unpinAchievement(
  achievementId: string,
): Promise<BaseServerResponse<null> | { error: Record<string, string> }> {
  try {
    const { data: response } = await $api.delete<BaseServerResponse<null>>(`/users/achievements/${achievementId}`)
    return response
  } catch (error: unknown) {
    serverLog('UNPIN_ACHIEVEMENT_ERROR', error, true)
    return {
      error: parseAxiosError(error),
      success: false,
    }
  }
}
