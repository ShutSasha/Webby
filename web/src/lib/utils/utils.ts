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

  console.log(`\n=== [${label}] ===`)
  console.dir(error.response?.data, { depth: null, colors: true })
  console.log(`=== === === ===\n`)

  if (detailed) {
    console.log(`--- [${label}] DETAILED ---`)
    // Log all properties except 'stack'
    const { stack, ...props } = error.toJSON() as any
    console.log(props)
    console.log(`--- --- --- ---\n`)
  }
}
