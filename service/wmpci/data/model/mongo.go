package model

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo 单纯的处理 WMPDataModelInterface 数据添加 `_id` 属性标签，输出 document interface{}
func Mongo(d WMPDataModelInterface) interface{} {
	// 获取原始结构体的类型
	originalVal := reflect.ValueOf(d)
	originalType := originalVal.Type()

	// 创建一个新的结构体类型
	newFields := []reflect.StructField{}

	// 将原始结构体的字段添加到新的字段列表中
	for i := 0; i < originalType.NumField(); i++ {
		field := originalType.Field(i)
		newFields = append(newFields, field)
	}

	// 添加新的字段
	newField := reflect.StructField{
		Name: "MongoId",
		Type: reflect.TypeOf(d.DataModel().Identifier),
		Tag:  reflect.StructTag((`bson:"_id"`)),
	}
	newFields = append(newFields, newField)
	newStructType := reflect.StructOf(newFields)
	// 创建新的结构体实例
	newStructValue := reflect.New(newStructType).Elem()

	// 将原始结构体的值复制到新的结构体中
	for i := 0; i < originalVal.NumField(); i++ {
		newStructValue.Field(i).Set(originalVal.Field(i))
	}

	// 设置新字段的值
	newStructValue.Field(len(newFields) - 1).Set(reflect.ValueOf(d.DataModel().Identifier))

	return newStructValue.Interface()

}

func WMPDataModelMongoFindFunc(ctx context.Context, client *mongo.Client, dn, cn string, resultType WMPDataModelInterface, filter bson.D) ([]WMPDataModelInterface, int64, error) {
	if resultType == nil {
		return nil, 0, fmt.Errorf("resultType is null interface")
	}
	coll := client.Database(dn).Collection(cn)
	var count int64

	if cursor, err := coll.Find(ctx, filter, nil); err != nil {
		return nil, 0, err
	} else {
		var results []WMPDataModelInterface
		// 使用反射创建 resultType 的实例
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			// 创建一个新的实例
			resultValue := reflect.New(reflect.TypeOf(resultType)).Interface()

			// 解码到新实例
			if err := cursor.Decode(resultValue); err != nil {
				return nil, count, nil
			}

			// 将结果添加到切片中
			results = append(results, resultValue.(WMPDataModelInterface))
			count = count + 1
		}

		return results, count, nil
	}
}

func WMPDataModelMongoInsertFunc(ctx context.Context, client *mongo.Client, dn, cn string, data []WMPDataModelInterface) (int64, error) {
	var count int64
	coll := client.Database(dn).Collection(cn)

	for _, item := range data {
		if _, err := coll.InsertOne(ctx, Mongo(item)); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				if _, err = coll.ReplaceOne(ctx, bson.D{{Key: "_id", Value: item.DataModel().AssetId}}, Mongo(item)); err != nil {
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

func WMPDataModelMongoUpdateFunc(ctx context.Context, client *mongo.Client, dn, cn string, data []WMPDataModelInterface) (int64, error) {
	var count int64
	coll := client.Database(dn).Collection(cn)

	for _, item := range data {
		if updateResult, err := coll.UpdateOne(ctx, bson.D{{Key: "_id", Value: item.DataModel().Identifier}}, bson.D{{Key: "$set", Value: item}}); err != nil {
			return 0, err
		} else if updateResult.MatchedCount == 0 {
			if _, err := coll.InsertOne(ctx, item); err != nil {
				return 0, err
			}
		}
		count = count + 1
	}

	return count, nil
}

func WMPDataModelMongoDeleteFunc(ctx context.Context, client *mongo.Client, dn, cn string, filter bson.D) (int64, error) {
	coll := client.Database(dn).Collection(cn)
	if result, err := coll.DeleteMany(ctx, filter); err != nil {
		return 0, err
	} else {
		return result.DeletedCount, nil
	}
}
