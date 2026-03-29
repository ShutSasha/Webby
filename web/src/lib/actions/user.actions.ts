'use server'

import $api from '@/lib/api/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/utils'
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

export async function getUserFollows(id: string): Promise<BaseServerResponse<GetUserFollowsResponse[]>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<GetUserFollowsResponse[]>>(`${endpoint}/${id}/follows`)

    return response
  } catch (error: unknown) {
    serverLog('FETCH_USER_FOLLOWS_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: `Failed to get user follows`,
      errors: parseAxiosError(error),
    }
  }
}

export async function getUserFollowers(id: string): Promise<BaseServerResponse<GetUserFollowersResponse[]>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<GetUserFollowersResponse[]>>(
      `${endpoint}/${id}/followers`,
    )

    return response
  } catch (error: unknown) {
    serverLog('FETCH_USER_FOLLOWERS_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: `Failed to get user followers`,
      errors: parseAxiosError(error),
    }
  }
}

type TargetType = 'Video' | 'User'

export async function leaveComplaint(
  targetUserId: string,
  reasonType: string,
  targetType: TargetType,
  additionalInfo: string,
): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/complaints`, {
      targetId: targetUserId,
      reasonType,
      targetType,
      additionalInfo,
    })

    return response
  } catch (error: unknown) {
    serverLog('LEAVE_COMPLAINT_ERROR', error, true)

    return {
      success: false,
      data: null,
      message: 'Create complaint error',
      errors: parseAxiosError(error),
    }
  }
}

export async function toggleFollow(targerId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`${endpoint}/${targerId}/follow`)

    return response
  } catch (error: unknown) {
    serverLog('TOGGLE_FOLLOW_ERROR', error, true)
    return {
      success: false,
      data: null,
      message: 'Toggle follow error',
      errors: parseAxiosError(error),
    }
  }
}

export async function checkFollowing(targerId: string): Promise<BaseServerResponse<{ isFollowing: boolean }>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<{ isFollowing: boolean }>>(
      `${endpoint}/is-following/${targerId}`,
    )

    return response
  } catch (error: unknown) {
    serverLog('IS_FOLLOWING_ERROR', error, true)
    return {
      data: null,
      success: false,
      message: 'Check follow error',
      errors: parseAxiosError(error),
    }
  }
}

export async function updateAboutField(
  userId: string,
  about: string,
): Promise<BaseServerResponse<User> | { success: boolean; errors: Record<string, string> }> {
  try {
    const { data: response } = await $api.patch<BaseServerResponse<User>>(endpoint, {
      userId,
      about,
    })

    return response
  } catch (error: unknown) {
    serverLog('UPDATE_ABOUT_ERROR', error, true)
    return {
      success: false,
      errors: parseAxiosError(error),
    }
  }
}
