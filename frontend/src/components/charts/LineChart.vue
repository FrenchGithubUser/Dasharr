<template>
  <apexchart height="350" ref="chart" :options="chartOptions" :series="series"></apexchart>
</template>

<script lang="ts">
import { defineComponent, watch } from 'vue'
import VueApexCharts from 'vue3-apexcharts'

export default defineComponent({
  name: 'LineChart',
  components: {
    apexchart: VueApexCharts,
  },
  props: {
    series: {
      type: Array,
      required: true,
    },
    xaxis: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      chartOptions: {
        chart: {
          height: 350,
          type: 'area',
        },
        dataLabels: {
          enabled: false,
        },
        stroke: {
          curve: 'smooth',
        },
        xaxis: this.xaxis,
        // tooltip: {
        //   x: {
        //     format: 'yy-MM-ddTHH:mm:ss',
        //   },
        // },
      },
    }
  },
  setup() {},
  watch: {
    xaxis: function (newVal) {
      this.$refs.chart.updateOptions({
        xaxis: newVal,
      })
    },
  },
})
</script>
