<template>
  <Card class="value-counter">
    <template #content>
      <!-- counter too laggy
      TODO: make it smoother -->
      <!-- <div>{{ value % 1 === 0 ? currentValue.toFixed(0) : currentValue.toFixed(2) }} {{ unit }}</div> -->
      <div>{{ value?.toFixed(2) }}</div>
      <div>{{ meaning }}</div>
    </template>
  </Card>
</template>

<script lang="ts">
import { watch, ref, onMounted } from 'vue'
import Card from 'primevue/card'

export default {
  components: {
    Card,
  },
  props: {
    value: {
      type: Number,
      required: true,
    },
    duration: {
      type: Number,
      required: true,
    },
    unit: {
      type: String,
      required: false,
    },
    meaning: {
      required: true,
    },
  },
  setup(props) {
    const currentValue = ref<number>(0)

    const startCounter = (targetValue: number) => {
      // currentValue.value = 0
      const step = targetValue / (props.duration / 100)

      // let currentStep = 0
      const interval = setInterval(() => {
        currentValue.value += step
        // currentStep += 10

        if (currentValue.value >= targetValue) {
          currentValue.value = targetValue
          clearInterval(interval)
        }
      }, 10)
    }
    onMounted(() => {
      startCounter(props.value)
    })
    watch(
      () => props.value,
      (newValue) => {
        startCounter(newValue)
      },
    )

    return {
      currentValue,
    }
  },
}
</script>

<style scoped>
.value-counter {
  text-align: center;
}
</style>
