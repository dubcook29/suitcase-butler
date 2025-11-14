package wmp_subdomain

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"github.com/suitcase/butler/wmpci"
	connect_builtin "github.com/suitcase/butler/wmpci/connector/connect/built-in"
)

var wmp_application_basic = wmpci.WMPBasic{
	Id:        "ffffffff-ffff-ffff-0004-b0bd42efd9ae",
	Name:      "wmp application of subdomain",
	Version:   "0.0.1",
	PushState: "opensource",
}

var wmp_application_config = map[string]wmpci.WMPCustom{
	"threads": {
		Name:        "threads",
		Value:       10,
		Description: "Set the max thread for subdomain query requests",
		Required:    true,
	},
	"timeout": {
		Name:        "timeout",
		Value:       30,
		Description: "Set the max timeout for subdomain query requests",
		Required:    true,
	},
	"maxEnumerationTime": {
		Name:        "maxEnumerationTime",
		Value:       10,
		Description: "Set the max enumeration time for subdomain query requests",
		Required:    true,
	},
}

var wmp_application_request = wmpci.WMPRequest{
	"domain_name": {"example.com"},
}

var wmp_application_response = wmpci.WMPResponse{
	"subdomain": {nil},
	"^":         {nil},
}

type WMPSubdomain struct{}

func initSubdomainTool(ctx context.Context, threads, timeout, maxEnumerationTime int, domain string) ([]string, error) {
	subfinderOpts := &runner.Options{
		Threads:            threads,            // Thread controls the number of threads to use for active enumerations
		Timeout:            timeout,            // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: maxEnumerationTime, // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
	}

	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create subfinder runner: %v", err)
	}

	output := &bytes.Buffer{}
	var sourceMap map[string]map[string]struct{}
	// To run subdomain enumeration on a single domain
	if sourceMap, err = subfinder.EnumerateSingleDomainWithCtx(ctx, domain, []io.Writer{output}); err != nil {
		return nil, fmt.Errorf("failed to enumerate single domain(%s): %v", domain, err)
	}

	var subdomainList []string
	for subdomain := range sourceMap {
		subdomainList = append(subdomainList, subdomain)
	}

	return subdomainList, nil
}

func (wmp *WMPSubdomain) WMPService(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error) {
	var threads, timeout, maxEnumerationTime int

	if v, ok := wmp_application_config["threads"]; ok {
		if v, ok := v.Value.(int); ok {
			threads = int(v)
		}
	}

	if v, ok := wmp_application_config["timeout"]; ok {
		if v, ok := v.Value.(int); ok {
			timeout = int(v)
		}
	}

	if v, ok := wmp_application_config["maxEnumerationTime"]; ok {
		if v, ok := v.Value.(int); ok {
			maxEnumerationTime = int(v)
		}
	}

	wmp_response_result := make(wmpci.WMPResponse)

	for key, value := range request {
		switch key {
		case "domain_name", "domain":

			for _, v := range value {
				if domain, ok := v.(string); ok {
					if subdomainList, err := initSubdomainTool(ctx, threads, timeout, maxEnumerationTime, domain); err != nil {
						return nil, err
					} else {
						for _, subdomain := range subdomainList {
							wmp_response_result, _ = wmp_response_result.Add("subdomain", subdomain)
							wmp_response_result, _ = wmp_response_result.Add("^", fmt.Sprintf("%s:%s", "domain", subdomain))
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

func (wmp *WMPSubdomain) WMPConfig(custom map[string]wmpci.WMPCustom) (bool, error) {
	for k, v := range custom {
		wmp_application_config[k] = v
	}

	return true, nil
}

func (wmp *WMPSubdomain) WMPRegist() (connect_builtin.BuiltinServer, wmpci.WMPRegistrars) {
	return wmp, wmpci.WMPRegistrars{
		WMPBasic:    wmp_application_basic,
		WMPRequest:  wmp_application_request,
		WMPResponse: wmp_application_response,
		WMPCustom:   wmp_application_config,
	}
}
