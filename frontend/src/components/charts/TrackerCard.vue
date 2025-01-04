<template>
  <Card class="tracker-card">
    <template #content>
      <div class="logo-wrapper">
        <img :src="'/images/' + trackerName + '.png'" width="250px" class="logo" :alt="trackerName" />
      </div>
      <Card>
        <template #content>
          <div class="explanation">Amounts increase during the selected period</div>
          <div class="counters">
            <ValueCounter v-for="(stat, label) in statsToDisplay" :key="label" :value="stat" :meaning="label" :unit="['uploaded_amount', 'downloaded_amount'].indexOf(label.toString()) > -1 ? 'GiB' : ''" :duration="500" />
          </div>
        </template>
      </Card>
      <Card class="graphs">
        <template #content>
          <div class="explanation">Amounts evolution during the selected period</div>
          <LineChart :series="uploadDetail.series" :xaxis="uploadDetail.xaxis" />
        </template>
      </Card>
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
    trackerName: {
      type: String,
    },
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
      return {
        series: [
          {
            data: this.statsDetailed.map((stat: any) => stat.uploaded_amount.toFixed(1)),
            name: 'uploaded_amount',
          },
        ],
        xaxis: {
          type: 'datetime',
          categories: this.statsDetailed.map((stat: any) => stat.collected_at),
        },
      }
    },
  },
  setup() {},
}
</script>

<style scoped>
.logo-wrapper {
  text-align: center;
}
.logo {
  margin-bottom: 10px;
}
.counters {
  display: flex;
  justify-content: space-around;
}
.explanation {
  margin-bottom: 10px;
  font-weight: bold;
}
.graphs {
  margin-top: 20px;
}
</style>
