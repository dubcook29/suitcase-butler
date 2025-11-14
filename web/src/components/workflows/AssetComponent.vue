<template>
    
    <v-container>
        <v-banner v-if="bannerShow" icon="$warning" color="warning" class="mb-3">
            <template v-slot:text>
                {{ bannerMessages }}
            </template>

            <!-- <template v-slot:actions>
                <v-btn color="deep-purple-accent-4" @click="close">
                    Close
                </v-btn>
            </template> -->
        </v-banner>
        <v-row>
            <v-col cols="12" md="4">
                <v-text-field v-model="assetData.group_id" variant="outlined" label="Group ID" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field v-model="assetData.asset_id" variant="outlined" label="Asset ID" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field v-model="assetData.domain_name" variant="outlined" label="Domains" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field v-model="assetData.ip_address" variant="outlined" label="IPs" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field v-model="assetData.as_number" variant="outlined" label="ASNs" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field v-model="assetData.org_name" variant="outlined" label="Organizations" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field v-model="assetData.cloud" variant="outlined" label="Clouds" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field v-model="assetData.cdn" variant="outlined" label="CDNs" readonly></v-text-field>
            </v-col>

        </v-row>

        <v-row justify="end">

            <v-col cols="auto">
                <v-btn prepend-icon="mdi-notebook-outline" @click="toggleNotes">
                    NOTES
                </v-btn>
            </v-col>

            <v-col cols="auto">
                <v-btn prepend-icon="mdi-file-chart-outline" @click="toggleReport">
                    REPORT
                </v-btn>
            </v-col>

            <v-col cols="auto">
                <v-btn color="cyan" prepend-icon="mdi-square-edit-outline" @click="toggleAssetEditDialogState">
                    EDIT
                </v-btn>
            </v-col>

            <v-col cols="auto">
                <v-btn color="error" prepend-icon="mdi-delete-outline" @click="onDeleteAsset">
                    DELETE
                </v-btn>
            </v-col>
        </v-row>

        <v-expand-transition v-if="isReport">
            <v-row>
                <v-tabs color="deep-purple-accent-4" center-active>
                    <v-tab v-for="wmp_name in assetData.result_wmp_data_lists" :value=wmp_name :key=wmp_name
                        @click="loadingWMPData(wmp_name)">
                        {{ wmp_name }}
                    </v-tab>
                </v-tabs>
                <v-card style="padding: 20px; margin-top: 10px;" width="100%">
                    <WmpdataComponent :data="wmpdata">

                    </WmpdataComponent>
                </v-card>
            </v-row>
        </v-expand-transition>

    </v-container>

    <v-dialog v-model="assetEditDialogState" max-width="800px" persistent>
        <v-card>
            <v-card-title class="d-flex justify-space-between align-center">
                <div>Asset Edit</div>
                <v-icon @click="closeAssetEditDialogState">mdi-close</v-icon>
            </v-card-title>

            <v-divider></v-divider>

            <v-card-text>
                <AssetEditComponent loadingMode="loading" :loadingData="assetEditDialogDatas" @cancel="closeAssetEditDialogState" @saved="syncAssetEditDialogData"></AssetEditComponent>
            </v-card-text>
            <v-divider></v-divider>

            <v-card-actions>
                <v-spacer />
                <v-btn color="primary" @click="closeAssetEditDialogState">取消</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup>
import { ref,reactive,computed } from 'vue';
import { useRoute } from 'vue-router';
import { GetAssetWMPData } from '@/api/wmpdata';
import { GetAssetData,DelAssetData } from '@/api/asset';
import AssetEditComponent from "./AssetEditComponent.vue";
import WmpdataComponent from '../wmpdata/WmpdataComponent.vue';

// asset note data loading

const isNotes = ref(false);
function toggleNotes() {
  isNotes.value = !isNotes.value;
  if (isNotes.value) isReport.value = false;
}

// asset wmpdata loading

const wmpdata = ref([])
const isReport = ref(false);
function toggleReport() {
  isReport.value = !isReport.value;
  if (isReport.value) isNotes.value = false;
}

const getcurrentwindowswmpdata = computed(() => {
    console.log(`loading wmp-data : ${JSON.stringify(wmpdata.value)}`)
    return wmpdata.value;
})


async function loadingWMPData(wmp_name) {
    const response = await GetAssetWMPData(assetData.asset_id,wmp_name);
    console.log("get asset wmpdata", response.msg, response.data);
    // Object.assign(wmpdata, response.data);
    wmpdata.value = response.data;
}

// review design

// asset data loading 

const props = defineProps({
    asset_data: { type: Object ,default: null},
});

const assetData = reactive(props.asset_data == null ? {} : props.asset_data);

const assetIdReader = computed(() => {
    const route = useRoute()
    const assetId = computed(() => route.params?.asset_id ?? route.params?.assetId ?? null)
    return assetId.value
});

if (props.asset_data == null  && assetIdReader.value != null) {
    loadingAssetData(assetIdReader.value)
} else {
    Object.assign(assetData, {
        "group_id": "************",
        "asset_id": "************",
        "domain_name": "************.***",
        "ip_address": "************",
        "org_name": "",
        "as_number": "",
        "cloud": "oracle(cloud);aws(cloud)",
        "cdn": "fastly(cdn)",
        "other_input_value": "",
        "result_wmp_data_lists": [],
        "created_at": "2025-08-21T08:19:06.593Z",
        "updated_at": "2025-08-21T08:19:06.593Z",
        "deleted_at": "0001-01-01T00:00:00Z"
    });
};
    
const bannerShow = ref(false);
const bannerMessages = ref("nobady")

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

// manager asset edit dialog state

const assetEditDialogState = ref(false);
const assetEditDialogDatas = reactive({});

function toggleAssetEditDialogState() {
    assetEditDialogState.value = !assetEditDialogState.value;
    if (assetEditDialogState.value === true) {
        Object.assign(assetEditDialogDatas,assetData);
    } else {
        Object.assign(assetEditDialogDatas,{});
    }

    console.log(`asset edit dialog state (${assetEditDialogState.value}) data: ${JSON.stringify(assetEditDialogDatas)}`);
    console.log(`asset data reader: ${JSON.stringify(assetData)}}`)
}

function closeAssetEditDialogState() {
    Object.assign(assetEditDialogDatas,{});
    assetEditDialogState.value = false;
}

function syncAssetEditDialogData(data) {
    Object.assign(assetData,data);
}

// asset delete button

function onDeleteAsset(){
    deleteAssetData()
}

async function deleteAssetData() {
    const response = await DelAssetData(assetData.asset_id);
    if (response.code === 200) {
        bannerShow.value = true;
        bannerMessages.value = `deleted asset successfully: ${response.msg}(${response.code})`
        console.log(`deleted asset data successed: ${response.msg}`);
        Object.assign(assetData, {})
    } else {
        bannerShow.value = true;
        bannerMessages.value = `deleted asset error: ${response.msg}(${response.code})`
        console.log(`deleted asset data failed: (${response.code})${response.msg}`);
    }
    // loadingAssetData(assetData.asset_id);
}



</script>