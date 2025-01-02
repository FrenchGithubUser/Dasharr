<template>
  <Card class="value-counter">
    <template #content>
      <div>{{ currentValue % 1 === 0 ? currentValue.toFixed(0) : currentValue.toFixed(5) }} {{ unit }}</div>
      <div>{{ meaning }}</div>
    </template>
  </Card>
</template>

<script lang="ts">
import { watch, ref } from 'vue'
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
      type: String,
      required: true,
    },
    data: {
      type: Boolean,
      required: true,
    },
  },
  setup(props) {
    const currentValue = ref<number>(0)

    const startCounter = (targetValue: number) => {
      const step = targetValue / (props.duration / 10)
      console.log(step)

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
