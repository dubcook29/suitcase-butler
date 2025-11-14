package wmpdata

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/suitcase/butler/wmpci/data/model"
)

type Tech struct {
	model.WMPDataModelBasicStructure `json:"wmpdatamodelbasicstructure" bson:"wmpdatamodelbasicstructure"`
	PortServiceId                    string `json:"port_service_id,omitempty" bson:"port_service_id"`

	Method      string  `json:"method,omitempty" bson:"method"`
	Name        string  `json:"name,omitempty" bson:"name"`               // name
	Types       string  `json:"types,omitempty" bson:"types"`             // language/cms/server/cloud/cdn/service/
	Description string  `json:"description,omitempty" bson:"description"` // describe the fingerprint
	Signature   string  `json:"signature,omitempty" bson:"signature"`     //
	Accuracy    float64 `json:"accuracy,omitempty" bson:"accuracy"`       // 0 ～ 10.0
}

func (t Tech) Key() string {
	return "tech"
}

func (t Tech) Exchange(data map[string][]interface{}) map[string][]interface{} {
	if val, ok := data[t.Key()]; ok {
		data[t.Key()] = append(val, t)
	} else {
		data[t.Key()] = []interface{}{t}
	}
	return data
}

func (t Tech) DataModel() model.WMPDataModelBasicStructure {
	return t.WMPDataModelBasicStructure
}

func (t Tech) JSON() (json.RawMessage, error) {
	return json.Marshal(&t)
}

func (t Tech) New(assetId string) Tech {
	t.WMPDataModelBasicStructure.CreatedAt = time.Now()
	t.WMPDataModelBasicStructure.UpdatedAt = time.Now()
	t.WMPDataModelBasicStructure.AssetId = assetId
	t.WMPDataModelBasicStructure.Identifier = uuid.NewString()
	t.WMPDataModelBasicStructure.DataModelKey = t.Key()
	return t
}
