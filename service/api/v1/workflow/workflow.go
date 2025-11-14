package workflow_api_service

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suitcase/butler/api/v1/responses"
	"github.com/suitcase/butler/grid"
	"github.com/suitcase/butler/workflow/manager"
)

type WorkflowAPIService struct {
	resp            responses.DefaultResponses
	workflowmanager *manager.Manager
}

func (api *WorkflowAPIService) InitialAPIService(g *gin.RouterGroup) {
	workflow_api_service := g.Group("/workflow")

	workflow_api_service.GET("", api.WorkflowQueueHealth)
	workflow_api_service.GET("/", api.WorkflowTaskRuntimeQueryerAPI)
	workflow_api_service.GET("/:task_id", api.WorkflowtaskRuntimeRestoreAPI)
	workflow_api_service.DELETE("/:task_id", api.WorkflowtaskRuntimeDeletedAPI)

	workflow_api_service.GET("/:task_id/add/", api.WorkflowtaskRuntimeAdds)
	workflow_api_service.GET("/:task_id/add/:asset_id", api.WorkflowtaskRuntimeAdd)
	workflow_api_service.GET("/:task_id/start/", api.WorkflowtaskRuntimeStarts)
	workflow_api_service.GET("/:task_id/start/:asset_id", api.WorkflowtaskRuntimeStart)
	workflow_api_service.GET("/:task_id/stop/", api.WorkflowtaskRuntimeStops)
	workflow_api_service.GET("/:task_id/stop/:asset_id", api.WorkflowtaskRuntimeStop)
	workflow_api_service.GET("/:task_id/delete/", api.WorkflowtaskRuntimeDeletes)
	workflow_api_service.GET("/:task_id/delete/:asset_id", api.WorkflowtaskRuntimeDelete)
	workflow_api_service.GET("/:task_id/restore/", api.WorkflowtaskRuntimeRestores)
	workflow_api_service.GET("/:task_id/restore/:asset_id", api.WorkflowtaskRuntimeRestore)
	workflow_api_service.GET("/:task_id/health/", api.WorkflowtaskRuntimeHealths)
	workflow_api_service.GET("/:task_id/health/:asset_id", api.WorkflowtaskRuntimeHealth)

	task_api_servicd := g.Group("/task")
	task_api_servicd.GET("/", api.WorkflowtaskNotRuntimeQuerys)
	task_api_servicd.POST("/", api.WorkflowtaskNotRuntimeModifys)
	task_api_servicd.DELETE("/", api.WorkflowtaskNotRuntimeDeletes)
	task_api_servicd.GET("/:task_id", api.WorkflowtaskNotRuntimeQuery)
	task_api_servicd.POST("/:task_id", api.WorkflowtaskNotRuntimeModify)
	task_api_servicd.DELETE("/:task_id", api.WorkflowtaskNotRuntimeDelete)
	task_api_servicd.GET("/:task_id/asset/", api.WorkflowtaskAssetQuery)
	task_api_servicd.GET("/:task_id/notasset/", api.WorkflowtaskAssetNoBeInQuery)
	task_api_servicd.POST("/:task_id/asset/:asset_id", api.WorkflowtaskAssetAdd)
	task_api_servicd.DELETE("/:task_id/asset/:asset_id", api.WorkflowtaskAssetDelete)

}

func (api *WorkflowAPIService) InitialAPIServiceAsWorkflowManager(s *grid.GridManager) {
	api.workflowmanager = manager.NewWorkflowManager(context.TODO(), 10)
	if err := api.workflowmanager.InitialGridManager(s).InitialWorkflowSessions(10).SelfCheck(); err != nil {
		panic(err)
	}
}

// |GET   |	/workflow	quer all workflow session health（workflow task runtime）
func (api *WorkflowAPIService) WorkflowQueueHealth(ctx *gin.Context) {
	if data := api.workflowmanager.WorkflowSessionsStatus(); data != nil {
		api.resp.Successed(ctx, data)
		return
	} else {
		api.resp.Successed(ctx, nil)
		return
	}
}

// |GET   |	/workflow/	quer all task（workflow task runtime）
func (api *WorkflowAPIService) WorkflowTaskRuntimeQueryerAPI(ctx *gin.Context) {
	if data := api.workflowmanager.WorkflowTaskRuntimeQueryer(); data != nil {
		api.resp.Successed(ctx, data)
		return
	} else {
		api.resp.Successed(ctx, nil)
		return
	}
}

// |GET   |	/workflow/:task_id	restart task（workflow task runtime）
func (api *WorkflowAPIService) WorkflowtaskRuntimeRestoreAPI(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be restore without specifying task_id"))
		return
	} else {
		if err := api.workflowmanager.WorkflowTaskRuntimeRestore(task_id); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, nil)
			return
		}
	}
}

