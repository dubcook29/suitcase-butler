<template>
    <v-container no-gutters style="padding: 20px; margin-top: 10px; width:100%" >
        <v-row align="center">
            <v-banner v-if="bannerShow" icon="$warning" color="warning" class="mb-3">
                <template v-slot:text>
                    {{ bannerMessages }}
                </template>
            </v-banner>
            <v-data-iterator class="w-100" :items="assetData" :page="pageManager.page">
                <template v-slot:default="{ items }">
                    <v-toolbar title="Asset Lists" style="margin-bottom: 10px; padding-right: 10px;">
                        <template v-slot:append>
                            <v-btn prepend-icon="mdi-plus-box-outline" @click="onAddAsset">Add</v-btn>
                            <v-btn v-if="!isTaskLoading" prepend-icon="mdi-refresh"
                                @click="getAllAssetData">Refresh</v-btn>
                            <v-switch v-if="isTaskLoading" v-model="isTaskAsset" :label="loadAssetDataClass"
                                hide-details inset></v-switch>
                        </template>
                    </v-toolbar>

                    <template v-for="(item, i) in items" :key="i">
                        <v-card class="pa-4 mb-3 w-100">
                            <v-row align="center">
                                <v-col>
                                    <v-row align="center">

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.group_id" variant="outlined"
                                                label="Group ID" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.asset_id" variant="outlined"
                                                label="Asset ID" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.domain_name" variant="outlined"
                                                label="Domains" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.ip_address" variant="outlined" label="IPs"
                                                readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.as_number" variant="outlined" label="ASNs"
                                                readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.org_name" variant="outlined"
                                                label="Organizations" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.cloud" variant="outlined" label="Clouds"
                                                readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.cdn" variant="outlined" label="CDNs"
                                                readonly></v-text-field>
                                        </v-col>

                                    </v-row>
                                </v-col>

                                <v-col cols="12" md="3" class="d-flex flex-column" style="gap:8px;">
                                    <v-btn variant="outlined" color="primary" prepend-icon="mdi-circle-edit-outline"
                                        @click="onEditAsset(item.raw)">
                                        Asset Edit
                                    </v-btn>
                                    <v-btn variant="outlined" color="secondary" prepend-icon="mdi-open-in-new"
                                        @click="onOpenAsset(item.raw.asset_id)">
                                        Open Asset
                                    </v-btn>
                                    <v-btn v-if="isTaskLoading && isTaskAssetData" variant="outlined" color="secondary"
                                        prepend-icon="mdi-plus" @click="onAddAssetToWorkflowTask(item.raw.asset_id)">
                                        workflowtask
                                    </v-btn>

                                    <v-btn v-if="isTaskLoading && !isTaskAssetData" variant="outlined" color="secondary"
                                        prepend-icon="mdi-close" @click="onDelAssetToWorkflowTask(item.raw.asset_id)">
                                        workflowtask
                                    </v-btn>

                                    <v-menu v-if="isTaskLoading" :close-on-content-click="false" offset-y>
                                        <template #activator="{ props: menuProps }">
                                            <v-btn v-bind="menuProps" color="success">
                                                <v-icon>{{ runtimeStatusIcon(item.raw.asset_id) }}</v-icon>
                                                <v-divider vertical style="height: 24px; margin: 0 8px;"></v-divider>
                                                Runtime Actions
                                                <v-icon right>mdi-menu-down</v-icon>
                                            </v-btn>
                                        </template>

                                        <v-list>
                                            <v-list-item prepend-icon="mdi-record-circle-outline"
                                                @click="onWorkflowRuntimeAction('ready', item.raw.asset_id)"
                                                :disabled="!(assetRuntimeStatus(item.raw.asset_id) === -1)"><v-list-item-title>READY</v-list-item-title></v-list-item>
                                            <v-list-item prepend-icon="mdi-step-forward"
                                                @click="onWorkflowRuntimeAction('start', item.raw.asset_id)"
                                                :disabled="!([0, 1, 3, 4].includes(assetRuntimeStatus(item.raw.asset_id)))"><v-list-item-title>START</v-list-item-title></v-list-item>
                                            <v-list-item prepend-icon="mdi-pause"
                                                @click="onWorkflowRuntimeAction('stop', item.raw.asset_id)"
                                                :disabled="!(assetRuntimeStatus(item.raw.asset_id) === 2)"><v-list-item-title>STOP</v-list-item-title></v-list-item>
                                            <v-list-item prepend-icon="mdi-exit-to-app"
                                                @click="onWorkflowRuntimeAction('exit', item.raw.asset_id)"
                                                :disabled="!([0, 3, 4].includes(assetRuntimeStatus(item.raw.asset_id)))"><v-list-item-title>EXIT</v-list-item-title></v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-col>
                            </v-row>
                        </v-card>
                    </template>

                    <v-pagination v-model="pageManager.page" :length="pageManager.total"></v-pagination>
                </template>
            </v-data-iterator>

        </v-row>

        <v-dialog v-model="isAddDialog" max-width="800px" persistent>
            <v-card>
                <v-card-title class="d-flex justify-space-between align-center">
                    <div>Asset insert</div>
                    <v-icon @click="onAddAsset">mdi-close</v-icon>
                </v-card-title>

                <v-divider></v-divider>

                <v-card-text>
                    <AssetEditComponent loadingMode="insert" @cancel="onAddAsset"></AssetEditComponent>
                </v-card-text>
                <v-divider></v-divider>

                <v-card-actions>
                    <v-spacer />
                    <v-btn color="primary" @click="onAddAsset">cancel</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>

        <v-dialog v-model="isEditDialog" max-width="800px" persistent>
            <v-card>
                <v-card-title class="d-flex justify-space-between align-center">
                    <div>Asset Edit</div>
                    <v-icon @click="onEditAsset">mdi-close</v-icon>
                </v-card-title>

                <v-divider></v-divider>

                <v-card-text>
                    <AssetEditComponent loadingMode="loading" :loadingData="isEditAssetData" @cancel="onEditAsset">
                    </AssetEditComponent>
                </v-card-text>
                <v-divider></v-divider>

                <v-card-actions>
                    <v-spacer />
                    <v-btn color="primary" @click="onEditAsset">cancel</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </v-container>

