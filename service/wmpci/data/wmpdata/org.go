package wmpdata

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/suitcase/butler/wmpci/data/model"
)

type Org struct {
	model.WMPDataModelBasicStructure `json:"wmpdatamodelbasicstructure" bson:"wmpdatamodelbasicstructure"`

	OrgName       string   `json:"org_name,omitempty" bson:"org_name"`
	FullNames     []string `json:"full_names,omitempty" bson:"full_names"`
	Description   []string `json:"description,omitempty" bson:"description"`
	OfficeAddress []string `json:"office_address,omitempty" bson:"office_address"`
	OfficePhones  []string `json:"office_phones,omitempty" bson:"office_phones"`
	OfficeEmails  []string `json:"office_emails,omitempty" bson:"office_emails"`
	OfficeStaffs  []string `json:"office_staffs,omitempty" bson:"office_staffs"`
	RegisterInfos []string `json:"register_infos,omitempty" bson:"register_infos"`
}

func (o Org) Key() string {
	return "org"
}

func (o Org) Exchange(data map[string][]interface{}) map[string][]interface{} {
	if val, ok := data[o.Key()]; ok {
		data[o.Key()] = append(val, o)
	} else {
		data[o.Key()] = []interface{}{o}
	}
	return data
}

func (o Org) DataModel() model.WMPDataModelBasicStructure {
	return o.WMPDataModelBasicStructure
}

func (o Org) JSON() (json.RawMessage, error) {
	return json.Marshal(&o)
}

func (o Org) New(assetId string) Org {
	o.WMPDataModelBasicStructure.CreatedAt = time.Now()
	o.WMPDataModelBasicStructure.UpdatedAt = time.Now()
	o.WMPDataModelBasicStructure.AssetId = assetId
	o.WMPDataModelBasicStructure.Identifier = uuid.NewString()
	o.WMPDataModelBasicStructure.DataModelKey = o.Key()
	return o
}
