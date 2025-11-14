package manager

import (
	"context"
	"errors"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/db"
	workflowtask "github.com/suitcase/butler/workflow/task"
	"go.mongodb.org/mongo-driver/bson"
)

// WorkflowTaskRuntimeRestore: restore workflowtask from the database and load as runtime status
func (w *Manager) WorkflowTaskRuntimeRestore(task_id string) error {

	if tasks, count, err := workflowtask.WorkflowTasksFindAll(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "_id", Value: task_id}}); err != nil {
		return err
	} else if count == 0 {
		return errors.New("[workflow manager] this task cannot be found" + task_id)
	} else {
		for _, in := range tasks {
			if err := w.workflowTaskOnlineWithLock(in); err != nil {
				return err
			}
		}

	}

	return nil
}

// WorkflowTaskRuntimeDeleted: remove a specific workflowtask from the runtime state
func (w *Manager) WorkflowTaskRuntimeDeleted(task_id string) error {

	if err := w.workflowTaskOfflineWithLock(task_id); err != nil {
		return err
	} else {
		return nil
	}
}

// WorkflowTaskRuntimeQueryer: get the information data of Runtime workflowtask for display
func (w *Manager) WorkflowTaskRuntimeQueryer(task_id ...string) []workflowtask.WorkflowTasks {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var out []workflowtask.WorkflowTasks

	if task_id != nil {
		for _, id := range task_id {
			if v, ok := w.workflowTaskQueue[id]; ok {
				out = append(out, v.WorkflowTasks().Model())
			}
		}
	} else {
		for _, v := range w.workflowTaskQueue {
			out = append(out, v.WorkflowTasks().Model())
		}
	}

	return out
}

func (w *Manager) workflowTaskOnlineWithLock(in workflowtask.WorkflowTasks) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.workflowTaskQueue[in.TaskID]; ok {
		return errors.New("[workflow manager] the task queue has duplicate conflicts [" + in.TaskID + "::" + in.TaskName + "]")
	} else {
		if w.nowWorkflowTask < w.maxWorkflowTask {
			w.wg.Add(1)
			w.workflowTaskQueue[in.TaskID] = new(workflowtask.WorkflowTaskRuntime).Initial(w.ctx, 10).InitialGridManager(w.gridManager).InitialWorkflowSessions(w.workflowSessions).InitialWorkflowTasks(&in) // online
			w.nowWorkflowTask = w.nowWorkflowTask + 1
			logrus.Debugf("[workflow manager] the workflow task runtime has been online, now (%d/%d)", w.nowWorkflowTask, w.maxWorkflowTask)
			if _, err := in.UpdateOne(context.TODO(), db.GetCurrentMongoClient()); err != nil {
				return err
			}
		} else {
			return errors.New("[workflow manager] maximum number of task runtime queues exceeded, now (" + strconv.Itoa(w.nowWorkflowTask) + "/" + strconv.Itoa(w.maxWorkflowTask) + ")")
		}
	}
	return nil
}

func (w *Manager) workflowTaskOfflineWithLock(task_id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if runtime, ok := w.workflowTaskQueue[task_id]; ok {
		err := runtime.ExitRuntime()
		delete(w.workflowTaskQueue, task_id)
		if w.nowWorkflowTask > 0 {
			defer w.wg.Done()
			w.nowWorkflowTask = w.nowWorkflowTask - 1 // offline
			logrus.Debugf("[workflow manager] the workflow task runtime has been offline, now (%d/%d)", w.nowWorkflowTask, w.nowWorkflowTask)
		}
		if err != nil {
			return err
		}
	} else {
		return errors.New("[workflow manager] not available in the task runtime queue [" + task_id + "]")
	}

	return nil
}