</template>

<script setup>
import { GetAllAssetData } from "@/api/asset";
import { ref, reactive, watch, toRef, computed } from "vue";
import AssetEditComponent from "./AssetEditComponent.vue";


// review 

// 1. loading referer : 
//  a. page open: show all asset data or filter asset data show
//  b. page loading: loading asset data
const props = defineProps({
    
    isTaskAssetData: { type: Boolean, default: false },
    taskLoadingAssets: { type: Boolean, default: false },
    assetData: { type: Array, default: [] }, 
    workflowRuntimeStatus: { type: {}, default: {} },
});

const emits = defineEmits([
    'switchTaskAssetLoadTypes',
    'workflowRuntimeActions',
    'workflowtaskAddAsset',
    'workflowtaskDelAsset'
])

const assetData = ref([])
const currentWorkflowRuntimeStatus = reactive({})
const isAddDialog = ref(false)
const isEditDialog = ref(false)
const isEditAssetData = reactive({})
const isTaskAsset = ref(props.isTaskAssetData)
const bannerShow = ref(false)
const bannerMessages = ref("nobady")

async function initialComponent() {
    if (props.taskLoadingAssets) {
        assetData.value = props.assetData;
        listenToParentAssetData();
        listenToParentWorkflowRuntimeStatus();
        initialTaskSwitchMode();
        if (props.assetData.length <= 0) {
            bannerShow.value = false
            bannerMessages.value = `asset data loading is not show data`
        }

    } else {
        await getAllAssetData();
    }
}

async function getAllAssetData() {
    const response = await GetAllAssetData();
    console.log("get all asset data: ", response.data)
    if (response.data != null) {
        assetData.value = response.data;
        return
    }
}

const pageManager = reactive({
    limit: 5,
    page: 1,
    total: 5,
})

const loadAssetDataClass = computed(()=>{
    return isTaskAsset.value ? "Not Task Asset " : "In Task Asset"
})

const isTaskLoading = computed(() => {
    return props.taskLoadingAssets;
})

