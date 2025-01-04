import { ref } from 'vue'
import { fetchData } from '../api/api.ts'

export function useApi<T>() {
  const data = ref<T | null>(null)
  const loading = ref<boolean>(false)
  const error = ref<T | null>(null)

  const getUserStats = async (date_from: string, date_to: string, tracker_ids: string) => {
    try {
      const response = await fetchData(`/stats?date_from=${date_from}&date_to=${date_to}&tracker_ids=${tracker_ids}`)
      data.value = response as T
      return response
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }
  const getTrackerMap = async () => {
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

  return { data, loading, error, getUserStats, getTrackerMap }
}
