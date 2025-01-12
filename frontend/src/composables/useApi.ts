import { ref } from 'vue'
import { fetchData, sendData } from '../api/api.ts'

export function useApi<T>() {
  const data = ref<T | null>(null)
  const loading = ref<boolean>(false)
  const error = ref<T | null>(null)

  const getUserStats = async (date_from: string, date_to: string, indexer_ids: string) => {
    try {
      const response = await fetchData(`/stats?date_from=${date_from}&date_to=${date_to}&indexer_ids=${indexer_ids}`)
      data.value = response as T
      return response
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }
  //todo : call this on first page load and make the result accessible everywhere instead
  const getIndexerMap = async () => {
    try {
      const response = await fetchData(`/prowlarrConfig`)
      data.value = response as T
      return response
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }
  const getConfig = async () => {
    try {
      const response = await fetchData(`/config`)
      data.value = response as T
      return response
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }
  const saveCredentials = async (jsonBody: object) => {
    try {
      const response = await sendData(`/saveCredentials`, jsonBody)
      data.value = response as T
      return response
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }
  const savedCredentials = async () => {
    try {
      const response = await fetchData(`/savedCredentials`)
      data.value = response as T
      return response
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }

  return { data, loading, error, getUserStats, getIndexerMap, getConfig, saveCredentials, savedCredentials }
}
