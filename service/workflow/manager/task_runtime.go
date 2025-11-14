package manager

import (
	"errors"

	workflowtask "github.com/suitcase/butler/workflow/task"
	"github.com/suitcase/butler/workflow/workflow"
)

// WorkflowTaskRuntimeActionAdd: Specify that the runtime workflowtask adds an asset (workflow), if asset_id is not found, an error will be returned.
// Specifying an asset_id retrieves specific asset data from the database and loads it into workflow to wait for it to start running.
// If asset_id is not specified, all asset_ids in the task are retrieved from the database and loads it into workflow to wait for it to start running.
func (w *Manager) WorkflowTaskRuntimeActionAdd(task_id string, asset_id ...string) error {
	if task_runtime := w.workflowTaskRuntime(task_id); task_runtime == nil {
		return errors.New("the [" + task_id + "] task you specified not is runtime or does not exist")
	} else {
		return task_runtime.WorkflowTaskRuntimeAdd(asset_id...)
	}
}

// WorkflowTaskRuntimeActionStart: Specify that the runtime workflowtask start an asset (workflow), if asset_id is not found, an error will be returned.
// If asset_id is specified, it will be searched for in the runtime workflow_queue and executed (start),
// If asset_id is not specified, it looks for all asset_ids of the workflowtask in the runtime workflow_queue and executes (start) them,
func (w *Manager) WorkflowTaskRuntimeActionStart(task_id string, asset_id ...string) error {
	if task_runtime := w.workflowTaskRuntime(task_id); task_runtime == nil {
		return errors.New("the [" + task_id + "] task you specified not is runtime or does not exist")
	} else {
		return task_runtime.WorkflowTaskRuntimeStart(asset_id...)
	}
}

// WorkflowTaskRuntimeActionStop: Specify that the runtime workflowtask stops an asset (workflow), if asset_id is not found, an error will be returned.
func (w *Manager) WorkflowTaskRuntimeActionStop(task_id string, asset_id ...string) error {
	if task_runtime := w.workflowTaskRuntime(task_id); task_runtime == nil {
		return errors.New("the [" + task_id + "] task you specified not is runtime or does not exist")
	} else {
		return task_runtime.WorkflowTaskRuntimeStop(asset_id...)
	}
}

// WorkflowTaskRuntimeActionRestart: Specify that the runtime workflowtask restart an asset (workflow), if asset_id is not found, an error will be returned.
func (w *Manager) WorkflowTaskRuntimeActionRestart(task_id string, asset_id ...string) error {
	if task_runtime := w.workflowTaskRuntime(task_id); task_runtime == nil {
		return errors.New("the [" + task_id + "] task you specified not is runtime or does not exist")
	} else {
		return task_runtime.WorkflowTaskRuntimeRestart(asset_id...)
	}
}

// WorkflowTaskRuntimeActionDeleted: Specify that the runtime workflowtask deleted an asset (workflow), if asset_id is not found, an error will be returned.
func (w *Manager) WorkflowTaskRuntimeActionDeleted(task_id string, asset_id ...string) error {
	if task_runtime := w.workflowTaskRuntime(task_id); task_runtime == nil {
		return errors.New("the [" + task_id + "] task you specified not is runtime or does not exist")
	} else {
		return task_runtime.WorkflowTaskRuntimeDelete(asset_id...)
	}
}

func (w *Manager) WorkflowTaskRuntimeActionHealth(task_id string, asset_id ...string) ([]workflow.WorkflowRuntimeHealth, error) {
	if task_runtime := w.workflowTaskRuntime(task_id); task_runtime == nil {
		return nil, errors.New("the [" + task_id + "] task you specified not is runtime or does not exist")
	} else {
		return task_runtime.WorkflowTaskRuntimeHealth(asset_id...), nil
	}
}

// workflowTaskRuntime: get a workflowTaskRuntime object
// This is usually used to fetch a runtime workflowtask from the workflow_manager and execute the associated grass walk.
func (w *Manager) workflowTaskRuntime(task_id string) *workflowtask.WorkflowTaskRuntime {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if v, ok := w.workflowTaskQueue[task_id]; ok {
		return v
	}
	return nil
}
