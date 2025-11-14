package v1

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	asset_api_service "github.com/suitcase/butler/api/v1/asset"
	scheduler_api_service "github.com/suitcase/butler/api/v1/scheduler"
	wmpci_api_service "github.com/suitcase/butler/api/v1/wmpci"
	workflow_api_service "github.com/suitcase/butler/api/v1/workflow"
)

func ButlerAPIServiceStarter(host, port string) {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If it is a preflight request, return directly
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	var workflow_api_service = new(workflow_api_service.WorkflowAPIService)
	var wmpci_api_service = new(wmpci_api_service.WMPCIAPIService)
	var asset_api_service = new(asset_api_service.AssetAPIService)
	var scheduler_api_service = new(scheduler_api_service.SchedulerAPIService)

	workflow_api_service.
		InitialAPIServiceAsWorkflowManager(
			scheduler_api_service.InitialAPIServiceAsScheudlerManager(
				"./", wmpci_api_service.InitialAPIServiceAsWMPRegistrarManager("./"),
			),
		)

	router_v1 := router.Group("/api/v1")

	workflow_api_service.InitialAPIService(router_v1)
	wmpci_api_service.InitialAPIService(router_v1)
	asset_api_service.InitialAPIService(router_v1)
	scheduler_api_service.InitialAPIService(router_v1)

	if err := router.Run(net.JoinHostPort(host, port)); err != nil {
		panic(err)
	}
}
