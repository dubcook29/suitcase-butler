package wmp_whois

import (
	"context"

	"github.com/likexian/whois"
	"github.com/suitcase/butler/wmpci"
	connect_builtin "github.com/suitcase/butler/wmpci/connector/connect/built-in"
)

var wmp_application_basic = wmpci.WMPBasic{
	Id:        "ffffffff-ffff-ffff-0001-b0bd42efd9ae",
	Name:      "wmp application of whois",
	Version:   "0.0.1",
	PushState: "opensource",
}

var wmp_application_config = map[string]wmpci.WMPCustom{
	"resend": {
		Name:        "resend",
		Value:       3,
		Description: "Set the number of retries for CDN Check query requests",
	},
}

var wmp_application_request = wmpci.WMPRequest{
	"domain": {"example.com"},
	"asn":    {"AS15133"},
	"ip":     {"23.192.228.84"},
}

var wmp_application_response = wmpci.WMPResponse{
	"domain_whois":     {""},
	"domain_whois_raw": {""},
	"asn_whois":        {""},
	"asn_whois_raw":    {""},
	"ip_whois":         {""},
	"ip_whois_raw":     {""},
}

type WMPWhois struct {
	assetId string
	resend  int
}

func lookupwhoisAgentClient(data []interface{}, resend int) ([]string, error) {
	var result []string
	for _, v := range data {
		if w, err := whoisAgentClient(v.(string), resend); err != nil {
			return nil, err
		} else {
			result = append(result, w)
		}
	}
	return result, nil
}

func whoisAgentClient(l string, resend int) (string, error) {
	whoisInquire := func(client *whois.Client, l string) (r string, e error) {
		for i := 0; i < resend; i++ {
			if r, e = client.Whois(l); e != nil {
				continue
			} else {
				return
			}
		}
		return
	}

	client := whois.NewClient()
	client.SetDisableStats(true)
	client.SetDisableReferral(true)
	return whoisInquire(client, l)
}

// whois.markmonitor.com
// whois.verisign-grs.com
// whois.arin.net

func (wmp *WMPWhois) WMPService(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error) {

	if v, ok := wmp_application_config["resend"].Value.(int64); ok {
		wmp.resend = int(v)
	} else {
		wmp.resend = 10
	}

	if aid, ok := request["asset_id"]; ok {
		wmp.assetId = aid[0].(string)
	}

	wmp_response_result := make(wmpci.WMPResponse)

	for key, value := range request {
		switch key {
		case "domain", "domain_name":

			if whoisRawResult, err := lookupwhoisAgentClient(value, wmp.resend); err != nil {
				return nil, err
			} else {

				for _, whoisRaw := range whoisRawResult {
					if wmp_response_result, err = wmp_response_result.Add("domain_whois_raw", whoisRaw); err != nil {
						return nil, err
					}

					// if whoisStructResult, err := wmp.domainWhoisFormatFunc(whoisRaw); err != nil {
					// 	return nil, err
					// } else {
					// 	wmp_response_result = whoisStructResult.Exchange(wmp_response_result)
					// }
				}
			}

		case "asn", "as":

			for _, v := range value {
				if w, err := whoisAgentClient(v.(string), wmp.resend); err != nil {

					return nil, err
				} else {
					if wmp_response_result, err = wmp_response_result.Add("as_whois_raw", w); err != nil {
						return nil, err
					}
				}
			}

		case "ip", "ip_address":

			for _, v := range value {
				if w, err := whoisAgentClient(v.(string), wmp.resend); err != nil {
					return nil, err
				} else {
					if wmp_response_result, err = wmp_response_result.Add("ip_whois_raw", w); err != nil {
						return nil, err
					}
				}
			}

		default:
			continue
		}
	}

	return wmp_response_result, nil
}

func (wmp *WMPWhois) WMPConfig(custom map[string]wmpci.WMPCustom) (bool, error) {
	for k, v := range custom {
		wmp_application_config[k] = v
	}

	return true, nil
}

func (wmp *WMPWhois) WMPRegist() (connect_builtin.BuiltinServer, wmpci.WMPRegistrars) {
	return wmp, wmpci.WMPRegistrars{
		WMPBasic:    wmp_application_basic,
		WMPRequest:  wmp_application_request,
		WMPResponse: wmp_application_response,
		WMPCustom:   wmp_application_config,
	}
}