// |POST  |	/workflow/:task_id	modify task（workflow task runtime）

// |DELETE|	/workflow/:task_id  delete task（workflow task runtime）
func (api *WorkflowAPIService) WorkflowtaskRuntimeDeletedAPI(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be restore without specifying task_id"))
		return
	} else {
		if err := api.workflowmanager.WorkflowTaskRuntimeDeleted(task_id); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, nil)
			return
		}
	}
}

// |GET   |	/workflow/:task_id/query/
// |GET   |	/workflow/:task_id/query/:asset_id

// |GET   | /workflow/:task_id/add/		add all workflow
func (api *WorkflowAPIService) WorkflowtaskRuntimeAdds(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else {
		if err := api.workflowmanager.WorkflowTaskRuntimeActionAdd(task_id, nil...); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, nil)
			return
		}
	}
}

// |GET   | /workflow/:task_id/add/:asset_id		add workflow
func (api *WorkflowAPIService) WorkflowtaskRuntimeAdd(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else if asset_id, ok := ctx.Params.Get("asset_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying asset_id"))
		return
	} else {
		if err := api.workflowmanager.WorkflowTaskRuntimeActionAdd(task_id, asset_id); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, nil)
			return
		}
	}
}

// |GET   |	/workflow/:task_id/start/				start all workflow
func (api *WorkflowAPIService) WorkflowtaskRuntimeStarts(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else {
		if err := api.workflowmanager.WorkflowTaskRuntimeActionStart(task_id, nil...); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, nil)
			return
		}
	}
}

// |GET   |	/workflow/:task_id/start/:asset_id		start workflow
func (api *WorkflowAPIService) WorkflowtaskRuntimeStart(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else if asset_id, ok := ctx.Params.Get("asset_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying asset_id"))
		return
	} else {
		if err := api.workflowmanager.WorkflowTaskRuntimeActionStart(task_id, asset_id); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, nil)
			return
		}
	}
}

// |GET   |	/workflow/:task_id/stop/				stop all workflow
func (api *WorkflowAPIService) WorkflowtaskRuntimeStops(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else {
		if err := api.workflowmanager.WorkflowTaskRuntimeActionStop(task_id, nil...); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, nil)
			return
		}
	}
}

// |GET   |	/workflow/:task_id/stop/:asset_id		stop workflow
func (api *WorkflowAPIService) WorkflowtaskRuntimeStop(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else if asset_id, ok := ctx.Params.Get("asset_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying asset_id"))
		return
	} else {
		if err := api.workflowmanager.WorkflowTaskRuntimeActionStop(task_id, asset_id); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, nil)
			return
		}
	}
}

// |GET   |	/workflow/:task_id/delete/				delete all workflow
func (api *WorkflowAPIService) WorkflowtaskRuntimeDeletes(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else {
		if err := api.workflowmanager.WorkflowTaskRuntimeActionDeleted(task_id, nil...); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, nil)
			return
		}
	}
}

// |GET   |	/workflow/:task_id/delete/:asset_id		delete workflow
func (api *WorkflowAPIService) WorkflowtaskRuntimeDelete(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else if asset_id, ok := ctx.Params.Get("asset_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying asset_id"))
		return
	} else {
		if err := api.workflowmanager.WorkflowTaskRuntimeActionDeleted(task_id, asset_id); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, nil)
			return
		}
	}
}

// WorkflowTaskRuntimeActionHealth
// |GET   |	/workflow/:task_id/health/				get all workflowtask health
func (api *WorkflowAPIService) WorkflowtaskRuntimeHealths(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else {
		if data, err := api.workflowmanager.WorkflowTaskRuntimeActionHealth(task_id, nil...); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, data)
			return
		}
	}
}

// |GET   |	/workflow/:task_id/health/:asset_id		get workflowtask health
func (api *WorkflowAPIService) WorkflowtaskRuntimeHealth(ctx *gin.Context) {
	if task_id, ok := ctx.Params.Get("task_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying task_id"))
		return
	} else if asset_id, ok := ctx.Params.Get("asset_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be action without specifying asset_id"))
		return
	} else {
		if data, err := api.workflowmanager.WorkflowTaskRuntimeActionHealth(task_id, asset_id); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.Successed(ctx, data)
			return
		}
	}
}

// |GET   |	/workflow/:task_id/restore/				restore all workflow
func (api *WorkflowAPIService) WorkflowtaskRuntimeRestores(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, nil)
	// TODO Functions to be developed
}

// |GET   |	/workflow/:task_id/restore/:asset_id	restore workflow
func (api *WorkflowAPIService) WorkflowtaskRuntimeRestore(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, nil)
	// TODO Functions to be developed
}
