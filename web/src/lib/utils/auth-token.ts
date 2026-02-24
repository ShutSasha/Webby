let globalAccessToken: string | null = null

export const setGlobalToken = (token: string | null) => {
  globalAccessToken = token
}

export const getGlobalToken = () => globalAccessToken
