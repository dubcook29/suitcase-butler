package whois

import (
	"context"
	"fmt"

	"github.com/suitcase/butler/wmpci"
	connect_builtin "github.com/suitcase/butler/wmpci/connector/connect/built-in"

	whois_handle "github.com/likexian/whois"
)

type WMPWhois struct {
	assetId string
	resend  int
}

// import wmp application config
func (wmp WMPWhois) importWmpApplicationConfig() WMPWhois {
	if v, ok := WMPApplicationConfigDefault["resend"].Value.(int64); ok {
		wmp.resend = int(v)
	} else {
		wmp.resend = 10
	}
	return wmp
}

// import wmp application request body from wmpci request
func (wmp WMPWhois) importWmpApplicationRequest(request wmpci.WMPRequest) WMPWhois {
	if aid, ok := request["asset_id"]; ok {
		wmp.assetId = aid[0].(string)
		delete(request, "asset_id")
	}
	return wmp
}

func whoisQueryHandle(where string, resend int) (string, error) {

	var client = whois_handle.NewClient().SetDisableStats(true).SetDisableReferral(true)

	var lastError error

	for i := 0; i < resend; i++ {
		if result, err := client.Whois(where); err != nil {
			lastError = err
			continue
		} else {
			return result, nil
		}
	}

	return "", lastError

}

func (wmp WMPWhois) WMPService(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error) {

	wmp = wmp.importWmpApplicationConfig().importWmpApplicationRequest(request)

	response_result := make(wmpci.WMPResponse)

	lookQuery := func(data []interface{}) ([]string, error) {
		var result []string
		for _, v := range data {
			if w, err := whoisQueryHandle(v.(string), wmp.resend); err != nil {
				return nil, err
			} else {
				result = append(result, w)
			}
		}
		return result, nil
	}

	// lookup request all key and value, and loading wmp function
	// Keys and values ​​that cannot be understood and processed can be skipped
	for key, value := range request {
		fmt.Printf("[!] read <%s> : %+v form request\n", key, value)
		switch key {
		case "domain", "domain_name":
			if result, err := lookQuery(value); err != nil {

			} else {
				for _, raw := range result {
					if response_result, err = response_result.Add("domain_whois_raw", raw); err != nil {
						return nil, err
					}
				}
			}
		case "asn", "as":
			if result, err := lookQuery(value); err != nil {

			} else {
				for _, raw := range result {
					if response_result, err = response_result.Add("as_whois_raw", raw); err != nil {
						return nil, err
					}
				}
			}
		case "ip", "ip_address":
			if result, err := lookQuery(value); err != nil {

			} else {
				for _, raw := range result {
					if response_result, err = response_result.Add("ip_whois_raw", raw); err != nil {
						return nil, err
					}
				}
			}
		default:
			continue
		}
	}

	return response_result, nil
}

func (wmp WMPWhois) WMPConfig(custom map[string]wmpci.WMPCustom) (bool, error) {
	for k, v := range custom {
		WMPApplicationConfigDefault[k] = v
	}
	return true, nil
}

func (wmp WMPWhois) WMPRegist() (connect_builtin.BuiltinServer, wmpci.WMPRegistrars) {
	return wmp, wmpci.WMPRegistrars{
		WMPBasic:    WMPApplicationBasic,
		WMPRequest:  WMPApplicationRequestDefault,
		WMPResponse: WMPApplicationResponseDefault,
		WMPCustom:   WMPApplicationConfigDefault,
	}
}
