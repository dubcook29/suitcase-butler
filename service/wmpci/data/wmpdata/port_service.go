package wmpdata

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/suitcase/butler/wmpci/data/model"
)

type PortService struct {
	model.WMPDataModelBasicStructure `json:"wmpdatamodelbasicstructure" bson:"wmpdatamodelbasicstructure"`
	ServiceId                        string              `json:"service_id" bson:"service_id"`   // unique used to label services
	Address                          string              `json:"address" bson:"address"`         // e.g. 192.168.3.1, localhost, google.com ...
	PortNumber                       int                 `json:"port_number" bson:"port_number"` // e.g. 22, 21, 80, 443, 3389, 554 ...
	Protocol                         string              `json:"protocol" bson:"protocol"`       // e.g. SSH, FTP, HTTP, HTTPS, RDP, RTSP ...
	Product                          []string            `json:"product" bson:"product"`         // e.g. Centos-SSH, WordPress, YApi, DVR, ...
	Banner                           string              `json:"banner" bson:"banner"`
	Headers                          map[string][]string `json:"headers" bson:"headers"`
	Url                              string              `json:"uri" bson:"uri"`           // e.g. [protocol]://[address]:[port] is default format
	Redirect                         string              `json:"redirect" bson:"redirect"` // the location of the last resource
}

func (p PortService) Key() string {
	return "port_service"
}

func (p PortService) Exchange(data map[string][]interface{}) map[string][]interface{} {
	if val, ok := data[p.Key()]; ok {
		data[p.Key()] = append(val, p)
	} else {
		data[p.Key()] = []interface{}{p}
	}
	return data
}

func (p PortService) DataModel() model.WMPDataModelBasicStructure {
	return p.WMPDataModelBasicStructure
}

func (p PortService) JSON() (json.RawMessage, error) {
	return json.Marshal(&p)
}

func (p PortService) New(assetId string) PortService {
	p.WMPDataModelBasicStructure.CreatedAt = time.Now()
	p.WMPDataModelBasicStructure.UpdatedAt = time.Now()
	p.WMPDataModelBasicStructure.AssetId = assetId
	p.WMPDataModelBasicStructure.Identifier = uuid.NewString()
	p.WMPDataModelBasicStructure.DataModelKey = p.Key()
	return p
}
