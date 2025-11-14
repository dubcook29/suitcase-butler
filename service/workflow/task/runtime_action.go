package workflowtask

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/data/meta"
	"github.com/suitcase/butler/db"
	"github.com/suitcase/butler/workflow/workflow"
	"go.mongodb.org/mongo-driver/bson"
)

// WorkflowTaskRuntimeAdd: Runtime task Adds the specified or all workflows to be started
func (w *WorkflowTaskRuntime) WorkflowTaskRuntimeAdd(asset_id ...string) error {
	for _, id := range w.WorkflowTaskAssetsNotInSessions(asset_id...) {
		if assetResult, count, err := meta.AssetMetaDataMongoFindFunc(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "_id", Value: id}}); err != nil {
			return err
		} else if count > 0 {
			for _, v := range assetResult {
				if err := w.AddWorkflowSession(v); err != nil {
					return err
				} else {
					if _, err := w.workflowTasks.AddAsset(id).UpdateOne(w.ctx, db.GetCurrentMongoClient()); err != nil {
						return err
					}
				}
			}
		} else {
			return errors.New("[workflow task runtime] the asset you specified cannot be found in the database: " + id)
		}
	}
	return nil
}

// WorkflowTaskRuntimeStart: Runtime task Starts the specified or all workflows to be started
func (w *WorkflowTaskRuntime) WorkflowTaskRuntimeStart(asset_id ...string) error {
	for _, id := range w.WorkflowTaskAssetsInSessions(asset_id...) {
		if err := w.workflowSessions.WorkflowStartOperation(id); err != nil {
			return err
		}
	}
	return nil
}

// WorkflowTaskRuntimeStop: Runtime task Stops the specified or all workflows to be started
func (w *WorkflowTaskRuntime) WorkflowTaskRuntimeStop(asset_id ...string) error {
	for _, id := range w.WorkflowTaskAssetsInSessions(asset_id...) {
		if err := w.workflowSessions.WorkflowStopOperation(id); err != nil {
			return err
		}
	}
	return nil
}

// TODO Not tested, not recommended
func (w *WorkflowTaskRuntime) WorkflowTaskRuntimeRestart(asset_id ...string) error {
	for _, id := range w.WorkflowTaskAssetsInSessions(asset_id...) {
		if err := w.Restore(id); err != nil {
			return err
		}
	}
	return nil
}

// TODO Not tested, not recommended
func (w *WorkflowTaskRuntime) Restore(asset_id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if assets, count, err := meta.AssetMetaDataMongoFindFunc(w.ctx, db.GetCurrentMongoClient(), bson.D{{Key: "_id", Value: asset_id}}); err != nil {
		logrus.Errorf("[task] asset find failed: %v", err)
		return err
	} else if count > 0 {

		for _, asset := range assets {
			asset.UpdatedAt = time.Now()

			if err := w.AddWorkflowSession(asset); err != nil {
				logrus.Errorf("[workflow task runtime] failed to restore workflow (%v) from cache, [%+v]", err, asset)
				return err
			} else {
				logrus.Debugf("[workflow task runtime] successes to add a new workflow [%+v]", asset)

				if err := w.workflowSessions.WorkflowStartOperation(asset_id); err != nil {
					return err
				} else {
				}
			}

		}

	}

	return nil
}

// deleted the workflow from `workflow.WorkflowSessions`
func (w *WorkflowTaskRuntime) WorkflowTaskRuntimeDelete(asset_id ...string) error {
	for _, id := range w.WorkflowTaskAssetsInSessions(asset_id...) {
		if err := w.workflowSessions.WorkflowDeleteOperation(id); err != nil {
			return err
		}
	}
	return nil
}

// get the health of the workflow in the workflow task from `workflow.WorkflowSessions`
func (w *WorkflowTaskRuntime) WorkflowTaskRuntimeHealth(asset_id ...string) []workflow.WorkflowRuntimeHealth {
	var healths []workflow.WorkflowRuntimeHealth
	for _, id := range w.WorkflowTaskAssetsInSessions(asset_id...) {
		healths = append(healths, w.workflowSessions.WorkflowRuntimeHealth(id)...)
	}
	return healths
}
