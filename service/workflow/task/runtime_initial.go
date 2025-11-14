package workflowtask

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/data/meta"
	"github.com/suitcase/butler/db"
	"github.com/suitcase/butler/grid"
	"github.com/suitcase/butler/workflow/workflow"
)

func (w *WorkflowTaskRuntime) Initial(ctx context.Context, thread int) *WorkflowTaskRuntime {
	w.ctx, w.done = context.WithCancel(ctx)

	w.wg = sync.WaitGroup{}

	w.mu = sync.RWMutex{}

	w.backCallChannel = make(chan string, thread)

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.BackCallChannel()
	}()

	return w
}

func (w *WorkflowTaskRuntime) InitialGridManager(grider *grid.GridManager) *WorkflowTaskRuntime {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.gridManager = grider
	return w
}

func (w *WorkflowTaskRuntime) InitialWorkflowSessions(workflowSessions *workflow.WorkflowSessions) *WorkflowTaskRuntime {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.workflowSessions = workflowSessions
	return w
}

func (w *WorkflowTaskRuntime) InitialWorkflowTasks(workflowTasks *WorkflowTasks) *WorkflowTaskRuntime {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.workflowTasks = workflowTasks.Runtime()
	return w
}

func (w *WorkflowTaskRuntime) SelfCheck() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.workflowTasks.ScheudlerID == "" {
		return fmt.Errorf("[workflow runtime task] scheduling plan is not designed and cannot be executed")
	} else if w.gridManager == nil {
		return fmt.Errorf("[workflow runtime task] the scheduler cannot recognize and cannot perform scheduled tasks")
	} else if w.workflowSessions == nil {
		return fmt.Errorf("[workflow runtime task] workflow sessions cannot be found, workflow system is not available")
	} else if w.workflowTasks == nil {
		return fmt.Errorf("[workflow runtime task] workflow task is empty and cannot be executed")
	}

	return nil
}

func (w *WorkflowTaskRuntime) Close() {
	w.done()
	w.wg.Wait()
	close(w.backCallChannel)
}

func (w *WorkflowTaskRuntime) Wait() {
	w.wg.Wait()
	close(w.backCallChannel)
}

func (w *WorkflowTaskRuntime) BackCallChannel() {
	for iv := range w.backCallChannel {
		logrus.Debugf("[workflow runtime task] BackCallChannel New data received %v", iv)
		if asset, err := meta.AssetMetaDataConstruct(iv); err != nil {
			logrus.Errorf("[workflow runtime task] back call asset data is error: %v", err)
		} else {
			if err := asset.UpdateAsset(); err != nil {
				return
			}
			if err := w.AddWorkflowSession(asset); err != nil {
				logrus.Errorf("[workflow runtime task] back call asset add to workflow queue is error: %v", err)
			} else {
				if _, err := w.workflowTasks.AddAsset(asset.AssetId).UpdateOne(context.TODO(), db.GetCurrentMongoClient()); err != nil {
					logrus.Error("[workflow runtime task]", err.Error())
				}
			}
		}
	}
}
