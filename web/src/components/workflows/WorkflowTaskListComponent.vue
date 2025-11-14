<template>
    <div>
        <v-row align="center" no-gutters style="padding: 20px; margin-top: 10px; width:100%">
            <v-banner v-if="bannerShow" icon="$warning" color="warning" class="mb-3">
                <template v-slot:text>
                    {{ bannerMessages }}
                </template>
            </v-banner>
            <v-data-iterator class="w-100" :items="task_list" :page="pageManager.page">
                <template v-slot:default="{ items }">
                    <v-toolbar title="Workflow Task Lists" style="margin-bottom: 10px; padding-right: 10px;">
                        <template v-slot:append>
                            <v-btn prepend-icon="mdi-plus-box-outline" @click="addButtonClick">Add</v-btn>
                            <v-btn prepend-icon="mdi-refresh" @click="loadingAllWorkflowTaskData">Refresh</v-btn>
                        </template>
                    </v-toolbar>

                    <template v-for="(item, i) in items" :key="i">
                        <v-card class="pa-4 mb-3 w-100">
                            <v-row align="center">
                                <v-col>
                                    <v-row align="center">

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.task_id" variant="outlined" label="Task ID"
                                                readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.task_name" variant="outlined"
                                                label="Task Name" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.status" variant="outlined"
                                                label="Task Status" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.scheduler_id" variant="outlined"
                                                label="Scheduler" readonly></v-text-field>
                                        </v-col>

                                        <v-col cols="12" md="3">
                                            <v-text-field v-model="item.raw.task_description" variant="outlined"
                                                label="Task Description" readonly></v-text-field>
                                        </v-col>
                                    </v-row>
                                </v-col>

                                <v-col cols="12" md="3" class="d-flex flex-column" style="gap:8px;">
                                    <v-btn variant="outlined" color="primary" prepend-icon="mdi-circle-edit-outline"
                                        @click="editButtonClick(item.raw)">
                                        Task Edit
                                    </v-btn>
                                    <v-btn variant="outlined" color="secondary" prepend-icon="mdi-open-in-new"
                                        @click="openWorkflowTaskWindows(item.raw.task_id)">
                                        Open Task
                                    </v-btn>
                                    <v-btn variant="outlined" color="secondary" prepend-icon="mdi-power"
                                        @click="runtimeButtonClick(item.raw.task_id, item.raw.status == 2)">
                                        Runtime
                                    </v-btn>
                                </v-col>
                            </v-row>
                        </v-card>
                    </template>

                    <v-pagination v-model="pageManager.page" :length="pageManager.total"></v-pagination>
                </template>
            </v-data-iterator>

        </v-row>

        <v-dialog v-model="currentAddDialogState" max-width="800px" persistent>
            <v-card>
                <v-card-title class="d-flex justify-space-between align-center">
                    <div>Asset insert</div>
                    <v-icon @click="closeTaskEditDialogState">mdi-close</v-icon>
                </v-card-title>

                <v-divider></v-divider>

                <v-card-text>
                    <WorkflowTaskEditComponent :create="true" @cancel="closeTaskEditDialogState"
                        @saved="dialogEditComponetAddButton"></WorkflowTaskEditComponent>
                </v-card-text>
                <v-divider></v-divider>

                <v-card-actions>
                    <v-spacer />
                    <v-btn color="primary" @click="closeTaskEditDialogState">cancel</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>

        <v-dialog v-model="currentEditDialogState" max-width="800px" persistent>
            <v-card>
                <v-card-title class="d-flex justify-space-between align-center">
                    <div>Asset Edit</div>
                    <v-icon @click="closeTaskEditDialogState">mdi-close</v-icon>
                </v-card-title>

                <v-divider></v-divider>

                <v-card-text>
                    <WorkflowTaskEditComponent :task_data="currentEditTaskData" @cancel="closeTaskEditDialogState"
                        @saved="dialogEditComponetSaveButton"></WorkflowTaskEditComponent>
                </v-card-text>
                <v-divider></v-divider>

                <v-card-actions>
                    <v-spacer />
                    <v-btn color="primary" @click="closeTaskEditDialogState">cancel</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </div>

</template>

<script setup>
import { GetAllWorkflowTask,EditWorkflowTask,AddWorkflowTask,WorkflowRuntimeOptionsFunction } from "@/api/workflow"
import { ref, reactive, watch, toRef, computed } from "vue";
import WorkflowTaskEditComponent from "./WorkflowTaskEditComponent.vue";


// workflow task list component codeing ...
const props = defineProps({
    
    isTaskAssetData: { type: Boolean, default: false },
    taskLoadingAssets: { type: Boolean, default: false },
    workflowtaskData: { type: Array, default: [] }, 
});


const task_list = ref([]);

async function loadingAllWorkflowTaskData() {
    const response = await GetAllWorkflowTask();
    if (response.data != null) {
        console.log("get all workflow task data: ", response.data)
        task_list.value = response.data;
        return
    }
}

loadingAllWorkflowTaskData();


const pageManager = reactive({
    limit: 5,
    page: 1,
    total: 5,
})

watch(() => pageManager.page, (newPage) => {
    console.log(`switch to page ${newPage}`);
})

// asset data refersh -> page refersh
watch(() => task_list.value,(newTaskData) => {
    pageManager.total = Math.ceil(newTaskData.length / pageManager.limit);
    console.log(`asset data pagination refresh: ${pageManager.total}`)
})


const bannerShow = ref(false)
const bannerMessages = ref("nobady")

// 
const currentEditDialogState = ref(false)
const currentAddDialogState = ref(false)
const currentEditTaskData = reactive({})
async function editButtonClick(task) {
    
    Object.assign(currentEditTaskData, task)
    currentEditDialogState.value = !currentAddDialogState.value
}

async function addButtonClick() {
    currentAddDialogState.value = !currentAddDialogState.value
}

async function closeTaskEditDialogState() {
    currentEditDialogState.value = false;
    currentAddDialogState.value = false;
}

async function openWorkflowTaskWindows(task_id) {
    const domains = window.location.origin;
    window.open(`${domains}/workflow/task/${task_id}`,'_blank');
}

async function dialogEditComponetSaveButton(updated) {
    // Object.assign(task_data, updated)
    const response = await EditWorkflowTask(updated.task_id, updated);
    await loadingAllWorkflowTaskData();
    await closeTaskEditDialogState();
}

async function dialogEditComponetAddButton(updated) {
    // Object.assign(task_data, updated)
    const response = await AddWorkflowTask(updated);
    await loadingAllWorkflowTaskData();
    await closeTaskEditDialogState();
}

async function runtimeButtonClick(task_id, status) {
    console.log(`current runtime state toggle to ${!status}`)
    if (status) {
        const response = await WorkflowRuntimeOptionsFunction(task_id, null, "notruntime");
        console.log(`The workflow task ${task_id} perform the [notruntime] operation: ${JSON.stringify(response)}`)
    } else {
        const response = await WorkflowRuntimeOptionsFunction(task_id, null, "runtime");
        console.log(`The workflow task ${task_id} perform the [runtime] operation: ${JSON.stringify(response)}`)
    }


    await loadingAllWorkflowTaskData();
}
</script>