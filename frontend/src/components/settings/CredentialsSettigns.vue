<template>
  <div id="credentials-settings">
    <div class="note">Note: Indexers need to be setup in Prowlarr first in order to work with Dasharr, even if the credentials from Prowlarr are not used</div>
    <Card v-for="indexer in indexersConfig" :key="indexer['site_name']" class="indexer-card">
      <template #content>
        <div class="indexer">
          <div class="indexer-name">{{ indexer['site_name'] }}</div>
          <Button v-if="indexer['credentials']['method'] == 'built_in'" icon="pi pi-pencil" @click="selectIndexer(indexer)" />
          <div v-if="indexer['credentials']['method'] == 'prowlarr'">Credentials managed in Prowlarr</div>
        </div>
      </template>
    </Card>
    <Dialog v-model:visible="editCredentialsDialog" modal header="Edit credentials" @credentialsSaved="editCredentialsDialog = false"><EditCredentials :fields="selectedIndexer.fillableFields" :indexerName="selectedIndexer.name" /></Dialog>
  </div>
</template>

<script lang="ts">
import { useApi } from '@/composables/useApi'
import Button from 'primevue/button'
import Card from 'primevue/card'
import { onMounted, ref } from 'vue'
import EditCredentials from './EditCredentials.vue'
import { Dialog } from 'primevue'
export default {
  components: {
    // eslint-disable-next-line vue/no-reserved-component-names
    Button,
    Card,
    EditCredentials,
    // eslint-disable-next-line vue/no-reserved-component-names
    Dialog,
  },
  data() {
    return {
      editCredentialsDialog: false,
      selectedIndexer: {
        fillableFields: [] as string[],
        name: '',
      },
    }
  },
  methods: {
    selectIndexer(indexer: any) {
      this.selectedIndexer.fillableFields = Object.keys(indexer['login']['fields']).filter((key) => key !== 'extra')
      this.selectedIndexer.name = indexer['site_name']
      this.editCredentialsDialog = true
    },
  },

  setup() {
    const { getConfig } = useApi()
    const indexersConfig = ref<object>({})

    onMounted(() => {
      getConfig().then((data) => {
        indexersConfig.value = data
      })
    })
    return {
      indexersConfig,
    }
  },
}
</script>

<style scoped lang="scss">
#credentials-settings {
  overflow-y: scroll;
  height: 100%;
}
.note {
  font-weight: bold;
  margin-bottom: 15px;
}
.indexer-card {
  margin: 5px;
  margin-bottom: 10px;
}
.indexer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.indexer-name {
  font-weight: bold;
}
</style>
