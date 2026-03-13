'use server'

import $api from '@/app/api'
import { serverLog } from '@/lib/utils/utils'
import { Achievement } from '@/types/achivement'
import { BaseServerResponse } from '@/types/general'
import { User, UserFolowStats } from '@/types/user'

const endpoint = '/users'

type GetUserResponse = {
  user: User
  userFollowStats: UserFolowStats
  pinnedUserAchievements: Achievement[]
}

type GetUserFollowsResponse = {
  userId: string
  avatarUrl: string
  username: string
  followersCount: number
}

type GetUserFollowersResponse = GetUserFollowsResponse

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

    return response.data
  } catch (error: unknown) {
    serverLog('USER_UPLOAD_ERROR', error, true)
    return null
  }
}

export async function getUserFollows(id: string) {
  try {
    const { data: response } = await $api.get<BaseServerResponse<GetUserFollowsResponse[]>>(`${endpoint}/${id}/follows`)

    return response.data
  } catch (error: unknown) {
    serverLog('FETCH_USER_FOLLOWS_ERROR', error, true)
    return null
  }
}

export async function getUserFollowers(id: string) {
  try {
    const { data: response } = await $api.get<BaseServerResponse<GetUserFollowersResponse[]>>(
      `${endpoint}/${id}/followers`,
    )

    return response.data
  } catch (error: unknown) {
    serverLog('FETCH_USER_FOLLOWERS_ERROR', error, true)
    return null
  }
}
