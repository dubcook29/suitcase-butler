<template>
    <v-container v-if="true">
        <v-data-table :headers="headers" :items="data2" class="custom-data-table">
            <template v-slot:item="{ item }">
                <tr>
                    <wmp-tables :data="item"></wmp-tables>

                </tr>
            </template>
        </v-data-table>
    </v-container>
</template>

<script setup>
import { reactive, ref } from 'vue';
import WmpTables from './WmpTables.vue';


const data2 = ref([
        {
            "metadata": {
                "asset_id": "TestWMPRequestDataWriter_AssetId",
                "identifier": "9e715d18-4f7e-4875-9f51-3692d61c4ce6",
                "data_model_key": "dns",
                "created_at": "2025-08-20T11:43:37.177Z",
                "updated_at": "2025-08-20T11:43:37.177Z",
                "deleted_at": "0001-01-01T00:00:00Z"
            },
            "value": [
                {
                    "Key": "host",
                    "Value": "leavesongs.com"
                },
                {
                    "Key": "a",
                    "Value": [
                        "155.248.164.105"
                    ]
                },
                {
                    "Key": "mx",
                    "Value": [
                        "mx01.mail.icloud.com",
                        "mx02.mail.icloud.com"
                    ]
                },
                {
                    "Key": "ns",
                    "Value": [
                        "ns1.vultr.com",
                        "ns2.vultr.com"
                    ]
                },
                {
                    "Key": "txt",
                    "Value": [
                        "v=spf1 include:icloud.com ~all",
                        "apple-domain=eLjOdduEsCfjjOTj"
                    ]
                },
                {
                    "Key": "ptr",
                    "Value": null
                },
                {
                    "Key": "aaaa",
                    "Value": null
                },
                {
                    "Key": "cname",
                    "Value": null
                }
            ],
            "formats": "wmp_dns.DNSFormats",
            "from": "ffffffff-ffff-ffff-0002-b0bd42efd9ae"
        },
        {
            "metadata": {
                "asset_id": "TestWMPRequestDataWriter_AssetId",
                "identifier": "fcad7625-ecb8-438a-94ba-1abf535720fc",
                "data_model_key": "dns",
                "created_at": "2025-08-23T15:42:38.884Z",
                "updated_at": "2025-08-23T15:42:38.884Z",
                "deleted_at": "0001-01-01T00:00:00Z"
            },
            "value": [
                {
                    "Key": "host",
                    "Value": "leavesongs.com"
                },
                {
                    "Key": "a",
                    "Value": [
                        "155.248.164.105"
                    ]
                },
                {
                    "Key": "mx",
                    "Value": [
                        "mx01.mail.icloud.com",
                        "mx02.mail.icloud.com"
                    ]
                },
                {
                    "Key": "ns",
                    "Value": [
                        "ns1.vultr.com",
                        "ns2.vultr.com"
                    ]
                },
                {
                    "Key": "txt",
                    "Value": [
                        "apple-domain=eLjOdduEsCfjjOTj",
                        "v=spf1 include:icloud.com ~all"
                    ]
                },
                {
                    "Key": "ptr",
                    "Value": null
                },
                {
                    "Key": "aaaa",
                    "Value": null
                },
                {
                    "Key": "cname",
                    "Value": null
                }
            ],
            "formats": "wmp_dns.DNSFormats",
            "from": "ffffffff-ffff-ffff-0002-b0bd42efd9ae"
        },
        {
            "metadata": {
                "asset_id": "TestWMPRequestDataWriter_AssetId",
                "identifier": "d6fb3fea-eebe-4f52-b001-67fb28e3d0d3",
                "data_model_key": "dns",
                "created_at": "2025-10-13T12:03:15.39Z",
                "updated_at": "2025-10-13T12:03:15.39Z",
                "deleted_at": "0001-01-01T00:00:00Z"
            },
            "value": [
                {
                    "Key": "host",
                    "Value": "leavesongs.com"
                },
                {
                    "Key": "a",
                    "Value": [
                        "155.248.164.105"
                    ]
                },
                {
                    "Key": "mx",
                    "Value": [
                        "mx02.mail.icloud.com",
                        "mx01.mail.icloud.com"
                    ]
                },
                {
                    "Key": "ns",
                    "Value": [
                        "ns1.vultr.com",
                        "ns2.vultr.com"
                    ]
                },
                {
                    "Key": "txt",
                    "Value": [
                        "apple-domain=eLjOdduEsCfjjOTj",
                        "v=spf1 include:icloud.com ~all"
                    ]
                },
                {
                    "Key": "ptr",
                    "Value": null
                },
                {
                    "Key": "aaaa",
                    "Value": null
                },
                {
                    "Key": "cname",
                    "Value": null
                }
            ],
            "formats": "wmp_dns.DNSFormats",
            "from": "ffffffff-ffff-ffff-0002-b0bd42efd9ae"
        }
    ])

const headers = ref([])

addHeaders(data2.value[0])


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
console.log("print header list 11311=> ", JSON.stringify(headers.value, null, 4))

</script>

<style scoped>
.custom-data-table .v-data-table-header th {
    text-align: center;
    vertical-align: middle;
    font-weight: bold;
    background-color: #d3d3d3;
}
</style>