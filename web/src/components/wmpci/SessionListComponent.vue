<template>
    <v-container no-gutters style="padding: 20px; margin-top: 10px; width:100%" >
        <v-row align="center">
            
            <v-data-iterator class="w-100" :items="data" :page="pageManager.page">
                <template v-slot:default="{ items }">
                    <v-toolbar title="WMPCI Sessions" style="margin-bottom: 10px; padding-right: 10px;">
                        <template v-slot:append>
                            <v-btn prepend-icon="mdi-plus-box-outline" @click="null">Add</v-btn>
                        </template>
                    </v-toolbar>

                    <template v-for="(item, i) in items" :key="i">
                        <v-card class="pa-4 mb-3 w-100">
                            <v-row align="center">
                                <v-col>
                                    <v-row align="center">

                                        <v-col cols="12" md="4">
                                            <v-text-field v-model="item.raw.session_id" variant="outlined"
                                                label="Session ID" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="4">
                                            <v-text-field v-model="item.raw.last_connection" variant="outlined"
                                                label="Last Connection" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="4">
                                            <v-text-field v-model="item.raw.first_connection" variant="outlined"
                                                label="First Connection" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="4">
                                            <v-text-field v-model="item.raw.application_registration.wmp_basic.id" variant="outlined" label="WMP ID"
                                                readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="4">
                                            <v-text-field v-model="item.raw.application_registration.wmp_basic.name" variant="outlined" label="WMP Name"
                                                readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="4">
                                            <v-text-field v-model="item.raw.application_registration.wmp_basic.version" variant="outlined" label="WMP Version"
                                                readonly></v-text-field>
                                        </v-col>

                                        

                                    </v-row>
                                </v-col>

                                <v-col cols="12" md="3" class="d-flex flex-column" style="gap:8px;">
                                    <v-btn variant="outlined" color="secondary" prepend-icon="mdi-open-in-new"
                                        @click="newWindowForOpenSessionComponent(item.raw.session_id)">
                                        Open WMPCI
                                    </v-btn>
                                </v-col>
                            </v-row>
                        </v-card>
                    </template>

                    <v-pagination v-model="pageManager.page" :length="pageManager.total"></v-pagination>
                </template>
            </v-data-iterator>

        </v-row>


    </v-container>
</template>
<script setup>
import { GetCurrentAllSessionConnect } from '@/api/wmpci';
import { reactive, ref, computed, watch, onMounted } from 'vue';

const data = ref([])

async function loadingAllWMPCISessionConnect() {
    const response = await GetCurrentAllSessionConnect();
    if (Array.isArray(response.data) && response.data.length > 0) {
        data.value.splice(0, data.value.length, ...(response.data || []))
    }
}

const pageManager = reactive({
    limit: 5,
    page: 1,
    total: 1,
})

watch(() => data, (newData) => {
    pageManager.total = Math.ceil(newData.value.length / pageManager.limit);
}, { deep: true })

onMounted(() => {
    loadingAllWMPCISessionConnect();
})

async function newWindowForOpenSessionComponent(session_id) {
    const domains = window.location.origin;
    window.open(`${domains}/wmpci/sessions/${session_id}`,'_blank');
}

</script>