<template>

    <v-card variant="outlined" style="margin-top: 20px;">
        <v-card-title>
            <v-img v-if="data.wmp_basic.plugin_avatar && data.wmp_basic.plugin_avatar.length > 0"
                :src="`data:image/png;base64,${data.wmp_basic.plugin_avatar}`" max-width="100" max-height="100"
                contain></v-img>
            <span class="ml-3">Work Method Plugin</span>
        </v-card-title>

        <v-card-item class="w-100">
            <v-card-subtitle style="margin-bottom: 10px;">
                <span>Basic</span>
            </v-card-subtitle>
            <v-row>
                <v-col cols="12" md="4">
                    <v-text-field v-model="data.wmp_basic.id" variant="outlined" label="ID" readonly></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="data.wmp_basic.name" variant="outlined" label="Name" readonly></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="data.wmp_basic.version" variant="outlined" label="Version"
                        readonly></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="data.wmp_basic.push_state" variant="outlined" label="Push State"
                        readonly></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="data.wmp_basic.open_source" variant="outlined" label="Open Source"
                        readonly></v-text-field>
                </v-col>

                <v-col cols="12" md="4">
                    <v-text-field v-model="data.wmp_basic.copyright" variant="outlined" label="Copyright"
                        readonly></v-text-field>
                </v-col>

                <v-col cols="12" md="12">
                    <v-text-field v-model="data.wmp_basic.description" variant="outlined" label="Description"
                        readonly></v-text-field>
                </v-col>

                <v-col cols="12" md="12">
                    <v-text-field v-model="data.wmp_basic.cretificate" variant="outlined" label="Cretificate"
                        readonly></v-text-field>
                </v-col>
            </v-row>
        </v-card-item>


        <v-card-item class="w-100">
            <v-card-subtitle style="margin-bottom: 10px;">
                <span>Request</span>
            </v-card-subtitle>
            <v-row v-for="(value, key) in data.wmp_request" :key="key">
                <CustomeDataValueShowComponent cols="12" md="12" v-for="vv, kk in data.wmp_request" :label="kk" :value="vv" @update="(newdata) =>{
                    data.wmp_request[kk] = newdata
                }"></CustomeDataValueShowComponent>
            </v-row>
        </v-card-item>

        <v-card-item class="w-100">
            <v-card-subtitle style="margin-bottom: 10px;">
                <span>Response</span>
            </v-card-subtitle>
            <v-row v-for="(value, key) in data.wmp_response" :key="key">
                <CustomeDataValueShowComponent cols="12" md="12" v-for="vv, kk in data.wmp_response" :label="kk" :value="vv" @update="(newdata) =>{
                    data.wmp_response[kk] = newdata
                }"></CustomeDataValueShowComponent>
            </v-row>
        </v-card-item>

        <v-card-item class="w-100">
            <v-card-subtitle style="margin-bottom: 10px;">
                <span>Custome</span>
            </v-card-subtitle>
            <v-row>
                <CustomeDataValueShowComponent cols="12" md="12" v-for="vv, kk in data.wmp_custom" :label="kk" :value="vv.value" :description="vv.description" :required="vv.required" @update="(newdata) =>{
                    data.wmp_custom[kk].value = newdata
                }"></CustomeDataValueShowComponent>
            </v-row>
        </v-card-item>

    </v-card>


</template>

<script setup>
import { reactive, toRef, watch } from 'vue'
import CustomeDataValueShowComponent from './custom/CustomeDataValueShowComponent.vue'

const props = defineProps({
    wmp: { type: {}, default: {} },
})

const data = toRef(props.wmp);

</script>