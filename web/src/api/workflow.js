import { bulterAPI } from './index.js';

export function GetAllWorkflowTask() {
    return bulterAPI.get(`/task/`)
}

export function GetWorkflowTask(task_id) {
    return bulterAPI.get(`/task/${task_id}`)
}

export function EditWorkflowTask(task_id, data) {
    return bulterAPI.post(`/task/${task_id}`, data)
}

export function DelWorkflowTask(task_id) {
    return bulterAPI.delete(`/task/${task_id}`)
}

export function DelAllWorkflowTask() {
    return bulterAPI.delete(`/task/`)
}


export function AddWorkflowTask(data) {
    return bulterAPI.post(`/task/`, data)
}

// get asset data in all tasks with specified task_id
export function WorkflowTaskAllAssetData(task_id) {
    return bulterAPI.get(`/task/${task_id}/asset/`)
}

// get asset data in all non-tasks with specified task_id
export function WorkflowTaskNotAssetData(task_id) {
    return bulterAPI.get(`/task/${task_id}/notasset/`)
}

// add the specified asset to the specified task
export function AddAssetToWorkflowTask(task_id, asset_id) {
    return bulterAPI.post(`/task/${task_id}/asset/${asset_id}`)
}

// delete the specified asset from the specified task
export function DeleteAssetFromWorkflowTask(task_id, asset_id) {
    return bulterAPI.delete(`/task/${task_id}/asset/${asset_id}`)
}

// get the status attribute of the current workflow 
export function WorkflowRuntimeGetQueueStatus() {
    return bulterAPI.get(`/workflow`)
}

export function GetWorkflowRuntimeQueueStatus() {
    return bulterAPI.get(`/workflow`)
}

// get all current workflow runtime tasks
export function GetAllWorkflowRuntimeTask() {
    return bulterAPI.get(`/workflow/`)
}

// revert the specified workflow task to workflow runtime task
export function WorkflowRuntimeLoadTask(task_id) {
    return bulterAPI.get(`/workflow/${task_id}`)
}

// remove the specified task from the workflow runtime 
export function WorkflowRuntimeDeleteTask(task_id) {
    return bulterAPI.delete(`/workflow/${task_id}`)
}

// load the specified asset into the workflow runtime task
export function WorkflowRuntimeAsAddAsset(task_id, asset_id) {
    return bulterAPI.get(`/workflow/${task_id}/add/${asset_id}`)
}

// load the all asset into the workflow runtime task
export function WorkflowRuntimeAsAddAllAsset(task_id) {
    return bulterAPI.get(`/workflow/${task_id}/add/`)
}

// start the specified asset from the workflow runtime task
export function WorkflowRuntimeAsStartAsset(task_id, asset_id) {
    return bulterAPI.get(`/workflow/${task_id}/start/${asset_id}`)
}

// start the all asset from the workflow runtime task
export function WorkflowRuntimeAsStartAllAsset(task_id) {
    return bulterAPI.get(`/workflow/${task_id}/start/`)
}

// stop the specified asset from the workflow runtime task
export function WorkflowRuntimeAsStopAsset(task_id, asset_id) {
    return bulterAPI.get(`/workflow/${task_id}/stop/${asset_id}`)
}

// stop the all asset from the workflow runtime task
export function WorkflowRuntimeAsStopAllAsset(task_id) {
    return bulterAPI.get(`/workflow/${task_id}/stop/`)
}

// delete the specified asset from the workflow runtime task
export function WorkflowRuntimeAsDeleteAsset(task_id, asset_id) {
    return bulterAPI.get(`/workflow/${task_id}/delete/${asset_id}`)
}

// delete the all asset from the workflow runtime task
export function WorkflowRuntimeAsDeleteAllAsset(task_id) {
    return bulterAPI.get(`/workflow/${task_id}/delete/`)
}

// check health status the specified asset from the workflow runtime task
export function WorkflowRuntimeAsHealthAllAsset(task_id, asset_id) {
    return bulterAPI.get(`/workflow/${task_id}/health/${asset_id}`)
}

// 

export async function WorkflowRuntimeOptionsFunction(task_id, asset_id, method) {

    if (task_id === null || task_id === "" || task_id === undefined) {
        return {
            "code": 101,
            "err": "task_id is null"
        }
    }

    if (method === "ready") {
        if (asset_id === null || asset_id === "" || asset_id === undefined) {
            return await WorkflowRuntimeAsAddAllAsset(task_id);
        } else {
            return await WorkflowRuntimeAsAddAsset(task_id, asset_id);
        }
    } else if (method === "start") {
        if (asset_id === null || asset_id === "" || asset_id === undefined) {
            return await WorkflowRuntimeAsStartAllAsset(task_id);
        } else {
            return await WorkflowRuntimeAsStartAsset(task_id, asset_id);
        }
    } else if (method === "stop") {
        if (asset_id === null || asset_id === "" || asset_id === undefined) {
            return await WorkflowRuntimeAsStopAllAsset(task_id);
        } else {
            return await WorkflowRuntimeAsStopAsset(task_id, asset_id);
        }
    } else if (method === "exit") {
        if (asset_id === null || asset_id === "" || asset_id === undefined) {
            return await WorkflowRuntimeAsDeleteAllAsset(task_id);
        } else {
            return await WorkflowRuntimeAsDeleteAsset(task_id, asset_id);
        }
    } else if (method === "runtime") {
        return await WorkflowRuntimeLoadTask(task_id);
    } else if (method === "notruntime") {
        return await WorkflowRuntimeDeleteTask(task_id);
    } else {
        return {
            "code": 101,
            "err": "call method (" + method + ") does not exist"
        }
    }

}

export async function WorkflowRuntimeStatus() {
    return bulterAPI.get(`/workflow`)
}