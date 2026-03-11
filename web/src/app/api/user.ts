'use server'

import $api from '@/app/api'
import { clog, serverLog } from '@/lib/utils/utils'
import { BaseServerResponse } from '@/types/general'
import { User } from '@/types/user'

const endpoint = '/users'

//TODO: add types for other properites
type GetUserResponse = {
  user: User
  userFollowStats: { followers: 0; following: 0 }
  pinnedUserAchievements: []
}

export async function getUser(id: string): Promise<GetUserResponse | undefined | null> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<GetUserResponse>>(`${endpoint}/${id}`)
    return response.data
  } catch (error: unknown) {
    serverLog('USER_FETCH_ERROR', error, true)
  }
}
