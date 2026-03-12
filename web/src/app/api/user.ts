'use server'

import $api from '@/app/api'
import { clog, serverLog } from '@/lib/utils/utils'
import { Achievement } from '@/types/achivement'
import { BaseServerResponse } from '@/types/general'
import { User, UserFolowStats } from '@/types/user'

const endpoint = '/users'

type GetUserResponse = {
  user: User
  userFollowStats: UserFolowStats
  pinnedUserAchievements: Achievement[]
}

export async function getUser(id: string): Promise<GetUserResponse | undefined | null> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<GetUserResponse>>(`${endpoint}/${id}`)
    return response.data
  } catch (error: unknown) {
    serverLog('USER_FETCH_ERROR', error, true)
  }
}

export async function uploadNewUserPhoto(id: string, formData: FormData) {
  try {
    const { data: response } = await $api.patch<BaseServerResponse<User>>(
      `${endpoint}/update-user-avatar/${id}`,
      formData,
    )

    clog('upload new user photo response', response)
    return response.data
  } catch (error: unknown) {
    serverLog('USER_UPLOAD_ERROR', error, true)
    return null
  }
}
