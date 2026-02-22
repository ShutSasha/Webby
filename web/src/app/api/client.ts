import axios from 'axios'

import { clog } from '@/lib/utils/utils'

import { getEncryptedToken } from './auth'

const $apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
})

$apiClient.interceptors.request.use(async config => {
  const token = await getEncryptedToken()

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  clog('req CLIENT HTTP headers', config.headers)

  return config
})

export default $apiClient
