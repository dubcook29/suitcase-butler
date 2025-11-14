package workflow_api_service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suitcase/butler/data/meta"
	"github.com/suitcase/butler/db"
	workflowtask "github.com/suitcase/butler/workflow/task"
	"go.mongodb.org/mongo-driver/bson"
)

// |GET   |	/task/	查询所有task（database）
func (api *WorkflowAPIService) WorkflowtaskNotRuntimeQuerys(ctx *gin.Context) {
	if data, count, err := workflowtask.WorkflowTasksFindAll(context.TODO(), db.GetCurrentMongoClient(), bson.D{}); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows query", count), data)
		return
	}
}

// |POST  |	/task/	提交创建task（database）
func (api *WorkflowAPIService) WorkflowtaskNotRuntimeModifys(ctx *gin.Context) {
	var in workflowtask.WorkflowTasks
	if err := ctx.ShouldBindJSON(&in); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		if err := in.InitialTaskId().SelfCheck(); err != nil {
			api.resp.Failed(ctx, err)
			return
		}

		if count, err := workflowtask.WorkflowTasksInsert(context.TODO(), db.GetCurrentMongoClient(), []workflowtask.WorkflowTasks{in}); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows insert", count), nil)
			return
		}
	}
}

// |DELETE|	/task/	删除所有task（database）
func (api *WorkflowAPIService) WorkflowtaskNotRuntimeDeletes(ctx *gin.Context) {
	if count, err := workflowtask.WorkflowTasksDelete(context.TODO(), db.GetCurrentMongoClient(), bson.D{}); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows deleted", count), nil)
		return
	}
}

// |GET   |	/task/:task_id	查询指定task（database）
func (api *WorkflowAPIService) WorkflowtaskNotRuntimeQuery(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be query without specifying task_id"))
		return
	} else {
		if data, count, err := workflowtask.WorkflowTasksFindAll(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "task_id", Value: task_id}}); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows query", count), data)
			return
		}
	}
}

// |POST  |	/task/:task_id	编辑指定task（database）
func (api *WorkflowAPIService) WorkflowtaskNotRuntimeModify(ctx *gin.Context) {
	if _, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be modify without specifying task_id"))
		return
	} else {
		var in workflowtask.WorkflowTasks
		if err := ctx.ShouldBindJSON(&in); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			if count, err := in.UpdateOne(context.TODO(), db.GetCurrentMongoClient()); err != nil {
				api.resp.Failed(ctx, err)
				return
			} else {
				api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows modify", count), nil)
				return
			}
		}
	}
}

// |DELETE|	/task/:task_id	删除指定task（database）
func (api *WorkflowAPIService) WorkflowtaskNotRuntimeDelete(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be deleted without specifying task_id"))
		return
	} else {
		if count, err := workflowtask.WorkflowTasksDelete(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "task_id", Value: task_id}}); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows deleted", count), nil)
			return
		}
	}
}

// ** //
//	上述的 两组八个 API 接口方法，主要是用于管理和维护 NotRuntime 的 workflowtask 增、删、改、查
//
//	除了 task 上述之外，还需要考虑对 task 中的其他数据进行 增删改差，这其中就包括 task asset-list 和 run-log 的内容查询
// ** //

// |GET   |	/task/:task_id/asset/			查询所有task中的asset数据（database）
func (api *WorkflowAPIService) WorkflowtaskAssetQuery(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be query without specifying task_id"))
		return
	} else {
		if data, count, err := workflowtask.WorkflowTasksFindAll(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "task_id", Value: task_id}}); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else if count > 0 {
			asset_list := []string{}

			for _, task := range data {
				asset_list = append(asset_list, task.TaskAssetQueue...)
			}

			if asset_data, asset_count, err := meta.AssetMetaDataMongoFindFunc(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "asset_id", Value: bson.D{{Key: "$in", Value: asset_list}}}}); err != nil {
				api.resp.Failed(ctx, err)
				return
			} else {
				api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows query", asset_count), asset_data)
				return
			}

		} else {
			api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows query", count), data)
			return
		}
	}
}

// |GET   |	/task/:task_id/notasset/
func (api *WorkflowAPIService) WorkflowtaskAssetNoBeInQuery(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be query without specifying task_id"))
		return
	} else {
		if data, count, err := workflowtask.WorkflowTasksFindAll(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "task_id", Value: task_id}}); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else if count > 0 {
			asset_list := []string{}

			for _, task := range data {
				asset_list = append(asset_list, task.TaskAssetQueue...)
			}

			if asset_data, asset_count, err := meta.AssetMetaDataMongoFindFunc(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "asset_id", Value: bson.D{{Key: "$nin", Value: asset_list}}}}); err != nil {
				api.resp.Failed(ctx, err)
				return
			} else {
				api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows query", asset_count), asset_data)
				return
			}

		} else {
			api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows query", count), data)
			return
		}
	}
}

// |POST  | /task/:task_id/asset/:asset_id			添加一个新的Asset数据并将 asset_id 添加到task中
func (api *WorkflowAPIService) WorkflowtaskAssetAdd(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else if asset_id, ok := ctx.Params.Get("asset_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying asset_id"))
		return
	} else {
		// 查询 asset_id 是否存在
		if data, count, err := meta.AssetMetaDataMongoFindFunc(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "asset_id", Value: bson.D{{Key: "$in", Value: strings.Split(asset_id, ",")}}}}); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else if count > 0 {
			// asset_id is exist

			if tasks, err := api.workflowmanager.WorkflowTasks(task_id); err != nil {
				api.resp.Failed(ctx, err)
				return
			} else {
				for _, aid := range strings.Split(asset_id, ",") {
					tasks.AddAsset(aid)
				}
				if count, err := tasks.UpdateOne(context.TODO(), db.GetCurrentMongoClient()); err != nil {
					api.resp.Failed(ctx, err)
					return
				} else {
					api.resp.SuccessWithMessage(ctx, fmt.Sprintf("Successfully added asset [%s] into task [%s](%d) list", asset_id, task_id, count), nil)
					return
				}
			}

		} else {
			api.resp.SuccessWithMessage(ctx, fmt.Sprintf("cannot query asset [%s], cannot add this asset into task [%s] list", asset_id, task_id), data)
			return
		}
	}
}

// |DELETE|	/task/:task_id/asset/:asset_id			删除task中的asset数据（database）
func (api *WorkflowAPIService) WorkflowtaskAssetDelete(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else if asset_id, ok := ctx.Params.Get("asset_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying asset_id"))
		return
	} else {
		if tasks, err := api.workflowmanager.WorkflowTasks(task_id); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			for _, aid := range strings.Split(asset_id, ",") {
				tasks.DelAsset(aid)
			}
			if count, err := tasks.UpdateOne(context.TODO(), db.GetCurrentMongoClient()); err != nil {
				api.resp.Failed(ctx, err)
				return
			} else {
				api.resp.SuccessWithMessage(ctx, fmt.Sprintf("Successfully deleted asset [%s] into task [%s](%d) list", asset_id, task_id, count), nil)
				return
			}
		}

	}
}

// |GET   |	/task/:task_id/asset/:asset_id	查询指定task中的asset_id（database）
// |POST  |	/task/:task_id/asset/:asset_id	修改指定task中的asset_id（database）
// |DELETE|	/task/:task_id/asset/:asset_id	删除指定task中的asset_id（database）

// |GET   |	/task/:task_id/logs/			查询所有task中的logs数据（database）
// |DELETE|	/task/:task_id/logs/			删除所有task中的logs数据（database）
