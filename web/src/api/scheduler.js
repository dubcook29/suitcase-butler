import { bulterAPI } from './index.js';


export function GetScheduler() {
    return bulterAPI.get(`/scheduler/`)
}

export function GetSchedulerId(scheduler_id) {
    return bulterAPI.get(`/scheduler/${scheduler_id}`)
}
