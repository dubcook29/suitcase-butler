<template>
    <v-container>
        <v-row>
            <v-col cols="12" md="4">
                <v-text-field v-model="data.session_id" variant="outlined" label="ID" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field v-model="data.last_connection" variant="outlined" label="Name" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field v-model="data.first_connection" variant="outlined" label="Description"
                    readonly></v-text-field>
            </v-col>
        </v-row>

        <v-row justify="end">

            <v-col cols="auto">
                <v-btn color="cyan" prepend-icon="mdi-square-edit-outline" @click="">
                    EDIT
                </v-btn>
            </v-col>

            <v-col cols="auto">
                <v-btn color="error" prepend-icon="mdi-delete-outline" @click="">
                    DELETE
                </v-btn>
            </v-col>
        </v-row>

        <WMPInfoComponent v-if="application" :wmp="data.application_registration"></WMPInfoComponent>

        <v-row>
            <v-textarea label="data" :model-value="printJsonFormatData" auto-grow>
            </v-textarea>
        </v-row>

    </v-container>
</template>


<script setup>
import { reactive, ref, computed, watch, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { GetCurrentSessionConnect } from '@/api/wmpci';
import WMPInfoComponent from './WMPInfoComponent.vue';

const wmpciSessionIDReader = computed(() => {
    const route = useRoute()
    const id = computed(() => route.params?.session_id ?? route.params?.sessionId ?? null)
    return id.value
});

const printJsonFormatData = computed(() => {
    const jsonFormatData = JSON.stringify(data.application_registration,null,4);
    return jsonFormatData 
})

const data = reactive({});
const application = ref(false)

watch(() => data, (newdata) => {
    console.log(`updated scheduler data: ${JSON.stringify(newdata)}`)
});

watch(() => data.application_registration, (newdata) =>{
    console.log(`wmpci application registration updated ${JSON.stringify(newdata,null,4)}`);
}, { deep: true });

async function loadingWMPCISessionConnect(session_id) {
    const response = await GetCurrentSessionConnect(session_id);
    if (Array.isArray(response.data) && response.data.length > 0) {
        Object.assign(data, response.data[0]);
        application.value = true;
    }
    return
}

onMounted(() => {
    loadingWMPCISessionConnect(wmpciSessionIDReader.value);

})

</script>