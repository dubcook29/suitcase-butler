<template>
    <v-container>
        <v-row>
            <v-col cols="12" md="4">
                <v-text-field variant="outlined" v-model="task_data.task_id" label="Task ID" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field variant="outlined" v-model="task_data.task_name" label="Task Name"
                    readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field variant="outlined" v-model="task_data.status" label="Task Status" readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-text-field variant="outlined" v-model="task_data.task_description" label="Task Description"
                    readonly></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
                <v-select variant="outlined" label="Scheduler ID" v-model="task_data.scheduler_id" readonly></v-select>
            </v-col>
        </v-row>

        <v-row justify="end">

            <v-col cols="auto" v-if="currentRuntimeState">
                <v-btn variant="elevated" prepend-icon="mdi-record-circle-outline" :disabled="task_data.status !== 2"
                    @click="readyButtonClick">
                    READY
                </v-btn>
            </v-col>

            <v-col cols="auto" v-if="currentRuntimeState">
                <v-btn variant="elevated" prepend-icon="mdi-step-forward" :disabled="task_data.status !== 2"
                    @click="startButtonClick">
                    START
                </v-btn>
            </v-col>

            <v-col cols="auto" v-if="currentRuntimeState">
                <v-btn variant="elevated" prepend-icon="mdi-pause" :disabled="task_data.status !== 2"
                    @click="stopButtonClick">
                    STOP
                </v-btn>
            </v-col>

            <v-col cols="auto" v-if="currentRuntimeState">
                <v-btn variant="elevated" prepend-icon="mdi-exit-to-app" :disabled="task_data.status !== 2"
                    @click="exitButtonClick">
                    EXIT
                </v-btn>
            </v-col>

            <v-col cols="auto">
                <v-btn variant="elevated" prepend-icon="mdi-power" :disabled="task_data.status == 0"
                    @click="runtimeButtonClick">
                    {{ task_data.status == 1 ? "open" : "close" }}
                </v-btn>
            </v-col>

            <v-col cols="auto">
                <v-btn variant="elevated" prepend-icon="mdi-clipboard-list-outline" @click="assetListBtnClick">
                    Assets
                </v-btn>
            </v-col>

            <v-col cols="auto">
                <v-btn variant="elevated" color="cyan" prepend-icon="mdi-square-edit-outline" @click="editButtonClick">
                    EDIT
                </v-btn>
            </v-col>

            <v-col cols="auto">
                <v-btn variant="elevated" color="error" prepend-icon="mdi-delete-outline"
                    :disabled="currentRuntimeState" @click="deleteButtonClick">
                    DELETE
                </v-btn>
            </v-col>
        </v-row>

        <AssetListComponent v-if="assetListShowState" :taskLoadingAssets="true"
            :assetData="task_asset_lists" :isTaskAssetData=assetListShowNonTask
            :workflowRuntimeStatus="currentWorkflowRuntimeStatus" @switchTaskAssetLoadTypes="toggleAssetListShowNonTask"
            @workflowRuntimeActions="assetComponentRuntimeButtonClick" @workflowtaskAddAsset="workflowtaskInstall"
            @workflowtaskDelAsset="workflowtaskUninstall"></AssetListComponent>

    </v-container>

    <!-- dialog edit workflowtask -->
    <v-dialog v-model="currentEditDialogState" max-width="800px" persistent>
        <v-card>
            <v-card-title class="d-flex justify-space-between align-center">
                <div>Workflow Task Edit</div>
                <v-icon @click="dialogEditComponetCloseButtonClick">mdi-close</v-icon>
            </v-card-title>

            <v-divider></v-divider>

            <v-card-text>
                <WorkflowTaskEditComponent :task_data="task_data" @saved="dialogEditComponetSaveButton"
                    @cancel="dialogEditComponetCloseButtonClick">
                </WorkflowTaskEditComponent>
            </v-card-text>
            <v-divider></v-divider>

            <v-card-actions>
                <v-spacer />
                <v-btn color="primary" @click="dialogEditComponetCloseButtonClick">取消</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup>

