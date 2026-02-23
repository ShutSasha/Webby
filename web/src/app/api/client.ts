import axios from 'axios'
import { getSession } from 'next-auth/react'

import { clog } from '@/lib/utils/utils'

const $apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
})

$apiClient.interceptors.request.use(async config => {
  const session = await getSession()

  if (session?.user?.accessToken) {
    config.headers.Authorization = `Bearer ${session.user.accessToken}`
  }

  clog('req CLIENT HTTP headers', config.headers)

  return config
})

export default $apiClient
