<template>
    <CustomeDataValueShowComponent label="verify array" :value="value" @update="updated">
    </CustomeDataValueShowComponent>
    <v-divider></v-divider>

    {{ value }}
</template>

<script setup>
import { ref, watch } from 'vue';
import CustomeDataValueShowComponent from './CustomeDataValueShowComponent.vue'


// const value = ref(1)
// const value = ref('s')
// const value = ref(false)
const value = ref([1,'s',false])
// const value = ref({
//     "name":"Bulter",
//     "version": 10,
//     "opensource":false,
// })


function updated(newdata) {
    console.log(`value data update to ${newdata}`);
    value.value = newdata
}
</script>