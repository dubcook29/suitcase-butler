<template>
    <v-container>
        <v-row>
            <v-col cols="12" md="4">
                <v-text-field v-model="data.scheduler.grid_id" variant="outlined" label="ID" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field v-model="data.scheduler.grid_name" variant="outlined" label="Name"
                    readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field v-model="data.scheduler.grid_descriptive" variant="outlined" label="Description"
                    readonly></v-text-field>
            </v-col>

        </v-row>

        <v-row justify="end">

            <v-col cols="auto">
                <v-btn color="cyan" prepend-icon="mdi-square-edit-outline" @click="toggleAssetEditDialogState" disabled>
                    EDIT
                </v-btn>
            </v-col>

            <v-col cols="auto">
                <v-btn color="error" prepend-icon="mdi-delete-outline" @click="onDeleteAsset" disabled>
                    DELETE
                </v-btn>
            </v-col>
        </v-row>

        <v-row>
            <v-textarea label="data.scheduler.grid_name" :model-value="printJsonFormatData" auto-grow>
            </v-textarea>
        </v-row>

    </v-container>
</template>


<script setup>
import { reactive, ref, computed, watch } from 'vue';
import { useRoute } from 'vue-router';
import { GetSchedulerId } from '@/api/scheduler';

const schedulerIDReader = computed(() => {
    const route = useRoute()
    const id = computed(() => route.params?.scheduler_id ?? route.params?.schedulerId ?? null)
    return id.value
});

const printJsonFormatData = computed(() => {
    const jsonFormatData = JSON.stringify(data.scheduler.grid_tasks,null,4);
    console.log(`${jsonFormatData}`)
    return jsonFormatData 
})

const data = reactive({
    scheduler: {}

});

watch(() => data.scheduler, (newdata) => {
    console.log(`updated scheduler data: ${JSON.stringify(newdata)}`)
});


const state = reactive({


});

async function loadingScheduler(scheduler_id) {
    const response = await GetSchedulerId(scheduler_id);
    Object.assign(data.scheduler, response.data);
    // data.scheduler = response.data;
    console.log(`loading ${scheduler_id} scheduler data: ${JSON.stringify(data.scheduler)}`)
}

loadingScheduler(schedulerIDReader.value)
</script>