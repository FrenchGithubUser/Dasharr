import type { AxiosInstance } from 'axios' // Import AxiosInstance as a type-only import
import axios from 'axios'
import router from '../router'

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || window.location.origin + '/api'

const apiClient: AxiosInstance = axios.create({
  baseURL: BACKEND_URL,
})

apiClient.interceptors.request.use((config) => {
  config.headers['X-API-Key'] = localStorage.getItem('api-key')
  return config
})

apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response.status !== 200 && error.response.data.message == 'invalid_api_key') {
      router.push({ name: 'Register' })
    }
    return Promise.reject(error)
  },
)

export const fetchData = async (endpoint: string) => {
  try {
    const response = await apiClient.get(endpoint)
    return response.data
  } catch (error) {
    throw error
  }
}

export const sendData = async (endpoint: string, jsonBody: object) => {
  try {
    const response = await apiClient.post(endpoint, jsonBody)
    return response.data
  } catch (error) {
    throw error
  }
}
