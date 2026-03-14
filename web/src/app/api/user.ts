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

export async function uploadNewUserPhoto(id: string, formData: FormData): Promise<User | null> {
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

export async function getUserFollows(id: string): Promise<GetUserFollowsResponse[] | null> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<GetUserFollowsResponse[]>>(`${endpoint}/${id}/follows`)

    return response.data
  } catch (error: unknown) {
    serverLog('FETCH_USER_FOLLOWS_ERROR', error, true)
    return null
  }
}

export async function getUserFollowers(id: string): Promise<GetUserFollowersResponse[] | null> {
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

export async function leaveComplaint(
  authorId: string,
  targetUserId: string,
  reasonType: string,
  additionalInfo: string,
): Promise<BaseServerResponse<null> | null> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/complaints`, {
      authorId,
      targetUserId,
      reasonType,
      additionalInfo,
    })

    return response
  } catch (error: unknown) {
    serverLog('LEAVE_COMPLAINT_ERROR', error, true)
    return null
  }
}

export async function toggleFollow(targerId: string): Promise<BaseServerResponse<null> | null> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`${endpoint}/${targerId}/follow`)

    return response
  } catch (error: unknown) {
    serverLog('TOGGLE_FOLLOW_ERROR', error, true)
    return null
  }
}

export async function checkFollowing(targerId: string): Promise<BaseServerResponse<{ isFollowing: boolean }> | null> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<{ isFollowing: boolean }>>(
      `${endpoint}/is-following/${targerId}`,
    )

    return response
  } catch (error: unknown) {
    serverLog('IS_FOLLOWING_ERROR', error, true)
    return null
  }
}
