package connect_builtin_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/suitcase/butler/wmpci"
	"github.com/suitcase/butler/wmpci/connector"
	connect_builtin "github.com/suitcase/butler/wmpci/connector/connect/built-in"
	wmpexample "github.com/suitcase/butler/wmpci/example/wmp_example"
)

func TestClientConnection(t *testing.T) {
	ctx, done := context.WithCancel(context.Background())

	session := connector.NewWMPCISession("./.wmpci/")

	if _, err := connector.DefaultWMPConnectorSupports.GetConnectorConnectCustomConfig("built-in"); err != nil {
		panic(err)
	} else {
		custom := map[string]wmpci.WMPCustom{
			"WMPCI": {
				Value: new(wmpexample.WMPCIApplication),
			},
		}
		if conn, err := connector.DefaultWMPConnectorSupports.ConnectorGenerator("built-in"); err != nil {
			panic(err)
		} else {
			if msg, ok, err := session.Connect(context.TODO(), conn, custom); err != nil {
				panic(err)
			} else {
				fmt.Println(string(msg), ok, "\nwmp application registration, the registration information is as follows:")
				session.Registration.PrintAllWmpregistration()
			}
		}
	}

	if response, err := session.Service(ctx, wmpci.WMPRequest{
		"domain": {
			"example.com",
			"google.com",
		},
	}); err != nil {
		panic(err)
	} else {
		fmt.Println("wmpci service called successfully, the response information is as follows:")
		for k, v := range response {
			fmt.Printf("%v:\n", k)
			for _, vv := range v {
				fmt.Printf("\t%+v\n", vv)
			}
		}
	}

	done()
}

func TestClientConnection1(t *testing.T) {

	ctx, done := context.WithCancel(context.Background())

	session := connector.NewWMPCISession("./.wmpci/")

	connector.DefaultWMPConnectorSupports.AddConnectorGenerator("builtin", func() connector.WMPConnector {
		return new(connect_builtin.BuiltinClient)
	})

	if _, err := connector.DefaultWMPConnectorSupports.GetConnectorConnectCustomConfig("builtin"); err != nil {
		panic(err)
	} else {
		custom := map[string]wmpci.WMPCustom{
			"WMPCI": {
				Value: new(wmpexample.WMPCIApplication),
			},
		}
		if conn, err := connector.DefaultWMPConnectorSupports.ConnectorGenerator("builtin"); err != nil {
			panic(err)
		} else {
			if msg, ok, err := session.Connect(context.TODO(), conn, custom); err != nil {
				panic(err)
			} else {
				fmt.Println(string(msg), ok, "\nwmp application registration, the registration information is as follows:")
				session.Registration.PrintAllWmpregistration()
			}
		}
	}

	if response, err := session.Service(ctx, wmpci.WMPRequest{
		"domain": {
			"example.com",
			"google.com",
		},
	}); err != nil {
		panic(err)
	} else {
		fmt.Println("wmpci service called successfully, the response information is as follows:")
		for k, v := range response {
			fmt.Printf("%v:\n", k)
			for _, vv := range v {
				fmt.Printf("\t%+v\n", vv)
			}
		}
	}

	done()
}
