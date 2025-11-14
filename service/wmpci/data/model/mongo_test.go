package model_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/db"
	"github.com/suitcase/butler/wmpci/data/model"
	"github.com/suitcase/butler/wmpci/data/wmpdata"
	"go.mongodb.org/mongo-driver/bson"
)

func TestWMPDataModelMongoUpdate(t *testing.T) {
	testing.Init()

	var data = make(map[string][]interface{})

	wmpdata.DNS{
		WMPDataModelBasicStructure: model.WMPDataModelBasicStructure{
			AssetId:      "wmpdata",
			Identifier:   "wmpdata_dns_1",
			DataModelKey: "dns",
		},
		Host: "youtube.de",
		A:    []string{"31.13.73.19"},
		NS:   []string{"ns1.google.com.", "ns2.google.com.", "ns3.google.com.", "ns4.google.com."},
		AAAA: []string{time.Now().String()},
	}.Exchange(data)

	for k, item := range data {
		var dataItems []model.WMPDataModelInterface
		for _, v := range item {
			if data, ok := v.(model.WMPDataModelInterface); ok {
				dataItems = append(dataItems, data)
			} else {
				logrus.Debugf("[workflow] wmp response data format failed")
			}
		}
		model.WMPDataModelMongoUpdateFunc(context.TODO(), db.GetCurrentMongoClient(), "wmpdata", k, dataItems)
	}
}

func TestWMPDataModelMongoInsert(t *testing.T) {
	testing.Init()

	var data = make(map[string][]interface{})

	wmpdata.DNS{
		WMPDataModelBasicStructure: model.WMPDataModelBasicStructure{
			AssetId:      "wmpdata",
			Identifier:   "wmpdata_dns_1" + uuid.NewString(),
			DataModelKey: "dns",
		},
		Host: "youtube.com",
		A:    []string{"31.13.73.9"},
		NS:   []string{"ns1.google.com.", "ns2.google.com.", "ns3.google.com.", "ns4.google.com."},
		AAAA: []string{"2001::1"},
	}.Exchange(data)

	for k, item := range data {
		var dataItems []model.WMPDataModelInterface
		for _, v := range item {
			if data, ok := v.(model.WMPDataModelInterface); ok {
				dataItems = append(dataItems, data)
			} else {
				logrus.Debugf("[workflow] wmp response data format failed")
			}
		}
		model.WMPDataModelMongoInsertFunc(context.TODO(), db.GetCurrentMongoClient(), "wmpdata", k, dataItems)
	}
}

func TestWMPDataModelMongoFindFunc(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	db.InitLoggerStdout()
	if err := db.InitialRuntimeDBConnect(context.TODO(), "localhost", "27017", "", ""); err != nil {
		panic(err)
	}

	db.SetRuntimeDevMode(true, "./")

	fmt.Println(model.WMPDataModelMongoFindFunc(context.TODO(), db.GetCurrentMongoClient(), "wmpdata", "dns", model.DefaultWMPDataModel{}, bson.D{}))
}
