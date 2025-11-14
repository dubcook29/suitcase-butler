package wmpci

import (
	"encoding/json"
	"fmt"

	"github.com/suitcase/butler/wmpci/data/model"
)

// Deprecated: This function is no longer in use and may be removed in the future.
func (r WMPRequest) Add(key string, value interface{}) (WMPRequest, error) {
	if verifyWMPDataExchangeTypes(value) {
		r[key] = append(r[key], value)
		return r, nil
	} else {
		return r, fmt.Errorf("[wmprequest] the value (%+v) entered is not an acceptable value type: %T", value, value)
	}
}

// Deprecated: This function is no longer in use and may be removed in the future.
func (r WMPResponse) Add(key string, value interface{}) (WMPResponse, error) {
	if verifyWMPDataExchangeTypes(value) {
		r[key] = append(r[key], value)
		return r, nil
	} else {
		return r, fmt.Errorf("[wmpresponse] the value (%+v) entered is not an acceptable value type: %T", value, value)
	}
}

// Deprecated: This function is no longer in use and may be removed in the future.
func verifyWMPDataExchangeTypes(v interface{}) bool {
	switch v.(type) {
	case int, float64, string, nil:
		return true
	case []byte:
		return true
	case json.RawMessage:
		return isValidJSON(v.([]byte))
	case model.WMPDataModelInterface:
		return true
	case interface{}:
		return true
	default:
		return false
	}
}

func isValidJSON(data []byte) bool {
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}
