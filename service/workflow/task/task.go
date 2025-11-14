package workflowtask

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Disabled uint = iota
	NotRuntime
	Runtime
)

type WorkflowTasks struct {
	TaskID          string        `json:"task_id,omitempty" bson:"task_id"`
	TaskName        string        `json:"task_name" bson:"task_name"`
	TaskDescription string        `json:"task_description,omitempty" bson:"task_description"`
	ScheudlerID     string        `json:"scheduler_id" bson:"scheduler_id"`
	Status          uint          `json:"status,omitempty" bson:"status"`
	CreatedAt       int64         `json:"created_at,omitempty" bson:"created_at"`
	StartedAt       int64         `json:"started_at,omitempty" bson:"started_at"`
	CloseddAt       int64         `json:"closed_at,omitempty" bson:"closed_at"`
	TaskRunlogs     []TaskRunLogs `json:"task_run_logs,omitempty" bson:"task_run_logs"`
	TaskAssetQueue  []string      `json:"task_asset_queue,omitempty" bson:"task_asset_queue"`
}

type TaskRunLogs struct {
	Time      int64  `json:"time" bson:"time"`
	AssetId   string `json:"asset_id" bson:"asset_id"`
	Messages  string `json:"messages,omitempty" bson:"messages,omitempty"`
	Excetpion string `json:"excetpion,omitempty" bson:"excetpion,omitempty"`
}

// initial task_id and created_at
func (w *WorkflowTasks) InitialTaskId() *WorkflowTasks {
	if w.TaskID == "" {
		w.TaskID = uuid.NewString()
	}

	if w.CreatedAt == 0 {
		w.CreatedAt = time.Now().Unix()
	}

	return w
}

// initial scheduler_id
func (w *WorkflowTasks) InitialSchedulerId(scheduler_id string) *WorkflowTasks {
	w.ScheudlerID = scheduler_id
	return w
}

// self-check
func (w *WorkflowTasks) SelfCheck() error {
	if w.TaskID == "" {
		return errors.New("[task] taks_id is empty or wrong")
	} else if w.ScheudlerID == "" {
		return errors.New("[task] scheduler_id is empty or wrong")
	}
	return nil
}

// Model output WorkflowTasks for mongodb storage
func (w *WorkflowTasks) Model() WorkflowTasks {
	return *w
}

func (w *WorkflowTasks) Runtime() *WorkflowTasks {

	w.Status = Runtime
	return w
}

func (w *WorkflowTasks) NotRuntime() *WorkflowTasks {
	w.Status = NotRuntime
	return w
}

func (w *WorkflowTasks) DisabledRuntime() *WorkflowTasks {
	w.Status = Disabled
	return w
}

// add new asset_id into TaskAssetQueue
func (w *WorkflowTasks) AddAsset(asset_id string) *WorkflowTasks {
	w.TaskAssetQueue = dataDuplication(append(w.TaskAssetQueue, asset_id))
	return w
}

func dataDuplication(in []string) []string {
	var duplication map[string]bool = make(map[string]bool)
	for _, v := range in {
		duplication[v] = true
	}

	var out []string
	for k := range duplication {
		out = append(out, k)
	}

	return out
}

// delete a asset_id from TaskAssetQueue
func (w *WorkflowTasks) DelAsset(asset_id string) *WorkflowTasks {
	w.TaskAssetQueue = dataRemoval(w.TaskAssetQueue, asset_id)
	return w
}

func dataRemoval(in []string, key string) []string {
	var duplication map[string]bool = make(map[string]bool)
	for _, v := range in {
		duplication[v] = true
	}

	delete(duplication, key)

	var out []string
	for k := range duplication {
		out = append(out, k)
	}

	return out
}

func (w *WorkflowTasks) GetAllTaskAssetQueue() []string {
	return w.TaskAssetQueue
}

func (w *WorkflowTasks) AddMessage(asset_id, text string) *WorkflowTasks {
	// TODO reserved, not deprecated
	w.TaskRunlogs = append(w.TaskRunlogs, TaskRunLogs{
		Time:     time.Now().Unix(),
		AssetId:  asset_id,
		Messages: text,
	})

	return w
}

func (w *WorkflowTasks) AddExcetpion(asset_id string, text string) *WorkflowTasks {
	// TODO reserved, not deprecated
	w.TaskRunlogs = append(w.TaskRunlogs, TaskRunLogs{
		Time:      time.Now().Unix(),
		AssetId:   asset_id,
		Excetpion: text,
	})

	return w
}

func (w *WorkflowTasks) InsertOne(ctx context.Context, client *mongo.Client) (int64, error) {
	// TODO reserved, not deprecated
	return WorkflowTasksInsert(ctx, client, []WorkflowTasks{*w})
}

func (w *WorkflowTasks) UpdateOne(ctx context.Context, client *mongo.Client) (int64, error) {
	coll := getCollectionHander(client)
	if updateResult, err := coll.UpdateOne(ctx, bson.D{{Key: "_id", Value: w.TaskID}}, bson.D{{Key: "$set", Value: w}}); err != nil {
		return 0, err
	} else {
		return updateResult.ModifiedCount, nil
	}
}

func (w *WorkflowTasks) Deleted(ctx context.Context, client *mongo.Client) (int64, error) {
	// TODO reserved, not deprecated
	return WorkflowTasksDelete(ctx, client, bson.D{{Key: "_id", Value: w.TaskID}})
}
