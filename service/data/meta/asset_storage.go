package meta

import (
	"context"

	"github.com/suitcase/butler/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollectionHander(client *mongo.Client) *mongo.Collection {
	return client.Database("meta").Collection("asset")
}

func AssetMetaDataMongoFindFunc(ctx context.Context, client *mongo.Client, filter bson.D) ([]AssetMetaData, int64, error) {
	coll := getCollectionHander(client)

	if cursor, err := coll.Find(ctx, filter, nil); err != nil {
		return nil, 0, err
	} else {

		var count int64
		if count, err = coll.CountDocuments(ctx, filter); err != nil {
			return nil, 0, err
		}

		defer cursor.Close(ctx)
		var result []AssetMetaData
		if err := cursor.All(ctx, &result); err != nil {
			return nil, 0, err
		} else {
			return result, count, nil
		}
	}
}

func AssetMetaDataMongoInsertFunc(ctx context.Context, client *mongo.Client, data []AssetMetaData) (int64, error) {
	coll := getCollectionHander(client)
	var count int64

	for _, item := range data {
		if _, err := coll.InsertOne(ctx, db.MongoId(item, item.AssetId)); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				if _, err = coll.ReplaceOne(ctx, bson.D{{Key: "_id", Value: item.AssetId}}, db.MongoId(item, item.AssetId)); err != nil {
					return count, err
				}
			} else {
				return count, err
			}
		}
		count = count + 1
	}

	return count, nil
}

func AssetMetaDataMongoUpdateFunc(ctx context.Context, client *mongo.Client, data []AssetMetaData) (int64, error) {
	coll := getCollectionHander(client)
	var count int64

	for _, item := range data {
		if updateResult, err := coll.UpdateOne(ctx, bson.D{{Key: "_id", Value: item.AssetId}}, bson.D{{Key: "$set", Value: db.MongoId(item, item.AssetId)}}); err != nil {
			return count, err
		} else if updateResult.MatchedCount == 0 {
			if _, err := coll.InsertOne(ctx, db.MongoId(item, item.AssetId)); err != nil {
				return count, err
			}
		}
		count = count + 1
	}

	return count, nil
}

func AssetMetaDataMongoDeleteFunc(ctx context.Context, client *mongo.Client, filter bson.D) (int64, error) {
	coll := getCollectionHander(client)
	if result, err := coll.DeleteMany(ctx, filter); err != nil {
		return 0, err
	} else {
		return result.DeletedCount, nil
	}
}
