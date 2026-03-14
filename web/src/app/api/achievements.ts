'use server'

import $api from '@/app/api'
import { serverLog } from '@/lib/utils/utils'
import { BaseServerResponse } from '@/types/general'

const endpoint = '/achievements'

type Achievement = {
  achievementId: string
  code: string
  title: string
  description: string
  iconUrl: string
  userAchievements?: null
}

export async function getAllAchievements(): Promise<BaseServerResponse<Achievement[]> | null> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<Achievement[]>>(`${endpoint}`)
    return response
  } catch (error: unknown) {
    serverLog('FETCH_ALL_ACHIEVEMENTS_ERROR', error, true)
    return null
  }
}
