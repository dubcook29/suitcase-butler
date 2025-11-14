package meta

import (
	// Unexpected package references prohibited github.com/suitcase/butler/meta
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/db"
)

var sep string = ";"
var privateKeys []string = []string{
	"group_id", "GroupId", "groupid",
	"asset_id", "AssetId", "assetid",
	"org", "org_name",
	"asn", "as", "as_number",
	"ip", "ips", "address", "ip_address", "IpAddress", "ipaddress",
	"domain", "domain_name", "DomainName", "domainname",
	"input", "other", "cloud", "cdn",
	"scheduler_id", "task_id",
}

type AssetMetaData struct {
	GroupId            string    `json:"group_id" bson:"group_id"`                   // asset group uuid
	AssetId            string    `json:"asset_id" bson:"asset_id"`                   // unique used to label asset
	DomainName         string    `json:"domain_name" bson:"domain_name"`             // e.g. google.com
	IpAddress          string    `json:"ip_address" bson:"ip_address"`               // e.g. 142.250.189.206
	OrgName            string    `json:"org_name" bson:"org_name"`                   // e.g. Google, Inc.
	ASNumber           string    `json:"as_number" bson:"as_number"`                 // e.g. AS15169
	Cloud              string    `json:"cloud" bson:"cloud"`                         // cloud marking
	CDN                string    `json:"cdn" bson:"cdn"`                             // cdn marking
	OtherInputValue    string    `json:"other_input_value" bson:"other_input_value"` // e.g. twitter:@elonmusk
	ResultWmpDataLists []string  `json:"result_wmp_data_lists" bson:"result_wmp_data_lists"`
	CreatedAt          time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt          time.Time `json:"deleted_at" bson:"deleted_at"`
}

func (a *AssetMetaData) Initial() {
	a.CreatedAt = time.Now()
	a.AssetId = uuid.NewString()
}

// UpdateAsset  update latest asset data into database
func (a *AssetMetaData) UpdateAsset() error {
	a.UpdatedAt = time.Now()
	if _, err := AssetMetaDataMongoUpdateFunc(context.TODO(), db.GetCurrentMongoClient(), []AssetMetaData{*a}); err != nil {
		logrus.Errorf("[asset] asset meta data mongo update or insert is error: %v", err)
		return err
	} else {
		logrus.Debugf("[asset] asset meta data mongo update or insert is successed: %+v", a)
		return nil
	}
}

// AddWMPData update a wmp name into asset lists
func (a *AssetMetaData) AddWMPData(n string) {
	defer a.UpdateAsset()
	a.ResultWmpDataLists = uniqueStringsPreserveOrder(append(a.ResultWmpDataLists, n))
}

func uniqueStringsPreserveOrder(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	result := make([]string, 0, len(input))
	for _, s := range input {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}

func (a *AssetMetaData) PrivateKeys() []string {
	return privateKeys
}

func (a *AssetMetaData) IsPrivateKeys(s string) bool {
	for _, v := range privateKeys {
		if v == s {
			return true
		}
	}
	return false
}

func (a *AssetMetaData) Reader(k string) []interface{} {
	switch strings.ToLower(k) {
	case "asset_id":
		return stringItemsToInterfaceItems(strings.Split(a.AssetId, sep))
	case strings.ToLower("DomainName"), "domain", "domain_name":
		if a.DomainName != "" {
			return stringItemsToInterfaceItems(strings.Split(a.DomainName, sep))
		}
	case strings.ToLower("ASNumber"), "asn", "as", "as_number":
		if a.ASNumber != "" {
			return stringItemsToInterfaceItems(strings.Split(a.ASNumber, sep))
		}
	case strings.ToLower("OrgName"), "org", "org_name":
		if a.OrgName != "" {
			return stringItemsToInterfaceItems(strings.Split(a.OrgName, sep))
		}
	case strings.ToLower("IpAddress"), "ip", "ips", "address", "ip_address":
		if a.IpAddress != "" {
			return stringItemsToInterfaceItems(strings.Split(a.IpAddress, sep))
		}
	case strings.ToLower("other"), "input":
		if a.OtherInputValue != "" {
			return stringItemsToInterfaceItems([]string{a.OtherInputValue})
		}
	case "cdn":
		if a.CDN != "" {
			return stringItemsToInterfaceItems(strings.Split(a.CDN, sep))
		}
	case "cloud":
		if a.Cloud != "" {
			return stringItemsToInterfaceItems(strings.Split(a.Cloud, sep))
		}
	}
	return nil
}

func (a *AssetMetaData) Writer(k string, in []interface{}) error {
	if in == nil {
		logrus.Debugf("[asset] the data of (%s) is null", k)
		return nil
	}
	defer a.UpdateAsset()

	switch strings.ToLower(k) {
	case strings.ToLower("DomainName"), "domain", "domain_name":

		if splits, err := readUniqueDataFromSlice(a.DomainName, in); err != nil {
			return err
		} else if splits != nil {
			a.DomainName = strings.Join(splits, sep)
		}

	case strings.ToLower("ASNumber"), "asn", "as", "as_number":

		if splits, err := readUniqueDataFromSlice(a.ASNumber, in); err != nil {
			return err
		} else if splits != nil {
			a.ASNumber = strings.Join(splits, sep)
		}

	case strings.ToLower("OrgName"), "org", "org_name":

		if splits, err := readUniqueDataFromSlice(a.OrgName, in); err != nil {
			return err
		} else if splits != nil {
			a.OrgName = strings.Join(splits, sep)
		}

	case strings.ToLower("IpAddress"), "ip", "ips", "address", "ip_address":

		if splits, err := readUniqueDataFromSlice(a.IpAddress, in); err != nil {
			return err
		} else if splits != nil {
			a.IpAddress = strings.Join(splits, sep)
		}

	case "cdn":

		if splits, err := readUniqueDataFromSlice(a.CDN, in); err != nil {
			return err
		} else if splits != nil {
			a.CDN = strings.Join(splits, sep)
		}

	case "cloud":

		if splits, err := readUniqueDataFromSlice(a.Cloud, in); err != nil {
			return err
		} else if splits != nil {
			a.Cloud = strings.Join(splits, sep)
		}
	}

	return nil
}

func stringItemsToInterfaceItems(list []string) []interface{} {
	if len(list) <= 0 {
		return nil
	}
	var out []interface{}
	for _, v := range list {
		out = append(out, strings.ReplaceAll(v, " ", ""))
	}
	return out
}

func removeDuplicates1(strings []string) []string {
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

// Read unique data from []interface{}
func readUniqueDataFromSlice(raw string, in []interface{}) ([]string, error) {
	var data []string

	dataInsert := func(ins string) {
		if strings.ReplaceAll(ins, " ", "") != "" {
			data = append(data, ins)
		}
	}

	if strings.ReplaceAll(raw, " ", "") != "" {
		data = strings.Split(raw, sep)
	}

	for _, item := range in {
		if v, ok := item.(string); ok {
			dataInsert(v)
		} else if v, ok := item.(int); ok {
			dataInsert(strconv.Itoa(v))
		} else if v, ok := item.(int32); ok {
			dataInsert(strconv.Itoa(int(v)))
		} else if v, ok := item.(int64); ok {
			dataInsert(strconv.Itoa(int(v)))
		} else if v, ok := item.(uint); ok {
			dataInsert(strconv.Itoa(int(v)))
		} else {
			return nil, fmt.Errorf("[asset] data type is not a legitimate String (%T)", item)
		}
	}
	return removeDuplicates1(data), nil
}
