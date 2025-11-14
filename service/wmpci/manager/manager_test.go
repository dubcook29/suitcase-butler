package manager

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/suitcase/butler/wmpci"
	"github.com/suitcase/butler/wmpci/connector"
	connect_jsonrpc "github.com/suitcase/butler/wmpci/connector/connect/json-rpc"
	wmpexample "github.com/suitcase/butler/wmpci/example/wmp_example"
)

func TestConnectorSessionFunctionForJSONRPC(t *testing.T) {

	var wg sync.WaitGroup
	ctx, done := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		connect_jsonrpc.QuickStartSimpleRPCOfBuiltinClient(ctx, "", "1080", "wmp", new(wmpexample.WMPCIApplication))
	}()

	// initial wmp session manager
	session_manager := InitialSessionManager(ctx, "./")

	// 1. get a list of connectors
	if lists, err := session_manager.ConnectorSupportLists(); err != nil {
		panic(err)
	} else {
		fmt.Println("system connector support lists:", lists)
	}

	// 2. select a connector and get connector binding parameters
	if custom, err := session_manager.SelectConnectorConnectionCumstom("json-rpc"); err != nil {
		panic(err)
	} else {
		fmt.Println("`json-rpc` connector custom query ok")
		// 3. bind connection parameters and create session
		if welcome, err := session_manager.ConnectorConnectionSession("json-rpc", custom); err != nil {
			fmt.Println("`json-rpc` connector session created failed")
			panic(err)
		} else {
			fmt.Println("`json-rpc` connector session created ok,", welcome)
		}
	}

	// print all session information
	if data, err := session_manager.SessionMap(); err != nil {
		panic(err)
	} else {
		fmt.Println(string(data))

	}

	for _, session := range session_manager.sessions {

		// verify config interface functionality
		test_sessionConfig(session, map[string]wmpci.WMPCustom{
			"whoami": {
				Name:  "whoami",
				Value: "root",
			},
		})

		// Verify service interface functionality
		test_sessionService(session, map[string][]interface{}{
			"domain": {
				"example.com",
				"google.com",
			},
		})

	}

	done()
	wg.Wait()
}

func TestConnectorSessionFunctionForBuiltin(t *testing.T) {

	ctx, done := context.WithCancel(context.Background())

	// initial wmp session manager
	session_manager := InitialSessionManager(ctx, "./")

	// 1. get a list of connectors
	if lists, err := session_manager.ConnectorSupportLists(); err != nil {
		panic(err)
	} else {
		fmt.Println("system connector support lists:", lists)
	}

	// 2. select a connector and get connector binding parameters
	if _, err := session_manager.SelectConnectorConnectionCumstom("built-in"); err != nil {
		panic(err)
	} else {
		fmt.Println("`built-in` connector custom query ok")
		// 3. bind connection parameters and create session
		custom := map[string]wmpci.WMPCustom{
			"WMPCI": {
				Value: new(wmpexample.WMPCIApplication),
			},
		}
		if welcome, err := session_manager.ConnectorConnectionSession("built-in", custom); err != nil {
			fmt.Println("`built-in` connector session created failed")
			panic(err)
		} else {
			fmt.Println("`built-in` connector session created ok,", welcome)
		}
	}

	// print all session information
	if data, err := session_manager.SessionMap(); err != nil {
		panic(err)
	} else {
		fmt.Println(string(data))

	}

	for _, session := range session_manager.sessions {

		// verify config interface functionality
		test_sessionConfig(session, map[string]wmpci.WMPCustom{
			"whoami": {
				Name:  "whoami",
				Value: "root",
			},
		})

		// Verify service interface functionality
		test_sessionService(session, map[string][]interface{}{
			"domain": {
				"example.com",
				"google.com",
			},
		})

	}

	done()
}

// The verification here is mainly to confirm that the logical relationship between
// "manager -> session -> connector -> wmp_application" is normal and can complete
// the normal wmp application call.

func test_sessionConfig(session *connector.WMPCISessions, update_config map[string]wmpci.WMPCustom) {
	fmt.Println("=== RUN: start verify session config function")
	defer fmt.Println("--- PASS: end verify session config function")

	old_config := session.Registration.RegistWMPCustom
	compareMaps(old_config, update_config)

	if ok, err := session.Config(context.TODO(), update_config); err != nil {
		panic(err)
	} else if ok {
		fmt.Println("session config update successfully")
	} else {
		fmt.Println("session config update failed")
	}

}

func compareMaps(old_config, new_config map[string]wmpci.WMPCustom) {
	keys := make(map[string]struct{})

	for k := range old_config {
		keys[k] = struct{}{}
	}
	for k := range new_config {
		keys[k] = struct{}{}
	}

	for k := range keys {
		val1, ok1 := old_config[k]
		val2, ok2 := new_config[k]

		if ok1 && ok2 {
			fmt.Printf("\t[%s] <updated> (%v) >>> (%v)  \n", k, val1.Value, val2.Value)
		} else if ok1 {
			fmt.Printf("\t[%s] <keep-it> (%v) >>> (%v) \n", k, val1.Value, val1.Value)
		} else if ok2 {
			fmt.Printf("\t[%s] <created> (null) >>> (%v) \n", k, val2.Value)
		}
	}
}

func test_sessionService(session *connector.WMPCISessions, request map[string][]interface{}) {
	fmt.Println("=== RUN: start verify session service function")
	defer fmt.Println("--- PASS: end verify session service function")
	if response, err := session.Service(context.TODO(), request); err != nil {
		panic(err)
	} else {
		for index, value := range response {
			fmt.Printf("\t=== print the value of [%s] in the response\n", index)
			for i, v := range value {
				fmt.Printf("\t\t[%d] %+v\n", i, v)
			}
			fmt.Printf("\t--- close print the value of [%s] in the response\n", index)
		}
	}
}
