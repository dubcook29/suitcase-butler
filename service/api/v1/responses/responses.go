package responses

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type DefaultResponses struct {
}

func (resp DefaultResponses) Successed(ctx *gin.Context, data any) {

	ctx.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"msg":   "successed",
		"data":  data,
		"count": calcCount(data),
	})

}

func (resp DefaultResponses) SuccessWithMessage(ctx *gin.Context, msg string, data any) {

	ctx.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"msg":   msg,
		"data":  data,
		"count": calcCount(data),
	})

}

func (resp DefaultResponses) Failed(ctx *gin.Context, err error) {

	ctx.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{
		"code":  http.StatusNotImplemented,
		"msg":   err.Error(),
		"err":   err.Error(),
		"data":  nil,
		"count": -1,
	})
}

func (resp DefaultResponses) Aborted(ctx *gin.Context) {

	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"code":  http.StatusInternalServerError,
		"msg":   "Internal Server error",
		"err":   "Internal Server error",
		"data":  nil,
		"count": -1,
	})

}

func (resp DefaultResponses) Null(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
		"code":  http.StatusServiceUnavailable,
		"msg":   "Service Unavailable",
		"err":   "Service Unavailable",
		"data":  nil,
		"count": -1,
	})
}

func calcCount(data any) int {

	if v := reflect.ValueOf(data); v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		return v.Len()
	} else {
		return 0
	}
}
