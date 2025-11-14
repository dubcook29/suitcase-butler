package cache

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/db"
	"go.mongodb.org/mongo-driver/bson"
)

type wmpDataCache struct {
	ID                    string                              `bson:"_id"`
	AssetId               string                              `bson:"asset_id"`
	SchedulingGridId      string                              `bson:"scheduling_grid_id"`
	WMPResponseDataBuffer map[string]map[string][]interface{} `bson:"wmp_response_data_buffer"`
	ExceptionDataBuffer   []string                            `bson:"exception_data_buffer"`
}

func (w *WorkflowDataBuffer) WorkflowDataBufferMongoUpdate(grid string) (wmpDataCache, int64, error) {
	coll := db.ConnectToDatabaseAndCollection(nil, "wmp", "wmp_data_cache")

	key := w.asset.AssetId + ":" + grid

	cache := wmpDataCache{
		ID:                    w.asset.AssetId + ":" + grid,
		AssetId:               w.asset.AssetId,
		SchedulingGridId:      grid,
		WMPResponseDataBuffer: w.wmpResponseDataBuffer,
		ExceptionDataBuffer:   w.exceptionDataBuffer,
	}

	if updateResult, err := coll.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: key}}, bson.D{{Key: "$set", Value: cache}}); err != nil {
		logrus.Errorf("[Buffer] WMP DMS Bind Data Update is error: %v", err)
		return cache, 0, err
	} else if updateResult.MatchedCount == 0 {
		if _, err := coll.InsertOne(context.TODO(), cache); err != nil {
			logrus.Errorf("[Buffer] WMP DMS Bind Data Update  (insert) is error: %v", err)
			return cache, 0, err
		} else {
			logrus.Debug("[Buffer] WMP DMS Bind Data Update (insert) is successed.")
			return cache, 1, nil
		}
	} else {
		logrus.Debug("[Buffer] WMP DMS Bind Data Update is successed.")
		return cache, updateResult.UpsertedCount, nil
	}
}

func (w *WorkflowDataBuffer) WorkflowDataBufferMongoFind(grid string) (int64, error) {
	coll := db.ConnectToDatabaseAndCollection(nil, "wmp", "wmp_data_cache")

	key := w.asset.AssetId + ":" + grid

	if cursor, err := coll.Find(context.TODO(), bson.D{{Key: "_id", Value: key}}); err != nil {
		logrus.Errorf("[Buffer] workflow data cache find failed: %v(%s)", err, key)
		return 0, err
	} else {
		defer cursor.Close(context.TODO())
		var result []wmpDataCache
		if err := cursor.All(context.TODO(), &result); err != nil {
			logrus.Errorf("[Buffer] workflow data cache decoding failed: %v(%s)", err, key)
			return 0, err
		} else if len(result) > 0 {
			w.wmpResponseDataBuffer = result[0].WMPResponseDataBuffer
			w.exceptionDataBuffer = result[0].ExceptionDataBuffer
			return 1, nil
		} else {
			return 0, nil
		}
	}
}
