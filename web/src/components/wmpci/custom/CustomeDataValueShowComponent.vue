<template>
    <v-col>
        <v-text-field v-if="valueDataTypes === 'string'" v-model="localValue" variant="outlined"
            :label="localLabel" :hint="localHint" persistent-hint :required="localRequired"></v-text-field>
        <v-number-input v-if="valueDataTypes === 'int'" v-model="localValue" variant="outlined"
            :label="localLabel" :hint="localHint" persistent-hint :required="localRequired"></v-number-input>
        <v-switch v-if="valueDataTypes === 'bool'" v-model="localValue" inset 
            :label="localLabel" :hint="localHint" persistent-hint :required="localRequired"></v-switch>

        <v-card v-if="valueDataTypes === 'array' || valueDataTypes === 'json'" style="width: 100%;"
            :subtitle="localLabel">
            <CustomeDataValueShowComponent cols="12" md="12" v-for="vv, kk in localValue" :value="localValue[kk]"
                :key="kk" :label="`${localLabel} / ${kk}`" @update="(newdata) => {
                    console.log(`print array data <${localLabel} / ${kk}> ${localValue[kk]} update to ${newdata}`);
                    localValue[kk] = newdata;
                }"></CustomeDataValueShowComponent>

            <v-divider></v-divider>

            <v-card-actions class="justify-end align-center">
                <v-btn color="success" prepend-icon="mdi-code-json">JSON Editer</v-btn> 
                
                <v-divider vertical style="height: 30px;"></v-divider>

                <v-menu v-if="valueDataTypes === 'array'" style="width: 100%;" :close-on-content-click="false" offset-y>
                    <template #activator="{ props: menuProps }">
                        <v-btn v-bind="menuProps" color="success">
                            <v-icon>mdi-plus-box-outline</v-icon>
                            <v-icon right>mdi-menu-down</v-icon>
                        </v-btn>
                    </template>

                    <v-list>
                        <v-list-item prepend-icon="mdi-code-array" @click="() => {
                            localValue.push([])
                        }"><v-list-item-title>Array</v-list-item-title></v-list-item>
                        <v-list-item prepend-icon="mdi-code-json" @click="() => {
                            localValue.push([])
                        }"><v-list-item-title>JSON</v-list-item-title></v-list-item>
                        <v-list-item prepend-icon="mdi-code-string" @click="() => {
                            localValue.push('')
                        }"><v-list-item-title>string</v-list-item-title></v-list-item>
                        <v-list-item prepend-icon="mdi-check" @click="() => {
                            localValue.push(false)
                        }"><v-list-item-title>boolean</v-list-item-title></v-list-item>
                        <v-list-item prepend-icon="mdi-numeric" @click="() => {
                            localValue.push(0)
                        }"><v-list-item-title>number</v-list-item-title></v-list-item>
                    </v-list>
                </v-menu>
            </v-card-actions>
        </v-card>
    </v-col>
</template>

<script setup>
import { onMounted, reactive, ref, toRef, watch } from 'vue';
import CustomeDataValueShowComponent from './CustomeDataValueShowComponent.vue'


const props = defineProps({
    label: { type: String },
    value: { default: null },
    description: { type: String, default: null },
    required: { type: Boolean, default: false },
})

const emits = defineEmits(
    ['update']
)

console.log(`key ${props.key} -> value ${props.value}`)

const valueDataTypes = ref('null')
const state = ref(false)

function verifyDataType(value) {
    if (value === null) {
        return null
    } else if ( Array.isArray(value) ){
        valueDataTypes.value = 'array'
    } else if (typeof value === 'object' ){
        valueDataTypes.value = 'json'
    } else if (typeof value === 'boolean') {
        valueDataTypes.value = 'bool'
    } else if (typeof value === 'number') {
        valueDataTypes.value = 'int'
    } else if (typeof value === 'string') {
        valueDataTypes.value = 'string'
    }
}

const localLabel = toRef(props.label)
const localValue = toRef(props.value)
const localHint = toRef(props.description)
const localRequired = toRef(props.required)

watch(()=>localValue.value ,(newdata) =>{
    emits('update',newdata);
},{deep:true})

verifyDataType(localValue.value);

</script>