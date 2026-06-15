'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'

export type CreateRightChoiceVotePayload = {
  voteText: string
  duration: number
  choices: string[]
}

export async function createRightChoiceVoteAction(
  roomId: string,
  payload: CreateRightChoiceVotePayload,
): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/rooms/${roomId}/votes/right-choice`, payload)

    return response
  } catch (error: unknown) {
    serverLog('CREATE_RIGHT_CHOICE_VOTE_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to create vote',
      errors: parseAxiosError(error),
    }
  }
}

export type RoomVote = {
  id: string
  voteText: string
  duration: number
  createdAt: string
  choices: string[]
  isLocked?: boolean
  rightChoice?: string
  winners?: string[]
}

export async function getRoomVotesAction(roomId: string): Promise<BaseServerResponse<RoomVote[]>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<RoomVote[]>>(`/rooms/${roomId}/votes`)

    return response
  } catch (error: unknown) {
    serverLog('GET_ROOM_VOTES_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve room votes',
      errors: parseAxiosError(error),
    }
  }
}

export async function castVoteAction(
  roomId: string,
  voteId: string,
  choice: string,
): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/rooms/${roomId}/votes/${voteId}/vote`, {
      choice,
    })

    return response
  } catch (error: unknown) {
    serverLog('CAST_VOTE_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to cast vote',
      errors: parseAxiosError(error),
    }
  }
}

export async function resolveVoteAction(
  roomId: string,
  voteId: string,
  rightChoice: string,
): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/rooms/${roomId}/votes/${voteId}/resolve`, {
      rightChoice,
    })

    return response
  } catch (error: unknown) {
    serverLog('RESOLVE_VOTE_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to resolve vote',
      errors: parseAxiosError(error),
    }
  }
}
