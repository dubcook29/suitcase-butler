package whois

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/suitcase/butler/wmpci"
	connect_jsonrpc "github.com/suitcase/butler/wmpci/connector/connect/json-rpc"
)

var test_wmpApplicationRequest = wmpci.WMPRequest{
	"domain": {"example.com"},
	"asn":    {"AS15133"},
	"ip":     {"23.192.228.84"},
}

func TestWMPService(t *testing.T) {
	var demo WMPWhois

	if response, err := demo.WMPService(context.TODO(), test_wmpApplicationRequest); err != nil {
		panic(err)
	} else {
		if data, err := json.MarshalIndent(response, "", "\t"); err != nil {
			panic(err)
		} else {
			fmt.Println("print current wmp service called response:\n", string(data))
		}
	}

}

func Test_whoisQueryHandle(t *testing.T) {
	var domains = []string{
		"google.com",
		"microsoft.com",
		"amazon.com",
	}

	var asns = []string{
		"AS15169",
		"AS8075",
		"AS16509",
	}

	var ips = []string{
		"142.251.32.36",
		"104.99.50.44",
		"108.139.7.233",
	}

	var wheres []string
	wheres = append(wheres, domains...)
	wheres = append(wheres, asns...)
	wheres = append(wheres, ips...)
	for _, where := range wheres {
		if whois_result_raw, err := whoisQueryHandle(where, 3); err != nil {
			t.Errorf("<%s> whois query handle call failed:%v\n", where, err)
		} else {
			t.Logf("<%s> whois query handle call successed:\n%v\n", where, whois_result_raw)
		}
	}
}

func Test_importWmpApplicationConfig(t *testing.T) {

}

func Test_importWmpApplicationRequest(t *testing.T) {

}

func TestQuickStartSimpleRPCService(t *testing.T) {
	var demo WMPWhois
	if builtin_service, _ := demo.WMPRegist(); builtin_service != nil {
		if err := connect_jsonrpc.QuickStartSimpleRPCOfBuiltinClient(context.TODO(), "", "8081", "wmp", builtin_service); err != nil {
			panic(err)
		}
	}
}
