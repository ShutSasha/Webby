/* eslint-disable */
import { AxiosError } from 'axios'

export type BaseServerResponse = {
  data: Record<string, any> | null
  errors: Record<string, string> | null
  message: string
  success: boolean
}

export const serverLog = (label: string, error: any, detailed?: boolean) => {
  if (!(error instanceof AxiosError)) {
    console.log(`\n=== [${label}] Non-Axios Error ===`)
    console.log(error)
    console.log(`=== === === ===\n`)
    return
  }

  const responseData = error.response?.data

  console.log(`\n=== [${label}] ===`)
  if (typeof responseData === 'string' && responseData.includes('<!DOCTYPE html>')) {
    console.log(`Server returned HTML (likely 404 or 500 error)`)
    console.log(`Status: ${error.response?.status} ${error.response?.statusText}`)

    console.log(`Snippet: ${responseData.substring(0, 200)}...`)
  } else {
    console.dir(responseData, { depth: null, colors: true })
  }

  console.log(`=== === === ===\n`)

  if (detailed) {
    console.log(`--- [${label}] DETAILED ---`)
    const { stack, ...props } = error.toJSON() as any

    console.log({
      url: error.config?.url,
      method: error.config?.method,
      ...props,
    })
    console.log(`--- --- --- ---\n`)
  }
}
