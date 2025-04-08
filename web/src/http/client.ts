import axios from 'axios'

let unauthorizedErrorHandler: (() => void) | undefined = undefined

export function setUnauthorizedErrorHandler(callback: typeof unauthorizedErrorHandler) {
  unauthorizedErrorHandler = callback
}

var httpClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 1000,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
})

httpClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response.status === 401) {
      unauthorizedErrorHandler?.()
      return
    }
    return Promise.reject(error)
  },
)

export default httpClient
