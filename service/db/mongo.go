package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func MongoInsertFunc(ctx context.Context, coll *mongo.Collection, id string, data interface{}) (int64, error) {

	if _, err := coll.InsertOne(ctx, data); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			if ReplaceOneResult, err := coll.ReplaceOne(ctx, bson.D{{Key: "_id", Value: id}}, data); err != nil {
				return ReplaceOneResult.MatchedCount, err
			} else {
				return ReplaceOneResult.MatchedCount, nil
			}
		} else {
			return 0, err
		}
	} else {
		return 1, nil
	}

}

func MongoUpdateFunc(ctx context.Context, coll *mongo.Collection, id string, data interface{}) (int64, error) {

	if updateResult, err := coll.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.D{{Key: "$set", Value: data}}); err != nil {
		return updateResult.MatchedCount, err
	} else if updateResult.MatchedCount == 0 {
		if _, err := coll.InsertOne(ctx, data); err != nil {
			return 0, err
		} else {
			return 1, nil
		}
	} else {
		return updateResult.MatchedCount, nil
	}

}

func MongoDeleteFunc(ctx context.Context, coll *mongo.Collection, filter bson.D) (int64, error) {
	if result, err := coll.DeleteMany(ctx, filter); err != nil {
		return 0, err
	} else {
		return result.DeletedCount, nil
	}
}

func MongoFindFunc(ctx context.Context, coll *mongo.Collection, filter bson.D) ([]interface{}, int64, error) {

	if cursor, err := coll.Find(ctx, filter, nil); err != nil {
		return nil, 0, err
	} else {

		var count int64
		if count, err = coll.CountDocuments(ctx, filter); err != nil {
			return nil, 0, err
		}

		defer cursor.Close(ctx)
		var result []interface{}
		if err := cursor.All(ctx, &result); err != nil {
			return nil, 0, err
		} else {
			return result, count, nil
		}
	}
}
