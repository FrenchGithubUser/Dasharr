<template>
  <Card class="tracker-card">
    <template #content>
      <img :src="'/images/' + stats.tracker_id + '.png'" width="200px" />
      <div class="stats">
        <ValueCounter v-for="(stat, label) in statsToDisplay" :key="label" :value="stat" :meaning="label" :unit="['uploaded', 'downloaded'].indexOf(label) > -1 ? 'GiB' : ''" :duration="500" />
      </div>
    </template>
  </Card>
</template>

<script lang="ts">
import Card from 'primevue/card'
import ValueCounter from './ValueCounter.vue'

export default {
  components: {
    Card,
    ValueCounter,
  },
  props: {
    stats: {
      type: Object,
      required: true,
    },
  },
  computed: {
    statsToDisplay() {
      return Object.fromEntries(Object.entries(this.stats).filter(([label, value]) => typeof value === 'number' && label != 'tracker_id'))
    },
  },
  setup() {},
}
</script>

<style scoped>
.stats {
  display: flex;
  justify-content: space-around;
}
</style>
