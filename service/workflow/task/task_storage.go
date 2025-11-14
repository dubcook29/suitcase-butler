package workflowtask

import (
	"context"

	"github.com/suitcase/butler/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollectionHander(client *mongo.Client) *mongo.Collection {
	return client.Database("workflow").Collection("task")
}

func WorkflowTasksFindAll(ctx context.Context, client *mongo.Client, filter bson.D) ([]WorkflowTasks, int64, error) {
	coll := getCollectionHander(client)

	if curosr, err := coll.Find(ctx, filter); err != nil {
		return nil, 0, err
	} else {
		var count int64
		if count, err = coll.CountDocuments(ctx, filter); err != nil {
			return nil, 0, err
		}

		var data []WorkflowTasks
		if err := curosr.All(ctx, &data); err != nil {
			return nil, 0, err
		} else {
			return data, count, nil
		}
	}
}

func WorkflowTasksInsert(ctx context.Context, client *mongo.Client, data []WorkflowTasks) (int64, error) {

	coll := getCollectionHander(client)

	var count int64

	for _, item := range data {
		if _, err := coll.InsertOne(ctx, db.MongoId(item, item.TaskID)); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				if _, err = coll.ReplaceOne(ctx, bson.D{{Key: "_id", Value: item.TaskID}}, db.MongoId(item, item.TaskID)); err != nil {
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

func WorkflowTasksUpdate(ctx context.Context, client *mongo.Client, data []WorkflowTasks) (int64, error) {

	coll := getCollectionHander(client)
	var count int64

	for _, item := range data {
		if updateResult, err := coll.UpdateOne(ctx, bson.D{{Key: "_id", Value: item.TaskID}}, bson.D{{Key: "$set", Value: db.MongoId(item, item.TaskID)}}); err != nil {
			return count, err
		} else if updateResult.MatchedCount == 0 {
			if _, err := coll.InsertOne(ctx, db.MongoId(item, item.TaskID)); err != nil {
				return count, err
			}
		}
		count = count + 1
	}

	return count, nil
}

func WorkflowTasksDelete(ctx context.Context, client *mongo.Client, filter bson.D) (int64, error) {
	coll := getCollectionHander(client)
	if result, err := coll.DeleteMany(ctx, filter); err != nil {
		return 0, err
	} else {
		return result.DeletedCount, nil
	}
}
