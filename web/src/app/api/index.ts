import axios from 'axios'
import { cookies } from 'next/headers'

export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:5000/api'

const $api = axios.create({
  baseURL: API_URL,
})

$api.interceptors.request.use(async config => {
  const token = (await cookies()).get('authjs.session-token')?.value

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  console.log('----- req SERVER HTTP headers -----', config.headers)

  return config
})

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
