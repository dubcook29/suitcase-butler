<template>
    <!-- this test demo -->
     <v-container>
        <v-row>
            <CustomeDataValueShowComponent cols="12" md="12" v-for="vv, kk in data" :label="kk" :value="vv.value" :description="vv.description" :required="vv.required" @update="(newdata) =>{
                console.log(`wmp custom component value update to ${newdata}`);
                data[kk].value = newdata
            }"></CustomeDataValueShowComponent>
        </v-row>

        
    </v-container>
    {{ data }}
</template>

<script setup>
import { reactive } from 'vue';
import CustomeDataValueShowComponent from './CustomeDataValueShowComponent.vue'



const data = reactive({
    "host": {
        "name": "host",
        "value": "localhost",
        "description": "json-rpc service hostname/ipaddress",
        "required": true
    },
    "port": {
        "name": "port",
        "value": "1080",
        "description": "json-rpc service port number",
        "required": true
    },
    "serviceName": {
        "name": "serviceName",
        "value": "wmp",
        "description": "json-rpc service root-name",
        "required": true
    }
});

function updated(newdata) {
    console.log(`wmp custom component value update to ${newdata}`);

}
</script>