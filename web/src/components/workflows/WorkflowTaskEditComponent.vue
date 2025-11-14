<template>
    <v-form>
        <v-container>
            <v-row>
                <v-col cols="12" md="4">
                    <v-text-field v-model="task_data_local.task_id" :counter="36" label="Task ID"
                        required></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="task_data_local.task_name" :counter="36" label="Task Name"
                        required></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-select v-model="task_data_local.status" :counter="2" label="Task Status"
                        :items="[0,1,2]"></v-select>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="task_data_local.task_description"
                        label="Task Description" required></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-select label="Scheduler ID" v-model="task_data_local.scheduler_id"
                        :items="scheduler_items"></v-select>
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
import { reactive, watch, toRefs } from 'vue'

const props = defineProps({
    task_data: { type: Object },
    create: { type: Boolean, default: false },
});

const isEditMode = !props.task_data
const isCreateMode = props.task_data
const task_data_local = reactive({ ...(props.task_data || {}) })

const scheduler_items = ['DufalueBuiltWMPSchedulingGrid'];

const emit = defineEmits(['saved', 'cancel'])

// refresh task_data in real time
watch(() => props.task_data, (newVal) => {
    Object.assign(task_data_local, newVal || {})
}, { deep: true, immediate: true })

// submit btn click event
function submit() {
    emit('saved', task_data_local)
}

// reset btn click event
function reset() {
    Object.assign(task_data_local, props.task_data)
}

// close btn click event
function close() {
    emit('cancel')
}


</script>
