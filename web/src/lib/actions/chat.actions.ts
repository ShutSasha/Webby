'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse, PaginatedData } from '@/types/general.types'

const endpoint = '/chats'

export async function sendMessageAction(chatId: string, content: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<null>>(`${endpoint}/${chatId}/messages`, { content })

    return response
  } catch (error: unknown) {
    serverLog('SEND_MESSAGE_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to send message',
      errors: parseAxiosError(error),
    }
  }
}

export type ChatUser = {
  id: string
  username: string
  avatarUrl: string
}

export type ChatMessage = {
  id: string
  sender: ChatUser
  content: string
  isEdited: boolean
  createdAt: string
}

export async function getChatMessagesAction(
  chatId: string,
  page: number = 1,
  limit: number = 20,
): Promise<BaseServerResponse<PaginatedData<ChatMessage>>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<PaginatedData<ChatMessage>>>(
      `/chats/${chatId}/messages?page=${page}&limit=${limit}`,
    )
    return response
  } catch (error: unknown) {
    serverLog('GET_CHAT_MESSAGES_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve chat messages',
      errors: parseAxiosError(error),
    }
  }
}
