import { bulterAPI } from './index.js';

export function AssetDataQueryer() {
    
}

// 
export function GetAllAssetData() {
    return bulterAPI.get('/asset/')
}

// 
export function GetAssetData(asset_id) {
    return bulterAPI.get(`/asset/${asset_id}`)
}

// 
export function AddAssetData(data) {
    return bulterAPI.post(`/asset/`,data)
}

export function EditAssetData(asset_id, data) {
    return bulterAPI.post(`/asset/${asset_id}`,data)
}

export function DelAllAssetData() {
    return bulterAPI.delete(`/asset/`)
}

export function DelAssetData(asset_id) {
    return bulterAPI.delete(`/asset/${asset_id}`)
}

