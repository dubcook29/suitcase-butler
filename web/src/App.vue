<script setup>
import { ref } from 'vue';
import { useTheme } from 'vuetify/lib/composables/theme';

const theme = useTheme()
const themeIcon = ref(theme.global.name.value !== "dark" ? "☀️" : "🌙")
const watermark = ref(true)

function toggleTheme() {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark';
  themeIcon.value = theme.global.name.value !== "dark" ? "☀️" : "🌙";
}
</script>

<template>
  <v-app>
    <v-layout>
      <v-app-bar :elevation="2">
        <template v-slot:prepend>
          <v-app-bar-nav-icon>
            <v-icon @click="$router.push('/')" style="color: red;">mdi-hat-fedora</v-icon>
          </v-app-bar-nav-icon>
        </template>

        <v-app-bar-title>Red-Team BULTER Application</v-app-bar-title>

        <template v-slot:append>
          <v-btn icon="mdi-theme-light-dark" @click="toggleTheme">
            {{ themeIcon }}
          </v-btn>
          <v-btn icon="mdi-translate"></v-btn>

          <v-btn icon="mdi-github" href="https://github.com" target="_blank" rel="noopener noreferrer"></v-btn>

          <v-btn icon="mdi-dots-vertical"></v-btn>

          <v-switch v-model="watermark" hide-details inset label="watermark"></v-switch>
        </template>

      </v-app-bar>

      <v-navigation-drawer expand-on-hover permanent rail>
        <v-list>
          <v-list-item prepend-avatar="/public/github-mark.svg" subtitle="github.com/dubcook29/suitcase-butler"
            title="Suitcase Butler"></v-list-item>
        </v-list>

        <v-divider></v-divider>

        <v-list density="compact" nav>

          <v-list-item prepend-icon="mdi-target" title="Assets" value="assets"
            @click="$router.push('/asset')"></v-list-item>

          <v-list-item prepend-icon="mdi-floor-plan" title="Scheduler" value="scheduler"
            @click="$router.push('/scheduler/')"></v-list-item>

          <v-list-group value="WMPCI">
            <template v-slot:activator="{ props }">
              <v-list-item v-bind="props" prepend-icon="mdi-toy-brick" title="WMPCI"></v-list-item>
            </template>

            <v-list-item title="Sessions" value="Sessions" @click="$router.push('/wmpci/sessions')"></v-list-item>
            <v-list-item title="Connector" value="Connector" @click="$router.push('/wmpci/connector')"></v-list-item>

          </v-list-group>

          <v-list-group value="Assets">
            <template v-slot:activator="{ props }">
              <v-list-item v-bind="props" prepend-icon="mdi-sitemap" title="Workflow"></v-list-item>
            </template>

            <v-list-item title="Task" value="Task" @click="$router.push('/workflow/task')"></v-list-item>

          </v-list-group>
          <v-divider></v-divider>


        </v-list>
      </v-navigation-drawer>

      <v-main style="overflow: auto; margin: 10px;">
        <div v-if="watermark" class="watermark">
          <span>BUTLER WEB UI (V0.0.1-Alpha)</span>
          <br></br>
        </div>
        <RouterView />
      </v-main>
    </v-layout>
  </v-app>
</template>

<style scoped>
.watermark {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  opacity: 0.1;
  font-size: 5rem; 
  color: #000; 
  pointer-events: none; 
  z-index: 9999;
  white-space: nowrap;
}
</style>