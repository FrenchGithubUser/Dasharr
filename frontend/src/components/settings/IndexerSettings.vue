<template>
  <div id="indexer-settings">
    <div class="note">Note: Indexers need to be setup in Prowlarr first in order to work with Dasharr, even if the credentials from Prowlarr are not used</div>
    <Card v-for="indexer in indexersConfig" :key="indexer['site_name']" class="indexer-card">
      <template #content>
        <div class="indexer">
          <div class="left">
            <div class="indexer-name">{{ indexer['site_name'] }}</div>
          </div>
          <div class="right">
            <Button v-if="indexer['credentials']['method'] == 'built_in'" icon="pi pi-pencil" @click="selectIndexer(indexer)" />
            <div v-if="indexer['credentials']['method'] == 'prowlarr'">Credentials managed in Prowlarr</div>
            <Chip label="Already setup" icon="pi pi-check" class="status" v-if="savedCredentialIndexers.some((object: any) => object.indexer_name === indexer['site_name'])" />
            <Chip label="Not setup" icon="pi pi-times" class="status" v-if="indexer['credentials']['method'] != 'prowlarr' && !savedCredentialIndexers.some((object: any) => object.indexer_name === indexer['site_name'])" />
            <ToggleSwitch class="toggle-switch" @change="updateEnbaledTrackers(indexer['site_name'])" :modelValue="enabledIndexers.includes(Object.keys(trackerMap).find((key) => trackerMap[key] === indexer['site_name'])!)" />
          </div>
        </div>
      </template>
    </Card>
    <Dialog v-model:visible="editCredentialsDialog" modal header="Edit credentials" @credentialsSaved="credentialsSaved"><EditCredentials :fields="selectedIndexer.fillableFields" :indexerName="selectedIndexer.name" /></Dialog>
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
  methods: {
    selectIndexer(indexer: any) {
      this.selectedIndexer.fillableFields = Object.keys(indexer['login']['fields']).filter((key) => key !== 'extra')
      this.selectedIndexer.name = indexer['site_name']
      this.editCredentialsDialog = true
    },
    credentialsSaved() {
      this.savedCredentialIndexers.push({ indexer_name: this.selectIndexer.name })
      this.editCredentialsDialog = false
    },
    updateEnbaledTrackers(siteName: string) {
      const indexerId = Object.keys(this.trackerMap).find((key) => this.trackerMap[key] === siteName) ?? ''
      const indexerEnabled = this.enabledIndexers.includes(indexerId)
      if (indexerEnabled) {
        this.enabledIndexers = this.enabledIndexers.filter((id) => id !== indexerId)
      } else {
        this.enabledIndexers.push(indexerId)
      }
      localStorage.setItem('enabledIndexers', JSON.stringify(this.enabledIndexers))
      console.log(this.enabledIndexers)
    },
  },

  setup() {
    const { getConfig, savedCredentials, getTrackerMap } = useApi()
    const indexersConfig = ref<object>({})
    const savedCredentialIndexers = ref<object[]>([])
    const enabledIndexers = ref<Array<string>>([])
    const trackerMap = ref<any>({})

    onMounted(() => {
      enabledIndexers.value = JSON.parse(localStorage.getItem('enabledIndexers') ?? '[]')
      getConfig().then((data) => {
        indexersConfig.value = data
      })
      savedCredentials().then((data) => {
        savedCredentialIndexers.value = data
      })
      getTrackerMap().then((data) => {
        trackerMap.value = data
      })
    })
    return {
      indexersConfig,
      savedCredentialIndexers,
      enabledIndexers,
      trackerMap,
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
