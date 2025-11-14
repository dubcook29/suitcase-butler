<template>
    <v-row style="padding:0px 10px; padding-top: 5px; background-color: azure;" :key="keys">
        <v-col cols="3">
            <v-sheet>{{ data.value }}</v-sheet>
        </v-col>
    </v-row>
</template>

<script setup>

import { reactive, watch, toRefs } from 'vue'

const props = defineProps({
    data: { type: Object },
    type: { type: String },
    keys: { type: String },
});

</script>