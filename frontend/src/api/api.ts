import type { AxiosInstance } from 'axios' // Import AxiosInstance as a type-only import
import axios from 'axios'

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL

const apiClient: AxiosInstance = axios.create({
  baseURL: BACKEND_URL,
  // headers: {
  //   Authorization: `Bearer ${API_TOKEN}`,
  // },
})

export const fetchData = async (endpoint: string) => {
  try {
    const response = await apiClient.get(endpoint)
    return response.data
  } catch (error) {
    throw error
  }
}
