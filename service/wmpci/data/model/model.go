package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// WMPDataModelInterface WMP data model interface
type WMPDataModelInterface interface {
	Key() string
	Exchange(data map[string][]interface{}) map[string][]interface{}
	DataModel() WMPDataModelBasicStructure
	JSON() (json.RawMessage, error)
}

type WMPDataModelBasicStructure struct {
	AssetId      string    `json:"asset_id" bson:"asset_id"`             // Represents the asset to which this data belongs
	Identifier   string    `json:"identifier" bson:"identifier"`         // Unique identifier representing the data
	DataModelKey string    `json:"data_model_key" bson:"data_model_key"` // The name (pronoun) that represents this data
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at" bson:"deleted_at"`
}

func NewWMPDataModelBasicStructure(assetId, dataModelKey string) WMPDataModelBasicStructure {
	return WMPDataModelBasicStructure{
		AssetId:      assetId,
		DataModelKey: dataModelKey,
		Identifier:   uuid.NewString(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Deprecated: This function is no longer in use and may be removed in the future.
func AddReservedKeysAsIPAddress(data map[string][]interface{}, value []string) map[string][]interface{} {
	if data == nil {
		data = make(map[string][]interface{})
	}

	if v, ok := data["ip_address"]; ok {
		data["ip_address"] = removeDuplicates2(append(v, stringItemsToInterfaceItems(value)...))
	} else {
		data["ip_address"] = stringItemsToInterfaceItems(removeDuplicates(value))
	}

	return data
}

// Deprecated: This function is no longer in use and may be removed in the future.
func AddReservedKeysAsDomainName(data map[string][]interface{}, value []string) map[string][]interface{} {
	if data == nil {
		data = make(map[string][]interface{})
	}

	if v, ok := data["domain_name"]; ok {
		data["domain_name"] = removeDuplicates2(append(v, stringItemsToInterfaceItems(value)...))
	} else {
		data["domain_name"] = stringItemsToInterfaceItems(removeDuplicates(value))
	}

	return data
}

// Deprecated: This function is no longer in use and may be removed in the future.
func stringItemsToInterfaceItems(list []string) []interface{} {
	if list == nil {
		return nil
	}
	var out []interface{}
	for _, v := range list {
		out = append(out, v)
	}
	return out
}

// Deprecated: This function is no longer in use and may be removed in the future.
func removeDuplicates(strings []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, str := range strings {
		if _, ok := seen[str]; !ok && str != "" {
			seen[str] = struct{}{}
			result = append(result, str)
		}
	}

	return result
}

// Deprecated: This function is no longer in use and may be removed in the future.
func removeDuplicates2(input []interface{}) []interface{} {
	uniqueMap := make(map[string]struct{})
	var result []interface{}

	for _, item := range input {
		itemStr := fmt.Sprintf("%v", item)
		if _, exists := uniqueMap[itemStr]; !exists {
			uniqueMap[itemStr] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
