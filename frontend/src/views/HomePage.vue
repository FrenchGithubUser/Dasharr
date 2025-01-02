<template>
  <div id="home-page">
    <div class="selectors section">
      <FloatLabel>
        <DatePicker v-model="selectedPeriod" dateFormat="yy-mm-dd" showIcon fluid iconDisplay="input" selectionMode="range" />
        <label for="over_label">Period</label>
      </FloatLabel>
      <Button type="button" label="Get Stats" :loading="loading" @click="fetchStats()" />
    </div>
    <div class="counters section">
      <ValueCounter :value="stats?.total_summary.downloaded" :duration="1000" unit="GiB" :data="true" meaning="downloaded" />
      <ValueCounter :value="stats?.total_summary.uploaded" :duration="1000" unit="GiB" :data="true" meaning="uploaded" />
      <ValueCounter :value="stats?.total_summary.seeding" :duration="1000" unit="torrents" :data="false" meaning="seeding" />
    </div>
  </div>
</template>

<script lang="ts">
import ValueCounter from '@/components/misc/ValueCounter.vue'
import DatePicker from 'primevue/datepicker'
import { Button, FloatLabel } from 'primevue'
import { useApi } from '../composables/useApi'
import { ref, onMounted } from 'vue'

export default {
  components: {
    ValueCounter,
    DatePicker,
    FloatLabel,
    // eslint-disable-next-line vue/no-reserved-component-names
    Button,
  },
  data() {
    return {
      loading: false,
    }
  },
  setup() {
    const { getUserStats } = useApi()
    const stats = ref(null)
    const loading = ref(true)
    const selectedPeriod = ref([new Date(Date.now() - 7 * 24 * 60 * 60 * 1000), new Date()])
    const selectedTrackers = ref([{ id: 1 }, { id: 2 }, { id: 3 }])

    const fetchStats = () => {
      const date_from = selectedPeriod.value[0].toISOString().split('T')[0]
      const date_to = selectedPeriod.value[1].toISOString().split('T')[0]
      const tracker_ids = selectedTrackers.value.map((tracker) => tracker.id).join(',')
      getUserStats(date_from, date_to, tracker_ids)
        .then((res) => (stats.value = res))
        .finally(() => (loading.value = false))
    }

    onMounted(() => {
      fetchStats()
    })

    return { stats, selectedPeriod, selectedTrackers, fetchStats }
  },
}
</script>

<style scoped>
#home-page {
  margin: 2em;
}
.section {
  display: flex;
  justify-content: space-around;
  margin-bottom: 1.5em;
}
</style>
