<template>
  <div id="settings-dialog">
    <div class="sidebar">
      <Button class="indexers" label="Indexers" @click="currentSection = 'indexers'" />
    </div>
    <div class="sections">
      <IndexerSettigns class="settings" v-if="currentSection == 'indexers'" />
    </div>
  </div>
</template>

<script lang="ts">
import Button from 'primevue/button'
import IndexerSettigns from './IndexerSettings.vue'
export default {
  components: {
    // eslint-disable-next-line vue/no-reserved-component-names
    Button,
    IndexerSettigns,
  },
  data() {
    return {
      currentSection: 'indexers',
    }
  },
}
</script>

<style scoped lang="scss">
#settings-dialog {
  display: flex;
  height: 60vh;
}
.sidebar {
  margin-right: 15px;
  width: 10vw;
}
.sections {
  width: 60vw;
}
</style>
