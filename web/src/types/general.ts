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

export { type FormActionState, type BaseServerResponse }
