import axios from 'axios'

export const API_URL = 'http://localhost:5000/api'

const $api = axios.create({
  baseURL: API_URL,
})

$api.interceptors.request.use(config => {
  // interceptor request logic
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
