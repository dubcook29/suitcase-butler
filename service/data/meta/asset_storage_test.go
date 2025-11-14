package meta

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/suitcase/butler/db"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAssetMetaDataMongoInsertFunc(t *testing.T) {

	// add new asset data into databases

	asset_tmp := AssetMetaData{
		AssetId:    "test_asset_1",
		DomainName: "leavesongs.com",
		CreatedAt:  time.Now(),
	}

	if count, err := AssetMetaDataMongoInsertFunc(context.TODO(), db.GetCurrentMongoClient(), []AssetMetaData{asset_tmp}); err != nil {
		panic(err)
	} else {
		fmt.Println("[asset] insert count successed:", count)
	}

}

func TestAssetMetaDataMongoUpdateFunc(t *testing.T) {

	var assetTempData = AssetMetaData{
		GroupId:         "TestGroupID",
		OtherInputValue: "",
		AssetId:         "a7513479-30d7-42fa-983d-1cd8b68ee65b",
		DomainName:      "google.com",
		IpAddress:       "",
		OrgName:         "",
		ASNumber:        "",
	}

	assetTempData.DomainName = "youtube.net"

	_, err := AssetMetaDataMongoUpdateFunc(context.TODO(), db.GetCurrentMongoClient(), []AssetMetaData{assetTempData})
	if err != nil {
		panic(err)
	}
}

func TestAssetMetaDataMongoFindFunc(t *testing.T) {

	data, _, err := AssetMetaDataMongoFindFunc(context.TODO(), db.GetCurrentMongoClient(), bson.D{})
	if err != nil {
		panic(err)
	}
	for _, v := range data {
		fmt.Printf("%+v\n", v)
	}

}
