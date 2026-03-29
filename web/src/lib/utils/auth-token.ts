import { Session } from 'next-auth'

let globalAccessToken: string | null = null
let updateSessionFn: (() => Promise<Session | null>) | null = null

export const setGlobalToken = (token: string | null) => {
  globalAccessToken = token
}

export const getGlobalToken = () => globalAccessToken

export const setUpdateSession = (fn: () => Promise<Session | null>) => {
  updateSessionFn = fn
}

export const triggerSessionUpdate = async () => {
  if (updateSessionFn) {
    const newSession = await updateSessionFn()
    return newSession
  }
  return null
}
