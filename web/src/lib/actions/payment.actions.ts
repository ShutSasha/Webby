'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { BaseServerResponse } from '@/types/general.types'

const endpoint = '/payments'

export type SubscriptionDetails = {
  expirationDate: string
  username: string
}

export async function initializePaymentAction(): Promise<BaseServerResponse<string>> {
  try {
    const { data: response } = await $api.post<BaseServerResponse<string>>(`${endpoint}`)

    return response
  } catch (error: unknown) {
    serverLog('INITIALIZE_PAYMENT_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to initialize payment',
      errors: parseAxiosError(error),
    }
  }
}

export async function getSubscriptionDetailsAction(
  paymentId: string,
): Promise<BaseServerResponse<SubscriptionDetails>> {
  try {
    const { data: response } = await $api.get<BaseServerResponse<SubscriptionDetails>>(
      `${endpoint}/subscriptions/${paymentId}`,
    )

    return response
  } catch (error: unknown) {
    serverLog('GET_SUBSCRIPTION_DETAILS_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve subscription details',
      errors: parseAxiosError(error),
    }
  }
}
