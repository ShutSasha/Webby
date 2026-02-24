let globalAccessToken: string | null = null
let updateSessionFn: (() => Promise<any>) | null = null

export const setGlobalToken = (token: string | null) => {
  globalAccessToken = token
}

export const getGlobalToken = () => globalAccessToken

export const setUpdateSession = (fn: () => Promise<any>) => {
  updateSessionFn = fn
}

export const triggerSessionUpdate = async () => {
  if (updateSessionFn) {
    const newSession = await updateSessionFn()
    return newSession
  }
  return null
}
