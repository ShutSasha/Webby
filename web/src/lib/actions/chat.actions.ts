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

export type ChatHistoryItem = {
  chatId: string
  user: ChatUser
  lastMessage: {
    content: string
    createdAt: string
  } | null
}

export async function getChatHistoryAction(
  page: number = 1,
  limit: number = 10,
  search?: string,
): Promise<BaseServerResponse<PaginatedData<ChatHistoryItem>>> {
  try {
    const queryParams = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    })

    if (search) {
      queryParams.append('search', search)
    }

    const { data: response } = await $api.get<BaseServerResponse<PaginatedData<ChatHistoryItem>>>(
      `/chats?${queryParams.toString()}`,
    )

    return response
  } catch (error: unknown) {
    serverLog('GET_CHAT_HISTORY_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve chat history',
      errors: parseAxiosError(error),
    }
  }
}

export type ChatDetails = {
  id: string
  createdAt: string
  user: ChatUser
}

export async function getChatByIdAction(chatId: string): Promise<BaseServerResponse<ChatDetails>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<ChatDetails>>(`${endpoint}/${chatId}`)

    return response
  } catch (error: unknown) {
    serverLog('GET_CHAT_BY_ID_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve chat details',
      errors: parseAxiosError(error),
    }
  }
}

export async function deleteChatAction(chatId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.delete<BaseServerResponse<null>>(`${endpoint}/${chatId}`)

    return response
  } catch (error: unknown) {
    serverLog('DELETE_CHAT_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to delete chat',
      errors: parseAxiosError(error),
    }
  }
}

export type CreateChatResponse = {
  id: string
  createdAt: string
}

export async function createPrivateChatAction(targetId: string): Promise<BaseServerResponse<CreateChatResponse>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<CreateChatResponse>>('/chats', {
      targetId,
    })

    return response
  } catch (error: unknown) {
    serverLog('CREATE_CHAT_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to create chat',
      errors: parseAxiosError(error),
    }
  }
}

export async function deleteMessageAction(chatId: string, messageId: string): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.delete<BaseServerResponse<null>>(
      `${endpoint}/${chatId}/messages/${messageId}`,
    )

    return response
  } catch (error: unknown) {
    serverLog('DELETE_MESSAGE_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to delete message',
      errors: parseAxiosError(error),
    }
  }
}

export async function editMessageAction(
  chatId: string,
  messageId: string,
  content: string,
): Promise<BaseServerResponse<null>> {
  try {
    const { data: response } = await $api.patch<BaseServerResponse<null>>(
      `${endpoint}/${chatId}/messages/${messageId}`,
      { content },
    )

    return response
  } catch (error: unknown) {
    serverLog('EDIT_MESSAGE_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to edit message',
      errors: parseAxiosError(error),
    }
  }
}
