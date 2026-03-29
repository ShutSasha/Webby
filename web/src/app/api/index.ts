import axios from 'axios'

import { getGlobalToken, setGlobalToken, triggerSessionUpdate } from '@/lib/utils/auth-token'
import { auth } from '@/workspace/auth'

export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:5000/api'

const $api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

$api.interceptors.request.use(async config => {
  let token: string | null = null

  if (typeof window === 'undefined') {
    const session = await auth()
    token = session?.user?.accessToken || null
  } else {
    token = getGlobalToken()
  }

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
})

$api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      const newSession = await triggerSessionUpdate()
      const newToken = newSession?.user?.accessToken

      if (newToken) {
        setGlobalToken(newToken)

        originalRequest.headers.Authorization = `Bearer ${newToken}`

        return $api(originalRequest)
      }
    }

    return Promise.reject(error)
  },
)

export default $api
