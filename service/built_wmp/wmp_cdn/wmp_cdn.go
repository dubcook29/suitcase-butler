package wmp_cdn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/projectdiscovery/cdncheck"
	"github.com/suitcase/butler/wmpci"
	connect_builtin "github.com/suitcase/butler/wmpci/connector/connect/built-in"
)

var wmp_application_basic = wmpci.WMPBasic{
	Id:        "ffffffff-ffff-ffff-0003-b0bd42efd9ae",
	Name:      "wmp application of cdn check",
	Version:   "0.0.1",
	PushState: "opensource",
}

var wmp_application_config = map[string]wmpci.WMPCustom{
	"resend": {
		Name:        "resend",
		Value:       3,
		Description: "Set the number of retries for CDN Check query requests",
	},
	"cdn_data_source": {
		Name:        "cdn_data_source",
		Value:       "https://raw.githubusercontent.com/projectdiscovery/cdncheck/refs/heads/main/cmd/generate-index/sources_data.json",
		Description: "Set the CDN Check data config source http url",
	},
}

var wmp_application_request = wmpci.WMPRequest{
	"domain_name": {"example.com"},
	"ip_address":  {"8.8.8.8"},
}

var wmp_application_response = wmpci.WMPResponse{
	"cdn":   {"Akamai"},
	"cloud": {"Akamai"},
}

var cdncheckClient *cdncheck.Client

func init() {

	default_cdn_data_source := "https://raw.githubusercontent.com/projectdiscovery/cdncheck/refs/heads/main/cmd/generate-index/sources_data.json"

	if client, err := initCdnCheckConfigSourceData(default_cdn_data_source, 3, nil); err != nil {
		cdncheckClient = cdncheck.New()
	} else {
		cdncheckClient = client
	}

}

// https://github.com/projectdiscovery/cdncheck
type WMPCDNCheck struct {
	assetId string
}

func initCdnCheckConfigSourceData(generatedDataUrl string, MaxRetries int, resolvers []string) (*cdncheck.Client, error) {
	response, err := http.Get(generatedDataUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var generatedData cdncheck.InputCompiled
	fullFileData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(fullFileData, &generatedData); err != nil {
		return nil, err
	}

	fmt.Println("cdn data source download successfuly from " + generatedDataUrl)

	cdncheck.DefaultCDNProviders = mapKeys(generatedData.CDN)
	cdncheck.DefaultWafProviders = mapKeys(generatedData.WAF)
	cdncheck.DefaultCloudProviders = mapKeys(generatedData.Cloud)

	return cdncheck.NewWithOpts(MaxRetries, nil)
}

func mapKeys(m map[string][]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func (wmp *WMPCDNCheck) WMPService(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error) {
	if aid, ok := request["asset_id"]; ok {
		wmp.assetId = aid[0].(string)
	}

	wmp_response_result := make(wmpci.WMPResponse)

	for key, value := range request {
		switch key {
		case "ip_address":
			for _, v := range value {
				if ips, ok := v.(string); ok {
					ip := net.ParseIP(ips)

					if matched, val, err := cdncheckClient.CheckCDN(ip); err != nil {
						return nil, err
					} else if matched {
						wmp_response_result, _ = wmp_response_result.Add("cdn", fmt.Sprintf("%s(CDN)", val))
					}

					if matched, val, err := cdncheckClient.CheckWAF(ip); err != nil {
						return nil, err
					} else if matched {
						wmp_response_result, _ = wmp_response_result.Add("cdn", fmt.Sprintf("%s(WAF)", val))
					}

					if matched, val, err := cdncheckClient.CheckCloud(ip); err != nil {
						return nil, err
					} else if matched {
						wmp_response_result, _ = wmp_response_result.Add("cloud", fmt.Sprintf("%s(Cloud)", val))
					}
				}
			}
		case "domain", "domain_name":
			for _, v := range value {
				if dn, ok := v.(string); ok {
					if matched, val, itemType, err := cdncheckClient.CheckSuffix(dn); err != nil {
						return nil, err
					} else if matched {
						switch itemType {
						case "cdn":
							wmp_response_result, _ = wmp_response_result.Add("cdn", fmt.Sprintf("%s(cdn)", val))
						case "waf":
							wmp_response_result, _ = wmp_response_result.Add("cdn", fmt.Sprintf("%s(waf)", val))
						case "cloud":
							wmp_response_result, _ = wmp_response_result.Add("cloud", fmt.Sprintf("%s(cloud)", val))
						}
					}
				}
			}
		default:
			continue
		}
	}

	return wmp_response_result, nil
}

func (wmp *WMPCDNCheck) WMPConfig(custom map[string]wmpci.WMPCustom) (bool, error) {
	for keys, v := range custom {
		switch keys {
		case "cdn_data_source":
			if oldCustom, ok := wmp_application_config[keys]; ok {
				if oldCustom.Value == v {
					continue
				} else {
					wmp_application_config[keys] = v
					if vv, ok := v.Value.(string); ok {
						if client, err := initCdnCheckConfigSourceData(vv, 3, nil); err != nil {
							cdncheckClient = cdncheck.New()
						} else {
							cdncheckClient = client
						}
					}
				}
			}
		default:
			if v.Value != nil {
				wmp_application_config[keys] = v
			}
		}
	}

	return true, nil
}

func (wmp *WMPCDNCheck) WMPRegist() (connect_builtin.BuiltinServer, wmpci.WMPRegistrars) {
	return wmp, wmpci.WMPRegistrars{
		WMPBasic:    wmp_application_basic,
		WMPRequest:  wmp_application_request,
		WMPResponse: wmp_application_response,
		WMPCustom:   wmp_application_config,
	}
}
