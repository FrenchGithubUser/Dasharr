<template>
  <div id="edit-credentials">
    <div class="form">
      <InputText v-for="field in fields" :key="field" v-model="filledFields[field]" :placeholder="field" class="item" />
      <Button type="submit" label="Submit" class="item" @click="submitCredentials" :loading="loading" />
    </div>
  </div>
</template>

<script lang="ts">
import { Button, InputText } from 'primevue'

import { useApi } from '@/composables/useApi'
export default {
  emits: ['credentialsSaved'],
  components: {
    InputText,
    // eslint-disable-next-line vue/no-reserved-component-names
    Button,
  },
  props: {
    fields: {
      type: Array as () => string[],
    },
    indexerName: {
      type: String,
    },
  },
  data() {
    return {
      filledFields: {} as { [key: string]: string },
      loading: false,
    }
  },
  methods: {
    submitCredentials() {
      this.loading = true
      this.saveCredentials({ ...this.filledFields, indexer: this.indexerName }).then(() => {
        this.loading = false
        //todo: check why this doesn't work
        this.$emit('credentialsSaved')
      })
    },
  },
  setup() {
    const { saveCredentials } = useApi()

    return { saveCredentials }
  },
}
</script>

<style scoped lang="scss">
.form {
  display: flex;
  flex-direction: column;
  .item {
    margin: 5px;
  }
}
</style>
