package wmp_subdomain_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/suitcase/butler/built_wmp/wmp_subdomain"
	"github.com/suitcase/butler/wmpci"
	"github.com/suitcase/butler/wmpci/connector"
	connect_builtin "github.com/suitcase/butler/wmpci/connector/connect/built-in"
)

func TestBuiltinApplication(t *testing.T) {

	session := connector.NewWMPCISession("./.wmpci")

	if data, ok, err := session.Connect(context.TODO(),
		new(connect_builtin.BuiltinClient),
		map[string]wmpci.WMPCustom{
			"WMPCI": {
				Value: new(wmp_subdomain.WMPSubdomain),
			},
		}); err != nil {
		fmt.Println("wmp application connection errors:")
		panic(err)
	} else if ok {
		fmt.Println("wmp application connection successfully:", string(data))
	} else {
		fmt.Println("wmp application connection failed")
	}

	if response, err := session.Service(context.TODO(), wmpci.WMPRequest{"domain_name": {"example.com"}, "ip_address": {"8.8.8.8"}}); err != nil {
		panic(err)
	} else {
		fmt.Println("wmp application service request successfully")
		fmt.Printf("%+v\n", response)
	}

	if ok, err := session.Config(context.TODO(), map[string]wmpci.WMPCustom{
		"test": {Value: "Hello, this is a test"},
	}); err != nil {
		panic(err)
	} else {
		fmt.Println("wmp application config refresh:", ok)
	}

}
