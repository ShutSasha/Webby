'use server'

import $api from '@/lib/config/api.config'
import { parseAxiosError, serverLog } from '@/lib/utils/general.utils'
import { Category } from '@/types/category.types'
import { BaseServerResponse, PaginatedData } from '@/types/general.types'

const endpoint = '/categories'

type GetCategoriesResponse = PaginatedData<Category>

export async function getCategories(
  query?: string,
  page?: number,
  pageSize?: number,
): Promise<BaseServerResponse<GetCategoriesResponse>> {
  try {
    const params = new URLSearchParams()

    if (query) params.append('search', query)
    if (page) params.append('page', page.toString())
    if (pageSize) params.append('limit', pageSize.toString())

    const queryString = params.toString()

    const url = queryString ? `${endpoint}?${queryString}` : endpoint

    const { data: response } = await $api.get<BaseServerResponse<GetCategoriesResponse>>(url)

    return response
  } catch (error: unknown) {
    serverLog('GET_CATEGORIES_ERROR', error, true)

    return {
      data: null,
      success: false,
      message: 'Failed to retrieve categories',
      errors: parseAxiosError(error),
    }
  }
}
