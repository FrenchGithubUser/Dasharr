<template>
  <div id="edit-credentials">
    <div class="form">
      <InputText v-for="field in orderedFields" :key="field" v-model="filledFields[field]" :placeholder="field" class="item" />
      <Button type="submit" label="Submit" class="item" @click="submitCredentials" :loading="loading" />
    </div>
  </div>
</template>

<script lang="ts">
import { Button, InputText } from 'primevue'

import { useApi } from '@/composables/useApi'
export default {
  emits: ['credentials-saved'],
  components: {
    InputText,
    // eslint-disable-next-line vue/no-reserved-component-names
    Button,
  },
  props: {
    fields: {
      type: Array as () => string[],
      required: true,
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
  computed: {
    orderedFields() {
      // should be flexible to add more ordered fields
      const order = ['username', 'password', 'twoFaCode']
      const fieldsCopy = [...this.fields]
      return fieldsCopy.sort((a, b) => {
        const aIndex = order.indexOf(a)
        const bIndex = order.indexOf(b)
        if (aIndex === -1 && bIndex === -1) return 0
        if (aIndex === -1) return 1
        if (bIndex === -1) return -1
        return aIndex - bIndex
      })
    },
  },
  methods: {
    submitCredentials() {
      this.loading = true
      console.log(this.filledFields)
      this.saveCredentials({ ...this.filledFields, indexer: this.indexerName }).then(() => {
        this.loading = false
        //todo: check why this doesn't work
        this.$emit('credentials-saved')
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
