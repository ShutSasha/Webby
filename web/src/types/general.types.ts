export type Optional<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>

interface FormActionState {
  success: boolean
  errors: Record<string, string> | null
}

interface BaseServerResponse<T = any> {
  data: T | null
  errors: Record<string, string> | null
  message: string
  success: boolean
}

type PaginatedData<T> = {
  items: T[]
  page: number
  pageSize: number
  totalCount: number
}

type Platform = 'Webby' | 'Youtube' | 'Twitch'

export { type FormActionState, type BaseServerResponse, type PaginatedData, type Platform }
