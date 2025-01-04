<template>
  <div id="register-page">
    <div class="form">
      <InputText type="text" v-model="apiKey" class="input" placeholder="API key" />
      <Button type="submit" label="Submit" @click="saveApiKey" />
    </div>
  </div>
</template>

<script lang="ts">
import { Button } from 'primevue'
import InputText from 'primevue/inputtext'

export default {
  // eslint-disable-next-line vue/no-reserved-component-names
  components: { InputText, Button },
  data() {
    return {
      apiKey: '',
    }
  },
  methods: {
    saveApiKey() {
      localStorage.setItem('api-key', this.apiKey)
      this.$router.push('/home')
    },
  },
  setup() {},
}
</script>

<style scoped>
#register-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
.form {
  display: flex;
  flex-direction: column;
}
.input {
  margin-bottom: 20px;
}
</style>
