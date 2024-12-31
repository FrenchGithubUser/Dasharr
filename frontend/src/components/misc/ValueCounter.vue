<template>
  <div class="widget value-counter">
    <div>{{ value % 1 === 0 ? currentValue.toFixed(0) : currentValue.toFixed(2) }} {{ unit }}</div>
    <div>{{ meaning }}</div>
  </div>
</template>

<script lang="ts">
import { watch, ref } from 'vue'

export default {
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
      required: true,
    },
    meaning: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const currentValue = ref<number>(0)

    const startCounter = (targetValue: number) => {
      const step = targetValue / (props.duration / 10)

      let currentStep = 0
      const interval = setInterval(() => {
        currentValue.value += step
        currentStep += 10

        if (currentValue.value >= targetValue) {
          currentValue.value = targetValue
          clearInterval(interval)
        }
      }, 10)
    }

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
