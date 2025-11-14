<template>
    <v-container no-gutters style="padding: 20px; margin-top: 10px; width:100%" >
        <v-row align="center">
            <v-data-iterator class="w-100" :items="data.scheduler" :page="pageManager.page">
                <template v-slot:default="{ items }">
                    <v-toolbar title="Asset Lists" style="margin-bottom: 10px; padding-right: 10px;">
                    </v-toolbar>

                    <template v-for="(item, i) in items" :key="i">
                        <v-card class="pa-4 mb-3 w-100">
                            <v-row align="center">
                                <v-col>
                                    <v-row align="center">

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.grid_id" variant="outlined"
                                                label="ID" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.grid_name" variant="outlined"
                                                label="Name" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.grid_descriptive" variant="outlined"
                                                label="Description" readonly></v-text-field>
                                        </v-col>

                                    </v-row>
                                </v-col>

                                <v-col cols="12" md="3" class="d-flex flex-column" style="gap:8px;">
                                    <v-btn variant="outlined" color="secondary" prepend-icon="mdi-open-in-new"
                                        @click="openSchedulerWindows(item.raw.grid_id)">
                                        Open Scheduler
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
import { reactive, ref, computed, watch } from 'vue';
import { useRoute } from 'vue-router';
import { GetScheduler } from '@/api/scheduler';

const printJsonFormatData = computed(() => {
    const jsonFormatData = JSON.stringify(data.scheduler.grid_tasks,null,4);
    console.log(`${jsonFormatData}`)
    return jsonFormatData 
})

const data = reactive({
    scheduler: []

});

watch(() => data.scheduler, (newdata) => {
    console.log(`updated scheduler data: ${JSON.stringify(newdata)}`)
});


const state = reactive({


});

async function loadingScheduler() {
    const response = await GetScheduler();
    Object.assign(data.scheduler, response.data);
    // data.scheduler = response.data;
    Object.assign(data.scheduler,(response.data || []))
    console.log(`loading all scheduler data: ${JSON.stringify(data.scheduler)}`)
}

loadingScheduler()

async function openSchedulerWindows(schedulerId) {
    const domains = window.location.origin;
    window.open(`${domains}/scheduler/${schedulerId}`,'_blank');
}

const pageManager = reactive({
    limit: 5,
    page: 1,
    total: 5,
})

watch(() => pageManager.page, (newPage) => {
    console.log(`switch to page ${newPage}`);
})

// asset data refersh -> page refersh
watch(() => data.scheduler,(newScheduler) => {
    pageManager.total = Math.ceil(newScheduler.length / pageManager.limit);
    console.log(`asset data pagination refresh: ${pageManager.total}`)
})


</script>