package scheduler_api_service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/suitcase/butler/api/v1/responses"
	"github.com/suitcase/butler/db"
	"github.com/suitcase/butler/grid"
	wmpcisessions "github.com/suitcase/butler/wmpci/manager"
)

type SchedulerAPIService struct {
	resp           responses.DefaultResponses
	gridManager    *grid.GridManager
	uploadfilePath string
}

func (api *SchedulerAPIService) InitialAPIService(g *gin.RouterGroup) {
	scheduler_api_service := g.Group("/scheduler")

	scheduler_api_service.GET("/:scheduler_id", api.SchedulerGridQueryAPI)
	scheduler_api_service.GET("/", api.SchedulerGridQueryAPI)

	scheduler_api_service.POST("/:scheduler_id", api.SchedulerGridModifyAPI)

	scheduler_api_service.DELETE("/:scheduler_id", api.SchedulerGridDeleteAPI)

	scheduler_api_service.POST("/add/:scheduler_id", api.SchedulerGridAddJSONAPI)
	scheduler_api_service.PUT("/add/:scheduler_id", api.SchedulerGridAddYAMLAPI)

}

func (api *SchedulerAPIService) InitialAPIServiceAsScheudlerManager(upfilepath string, sessions *wmpcisessions.WMPSessionManager) *grid.GridManager {

	s := new(grid.GridManager).Initial(db.CURRENT_CACHE_PATH, sessions)

	if _, err := s.AddGrid(&grid.Grid{
		GridId:          "DufalueBuiltWMPSchedulingGrid",
		GridName:        "DufalueBuiltWMPSchedulingGrid",
		GridDescriptive: "DufalueBuiltWMPSchedulingGrid",
		GridTasks: []grid.Task{
			{
				TaskId:   "DNS",
				TaskName: "DNS",
				WMPId:    "ffffffff-ffff-ffff-0002-b0bd42efd9ae",
				NextTask: []grid.Task{
					{
						TaskId:   "SUBDOMAIN",
						TaskName: "SUBDOMAIN",
						WMPId:    "ffffffff-ffff-ffff-0004-b0bd42efd9ae",
						NextTask: []grid.Task{},
					},
					{
						TaskId:   "cdn",
						TaskName: "cdn",
						WMPId:    "ffffffff-ffff-ffff-0003-b0bd42efd9ae",
						NextTask: []grid.Task{},
					},
				},
			},
			{
				TaskId:   "WHOIS",
				TaskName: "WHOIS",
				WMPId:    "ffffffff-ffff-ffff-0001-b0bd42efd9ae",
				NextTask: []grid.Task{},
			},
		},
	}); err != nil {
		panic(err)
	}

	api.uploadfilePath = upfilepath
	api.gridManager = s

	return s
}

// api PUT */wmpci/add/[:scheduler_id]
func (api SchedulerAPIService) SchedulerGridAddYAMLAPI(ctx *gin.Context) {

	scheduler_id := ctx.Param("scheduler_id")
	if scheduler_id == "" {
		api.resp.Failed(ctx, errors.New("cannot be deleted without specifying scheduler_id"))
		return
	}

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		api.resp.Failed(ctx, fmt.Errorf("filed to read upload data: %s", err.Error()))
		return
	}

	var in_data grid.Grid

	if err := in_data.Deserialization(body); err != nil {
		api.resp.Failed(ctx, err)
		return
	}

	if _, err := api.gridManager.AddGrid(&in_data); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.getAllSchedulingGrid(ctx)
		return
	}

}

// api POST */wmpci/add/[:scheduler_id]
func (api SchedulerAPIService) SchedulerGridAddJSONAPI(ctx *gin.Context) {

	scheduler_id := ctx.Param("scheduler_id")
	if scheduler_id == "" {
		api.resp.Failed(ctx, errors.New("cannot be deleted without specifying scheduler_id"))
		return
	}

	var in_data grid.Grid
	if err := ctx.BindJSON(&in_data); err != nil {
		api.resp.Failed(ctx, err)
		return
	}

	if _, err := api.gridManager.AddGrid(&in_data); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.getAllSchedulingGrid(ctx)
		return
	}

}

// api DELETE */scheduler/[:scheduler_id]
func (api SchedulerAPIService) SchedulerGridDeleteAPI(ctx *gin.Context) {

	scheduler_id := ctx.Param("scheduler_id")
	if scheduler_id == "" {
		api.resp.Failed(ctx, errors.New("cannot be deleted without specifying scheduler_id"))
		return
	} else {
		if _, err := api.gridManager.DelGrid(scheduler_id); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.getAllSchedulingGrid(ctx)
			return
		}
	}

}

// api POST */scheduler/[:scheduler_id]
func (api SchedulerAPIService) SchedulerGridModifyAPI(ctx *gin.Context) {

	scheduler_id := ctx.Param("scheduler_id")
	var in_data grid.Grid
	if err := ctx.ShouldBindJSON(&in_data); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else if scheduler_id == in_data.GridId {
		if _, err := api.gridManager.ModGrid(&in_data); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.getAllSchedulingGrid(ctx)
			return
		}

	} else {
		api.resp.Failed(ctx, errors.New("cannot be modified without specifying scheduler_id"))
		return
	}

}

// api GET */scheduler/[:scheduler_id]
func (api SchedulerAPIService) SchedulerGridQueryAPI(ctx *gin.Context) {

	scheduler_id := ctx.Param("scheduler_id")
	if scheduler_id == "" {
		api.getAllSchedulingGrid(ctx)
		return
	} else {
		api.getSchedulingGrid(ctx, scheduler_id)
		return
	}

}

func (api SchedulerAPIService) getAllSchedulingGrid(ctx *gin.Context) {
	if data, err := api.gridManager.GetAllOnlineGrid(); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.Successed(ctx, data)
		return
	}
}

func (api SchedulerAPIService) getSchedulingGrid(ctx *gin.Context, scheduler_id string) {
	if data, err := api.gridManager.NewGrid(scheduler_id); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.Successed(ctx, data)
		return
	}
}

// ensureDir check if the directory exists, and if it does not exist, create it
func ensureDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
