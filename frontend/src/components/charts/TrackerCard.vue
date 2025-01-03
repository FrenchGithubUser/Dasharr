<template>
  <Card class="tracker-card">
    <template #content>
      <img :src="'/images/' + statsSummary.tracker_id + '.png'" width="200px" class="logo" />
      <div class="counters">
        <ValueCounter v-for="(stat, label) in statsToDisplay" :key="label" :value="stat" :meaning="label" :unit="['uploaded_amount', 'downloaded_amount'].indexOf(label) > -1 ? 'GiB' : ''" :duration="500" />
      </div>
      <div class="graphs">
        <LineChart :series="uploadDetail.series" :xaxis="uploadDetail.xaxis" />
      </div>
    </template>
  </Card>
</template>

<script lang="ts">
import Card from 'primevue/card'
import ValueCounter from './ValueCounter.vue'
import LineChart from './LineChart.vue'

export default {
  components: {
    Card,
    ValueCounter,
    LineChart,
  },
  props: {
    statsSummary: {
      type: Object,
      required: true,
    },
    statsDetailed: {
      type: Array,
      required: true,
      default: () => {
        return []
      },
    },
  },
  computed: {
    statsToDisplay() {
      return Object.fromEntries(Object.entries(this.statsSummary).filter(([label, value]) => typeof value === 'number' && label != 'tracker_id'))
    },
    uploadDetail() {
      return { series: [{ data: this.statsDetailed.map((stat: object) => stat.uploaded_amount.toFixed(1)), name: 'uploaded_amount' }], xaxis: { type: 'datetime', categories: this.statsDetailed.map((stat: object) => stat.collected_at) } }
    },
  },
  setup() {},
}
</script>

<style scoped>
.counters {
  display: flex;
  justify-content: space-around;
}
.logo {
  margin-bottom: 10px;
}
</style>
