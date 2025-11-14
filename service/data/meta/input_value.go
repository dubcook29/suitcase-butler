package meta

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/suitcase/butler/data/verification"
)

// construct a general string as AssetMetaData
func AssetMetaDataConstruct(iv string) (AssetMetaData, error) {
	var key, value string
	if iv_list := strings.Split(iv, ":"); len(iv_list) >= 2 {
		key = iv_list[0]
		value = strings.Join(iv_list[1:], ":")
	} else {
		return AssetMetaData{}, errors.New("input is wrong")
	}

	asset := AssetMetaData{
		AssetId:   uuid.NewString(),
		CreatedAt: time.Now(),
	}

	switch strings.Trim(key, " ") {
	case "domain", "domain_name":
		if verification.IsValidDomain(value) {
			asset.DomainName = value
		}
	case "ips", "ip", "ipaddress", "ip_address":
		if verification.IsValidIPv4(value) {
			asset.IpAddress = value
		}
	case "asn", "as", "as_number":
		asset.ASNumber = value
	case "org", "org_name":
		asset.OrgName = value
	case "":
		return asset, errors.New("input key is empty")
	default:
		asset.OtherInputValue = fmt.Sprintf("%s:%s", key, value)
	}

	return asset, nil
}
