import type { AxiosInstance } from 'axios' // Import AxiosInstance as a type-only import
import axios from 'axios'
import router from '../router'

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL

const apiClient: AxiosInstance = axios.create({
  baseURL: BACKEND_URL,
  headers: {
    'X-API-Key': localStorage.getItem('api-key'),
  },
})

axios.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response && error.response.status !== 200) {
      router.push({ name: 'RegisterPage' })
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
