package manager

import (
	"context"
	"errors"
	"sync"

	"github.com/suitcase/butler/db"
	"github.com/suitcase/butler/grid"
	workflowtask "github.com/suitcase/butler/workflow/task"
	"github.com/suitcase/butler/workflow/workflow"
	"go.mongodb.org/mongo-driver/bson"
)

type Manager struct {
	wg                sync.WaitGroup
	mu                sync.RWMutex
	ctx               context.Context
	done              context.CancelFunc
	maxWorkflowTask   int
	nowWorkflowTask   int
	gridManager       *grid.GridManager                            // Maintain an independent scheduler grid
	workflowSessions  *workflow.WorkflowSessions                   // Maintain an independent workflow sessions
	workflowTaskQueue map[string]*workflowtask.WorkflowTaskRuntime // Maintain an independent workflow task runtime queue
}

func NewWorkflowManager(ctx context.Context, maxWorkflowTaskOnlineNumber int) *Manager {
	wfm := &Manager{
		wg:                sync.WaitGroup{},
		mu:                sync.RWMutex{},
		maxWorkflowTask:   maxWorkflowTaskOnlineNumber,
		nowWorkflowTask:   0,
		workflowTaskQueue: make(map[string]*workflowtask.WorkflowTaskRuntime),
		gridManager:       nil,
		workflowSessions:  nil,
	}

	wfm.ctx, wfm.done = context.WithCancel(ctx)

	return wfm
}

func (w *Manager) InitialWorkflowSessions(maxWorkflowOnlineNumber int) *Manager {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.workflowSessions = workflow.NewWorkflowSessions(w.ctx, maxWorkflowOnlineNumber)

	return w
}

func (w *Manager) InitialGridManager(grider *grid.GridManager) *Manager {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.gridManager = grider

	return w
}

func (w *Manager) SelfCheck() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.gridManager == nil {
		return errors.New("[workflow manager] <wmpmanager.GridManager> cannot find it")
	} else if w.workflowSessions == nil {
		return errors.New("[workflow manager] <workflow.WorkflowSessions> cannot find it")
	}

	return nil
}

func (w *Manager) Wait() {
	w.wg.Wait()
}

func (w *Manager) Close() {
	w.done()
	w.wg.Wait()
}

func (w *Manager) WorkflowSessionsStatus() map[string]interface{} {
	return w.workflowSessions.WorkflowSessionsStatus()
}

func (w *Manager) WorkflowTasks(task_id string) (*workflowtask.WorkflowTasks, error) {
	if task_runtime := w.workflowTaskRuntime(task_id); task_runtime != nil {
		return task_runtime.WorkflowTasks(), nil
	} else {
		// If there is no task in the workflow task runtime queue, load it from the database.
		if tasks, count, err := workflowtask.WorkflowTasksFindAll(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "_id", Value: task_id}}); err != nil {
			return nil, err
		} else if count == 0 {
			return nil, errors.New("[workflow manager] this task cannot be found" + task_id)
		} else {
			return &tasks[0], nil
		}
	}
}
