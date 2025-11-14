
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

const app = createApp(App)

// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const vuetify = createVuetify({
  components,
  directives,
  icons: {
    // https://pictogrammers.com/library/mdi/
    defaultSet: 'mdi',
    aliases,
    sets:{
        mdi,
    },
  }
})

app.use(router).use(vuetify).mount('#app')