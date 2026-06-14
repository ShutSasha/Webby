/* eslint-disable */
import { BaseServerResponse } from '@/types/general.types'
import axios, { AxiosError } from 'axios'
import clsx from 'clsx'
import { ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function parseAxiosError(error: unknown): Record<string, string> {
  if (!axios.isAxiosError<BaseServerResponse>(error)) {
    return { global: 'Unexpected server error' }
  }

  const data = error.response?.data

  if (!data) {
    return { global: 'No response from server' }
  }

  if (data.errors && Object.keys(data.errors).length > 0) {
    return data.errors
  }

  if (data.message) {
    return { global: data.message }
  }

  return { global: 'Unknown server error' }
}

export function extractServerMessage(errors: Record<string, string> | null | undefined): string | undefined {
  if (!errors) return undefined

  return errors.global || errors.message || Object.values(errors)[0]
}

export const serverLog = (label: string, error: any, detailed?: boolean) => {
  if (!(error instanceof AxiosError)) {
    console.log(`\n=== [${label}] Non-Axios Error ===`)
    console.log(error)
    console.log(`============\n`)
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

  console.log(`============\n`)

  if (detailed) {
    console.log(`--- [${label}] DETAILED ---`)
    const { stack, ...props } = error.toJSON() as any

    console.log({
      url: error.config?.url,
      method: error.config?.method,
      ...props,
    })
    console.log(`------------\n`)
  }
}

// console.log()
export function clog(label: string, data?: any) {
  if (!data) {
    console.log(`\n=== [${label}] ===`)
    return
  }

  console.log(`\n=== [${label}] ===`)
  console.log(data)
  console.log(`============\n`)
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export class ServerActionError extends Error {
  errors?: Record<string, string> | null

  constructor(message: string, errors?: Record<string, string> | null) {
    super(message)
    this.errors = errors
    this.name = 'ServerActionError'
  }
}

export const unwrapServerAction = <T>(response: BaseServerResponse<T>): T => {
  if (!response.success) {
    throw new ServerActionError(response.message || 'An error occurred', response.errors)
  }

  return response.data as T
}
