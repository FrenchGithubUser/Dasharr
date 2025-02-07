import './assets/main.scss'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import { definePreset } from '@primevue/themes'
import Tooltip from 'primevue/tooltip'

const app = createApp(App)

app.use(router)

const Noir = definePreset(Aura, {
  components: {
    chip: {
      colorScheme: {
        // light: {
        //   root: {
        //     shadow: '0px 0px 10px rgba(0, 0, 0, 0.2)',
        //   },
        // },
        dark: {
          root: {
            'background-color': 'white',
          },
        },
      },
    },
  },
  semantic: {
    primary: {
      50: '{amber.50}',
      100: '{amber.100}',
      200: '{amber.200}',
      300: '{amber.300}',
      400: '{amber.400}',
      500: '{amber.500}',
      600: '{amber.600}',
      700: '{amber.700}',
      800: '{amber.800}',
      900: '{amber.900}',
      950: '{amber.950}',
    },
  },
})

app.use(PrimeVue, {
  theme: {
    preset: Noir,
    options: {
      prefix: 'p',
      cssLayer: false,
      darkModeSelector: true,
    },
  },
})

app.directive('tooltip', Tooltip)

app.mount('#app')
