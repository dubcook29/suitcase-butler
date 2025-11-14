package meta

import (
	"context"
	"fmt"
	"testing"

	"github.com/suitcase/butler/db"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	db.InitialRuntimeDBConnect(context.TODO(), "localhost", "27017", "", "")
	db.SetRuntimeDevMode(true, "./")
}

func TestFindAll(t *testing.T) {

	if data, count, err := FindAll(context.TODO()); err != nil {
		panic(err)
	} else {
		fmt.Println("find data, total:", count)
		for _, asset := range data {
			fmt.Println(asset)
		}
	}
}

func TestFindMany(t *testing.T) {

	if data, count, err := FindMany(context.TODO(), bson.D{{Key: "ip_address", Value: bson.M{"$ne": ""}}}); err != nil {
		panic(err)
	} else {
		fmt.Println("find data, total:", count)
		for _, asset := range data {
			fmt.Println(asset)
		}
	}
}

func TestAssetMetaDataFromDatabase(t *testing.T) {
	asset := InitialAssetEmpty()
	if err := asset.AssetMetaDataOneFromDatabase("0fbe823b-70a4-4132-bc89-6f2ad99d6994"); err != nil {
		panic(err)
	} else {
		fmt.Println(*asset)
	}

}

func TestUpdatedAndDeleted(t *testing.T) {
	asset := InitialAssetEmpty()
	asset.DomainName = "bilibili.com"
	if asset.Exist(context.TODO()) {
		if err := asset.Updated(context.TODO()); err != nil {
			panic(err)
		}
		fmt.Println("asset update successed:", asset.AssetId)
	} else {
		if err := asset.Insert(context.TODO()); err != nil {
			panic(err)
		}
		fmt.Println("asset insert successed:", asset.AssetId)
	}

	asset.IpAddress = "8.8.8.8"
	if err := asset.Updated(context.TODO()); err != nil {
		panic(err)
	} else {
		fmt.Println("asset update successed:", asset.AssetId)
	}

	if data, count, err := FindMany(context.TODO(), bson.D{{Key: "_id", Value: asset.AssetId}}); err != nil {
		panic(err)
	} else {
		fmt.Println("find data, total:", count)
		for _, asset := range data {
			fmt.Println(asset)
		}
	}

	if err := asset.Deleted(context.TODO()); err != nil {
		panic(err)
	} else if !asset.Exist(context.TODO()) {
		fmt.Println("asset delete successed:", asset.AssetId)
	} else {
		fmt.Println("asset delete failed:", asset.AssetId)
	}
}

func TestExist(t *testing.T) {
	asset := InitialAssetEmpty()
	if asset.Exist(context.TODO()) {
		fmt.Println("asset is exist:", asset.AssetId)
	} else {
		fmt.Println("asset not is exist:", asset.AssetId)
	}
}
