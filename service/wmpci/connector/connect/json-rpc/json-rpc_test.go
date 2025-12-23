package connect_jsonrpc_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/suitcase/butler/wmpci"
	"github.com/suitcase/butler/wmpci/connector"

	connect_jsonrpc "github.com/suitcase/butler/wmpci/connector/connect/json-rpc"
	wmpexample "github.com/suitcase/butler/wmpci/example/wmp_example"
)

func TestClientConnection(t *testing.T) {
	var wg sync.WaitGroup
	ctx, done := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		connect_jsonrpc.QuickStartSimpleRPCOfBuiltinClient(ctx, "", "1080", "wmp", new(wmpexample.WMPCIApplication))
	}()

	time.Sleep(4 * time.Second)

	session := connector.NewWMPCISession("./.wmpci/")

	if custom, err := connector.DefaultWMPConnectorSupports.GetConnectorConnectCustomConfig("json-rpc"); err != nil {
		panic(err)
	} else {
		if conn, err := connector.DefaultWMPConnectorSupports.ConnectorGenerator("json-rpc"); err != nil {
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

	if response, err := session.Service(context.TODO(), wmpci.WMPRequest{
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
	wg.Wait()
}

func TestServerConnection(t *testing.T) {
	var wg sync.WaitGroup
	ctx, done := context.WithCancel(context.Background())
	wg.Add(1)

	connect_jsonrpc.QuickStartSimpleRPCOfBuiltinClient(ctx, "", "1080", "wmp", new(wmpexample.WMPCIApplication))

	time.Sleep(4 * time.Second)

	done()
	wg.Wait()
}

func TestWMPCISessionCompleteTesting(t *testing.T) {
	connector.WMPCISessionCompleteTesting(&connect_jsonrpc.WMPRPCClient{}, map[string]wmpci.WMPCustom{
		"host": {
			Name:        "host",
			Value:       "localhost",
			Required:    true,
			Description: "json-rpc service hostname/ipaddress",
		},
		"port": {
			Name:        "port",
			Value:       8081,
			Required:    true,
			Description: "json-rpc service port number",
		},
		"serviceName": {
			Name:        "serviceName",
			Value:       "wmp",
			Required:    true,
			Description: "json-rpc service root-name",
		},
	})
}
