'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'

export type Reaction = {
  id: string
  name: string
  cost: number
  stickerUrl: string
}

export async function getReactionsAction(): Promise<BaseServerResponse<Reaction[]>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<Reaction[]>>('/reactions')
    return response
  } catch (error: unknown) {
    serverLog('GET_REACTIONS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve reactions',
      errors: parseAxiosError(error),
    }
  }
}

export async function sendReactionAction(roomId: string, reactionId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`/rooms/${roomId}/react`, {
      reactionId,
    })

    return response
  } catch (error: unknown) {
    serverLog('SEND_REACTION_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to send reaction',
      errors: parseAxiosError(error),
    }
  }
}
