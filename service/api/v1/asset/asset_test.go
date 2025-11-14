package asset_api_service

import (
	"net"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

var test_api_service = gin.Default()
var test_asset_api_group *gin.RouterGroup
var test_serviced AssetAPIService

func init() {
	test_api_service.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	test_asset_api_group = test_api_service.Group("/api/v1/asset")

	test_serviced.InitialAPIService(test_asset_api_group)

}

func run(port int) {
	if err := test_api_service.Run(net.JoinHostPort("", strconv.Itoa(port))); err != nil {
		panic(err)
	}
}
func TestAssetAPIService(t *testing.T) {
	run(8080)
}
