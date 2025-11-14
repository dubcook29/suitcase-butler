import { bulterAPI } from './index.js';

export function GetCurrentAllSupportConnector() {
    return bulterAPI.get(`/wmpci/sessions/connector/`)
}

export function GetConnectorConfig(connector) {
    return bulterAPI.get(`/wmpci/sessions/connector/${connector}`)
}

export function PostConnectorConfig(connector, data) {
    return bulterAPI.post(`/wmpci/sessions/connector/${connector}`, data)
}

export function GetCurrentAllSessionConnect() {
    return bulterAPI.get(`/wmpci/sessions/`)
}

export function GetCurrentSessionConnect(session_id) {
    return bulterAPI.get(`/wmpci/sessions/${session_id}`)
}

export function DeleteCurrentSessionConnect(session_id) {
    return bulterAPI.delete(`/wmpci/sessions/${session_id}`)
}