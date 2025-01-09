<template>
  <div id="settings-dialog">
    <div class="sidebar">
      <Button class="credentials" label="Credentials" @click="currentSection = 'credentials'" />
    </div>
    <div class="sections">
      <CredentialsSettigns class="settings" v-if="currentSection == 'credentials'" />
    </div>
  </div>
</template>

<script lang="ts">
import Button from 'primevue/button'
import CredentialsSettigns from './CredentialsSettigns.vue'
export default {
  components: {
    // eslint-disable-next-line vue/no-reserved-component-names
    Button,
    CredentialsSettigns,
  },
  data() {
    return {
      currentSection: 'credentials',
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
