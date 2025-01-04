<template>
  <div id="home-page">
    <div class="selectors section">
      <FloatLabel>
        <DatePicker v-model="selectedPeriod" dateFormat="yy-mm-dd" showIcon fluid iconDisplay="input" selectionMode="range" />
        <label for="over_label">Period</label>
      </FloatLabel>
      <Button type="button" label="Get Stats" :loading="loading" @click="fetchStats()" />
    </div>
    <div class="global-counters section">
      <ValueCounter :value="stats ? stats.total_summary.downloaded_amount : 0" :duration="1000" unit="GiB" meaning="downloaded" />
      <ValueCounter :value="stats ? stats.total_summary.uploaded_amount : 0" :duration="1000" unit="GiB" meaning="uploaded" />
      <ValueCounter :value="stats ? stats.total_summary.seeding : 0" :duration="1000" unit="torrents" meaning="seeding" />
    </div>
    <div class="section tracker-details">
      <TrackerCard v-for="tracker in stats?.per_tracker_summary" :key="tracker.tracker_id" :trackerName="trackerMap[tracker.tracker_id]" :statsSummary="tracker" :statsDetailed="detailedStats(tracker.tracker_id)" />
    </div>
  </div>
</template>

<script lang="ts">
import ValueCounter from '@/components/charts/ValueCounter.vue'
import DatePicker from 'primevue/datepicker'
import { Button, FloatLabel } from 'primevue'
import { useApi } from '../composables/useApi'
import { ref, onMounted } from 'vue'
import TrackerCard from '@/components/charts/TrackerCard.vue'

export default {
  components: {
    ValueCounter,
    DatePicker,
    FloatLabel,
    // eslint-disable-next-line vue/no-reserved-component-names
    Button,
    TrackerCard,
  },
  data() {
    return {
      loading: false,
    }
  },
  methods: {
    detailedStats(tracker_id) {
      return this.stats ? this.stats.all.filter((stat: object) => stat.tracker_id === tracker_id) : []
    },
  },
  setup() {
    const { getUserStats, getTrackerMap } = useApi()
    const stats = ref(null)
    const loading = ref(true)
    const selectedPeriod = ref([new Date(Date.now() - 7 * 24 * 60 * 60 * 1000), new Date()])
    const selectedTrackers = ref([{ id: 30 }, { id: 32 }, { id: 19 }, { id: 2 }, { id: 5 }, { id: 62 }])
    const trackerMap = ref({})

    const fetchTrackerMap = () => {
      getTrackerMap().then((res) => {
        trackerMap.value = res
      })
    }

    const fetchStats = () => {
      const date_from = selectedPeriod.value[0].toISOString().split('T')[0] + ' 00:00:00'
      const date_to = selectedPeriod.value[1].toISOString().split('T')[0] + ' 23:59:59'
      const tracker_ids = selectedTrackers.value.map((tracker) => tracker.id).join(',')
      getUserStats(date_from, date_to, tracker_ids)
        .then((res) => (stats.value = res))
        .finally(() => (loading.value = false))
    }

    onMounted(() => {
      fetchTrackerMap()
      fetchStats()
    })

    return { stats, selectedPeriod, selectedTrackers, trackerMap, fetchStats }
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
.tracker-details {
  flex-direction: column;
}
.tracker-card {
  margin: 1em 0;
}
</style>
