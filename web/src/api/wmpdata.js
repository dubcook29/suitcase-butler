import { bulterAPI } from './index.js';

export function GetAssetAllWMPData(asset_id) {
    return bulterAPI.get(`/wmpci/wmpdata/${asset_id}/`)
}

export function GetAssetWMPData(asset_id,wmpdata_name) {
    return bulterAPI.get(`/wmpci/wmpdata/${asset_id}/${wmpdata_name}`)
}