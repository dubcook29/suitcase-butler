package wmp_dns

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/projectdiscovery/cdncheck"
	"github.com/suitcase/butler/wmpci"
	connect_builtin "github.com/suitcase/butler/wmpci/connector/connect/built-in"
	"github.com/suitcase/butler/wmpci/data/model"
)

var wmp_application_basic = wmpci.WMPBasic{
	Id:        "ffffffff-ffff-ffff-0002-b0bd42efd9ae",
	Name:      "wmp application of dns",
	Version:   "0.0.1",
	PushState: "opensource",
}

var wmp_application_config = map[string]wmpci.WMPCustom{
	"cdn_data_source": {
		Name:        "cdn_data_source",
		Value:       "https://raw.githubusercontent.com/projectdiscovery/cdncheck/refs/heads/main/cmd/generate-index/sources_data.json",
		Description: "Set the CDN Check data config source http url",
	},
	"nameservers": {
		Name:  "nameservers",
		Value: nameservers,
	},
	"retries": {
		Name:  "retries",
		Value: 3,
	},
}

var wmp_application_request = wmpci.WMPRequest{
	"domain_name": {"example.com"},
}

var wmp_application_response = wmpci.WMPResponse{
	"ip_address": {"93.184.215.14"},
	"^":          {nil},
	"dns":        {DNSFormats{}},
	"cdn":        {nil},
	"cloud":      {nil},
}

var nameservers []string = []string{
	"8.8.8.8",        // Google
	"1.1.1.1",        // Cloudflare
	"9.9.9.9",        // Quad9
	"208.67.222.222", // Cisco OpenDNS
	"84.200.69.80",   // DNS.WATCH
	"64.6.64.6",      // Neustar DNS
	"8.26.56.26",     // Comodo Secure DNS
	"205.171.3.65",   // Level3
	"134.195.4.2",    // OpenNIC
	"185.228.168.9",  // CleanBrowsing
	"76.76.19.19",    // Alternate DNS
	"37.235.1.177",   // FreeDNS
	"77.88.8.1",      // Yandex.DNS
	"94.140.14.140",  // AdGuard
	"38.132.106.139", // CyberGhost
	"74.82.42.42",    // Hurricane Electric
	"76.76.2.0",      // ControlD
}

type DNSFormats struct {
	Host  string   `json:"host" bson:"host"`
	A     []string `json:"a" bson:"a"`
	MX    []string `json:"mx" bson:"mx"`
	NS    []string `json:"ns" bson:"ns"`
	TXT   []string `json:"txt" bson:"txt"`
	PTR   []string `json:"ptr" bson:"ptr"`
	AAAA  []string `json:"aaaa" bson:"aaaa"`
	CNAME []string `json:"cname" bson:"cname"`
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

type WMPDnsResolve struct {
	assetId    string
	retries    int
	nameserver []string
}

func (wmp *WMPDnsResolve) WMPService(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error) {

	wmp.retries = 3
	wmp.nameserver = nameservers

	if aid, ok := request["asset_id"]; ok {
		wmp.assetId = aid[0].(string)
	}

	wmp_response_result := make(wmpci.WMPResponse)

	for key, value := range request {
		switch key {
		case "domain", "domain_name":
			for _, v := range value {
				if hostname, ok := v.(string); ok {

					if dnsResponse, err := NSLookup(hostname, wmp.retries, wmp.nameserver); err != nil {
						// return nil, err
					} else {
						wmp_data_dns := DNSFormats{
							Host:  hostname,
							MX:    dnsResponse.MX,
							NS:    dnsResponse.NS,
							A:     dnsResponse.A,
							AAAA:  dnsResponse.AAAA,
							CNAME: dnsResponse.CNAME,
							PTR:   dnsResponse.PTR,
							TXT:   dnsResponse.TXT,
						}

						// wmp_response_result["dns"] = append(wmp_response_result["dns"], wmp_data_dns)
						wmp_response_result.Add("dns", wmp_data_dns)
						wmp_response_result = model.AddReservedKeysAsIPAddress(wmp_response_result, dnsResponse.A)
						if matched, val, itemType, err := cdncheckClient.CheckDNSResponse(dnsResponse); err != nil {
							// return nil, err
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

				} else {
					return nil, fmt.Errorf("the entered [domain] parameter value cannot be parsed")
				}

			}
		}
	}

	return wmp_response_result, nil
}

func (wmp *WMPDnsResolve) WMPConfig(custom map[string]wmpci.WMPCustom) (bool, error) {
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

func (wmp *WMPDnsResolve) WMPRegist() (connect_builtin.BuiltinServer, wmpci.WMPRegistrars) {
	return wmp, wmpci.WMPRegistrars{
		WMPBasic:    wmp_application_basic,
		WMPRequest:  wmp_application_request,
		WMPResponse: wmp_application_response,
		WMPCustom:   wmp_application_config,
	}
}
