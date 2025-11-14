<template>
    
    <template v-for="line,table_name in (props.data || ['N/A'])" :key="table_name">
        {{ console.log(`table td print ${line}`) }}
        <template v-if="Array.isArray(line)">
            <td v-if="typeof line[0] === 'object'" v-for="_value,name in line[0]">
                <tr v-for="value in line">
                    {{ value[name] }}
                </tr>
            </td>
            <td v-else>
                <tr v-for="value in line">
                    {{ value }}
                </tr>
            </td>
        </template>

        <wmp-tables v-else-if="typeof line === 'object'" :data="line">

        </wmp-tables>

        <td v-else>
            {{ line }}
        </td>
    </template>
</template>

<script setup>
import WmpTables from './WmpTables.vue';

const props = defineProps({
    data: {
        type: Object,
        required: true
    }
})

console.log(`wmp table data is : ${JSON.stringify(props.data)}`)
</script>