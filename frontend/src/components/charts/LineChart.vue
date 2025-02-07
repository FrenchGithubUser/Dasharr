<template>
  <apexchart height="350" ref="chart" :options="chartOptions" :series="series" class="line-chart"></apexchart>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
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
    label: {
      type: String,
      required: true,
    },
  },
  data(props) {
    return {
      chartOptions: {
        chart: {
          height: 350,
          type: 'area',
          toolbar: {
            show: false,
          },
        },
        dataLabels: {
          enabled: false,
        },
        stroke: {
          curve: 'smooth',
        },
        title: {
          text: props.label,
          align: 'center',
        },

        xaxis: this.xaxis,
        // tooltip: {
        //   x: {
        //     format: 'yy-MM-ddTHH:mm:ss',
        //   },
        // },
        theme: {
          mode: 'dark',
          palette: 'palette8',
          monochrome: {
            enabled: true,
            color: '#F86624',
            shadeTo: 'dark',
            shadeIntensity: 0.95,
          },
        },
      },
    }
  },
  setup() {},
  watch: {
    xaxis: function (newVal) {
      ;(this.$refs.chart as typeof VueApexCharts).updateOptions({
        xaxis: newVal,
      })
    },
  },
})
</script>
