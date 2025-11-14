<template>
    <v-container v-if="true">
        <v-data-table :headers="headers" :items="data" class="custom-data-table">
            <template v-slot:item="{ item }">
                <tr>
                    <wmp-tables :data="item"></wmp-tables>

                </tr>
            </template>
        </v-data-table>
    </v-container>
</template>

<script setup>

import { ref, toRef, watch } from 'vue';
import WmpTables from './WmpTables.vue';

const props = defineProps({
    data: { type: Array, default: () => [] },
})

const data = ref([]);
const headers = ref([])

watch(() => props.data, (newData) => {
    console.log(`wmp-data value update to ${JSON.stringify(newData)}`)
    data.value.splice(0, data.value.length, ...(newData || []))
    headers.value.splice(0, headers.value.length, ...([])) // clear current headers
    if (newData.length > 0) {
        const line = newData[0]
        addHeaders(line); // refresh headers
    }
}, { immediate: true }); 


function addHeaders(item, propsHeader = null) {
    console.log(`headers open: ${JSON.stringify(item)} and add to ${JSON.stringify(propsHeader)}`)
    if (typeof item != 'object') {
        return
    }
    for (const key in item) {
        
        const value = item[key]
        const header = { title: key, value: key, key: key, align: 'center' }

        console.log(`header ${JSON.stringify(header)} created successed`)

        console.log(`value of headers :${JSON.stringify(value)}`)

        if (Array.isArray(value)) {
            addHeaders(value[0], header)
        } else if (typeof value === 'object') {
            addHeaders(value, header)
        } else {
            if (propsHeader != null) {
                propsHeader.children = propsHeader.children ? propsHeader.children : []
                propsHeader.children.push(header)
                continue
            } else {
                headers.value.push(header)
                continue
            }
        }

        if (propsHeader != null) {
            propsHeader.children = propsHeader.children ? propsHeader.children : []
            propsHeader.children.push(header)
            continue
        } else {
            headers.value.push(header)
            continue
        }

    }
}
</script>

<style scoped>
.custom-data-table .v-data-table-header th {
    text-align: center;
    vertical-align: middle;
    font-weight: bold;
    background-color: #d3d3d3;
}
</style>