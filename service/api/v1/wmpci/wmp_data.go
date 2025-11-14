package wmpci_api_service

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/suitcase/butler/data/meta"
	"github.com/suitcase/butler/db"
	"github.com/suitcase/butler/wmpci/data/model"
	"go.mongodb.org/mongo-driver/bson"
)

// |GET   | /wmpdata/:asset_id/
func (api *WMPCIAPIService) WMPDataQueryers(ctx *gin.Context) {
	if asset_id, ok := ctx.Params.Get("asset_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be query without specifying asset_id"))
		return
	} else {
		if data, count, err := meta.AssetMetaDataMongoFindFunc(context.TODO(), db.GetCurrentMongoClient(), bson.D{{Key: "_id", Value: asset_id}}); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else if count > 0 {
			var result_all_data []model.WMPDataModelInterface
			for _, asset := range data {
				for _, model_name := range asset.ResultWmpDataLists {
					if data, count, err := model.WMPDataModelMongoFindFunc(context.TODO(), db.GetCurrentMongoClient(), "wmpdata", model_name, model.DefaultWMPDataModel{}, bson.D{{Key: "metadata.asset_id", Value: asset_id}}); err != nil {
						api.resp.Failed(ctx, err)
						return
					} else if count > 0 {
						result_all_data = append(result_all_data, data...)
					}
				}

			}

			api.resp.Successed(ctx, result_all_data)
			return
		} else {
			api.resp.Successed(ctx, nil)
			return
		}
	}

}

// |GET   | /wmpdata/:asset_id/:model_name
func (api *WMPCIAPIService) WMPDataQueryer(ctx *gin.Context) {
	if asset_id, ok := ctx.Params.Get("asset_id"); !ok {
		api.resp.Failed(ctx, errors.New("cannot be query without specifying asset_id"))
		return
	} else if model_name, ok := ctx.Params.Get("model_name"); ok {
		if data, count, err := model.WMPDataModelMongoFindFunc(context.TODO(), db.GetCurrentMongoClient(), "wmpdata", model_name, model.DefaultWMPDataModel{}, bson.D{{Key: "metadata.asset_id", Value: asset_id}}); err != nil {
			api.resp.Failed(ctx, err)
			return
		} else if count > 0 {
			api.resp.Successed(ctx, data)
			return
		}
	} else {
		api.resp.Successed(ctx, nil)
		return
	}
}
