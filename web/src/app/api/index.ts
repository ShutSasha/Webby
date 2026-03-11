import axios from 'axios'
import { cookies } from 'next/headers'

import { clog } from '@/lib/utils/utils'

export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:5000/api'

const $api = axios.create({
  baseURL: API_URL,
})

$api.interceptors.request.use(async config => {
  clog('req SERVER HTTP headers', config.headers)

  return config
})

// TODO: catch 401 or error with invalid accessToken for signOut method
$api.interceptors.response.use(
  config => {
    return config
  },
  async error => {
    // interceptor response logic
    throw error
  },
)

export default $api
