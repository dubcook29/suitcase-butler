package asset_api_service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suitcase/butler/api/v1/responses"
	"github.com/suitcase/butler/data/meta"
	"github.com/suitcase/butler/db"
	"go.mongodb.org/mongo-driver/bson"
)

type AssetAPIService struct {
	resp responses.DefaultResponses
}

func (api *AssetAPIService) InitialAPIService(g *gin.RouterGroup) {
	asset_api_service := g.Group("/asset")
	asset_api_service.GET("/", api.AssetQuerys)
	asset_api_service.GET("/:asset_id", api.AssetQuery)
	asset_api_service.POST("/", api.AssetUpdated)
	asset_api_service.POST("/:asset_id", api.AssetModify)
	asset_api_service.DELETE("/", api.AssetDeleted)
	asset_api_service.DELETE("/:asset_id", api.AssetDeletes)
}

// GET /{api}/asset/	query all asset data
func (api *AssetAPIService) AssetQuerys(ctx *gin.Context) {
	if data, count, err := meta.AssetMetaDataMongoFindFunc(context.TODO(), db.GetCurrentMongoClient(), bson.D{}); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows asset data", count), data)
	}
}

// GET /{api}/asset/:asset_id	query the specified asset data
func (api *AssetAPIService) AssetQuery(ctx *gin.Context) {
	if asset_id, ok := ctx.Params.Get("asset_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be query without specifying asset_id"))
		return
	} else if data, count, err := meta.AssetMetaDataMongoFindFunc(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "asset_id", Value: bson.D{{Key: "$in", Value: strings.Split(asset_id, ",")}}}}); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows asset data", count), data)
	}
}

// POST /{api}/asset/	created an new asset data write to the database
func (api *AssetAPIService) AssetUpdated(ctx *gin.Context) {
	var in meta.AssetMetaData
	if err := ctx.ShouldBindJSON(&in); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		in.Initial()
		if count, err := meta.AssetMetaDataMongoInsertFunc(context.TODO(), db.GetCurrentMongoClient(), []meta.AssetMetaData{in}); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows asset data insert to databased", count), nil)
			return
		}
	}
}

// POST /{api}/asset/:asset_id	updated the specified asset data write to the database
func (api *AssetAPIService) AssetModify(ctx *gin.Context) {
	var in meta.AssetMetaData
	if err := ctx.ShouldBindJSON(&in); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		if err := in.UpdateAsset(); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else {
			api.resp.SuccessWithMessage(ctx, "total 1 rows asset data update to databased", nil)
			return
		}
	}
}

// DELETE /{api}/asset/	deleted all asset data
func (api *AssetAPIService) AssetDeleted(ctx *gin.Context) {
	if count, err := meta.AssetMetaDataMongoDeleteFunc(context.TODO(), db.GetCurrentMongoClient(), bson.D{}); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows asset data deleted from databased", count), nil)
	}
}

// DELETE /{api}/asset/:asset_id	deleted the specified asset data
func (api *AssetAPIService) AssetDeletes(ctx *gin.Context) {
	if asset_id, ok := ctx.Params.Get("asset_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be query without specifying asset_id"))
		return
	} else if count, err := meta.AssetMetaDataMongoDeleteFunc(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "asset_id", Value: asset_id}}); err != nil {
		api.resp.Failed(ctx, err)
		return
	} else {
		api.resp.SuccessWithMessage(ctx, fmt.Sprintf("total %d rows asset data deleted from databased", count), nil)
	}
}
