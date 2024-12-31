<template>
  <div id="app">
    <div class="counters">
      <ValueCounter :value="stats?.total_summary?.downloaded" :duration="1000" unit="MiB" meaning="downloaded" />
      <ValueCounter :value="stats?.total_summary?.uploaded" :duration="1000" unit="MiB" meaning="uploaded" />
      <ValueCounter :value="stats?.total_summary?.seeding" :duration="500" unit="" meaning="seeding" />
    </div>
    <LineChart />
  </div>
</template>

<script lang="ts">
import LineChart from '@/components/charts/LineChart.vue'
import ValueCounter from '@/components/misc/ValueCounter.vue'
import { useApi } from '../composables/useApi'
import { ref, onMounted } from 'vue'

export default {
  components: {
    LineChart,
    ValueCounter,
  },
  setup() {
    const { getUserStats } = useApi()
    const stats = ref(null)
    const loading = ref(true)

    onMounted(async () => {
      try {
        getUserStats('2024-01-01', '2024-12-31', '0,1,2').then((res) => {
          // const response = res
          stats.value = res
          console.log(stats.value)
        })
      } catch (error) {
        console.error('Error fetching user stats:', error)
        stats.value = null
      } finally {
        loading.value = false
      }
    })

    return { stats, loading }
  },
}
</script>

<style scoped>
.counters {
  display: flex;
  justify-content: space-around;
}
</style>
