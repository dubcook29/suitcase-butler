package meta

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/suitcase/butler/db"
	"go.mongodb.org/mongo-driver/bson"
)

func FindOne(ctx context.Context, asset_id string) ([]AssetMetaData, int64, error) {
	return nil, 0, nil
}

func FindMany(ctx context.Context, filter bson.D) ([]AssetMetaData, int64, error) {
	coll := getCollectionHander(db.GetCurrentMongoClient())

	if cursor, err := coll.Find(ctx, filter, nil); err != nil {
		return nil, -1, err
	} else {
		defer cursor.Close(ctx)
		var result []AssetMetaData
		if err := cursor.All(ctx, &result); err != nil {
			return nil, 0, err
		} else {
			return result, int64(len(result)), nil
		}
	}
}

func FindAll(ctx context.Context) ([]AssetMetaData, int64, error) {
	coll := getCollectionHander(db.GetCurrentMongoClient())

	if cursor, err := coll.Find(ctx, bson.D{}, nil); err != nil {
		return nil, -1, err
	} else {
		defer cursor.Close(ctx)
		var result []AssetMetaData
		if err := cursor.All(ctx, &result); err != nil {
			return nil, 0, err
		} else {
			return result, int64(len(result)), nil
		}
	}
}

func InitialAssetObject(iv []byte) *AssetMetaData {
	var asset = new(AssetMetaData)
	asset.CreatedAt = time.Now()
	asset.AssetId = uuid.NewString()
	if err := json.Unmarshal(iv, asset); err != nil {
		return nil
	} else {
		return asset
	}
}

func InitialAssetEmpty() *AssetMetaData {
	var asset = new(AssetMetaData)
	asset.CreatedAt = time.Now()
	asset.AssetId = uuid.NewString()
	return asset
}

func NewAssetMetaDataManyFromInput(iv []byte) ([]AssetMetaData, int64, error) {
	var data []AssetMetaData
	if err := json.Unmarshal(iv, &data); err != nil {
		return nil, 0, err
	}
	return data, int64(len(data)), nil
}

func (asset *AssetMetaData) NewAssetMetaDataOneFromInput(iv []byte) error {
	// 将用户的输入解析成 AssetMetaData 并转为实体，提供给外部调用
	asset.CreatedAt = time.Now()
	asset.AssetId = uuid.NewString()
	if err := json.Unmarshal(iv, asset); err != nil {
		return err
	}
	return nil
}

func (asset *AssetMetaData) AssetMetaDataOneFromDatabase(asset_id string) error {
	// 从数据库中寻找指定的 AssetMetaData 并转为实体，提供给外部调用
	if coll := getCollectionHander(db.GetCurrentMongoClient()); coll != nil {
		if findResult := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: asset_id}}); findResult.Err() != nil {
			return findResult.Err()
		} else {
			return findResult.Decode(asset)
		}
	} else {
		return fmt.Errorf("database connect collection hander call error")
	}
}

func (asset *AssetMetaData) Find(ctx context.Context) error {
	if coll := getCollectionHander(db.GetCurrentMongoClient()); coll != nil {
		if findResult := coll.FindOne(ctx, bson.D{{Key: "_id", Value: asset.AssetId}}); findResult.Err() != nil {
			return findResult.Err()
		} else {
			return findResult.Decode(asset)
		}
	} else {
		return fmt.Errorf("database connect collection hander call error")
	}
}

func (asset *AssetMetaData) Exist(ctx context.Context) bool {
	if coll := getCollectionHander(db.GetCurrentMongoClient()); coll != nil {
		if findResult := coll.FindOne(ctx, bson.D{{Key: "_id", Value: asset.AssetId}}); findResult.Err() != nil {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

func (asset *AssetMetaData) Insert(ctx context.Context) error {
	if coll := getCollectionHander(db.GetCurrentMongoClient()); coll != nil {
		if _, err := coll.InsertOne(ctx, db.MongoId(*asset, asset.AssetId)); err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("database connect collection hander call error")
	}
}

func (asset *AssetMetaData) Updated(ctx context.Context) error {
	asset.UpdatedAt = time.Now()
	if coll := getCollectionHander(db.GetCurrentMongoClient()); coll != nil {
		if _, err := coll.UpdateByID(ctx, asset.AssetId, bson.D{{Key: "$set", Value: db.MongoId(*asset, asset.AssetId)}}); err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("database connect collection hander call error")
	}
}

func (asset *AssetMetaData) Deleted(ctx context.Context) error {
	if coll := getCollectionHander(db.GetCurrentMongoClient()); coll != nil {
		if deleteResult := coll.FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: asset.AssetId}}); deleteResult.Err() != nil {
			return deleteResult.Err()
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("database connect collection hander call error")
	}

}
