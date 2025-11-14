package workflowtask

import (
	"context"
	"sync"

	"github.com/suitcase/butler/data/meta"
	"github.com/suitcase/butler/db"
	"github.com/suitcase/butler/grid"
	"github.com/suitcase/butler/workflow/workflow"
)

type WorkflowTaskRuntime struct {
	mu               sync.RWMutex
	wg               sync.WaitGroup
	ctx              context.Context
	done             context.CancelFunc
	workflowSessions *workflow.WorkflowSessions
	gridManager      *grid.GridManager
	backCallChannel  chan string
	workflowTasks    *WorkflowTasks
}

func (w *WorkflowTaskRuntime) WorkflowTasks() *WorkflowTasks {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.workflowTasks
}

func (w *WorkflowTaskRuntime) WorkflowTaskAssets() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.workflowTasks.GetAllTaskAssetQueue()
}

// get the asset_id of this task plan in the workflow session
func (w *WorkflowTaskRuntime) WorkflowTaskAssetsInSessions(asset_id ...string) []string {
	asset_ids := func() []string {
		if asset_id != nil {
			return asset_id
		}
		return w.WorkflowTaskAssets()
	}

	var assets []string
	for _, id := range asset_ids() {
		if w.workflowSessions.IsExist(id) {
			assets = append(assets, id)
		}
	}

	return assets
}

// get the asset_id of this task plan not in the workflow session
func (w *WorkflowTaskRuntime) WorkflowTaskAssetsNotInSessions(asset_id ...string) []string {
	asset_ids := func() []string {
		if asset_id != nil {
			return asset_id
		}
		return w.WorkflowTaskAssets()
	}

	var assets []string
	for _, id := range asset_ids() {
		if !w.workflowSessions.IsExist(id) {
			assets = append(assets, id)
		}
	}

	return assets
}

func (w *WorkflowTaskRuntime) AddWorkflowSession(asset meta.AssetMetaData) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if grider, err := w.gridManager.NewGrid(w.workflowTasks.ScheudlerID); err != nil {
		return err
	} else {
		return w.workflowSessions.AddWorkflow(asset.AssetId, workflow.Workflower(w.ctx, grider, asset, nil)) // BackCallChannel() -> AddWorkflowSession(asset meta.AssetMetaData)
	}
}

func (w *WorkflowTaskRuntime) ExitRuntime() error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if _, err := w.workflowTasks.NotRuntime().UpdateOne(context.TODO(), db.GetCurrentMongoClient()); err != nil {
		return err
	}

	return nil
}
