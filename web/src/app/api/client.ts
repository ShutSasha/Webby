import axios from 'axios'

import { getGlobalToken } from '@/lib/utils/auth-token'
import { clog } from '@/lib/utils/utils'

const $apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
})

$apiClient.interceptors.request.use(async config => {
  const token = getGlobalToken()

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  clog('req CLIENT HTTP headers', config.headers)

  return config
})

export default $apiClient
