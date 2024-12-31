<template>
  <apexchart :series="series" :options="chartOptions" type="line" height="350"></apexchart>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import VueApexCharts from 'vue3-apexcharts'

export default defineComponent({
  name: 'LineChart',
  components: {
    apexchart: VueApexCharts,
  },
  setup() {
    const series = ref([
      {
        name: 'Desktops',
        data: [10, 41, 35, 51, 49, 62, 69, 91, 148],
      },
    ])

    const chartOptions = ref({
      chart: {
        height: 350,
        type: 'line',
        zoom: {
          enabled: false,
        },
      },
      dataLabels: {
        enabled: false,
      },
      stroke: {
        curve: 'straight',
      },
      title: {
        text: 'Product Trends by Month',
        align: 'left',
      },
      grid: {
        row: {
          colors: ['#f3f3f3', 'transparent'],
          opacity: 0.5,
        },
      },
      xaxis: {
        categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep'],
      },
    })

    return { series, chartOptions }
  },
})
</script>
