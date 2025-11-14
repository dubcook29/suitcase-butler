package workflow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// WorkflowSessions Workflow Sessions is for managing running Workflows at runtime
type WorkflowSessions struct {
	wg                      sync.WaitGroup
	mu                      sync.RWMutex
	ctx                     context.Context
	done                    context.CancelFunc
	maxOnlineWorkflowNumber int
	onlineWorkflowNumber    int
	workflows               map[string]*Workflow
}

// NewWorkflowSessions initialized *WorkflowSessions
func NewWorkflowSessions(ctx context.Context, maxOnlineWorkflowNumber int) *WorkflowSessions {

	wqm := &WorkflowSessions{
		wg:                      sync.WaitGroup{},
		mu:                      sync.RWMutex{},
		onlineWorkflowNumber:    0,
		maxOnlineWorkflowNumber: maxOnlineWorkflowNumber,
		workflows:               make(map[string]*Workflow),
	}

	wqm.ctx, wqm.done = context.WithCancel(ctx)

	logrus.Infof("[workflow session] workflow sessions manager created successfully!")
	return wqm
}

// new workflow insert to the current Workflow Sessions
func (w *WorkflowSessions) AddWorkflow(asset_id string, flow *Workflow) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.workflows[asset_id]; ok {
		return fmt.Errorf("[workflow session] the workflow [%s] is exist, not add into sessions", asset_id)
	} else {
		logrus.Debugf("[workflow session] new workflow [%s] append into sessions", asset_id)
		w.workflows[asset_id] = flow
		return nil
	}
}

// delete Workflow from the current Workflow Sessions
func (w *WorkflowSessions) WorkflowDeleteOperation(asset_id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if flow, ok := w.workflows[asset_id]; ok {
		flow.Close()
		delete(w.workflows, asset_id)
		return nil
	} else {
		return fmt.Errorf("[workflow session] the workflow you are looking for [%s] does not exist", asset_id)
	}
}

// Start the Workflow with the current Workflow Sessions specified as asset_id.
// In principle, only one workflow with asset_id can be executed at the same time.
func (w *WorkflowSessions) WorkflowStartOperation(asset_id string) error {
	if flow, err := w.read(asset_id); err != nil {
		return err
	} else if flow.Status != Running {
		go w.workflowStartHandle(asset_id, flow)
	} else {
		return fmt.Errorf("[workflow session] the workflow [%s] is running, operation is prohibited", asset_id)
	}
	return nil
}

func (w *WorkflowSessions) workflowStartHandle(asset_id string, flow *Workflow) {
	if flow.Status != Stop {
		flow.Status = Waiting
	}

	w.online(1)
	defer w.offline()

	if flow.Starter(w.ctx) {
		if err := w.WorkflowDeleteOperation(asset_id); err != nil {
			logrus.Errorf("[workflow session] %v", err)
		}
	}
}

// WorkflowStopOperation Stop the Workflow with the current Workflow Sessions specified as asset_id.
func (w *WorkflowSessions) WorkflowStopOperation(asset_id string) error {
	if flow, err := w.read(asset_id); err != nil {
		return err
	} else if flow.Status == Running {
		flow.Close()
	} else {
		return fmt.Errorf("[workflow session] the workflow [%s] not is running, operation is prohibited", asset_id)
	}
	return nil
}

// It will read and return the status information of the current Workflow Sessions
func (w *WorkflowSessions) WorkflowSessionsStatus() map[string]interface{} {
	w.mu.RLock()
	defer w.mu.RUnlock()

	healthData := map[string]interface{}{
		"max":    w.maxOnlineWorkflowNumber,
		"online": w.onlineWorkflowNumber,
		"queue":  len(w.workflows),
		"queues": func() map[string]uint {
			data := make(map[string]uint)

			for k, flow := range w.workflows {
				data[k] = flow.Status
			}

			return data
		}(),
	}

	return healthData
}

// It will read and return the WorkflowRuntimeHealth of the asset_id specified by the current workflow sessions
func (w *WorkflowSessions) WorkflowRuntimeHealth(asset_id ...string) []WorkflowRuntimeHealth {
	w.mu.RLock()
	defer w.mu.RUnlock()
	var result []WorkflowRuntimeHealth

	if asset_id == nil {
		for _, flow := range w.workflows {
			result = append(result, flow.RuntimeHealth())
		}
	} else {
		for _, id := range asset_id {
			if flow, ok := w.workflows[id]; ok {
				result = append(result, flow.RuntimeHealth())
			} else {
				logrus.Errorf("[workflow session] the workflow you are looking for [%s] does not exist", id)
			}
		}
	}

	return result
}

// It will progressively increasing the value of onlineWorkflowNumber in WorkflowSessions, and if it exceeds maxOnlineWorkflowNumber, it will block and wait
func (w *WorkflowSessions) online(i int) {
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			if func(i int) bool {
				w.mu.Lock()
				defer w.mu.Unlock()
				if w.onlineWorkflowNumber < w.maxOnlineWorkflowNumber {
					w.wg.Add(1)
					w.onlineWorkflowNumber = w.onlineWorkflowNumber + i
					logrus.Debugf("[workflow session] workflow usage (%d/%d), %d workflow started", w.onlineWorkflowNumber, w.maxOnlineWorkflowNumber, i)
					return true
				}
				return false
			}(i) {
				return
			}
			time.Sleep(5 * time.Second)
		}
	}
}

// It will decrease progressively the value of onlineWorkflowNumber in WorkflowSessions
func (w *WorkflowSessions) offline() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.onlineWorkflowNumber > 0 {
		defer w.wg.Done()
		w.onlineWorkflowNumber = w.onlineWorkflowNumber - 1
		logrus.Debugf("[workflow session] workflow usage (%d/%d), one workflow exited", w.onlineWorkflowNumber, w.maxOnlineWorkflowNumber)
	}
}

// check if the asset(workflow) is in the workflow session
func (w *WorkflowSessions) IsExist(asset_id string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if _, ok := w.workflows[asset_id]; ok {
		return true
	} else {
		return false
	}
}

func (w *WorkflowSessions) read(asset_id string) (*Workflow, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if flow, ok := w.workflows[asset_id]; ok {
		return flow, nil
	} else {
		return nil, fmt.Errorf("[workflow session] the workflow you are looking for [%s] does not exist", asset_id)
	}
}

func (w *WorkflowSessions) Wait() {
	w.wg.Wait()
}

func (w *WorkflowSessions) Close() {
	w.done()
	w.Wait()
}
