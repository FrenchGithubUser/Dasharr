<template>
  <div id="indexer-settings">
    <div class="note">Indexers need to be setup in Prowlarr first in order to work with Dasharr, even if the credentials from Prowlarr are not used</div>
    <div class="note">Make sure the API keys configured in prowlarr have access to user stats</div>
    <div class="note">Enabling/disabling an indexer here will only affect its visibility on the dashboard, data will always be collected if it is enabled in config.json</div>
    <Card v-tooltip.top="{ value: 'To display this indexer in the dashboard, first enable it in config.json', disabled: indexer['enabled'] }" v-for="indexer in indexersConfig" :key="indexer['indexer_name']" :class="{ disabled: !indexer['enabled'], 'indexer-card': true }">
      <template #content>
        <div class="indexer">
          <div class="left">
            <div class="indexer-name">{{ indexer['indexer_name'] }}</div>
          </div>
          <div class="right">
            <Button v-if="indexer['credentials']['method'] == 'built_in'" icon="pi pi-pencil" @click="selectIndexer(indexer)" />
            <div v-if="indexer['credentials']['method'] == 'prowlarr'">Credentials managed in Prowlarr</div>
            <div v-else>
              <Chip label="Credentials saved" icon="pi pi-check" class="status" v-if="areIndexerCredentialsSaved(indexer)" />
              <Chip label="Credentials missing" icon="pi pi-times" class="status" v-if="!areIndexerCredentialsSaved(indexer)" />
            </div>
            <ToggleSwitch class="toggle-switch" @change="updateEnbaledIndexers(indexer['indexer_name'])" :modelValue="indexer['enabled'] ? enabledIndexers.includes(Object.keys(indexerMap).find((key) => indexerMap[key] === indexer['indexer_name'])!) : false" :disabled="!indexer['enabled']" />
          </div>
        </div>
      </template>
    </Card>
    <Dialog v-model:visible="editCredentialsDialog" modal header="Edit credentials" @credentialsSaved="credentialsSaved"><EditCredentials @credentials-saved="credentialsSaved" :fields="selectedIndexer.fillableFields" :indexerName="selectedIndexer.name" /></Dialog>
  </div>
</template>

<script lang="ts">
import { useApi } from '@/composables/useApi'
import { onMounted, ref } from 'vue'
import EditCredentials from './EditCredentials.vue'
import { Dialog } from 'primevue'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Chip from 'primevue/chip'
import ToggleSwitch from 'primevue/toggleswitch'

export default {
  components: {
    // eslint-disable-next-line vue/no-reserved-component-names
    Button,
    Card,
    EditCredentials,
    // eslint-disable-next-line vue/no-reserved-component-names
    Dialog,
    Chip,
    ToggleSwitch,
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
  computed: {
    areIndexerCredentialsSaved() {
      return (indexer: any) => this.savedCredentialIndexers.some((obj: any) => obj.indexer_name === indexer.indexer_name)
    },
  },
  methods: {
    selectIndexer(indexer: any) {
      this.selectedIndexer.fillableFields = Object.keys(indexer['login']['fields']).filter((key) => key !== 'extra')
      this.selectedIndexer.name = indexer['indexer_name']
      this.editCredentialsDialog = true
    },
    credentialsSaved() {
      this.savedCredentialIndexers.push({ indexer_name: this.selectedIndexer.name })
      console.log(this.savedCredentialIndexers)
      this.editCredentialsDialog = false
    },
    updateEnbaledIndexers(indexerName: string) {
      const indexerId = Object.keys(this.indexerMap).find((key) => this.indexerMap[key] === indexerName) ?? ''
      const indexerEnabled = this.enabledIndexers.includes(indexerId)
      if (indexerEnabled) {
        this.enabledIndexers = this.enabledIndexers.filter((id) => id !== indexerId)
      } else {
        this.enabledIndexers.push(indexerId)
      }
      localStorage.setItem('enabledIndexers', JSON.stringify(this.enabledIndexers))
    },
  },

  setup() {
    const { getConfig, savedCredentials, getIndexerMap } = useApi()
    const indexersConfig = ref<object>({})
    const savedCredentialIndexers = ref<object[]>([])
    const enabledIndexers = ref<Array<string>>([])
    const indexerMap = ref<any>({})

    onMounted(() => {
      enabledIndexers.value = JSON.parse(localStorage.getItem('enabledIndexers') ?? '[]')
      getConfig().then((data) => {
        indexersConfig.value = data
      })
      savedCredentials().then((data) => {
        savedCredentialIndexers.value = data
      })
      getIndexerMap().then((data) => {
        indexerMap.value = data
      })
    })

    return {
      indexersConfig,
      savedCredentialIndexers,
      enabledIndexers,
      indexerMap,
    }
  },
}
</script>

<style scoped lang="scss">
#indexer-settings {
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
  &.disabled {
    background-color: rgb(255, 0, 0) !important;
  }
}
.indexer {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .indexer-name {
    font-weight: bold;
  }
  .right {
    display: flex;
    align-items: center;
    .status {
      margin-left: 10px;
    }
    .toggle-switch {
      margin-left: 10px;
    }
  }
}
</style>
