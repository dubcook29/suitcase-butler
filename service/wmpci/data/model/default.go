package model

import (
	"encoding/json"
	"time"
)

type DefaultWMPDataModel struct {
	MetaData WMPDataModelBasicStructure `json:"metadata" bson:"metadata"`

	Value   interface{} `json:"value" bson:"value"`
	Formats string      `json:"formats" bson:"formats"`
	From    string      `json:"from" bson:"from"`
}

func (d DefaultWMPDataModel) Key() string {
	return d.MetaData.DataModelKey
}

func (d DefaultWMPDataModel) Exchange(data map[string][]interface{}) map[string][]interface{} {
	if val, ok := data[d.Key()]; ok {
		data[d.Key()] = append(val, d)
	} else {
		data[d.Key()] = []interface{}{d}
	}
	d.MetaData.UpdatedAt = time.Now()
	return data
}

func (d DefaultWMPDataModel) DataModel() WMPDataModelBasicStructure {
	return d.MetaData
}

func (d DefaultWMPDataModel) JSON() (json.RawMessage, error) {
	return json.Marshal(&d)
}
