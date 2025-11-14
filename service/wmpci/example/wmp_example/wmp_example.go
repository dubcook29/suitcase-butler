package wmpexample

import (
	"context"
	"net/http"

	"github.com/suitcase/butler/wmpci"
	connect_builtin "github.com/suitcase/butler/wmpci/connector/connect/built-in"
)

var wmp_application_basic = wmpci.WMPBasic{
	Id:          "0123456789-abcdefg",
	Name:        "wmp application of example",
	Version:     "0.0.1",
	PushState:   "opensource",
	Description: "This is a WMP application for demonstration and is only used during the test process to verify the logical correctness of wmpci",
	Certificate: "",
	Copyright: []wmpci.Author{
		{
			AuthorName:   "butler-example",
			AuthorEmail:  "root@example.com",
			AuthorAvatar: "",
			Sponsor:      nil,
		},
		{
			AuthorName:   "suitcase-example",
			AuthorEmail:  "suitcase@example.com",
			AuthorAvatar: "",
			Sponsor:      nil,
		},
	},
}

var wmp_application_config = map[string]wmpci.WMPCustom{
	"username": {
		Name:  "username",
		Value: "admin",
	},
	"password": {
		Name:  "password",
		Value: "123456",
	},
}

var wmp_application_request = wmpci.WMPRequest{
	"domain": {
		"yahoo.com",
		"google.com",
	},
}

var wmp_application_response = wmpci.WMPResponse{
	"username": {
		"admin",
		"root",
	},
}

type WMPCIApplication struct {
}

func (wmp *WMPCIApplication) WMPService(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error) {
	response := make(map[string][]interface{})

	response["status_code"] = []interface{}{}
	response["content_length"] = []interface{}{}
	response["request_url"] = []interface{}{}

	for k, v := range request {
		switch k {
		case "domain", "domain_name":
			for _, vv := range v {
				if domain, ok := vv.(string); ok {
					req, err := http.NewRequest("HEAD", "https://"+domain, nil)
					if err != nil {
						return nil, err
					}

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					response["status_code"] = append(response["status_code"], resp.Status)
					response["content_length"] = append(response["content_length"], resp.ContentLength)
					response["request_url"] = append(response["request_url"], resp.Request.URL.String())
				}
			}

		}
	}

	return response, nil
}

func (wmp *WMPCIApplication) WMPConfig(custom map[string]wmpci.WMPCustom) (bool, error) {
	for k, v := range custom {
		wmp_application_config[k] = v
	}
	return true, nil
}

func (wmp *WMPCIApplication) WMPRegist() (connect_builtin.BuiltinServer, wmpci.WMPRegistrars) {
	return wmp, wmpci.WMPRegistrars{
		WMPBasic:    wmp_application_basic,
		WMPRequest:  wmp_application_request,
		WMPResponse: wmp_application_response,
		WMPCustom:   wmp_application_config,
	}
}
