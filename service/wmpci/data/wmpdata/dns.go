package wmpdata

import (
	"encoding/json"
	"time"

	"github.com/suitcase/butler/wmpci/data/model"
)

type DNS struct {
	model.WMPDataModelBasicStructure `json:"wmpdatamodelbasicstructure" bson:"wmpdatamodelbasicstructure"`

	Host  string   `json:"host" bson:"host"`
	A     []string `json:"a" bson:"a"`
	MX    []string `json:"mx" bson:"mx"`
	NS    []string `json:"ns" bson:"ns"`
	TXT   []string `json:"txt" bson:"txt"`
	PTR   []string `json:"ptr" bson:"ptr"`
	AAAA  []string `json:"aaaa" bson:"aaaa"`
	CNAME []string `json:"cname" bson:"cname"`
}

func (d DNS) Key() string {
	return "dns"
}

func (d DNS) Exchange(data map[string][]interface{}) map[string][]interface{} {
	if val, ok := data[d.Key()]; ok {
		data[d.Key()] = append(val, d)
	} else {
		data[d.Key()] = []interface{}{d}
	}
	d.WMPDataModelBasicStructure.UpdatedAt = time.Now()
	return data
}

func (d DNS) DataModel() model.WMPDataModelBasicStructure {
	return d.WMPDataModelBasicStructure
}

func (d DNS) JSON() (json.RawMessage, error) {
	return json.Marshal(&d)
}

func (d DNS) New(assetId string) DNS {
	d.WMPDataModelBasicStructure = model.NewWMPDataModelBasicStructure(assetId, d.Key())
	return d
}
