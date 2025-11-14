package wmpci_api_service

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/suitcase/butler/api/v1/responses"
	"github.com/suitcase/butler/db"
	"github.com/suitcase/butler/wmpci"
	wmpcisessions "github.com/suitcase/butler/wmpci/manager"
)

type WMPCIAPIService struct {
	resp                responses.DefaultResponses
	wmpcisessionmanager *wmpcisessions.WMPSessionManager
	uploadfilePath      string
}

func (api *WMPCIAPIService) InitialAPIService(g *gin.RouterGroup) {
	wmpci_api_service := g.Group("/wmpci")

	wmpci_api_service.GET("/sessions/", api.GetSessions)
	wmpci_api_service.GET("/sessions/:session_id", api.GetSession)
	wmpci_api_service.DELETE("/sessions/:session_id", api.DeleteSession)
	wmpci_api_service.GET("/sessions/connector/", api.GetCurrentWMPConnectorSupports)
	wmpci_api_service.GET("/sessions/connector/:connector_id", api.GetCurrentWMPConnectorCumstom)
	wmpci_api_service.POST("/sessions/connector/:connector_id", api.ConnectorConnectionSession)

	wmpci_api_service.GET("/wmpdata/:asset_id/", api.WMPDataQueryers)
	wmpci_api_service.GET("/wmpdata/:asset_id/:model_name", api.WMPDataQueryer)
}

func (api *WMPCIAPIService) InitialAPIServiceAsWMPRegistrarManager(upfilepath string) *wmpcisessions.WMPSessionManager {

	api.uploadfilePath = upfilepath
	api.wmpcisessionmanager = wmpcisessions.InitialSessionManager(context.TODO(), db.CURRENT_CACHE_PATH)

	return api.wmpcisessionmanager
}

// GET  /wmpci/sessions/connector/
func (api *WMPCIAPIService) GetCurrentWMPConnectorSupports(ctx *gin.Context) {
	if data, err := api.wmpcisessionmanager.ConnectorSupportLists(); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.Successed(ctx, data)
		return
	}
}

// GET  /wmpci/sessions/connector/:connector_id
func (api *WMPCIAPIService) GetCurrentWMPConnectorCumstom(ctx *gin.Context) {
	connector_id := ctx.Param("connector_id")
	if connector_id == "" {
		api.resp.Failed(ctx, errors.New("cannot be find without specifying connector_id"))
		return
	} else if data, err := api.wmpcisessionmanager.SelectConnectorConnectionCumstom(connector_id); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.Successed(ctx, data)
		return
	}
}

// POST  /wmpci/sessions/connector/:connector_id
func (api *WMPCIAPIService) ConnectorConnectionSession(ctx *gin.Context) {
	connector_id := ctx.Param("connector_id")
	if connector_id == "" {
		api.resp.Failed(ctx, errors.New("cannot be find without specifying connector_id"))
		return
	}

	var custom_config map[string]wmpci.WMPCustom
	if err := ctx.BindJSON(&custom_config); err != nil {
		api.resp.Failed(ctx, err)
		return
	}

	if data, err := api.wmpcisessionmanager.ConnectorConnectionSession(connector_id, custom_config); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.Successed(ctx, data)
		return
	}
}

// GET  /wmpci/sessions/
func (api *WMPCIAPIService) GetSessions(ctx *gin.Context) {
	if session, err := api.wmpcisessionmanager.Sessions(); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.Successed(ctx, session)
		return
	}
}

// GET  /wmpci/sessions/:session_id
func (api *WMPCIAPIService) GetSession(ctx *gin.Context) {
	session_id := ctx.Param("session_id")
	if session_id == "" {
		api.resp.Failed(ctx, errors.New("cannot be find without specifying session_id"))
		return
	}

	if session, err := api.wmpcisessionmanager.Session(session_id); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.Successed(ctx, []any{session})
		return
	}
}

// DELETE  /wmpci/sessions/:session_id
func (api *WMPCIAPIService) DeleteSession(ctx *gin.Context) {
	session_id := ctx.Param("session_id")
	if session_id == "" {
		api.resp.Failed(ctx, errors.New("cannot be delete without specifying session_id"))
		return
	}

	if err := api.wmpcisessionmanager.CloseSession(session_id); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.Successed(ctx, nil)
		return
	}
}

func (api WMPCIAPIService) saveUploadedFile(ctx *gin.Context) (string, error) {

	if ctx.Request.Method == "PUT" {
		fileName := ctx.Query("filename")
		if fileName == "" {
			return "", errors.New("there is no file name, you need to enter \"filename\"")
		}

		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return "", fmt.Errorf("filed to read upload data: %s", err.Error())
		}

		err = ioutil.WriteFile(api.uploadfilePath+fileName, body, 0644)
		if err != nil {
			return "", fmt.Errorf("failed to save the file: %s", err.Error())
		}

		return api.uploadfilePath + fileName, nil

	} else if ctx.Request.Method == "POST" {
		file, err := ctx.FormFile("filename")
		if err != nil {
			return "", fmt.Errorf("filed to file upload: %s", err.Error())
		}

		if err := ctx.SaveUploadedFile(file, api.uploadfilePath+file.Filename); err != nil {
			return "", fmt.Errorf("failed to save the file: %s", err.Error())
		}

		return api.uploadfilePath + file.Filename, nil
	}

	return "", nil
}

// ensureDir check if the directory exists, and if it does not exist, create it
func ensureDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
