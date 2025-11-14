package wmpci

type WMPRegistrars struct {
	WMPBasic    WMPBasic                 `json:"wmp_basic" yaml:"wmp_basic" bson:"wmp_basic"`
	WMPRequest  map[string][]interface{} `json:"wmp_request" yaml:"wmp_request" bson:"wmp_request"`
	WMPResponse map[string][]interface{} `json:"wmp_response" yaml:"wmp_response" bson:"wmp_response"`
	WMPCustom   map[string]WMPCustom     `json:"wmp_custom" yaml:"wmp_custom" bson:"wmp_custom"`
}
