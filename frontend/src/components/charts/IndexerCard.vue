<template>
  <Card class="indexer-card">
    <template #content>
      <div class="logo-wrapper">
        <img :src="'/images/' + indexerName + '.png'" width="250px" class="logo" :alt="indexerName" />
      </div>
      <Card>
        <template #content>
          <div class="explanation">Amounts increase during the selected period</div>
          <div class="counters">
            <ValueCounter v-for="(stat, label) in statsToDisplay" :key="label" :value="stat" :meaning="label" :unit="['uploaded_amount', 'downloaded_amount'].indexOf(label.toString()) > -1 ? 'GiB' : ''" :duration="500" />
          </div>
        </template>
      </Card>
      <Card class="graphs-card">
        <template #content>
          <div class="explanation">Amounts evolution during the selected period</div>
          <div class="graphs">
            <LineChart :series="uploadDetail.series" :xaxis="uploadDetail.xaxis" label="uploaded_amount" />
            <LineChart :series="downloadDetail.series" :xaxis="downloadDetail.xaxis" label="downloaded_amount" />
          </div>
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
    indexerName: {
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
      return Object.fromEntries(Object.entries(this.statsSummary).filter(([label, value]) => typeof value === 'number' && label != 'indexer_id'))
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
    downloadDetail() {
      return {
        series: [
          {
            data: this.statsDetailed.map((stat: any) => stat.downloaded_amount.toFixed(1)),
            name: 'downloaded_amount',
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

<style scoped lang="scss">
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
.graphs-card {
  margin-top: 20px;
  .graphs {
    display: flex;
  }
  .line-chart {
    width: 50%;
  }
}
</style>