watch(() => pageManager.page, (newPage) => {
    console.log(`switch to page ${newPage}`);
})

// asset data refersh -> page refersh
watch(() => assetData.value,(newAssetData) => {
    pageManager.total = Math.ceil(newAssetData.length / pageManager.limit);
    console.log(`asset data pagination refresh: ${pageManager.total}`)
})


async function listenToParentAssetData() {
    // 这里是建立了一个 watch 用于观察 parent component 的 asset data 变化，并同步更新到 asset list component 中，这样做的目的是因为 workflowtask 的数据加载可能是动态的，是有必要保留的
    const parentComponentAssetData = toRef(props, 'assetData');
    watch(parentComponentAssetData.value, (newAssetData) => {
        console.log(`listen parent component asset data update: ${JSON.stringify(newAssetData)}`)
        if (Array.isArray(newAssetData)) {
            assetData.value.splice(0, assetData.value.length, ...(newAssetData || []));
        } else {
            assetData.value.splice(0, assetData.value.length, ...([]));
        }
    }, { immediate: true, deep: false });
}

async function listenToParentWorkflowRuntimeStatus() {
    watch(props.workflowRuntimeStatus, (newStatus) => {
        console.log(`listen workflow runtime status update: ${JSON.stringify(newStatus)}`)
        Object.assign(currentWorkflowRuntimeStatus, newStatus)
    }, { immediate: true, deep: false });
}

function assetRuntimeStatus(asset_id) {


    if (typeof currentWorkflowRuntimeStatus === 'object' && 'queues' in currentWorkflowRuntimeStatus) {
        if (currentWorkflowRuntimeStatus.queues && asset_id in currentWorkflowRuntimeStatus.queues) {
            console.log(`current asset ${asset_id} workflow runtime status: ${JSON.stringify(currentWorkflowRuntimeStatus)}`)
            return currentWorkflowRuntimeStatus.queues[asset_id];
        }
    }
    return -1;

}


function runtimeStatusIcon(asset_id) {
    console.log(`the asset ${asset_id} runtime status toggle to ${assetRuntimeStatus(asset_id)}`)
    switch (assetRuntimeStatus(asset_id)) {
        case -1:
            return "mdi-exit-to-app"
        case 0: //Created
            return "mdi-record-circle-outline"
        case 1: //Waiting
            return "mdi-cached"
        case 2: //Running
            return "mdi-step-forward"
        case 3: //Pause
        case 4: //Stop
            return "mdi-pause"
        case 5: //Done
            return "mdi-check"
        case 6: //Exit
            return "mdi-location-exit"
        case 7: //Exception
        case 8: //Exception
            return "mdi-alert-circle-outline"

    }
}

async function initialTaskSwitchMode() {
    // 这里是建立一个 watch 用于观察 parent component 的 asset types 变化，一般情况下，isTaskAsset 的变化同时会带动 asset data 的更新
    watch(() => isTaskAsset.value, (newState) => {
        console.log("component switch task asset load type: ", newState)
        emits('switchTaskAssetLoadTypes', newState)
    });
    watch(() => props.isTaskAssetData, (newState) =>{
        isTaskAsset.value = newState;
    })
}

function onWorkflowRuntimeAction(option, asset_id) {
    emits('workflowRuntimeActions', option, asset_id)
}

function onOpenAsset(asset_id) {
    const domains = window.location.origin;
    const assetId = asset_id;
    window.open(`${domains}/asset/${assetId}`,'_blank');
}

function onAddAsset() {
    isAddDialog.value = !isAddDialog.value
}

function onAddAssetToWorkflowTask(asset_id) {
    emits('workflowtaskAddAsset', asset_id)
}

function onDelAssetToWorkflowTask(asset_id) {
    emits('workflowtaskDelAsset', asset_id)
}
function onEditAsset(assetData) {
    if (isEditDialog.value === false) {
        Object.assign(isEditAssetData,assetData)
    } 
    if (isEditDialog.value === true) {
        Object.assign(isEditAssetData,{})
    }
    isEditDialog.value = !isEditDialog.value
}

initialComponent()

</script>