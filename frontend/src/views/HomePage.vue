<template>
  <div id="home-page">
    <div class="logo">Dasharr</div>
    <div class="selectors">
      <FloatLabel>
        <DatePicker v-model="selectedPeriod" dateFormat="yy-mm-dd" showIcon fluid iconDisplay="input" selectionMode="range" />
        <label for="over_label">Period</label>
      </FloatLabel>
      <div class="buttons">
        <Button icon="pi pi-cog" class="settings-btn" @click="settingsDialog = true" />
        <Button type="button" label="Get Stats" :loading="loading" @click="fetchStats()" />
      </div>
    </div>
    <div class="global-counters section">
      <ValueCounter :value="stats ? (stats as any).total_summary.downloaded_amount : 0" :duration="1000" unit="GiB" meaning="downloaded" />
      <ValueCounter :value="stats ? (stats as any).total_summary.uploaded_amount : 0" :duration="1000" unit="GiB" meaning="uploaded" />
      <ValueCounter :value="stats ? (stats as any).total_summary.seeding : 0" :duration="1000" unit="torrents" meaning="seeding" />
    </div>
    <div class="section tracker-details" v-if="stats">
      <TrackerCard v-for="tracker in (stats as any).per_tracker_summary" :key="tracker.tracker_id" :trackerName="trackerMap[tracker.tracker_id]" :statsSummary="tracker" :statsDetailed="detailedStats(tracker.tracker_id)" />
    </div>
    <Dialog v-model:visible="settingsDialog" modal header="Settings"><SettingsDialog /></Dialog>
  </div>
</template>

<script lang="ts">
import ValueCounter from '@/components/charts/ValueCounter.vue'
import DatePicker from 'primevue/datepicker'
import { Button, FloatLabel } from 'primevue'
import { useApi } from '../composables/useApi'
import { ref, onMounted } from 'vue'
import TrackerCard from '@/components/charts/TrackerCard.vue'
import Dialog from 'primevue/dialog'
import SettingsDialog from '@/components/settings/SettingsDialog.vue'

export default {
  components: {
    ValueCounter,
    DatePicker,
    FloatLabel,
    // eslint-disable-next-line vue/no-reserved-component-names
    Button,
    // eslint-disable-next-line vue/no-reserved-component-names
    Dialog,
    TrackerCard,
    SettingsDialog,
  },
  data() {
    return {
      settingsDialog: false,
    }
  },
  methods: {
    detailedStats(tracker_id: number) {
      return this.stats ? this.stats.all.filter((stat: any) => stat.tracker_id === tracker_id) : []
    },
  },
  setup() {
    const { getUserStats, getTrackerMap } = useApi()
    const stats = ref<any>(null)
    const loading = ref(true)
    const selectedPeriod = ref([new Date(Date.now() - 7 * 24 * 60 * 60 * 1000), new Date()])
    const trackerMap = ref<any>({})

    const fetchTrackerMap = () => {
      getTrackerMap().then((res) => {
        trackerMap.value = res
      })
    }

    const fetchStats = () => {
      const enabledIndexers = JSON.parse(localStorage.getItem('enabledIndexers') ?? '[]')
      loading.value = true
      const date_from = selectedPeriod.value[0].toISOString().split('T')[0] + ' 00:00:00'
      const date_to = selectedPeriod.value[1].toISOString().split('T')[0] + ' 23:59:59'
      const tracker_ids = enabledIndexers.map((id: string) => id).join(',')
      getUserStats(date_from, date_to, tracker_ids)
        .then((res) => (stats.value = res))
        .finally(() => (loading.value = false))
    }

    onMounted(() => {
      fetchTrackerMap()
      fetchStats()
    })

    return { loading, stats, selectedPeriod, trackerMap, fetchStats }
  },
}
</script>

<style scoped lang="scss">
#home-page {
  margin: 2em;
  margin-top: 1em;
}
.logo {
  font-size: 2em;
  font-weight: bold;
  margin-bottom: 1em;
}
.selectors {
  display: flex;
  justify-content: center;
  margin-bottom: 55px;
  .settings-btn {
    margin-right: 7px;
    margin-left: 15px;
  }
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
