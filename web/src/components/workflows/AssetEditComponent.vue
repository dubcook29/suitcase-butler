<template>
    <v-form>
        <v-banner v-if="bannerShow" lines="one" icon="$warning" color="warning">
            <template v-slot:text>
                {{ bannerMessages }}
            </template>

            <template v-slot:actions>
                <v-btn color="deep-purple-accent-4" @click="close">
                    Close
                </v-btn>
            </template>
        </v-banner>
        <v-container>

            <v-row>
                <v-col cols="12" md="4">
                    <v-text-field v-model="assetData.group_id" :counter="36" label="Group ID" required></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="assetData.domain_name" :counter="36" label="Domain Name"
                        required></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="assetData.ip_address" :counter="36" label="IP Address"
                        required></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="assetData.org_name" :counter="36" label="Org Name" required></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="assetData.as_number" :counter="36" label="ASNs" required></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="assetData.cloud" :counter="36" label="Clouds" required></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="assetData.cdn" :counter="36" label="CDNs" required></v-text-field>
                </v-col>

            </v-row>
        </v-container>

        <div class="d-flex flex-column">
            <v-btn class="mt-4" color="success" block @click="submit">
                Submit
            </v-btn>

            <v-btn class="mt-4" color="error" block @click="reset">
                Reset
            </v-btn>

            <v-btn class="mt-4" color="warning" block @click="close">
                Close
            </v-btn>
        </div>
    </v-form>

</template>

<script setup>
import { reactive, watch, toRefs, ref } from 'vue'
import { useRoute } from 'vue-router';
import { GetAssetData, AddAssetData, EditAssetData } from '@/api/asset';

const props = defineProps({
    loadingMode: { type: String, default: null }, // component loading mode
    loadingData: { type: Object, default: null }, // component loading data
    loadingAssetId: { type: String, default: null }, // component self-loading asset—id
});

const bannerShow = ref(false);
const bannerMessages = ref("nobady")

const assetData = reactive({})


console.log("loading mode => ", props.loadingMode);
console.log("loading data => ", JSON.stringify(props.loadingData, null, 2));
console.log("loading asset id => ", props.loadingAssetId);


function initLoading() {
    switch (props.loadingMode) {
        case "loading":
            if (props.loadingData != null) {
                Object.assign(assetData, props.loadingData);
            } else {
                console.error("init loading component error: ", props)
            }
            break;
        case "insert":
            break;
        default:
            if (props.loadingAssetId == null) {
                // get `asset_id` from the full url params 
                loadingAssetData(() => useRoute().params?.asset_id ?? useRoute().params?.assetId ?? null);
            } else if (props.loadingAssetId != null) {
                loadingAssetData(props.loadingAssetId);
            } else {
                bannerShow.value = true;
                bannerMessages.value = "loading error: loading mode is not recognized."
                return
            }
    }
}

// Standard Method: loadingAssetData
async function loadingAssetData(assetId) {
    if (assetId == null) {
        bannerShow.value = true;
        bannerMessages.value = "loading error: `asset_id` is null, Unable to get data from the database."
        return
    }

    const response = await GetAssetData(assetId);
    if (response.code == 200) {
        if (Array.isArray(response.data)) {
            if (response.data.length > 0) {
                Object.assign(assetData, response.data[0]);
                console.log(`loading data successfully: ${response.msg}(${response.code})`)
            } else {
                bannerShow.value = true;
                bannerMessages.value = `loading data is null: ${response.msg}(${response.code})`
                console.log(`loading data is null: ${response.msg}(${response.code})`)
                return
            }
        } else {
            bannerShow.value = true;
            bannerMessages.value = `loading data format error: ${response.msg}(${response.code})`
            console.log(`loading data format error: ${response.msg}(${response.code})`)
            return
        }
    } else {
        bannerShow.value = true;
        bannerMessages.value = `loading request/response error: ${response.msg}(${response.code})`
        console.log(`loading request/response error: ${response.msg}(${response.code})`)
        return
    }
}

async function updateAssetData(assetId) {
    const response = await EditAssetData(assetId, assetData);
    console.log("update asset data: ", response)
}

async function insertAssetData() {
    const response = await AddAssetData(assetData)
    console.log("insert asset data: ", response)
}


const emit = defineEmits(['saved', 'cancel'])


function onAddAsset() {
    insertAssetData();
}

function onModAsset(assetId) {
    updateAssetData(assetId)
}

// submit btn click event
function submit() {
    emit('saved', { ...assetData })
    switch (props.loadingMode) {
        case "insert":
            onAddAsset();
            break;
        case "loading":
            onModAsset(assetData.asset_id);
            break;
    }
    console.log('close current edit windows.')
    emit('cancel');
}

// reset btn click event
function reset() {
    Object.assign(assetData, props.loadingData);
}

// close btn click event
function close() {
    emit('cancel')
}

initLoading();

</script>
