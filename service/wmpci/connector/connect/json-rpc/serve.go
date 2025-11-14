package connect_jsonrpc

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/suitcase/butler/wmpci"
	connect_builtin "github.com/suitcase/butler/wmpci/connector/connect/built-in"
)

type WMPJSONRPCServer interface {
	WMPCallHandle(request wmpci.WMPRequest, response *wmpci.WMPResponse) error
	WMPRefreshForCustome(custom map[string]wmpci.WMPCustom, ok *bool) error
	Registrars(ok bool, registrars *wmpci.WMPRegistrars) error
}

// Quick demonstration of converting `BuiltinServer` implementations into WMP JSON-RPC (WMPJSONRPCServer) implementations
type WMPJSONRPCServerApplication struct {
	ctx             context.Context
	done            context.CancelFunc
	host            string
	port            string
	service_name    string
	registrarErrors error
	builtin_service connect_builtin.BuiltinServer
}

func Serve(ctx context.Context, host, port, sn string, wmp WMPJSONRPCServer) error {
	serve_listen, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}

	rpc_service := rpc.NewServer()
	if err := rpc_service.RegisterName(sn, wmp); err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		serve_listen.Close()
	}()

	for {
		fmt.Println("listen prot", net.JoinHostPort(host, port), "wait next request enter.")
		conn, err := serve_listen.Accept()
		if err != nil {
			return err
		}
		go rpc_service.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

// QuickStartSimpleRPCOfBuiltinClient Quickly start an JSON-RPC Service via `connect_builtin.BuiltinServer`
func QuickStartSimpleRPCOfBuiltinClient(ctx context.Context, host, port, service_name string, builtin_service connect_builtin.BuiltinServer) error {
	var ws = new(WMPJSONRPCServerApplication)
	if err := ws.initialContext(ctx).initialRPCConnect(host, port, service_name).initialBuiltinClient(builtin_service).errors(); err != nil {
		return err
	}

	return ws.StartService()
}

func (ws *WMPJSONRPCServerApplication) initialContext(ctx context.Context) *WMPJSONRPCServerApplication {
	ws.ctx, ws.done = context.WithCancel(ctx)
	ws.registrarErrors = nil
	return ws
}

func (ws *WMPJSONRPCServerApplication) initialBuiltinClient(builtin_service connect_builtin.BuiltinServer) *WMPJSONRPCServerApplication {
	ws.builtin_service = builtin_service
	return ws
}

func (ws *WMPJSONRPCServerApplication) initialRPCConnect(host, port, service_name string) *WMPJSONRPCServerApplication {
	ws.host = host
	ws.port = port
	ws.service_name = service_name
	return ws
}

func (ws *WMPJSONRPCServerApplication) errors() error {
	return ws.registrarErrors
}

func (ws *WMPJSONRPCServerApplication) StartService() error {
	return Serve(ws.ctx, ws.host, ws.port, ws.service_name, ws)
}

func (ws *WMPJSONRPCServerApplication) Close() {
	ws.done()
}

func (ws *WMPJSONRPCServerApplication) WMPCallHandle(request wmpci.WMPRequest, response *wmpci.WMPResponse) error {
	if resp, err := ws.builtin_service.WMPService(ws.ctx, request); err != nil {
		return err
	} else {
		*response = resp
	}
	return nil
}

func (ws *WMPJSONRPCServerApplication) WMPRefreshForCustome(custom map[string]wmpci.WMPCustom, ok *bool) error {
	if ok1, err := ws.builtin_service.WMPConfig(custom); err != nil {
		return err
	} else {
		*ok = ok1
	}
	return nil
}

func (ws *WMPJSONRPCServerApplication) Registrars(ok bool, registrars *wmpci.WMPRegistrars) error {
	_, *registrars = ws.builtin_service.WMPRegist()

	return nil
}
