<template>
     <v-container>
        <v-row>
            <CustomeDataValueShowComponent cols="12" md="12" v-for="vv, kk in data" :label="kk" :value="vv.value" :description="vv.description" :required="vv.required" @update="(newdata) =>{
                console.log(`wmp custom component value update to ${newdata}`);
                data[kk].value = newdata
            }"></CustomeDataValueShowComponent>
        </v-row>

        <v-row justify="end">
            <v-col cols="auto">
                <v-btn color="cyan" prepend-icon="mdi-square-edit-outline" @click="sumbitWMPCIConnectorConfig">
                    Submit
                </v-btn>
            </v-col>

            <v-col cols="auto">
                <v-btn color="error" prepend-icon="mdi-delete-outline" @click="">
                    cancel
                </v-btn>
            </v-col>
        </v-row>
    </v-container>
</template>


<script setup>
import { reactive, ref, computed, watch } from 'vue';
import { useRoute } from 'vue-router';
import { GetConnectorConfig,PostConnectorConfig } from '@/api/wmpci';
import CustomeDataValueShowComponent from './custom/CustomeDataValueShowComponent.vue';

const wmpciConnectorReader = computed(() => {
    const route = useRoute()
    const id = computed(() => route.params?.connector_name ?? route.params?.connectorName ?? null)
    return id.value
});

const data = reactive({});

watch(() => data, (newdata) => {
    console.log(`updated scheduler data: ${JSON.stringify(newdata)}`)
});


async function loadingWMPCIConnectorConfig(connector_name) {
    const response = await GetConnectorConfig(connector_name);
    Object.assign(data, response.data);

    console.log(`loading ${connector_name} wmpci conector data: ${JSON.stringify(data)}`)
}

loadingWMPCIConnectorConfig(wmpciConnectorReader.value)

async function sumbitWMPCIConnectorConfig() {
    const response = await PostConnectorConfig(wmpciConnectorReader.value, data)
    if (response.code !== 200) {
        log.alert(`create failed: ${response.data} / ${response.msg}`);
        return
    }
    const domains = window.location.origin;
    window.open(`${domains}/wmpci/sessions`);
}
</script>