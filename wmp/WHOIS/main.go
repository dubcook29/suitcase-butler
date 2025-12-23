package main

import (
	whois "butler_wmp_whois/package"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	connect_jsonrpc "github.com/suitcase/butler/wmpci/connector/connect/json-rpc"
)

func main() {
	var (
		service_address string
		service_port    string
		service_name    string
	)

	flag.StringVar(&service_address, "host", "localhost", "")
	flag.StringVar(&service_port, "port", "8081", "")
	flag.StringVar(&service_name, "name", "wmp_whois", "")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Ctrl+C -> SIGINT
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		fmt.Printf("signal received: %v, Wait for exit...\n", sig)
		cancel()
	}()

	go func(ctx context.Context) {
		var wmp_whois whois.WMPWhois
		if builtin_service, _ := wmp_whois.WMPRegist(); builtin_service != nil {
			fmt.Printf("wmpci json-rpc connector quick starting : the service name running on [%s:%s] is %s\n", service_address, service_port, service_name)
			if err := connect_jsonrpc.QuickStartSimpleRPCOfBuiltinClient(ctx, service_address, service_port, service_name, builtin_service); err != nil {
				panic("wmpci json-rpc connector quick start failed: " + err.Error())
			}
		} else {
			panic("built-in service initial failed ...")
		}
	}(ctx)

	<-ctx.Done()
}