import { ref, reactive, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import WorkflowTaskEditComponent from './WorkflowTaskEditComponent.vue';


import {
    GetWorkflowTask,
    WorkflowTaskAllAssetData,
    WorkflowTaskNotAssetData,
    AddAssetToWorkflowTask,
    DeleteAssetFromWorkflowTask,
    EditWorkflowTask,
    WorkflowRuntimeOptionsFunction,
    WorkflowRuntimeStatus,
    DelWorkflowTask
} from '@/api/workflow'
import AssetListComponent from './AssetListComponent.vue';



const props = defineProps({
    task_data: { type: Object, default: null },
});

const route = useRoute()
const taskId = computed(() => route.params?.task_id ?? route.params?.taskId ?? null)

console.log("task id => ", taskId.value);

const default_task_data = {
    "task_id": "Defult Task ID",
    "task_name": "Defult Task Name",
    "task_description": "Defult Task Description",
    "scheduler_id": "Defult Task Scheduler ID",
    "status": -65535,
    "task_asset_queue": [
        "TestWMPRequestDataWriter_AssetId"
    ]
};

const task_data = reactive(props.task_data == null ? {} : props.task_data);

async function getWorkflowTaskData(task_id) {
    const response = await GetWorkflowTask(task_id);

    if (Array.isArray(response.data) && response.data.length > 0) {
        Object.assign(task_data, (response.data[0] || {}));
    } else {
        Object.assign(task_data, {});
    }

    console.log(`get workflow task data ok ${JSON.stringify(task_data)}`)

}

if (props.task_data == null && taskId.value != null) {
    getWorkflowTaskData(taskId.value)
} else {
    Object.assign(task_data, default_task_data);
};

const task_asset_lists = ref([])

// refresh the asset table outside the task
async function getAllAssetDataNotTask(task_id) {
    const response = await WorkflowTaskNotAssetData(task_id);
    console.log("loading asset not in task: ", response.data);
    task_asset_lists.value.splice(0, task_asset_lists.value.length, ...(response.data || []))
}

// Refresh the asset table in the task
async function getAllAssetDataInTesk(task_id) {
    const response = await WorkflowTaskAllAssetData(task_id);
    console.log("loading asset in a task: ", response.data);
    if ((response.data || []).length === 0) {
        await getAllAssetDataNotTask(task_id);
        assetListShowNonTask.value = true;
    } else {
        task_asset_lists.value.splice(0, task_asset_lists.value.length, ...(response.data || []))
    }
}

// maintain and manage the display status of assets through button and mouse click operations
const assetListShowState = ref(false);
async function assetListBtnClick() {
    console.log(`asset list show state switch ${assetListShowState.value} -> ${!assetListShowState.value} ... `)
    if (!assetListShowState.value) {
        await assetListRefresh();
    }
    assetListShowState.value = !assetListShowState.value;
}

// maintain and manage the data types of asset lists (assets in task asset lists and assets in non-task asset lists)
const assetListShowNonTask = ref(false);
async function assetListRefresh() {
    console.log("asset list refersh ... ")
    assetListShowNonTask.value ? await getAllAssetDataNotTask(task_data.task_id) : await getAllAssetDataInTesk(task_data.task_id);
}

async function toggleAssetListShowNonTask(newState) {
    console.log("asset list show type toggle to ", newState);
    assetListShowNonTask.value = newState;
    await assetListRefresh();
}

async function workflowtaskInstall(asset_id) {
    console.log(`workflow task ${task_data.task_id} install asset ${asset_id}`);
    const response = await AddAssetToWorkflowTask(task_data.task_id, asset_id);
    if (response.code === 200) {
        await assetListRefresh();
    } 
    await getWorkflowTaskData(task_data.task_id);
}

async function workflowtaskUninstall(asset_id) {
    console.log(`workflow task ${task_data.task_id} uninstall asset ${asset_id}`);
    const response = await DeleteAssetFromWorkflowTask(task_data.task_id, asset_id);
    if (response.code === 200) {
        await assetListRefresh();
    }
    await getWorkflowTaskData(task_data.task_id);
}

// ready button
async function readyButtonClick() {
    const response = await WorkflowRuntimeOptionsFunction(task_data.task_id, null, "ready");
    console.log(`All assets in the workflow task ${task_data.task_name} perform the [ready] operation: ${JSON.stringify(response)}`)

}

// start button
async function startButtonClick() {
    const response = await WorkflowRuntimeOptionsFunction(task_data.task_id, null, "start");
    console.log(`All assets in the workflow task ${task_data.task_name}  perform the [start] operation: ${JSON.stringify(response)}`)
}

// stop button
async function stopButtonClick() {
    const response = await WorkflowRuntimeOptionsFunction(task_data.task_id, null, "stop");
    console.log(`All assets in the workflow task ${task_data.task_name} perform the [stop] operation: ${JSON.stringify(response)}`)
}

// exit button
async function exitButtonClick() {
    const response = await WorkflowRuntimeOptionsFunction(task_data.task_id, null, "exit");
    console.log(`All assets in the workflow task ${task_data.task_name} perform the [exit] operation: ${JSON.stringify(response)}`)
}


// runtime open/close button
const currentRuntimeState = computed(() => task_data.status == 2)
async function runtimeButtonClick() {
    console.log(`current runtime state toggle to ${!currentRuntimeState.value}`)
    if (task_data.status == 2) {
        const response = await WorkflowRuntimeOptionsFunction(task_data.task_id, null, "notruntime");
        console.log(`The workflow task ${task_data.task_name} perform the [notruntime] operation: ${JSON.stringify(response)}`)
    } else {
        const response = await WorkflowRuntimeOptionsFunction(task_data.task_id, null, "runtime");
        console.log(`The workflow task ${task_data.task_name} perform the [runtime] operation: ${JSON.stringify(response)}`)
    }

    await getWorkflowTaskData(task_data.task_id);

}

// asset list component runtime button click process
async function assetComponentRuntimeButtonClick(option, id) {
    const response = await WorkflowRuntimeOptionsFunction(task_data.task_id, id, option);
    console.log(`Asset ${id} in the workflow task ${task_data.task_name} perform the [exit] operation: ${JSON.stringify(response)}`)
}

// delete button , delete current workflow task
async function deleteButtonClick() {
    if (currentRuntimeState.value) {
        return
    }
    const response = await DelWorkflowTask(task_data.task_id);

}

// edit button, open edit dialog window
const currentEditDialogState = ref(false);
async function editButtonClick() {
    currentEditDialogState.value = !currentEditDialogState.value
}

// edit dialog , save button
async function dialogEditComponetSaveButton(updated) {
    Object.assign(task_data, updated)
    const response = await EditWorkflowTask(task_data.task_id, task_data);
    await dialogEditComponetCloseButtonClick();
}

// edit dialog , close button
async function dialogEditComponetCloseButtonClick() {
    currentEditDialogState.value = false
}

const currentWorkflowRuntimeStatus = reactive({});
async function autoSyncWorkflowRuntimeStatus() {
    const response = await WorkflowRuntimeStatus();
    if (response.code === 200) {
        Object.assign(currentWorkflowRuntimeStatus, response.data);
    }
}

setInterval(async ()=>{
    await autoSyncWorkflowRuntimeStatus();
    console.log(`workflow task component update runtime status: ${JSON.stringify(currentWorkflowRuntimeStatus)}`)
}, 5000);
</script>
