package connect_jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/suitcase/butler/wmpci"
)

type WMPRPCClient struct {
	host        string
	port        string
	serviceName string
	clientConn  *rpc.Client
}

func NewWMPRPCClient(host, port, sn string) (*WMPRPCClient, error) {
	client := &WMPRPCClient{
		host:        host,
		port:        port,
		serviceName: sn,
	}
	if r, err := Client(host, port); err != nil {
		return nil, err
	} else {
		client.clientConn = r
		return client, nil
	}
}

func (wc *WMPRPCClient) Close() error {
	return wc.clientConn.Close()
}

func (wc *WMPRPCClient) WMPCallHandle(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error) {
	methodName := wc.serviceName + ".WMPCallHandle"
	var response wmpci.WMPResponse
	if err := wc.clientConn.Call(methodName, request, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (wc *WMPRPCClient) WMPConfig(ctx context.Context, conf map[string]wmpci.WMPCustom) (bool, error) {
	methodName := wc.serviceName + ".WMPRefreshForCustome"
	var ok bool
	if err := wc.clientConn.Call(methodName, conf, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (wc *WMPRPCClient) WMPRegistrars() ([]byte, error) {
	methodName := wc.serviceName + ".Registrars"
	var wmp_registrars wmpci.WMPRegistrars
	if err := wc.clientConn.Call(methodName, true, &wmp_registrars); err != nil {
		return nil, err
	} else {
		return json.Marshal(wmp_registrars)
	}
}

func (wc *WMPRPCClient) WMPHealth(ctx context.Context, text ...string) (bool, error) {
	return true, nil
}

func (wc *WMPRPCClient) WMPConnectStatus() (bool, error) {
	if wc.clientConn != nil {
		return true, nil
	} else {
		return false, nil
	}
}

func (wc *WMPRPCClient) WMPConnection(conn ...any) ([]byte, bool, error) {
	if len(conn) >= 3 {
		for k, v := range conn {
			if vv, ok := v.(string); ok {
				switch k {
				case 0:
					wc.host = vv
				case 1:
					wc.port = vv
				case 2:
					wc.serviceName = vv
				}
			}
		}

		if client, err := Client(wc.host, wc.port); err != nil {
			return nil, false, err
		} else {
			wc.clientConn = client
			return []byte("connection successfully"), true, nil
		}
	} else {
		return nil, false, fmt.Errorf("connection Parameter parsing error: %v", conn...)
	}
}

// ConnectorParameter
func (wc *WMPRPCClient) WMPConnectorParameter() map[string]wmpci.WMPCustom {
	return map[string]wmpci.WMPCustom{
		"host": {
			Name:        "host",
			Value:       "localhost",
			Required:    true,
			Description: "json-rpc service hostname/ipaddress",
		},
		"port": {
			Name:        "port",
			Value:       1080,
			Required:    true,
			Description: "json-rpc service port number",
		},
		"serviceName": {
			Name:        "serviceName",
			Value:       "wmp",
			Required:    true,
			Description: "json-rpc service root-name",
		},
	}
}

func (wc *WMPRPCClient) WMPConnectionParameter(conn map[string]wmpci.WMPCustom) ([]byte, bool, error) {
	for k, v := range conn {
		switch k {
		case "host":
			if vv, ok := v.Value.(string); ok {
				wc.host = vv
			} else {
				return nil, false, fmt.Errorf("connection Parameter parsing error: [%s] %v (%T)", k, v.Value, v.Value)
			}
		case "port":
			if vv, ok := v.Value.(int); ok {
				wc.port = fmt.Sprint(vv)
			} else if vv, ok := v.Value.(int64); ok {
				wc.port = fmt.Sprint(vv)
			} else if vv, ok := v.Value.(string); ok {
				wc.port = fmt.Sprint(vv)
			} else if vv, ok := v.Value.(float32); ok {
				wc.port = fmt.Sprint(vv)
			} else if vv, ok := v.Value.(float64); ok {
				wc.port = fmt.Sprint(vv)
			} else {
				return nil, false, fmt.Errorf("connection Parameter parsing error: [%s] %v (%T)", k, v.Value, v.Value)
			}
		case "serviceName":
			if vv, ok := v.Value.(string); ok {
				wc.serviceName = vv
			} else {
				return nil, false, fmt.Errorf("connection Parameter parsing error: [%s] %v (%T)", k, v.Value, v.Value)
			}
		}
	}

	if client, err := Client(wc.host, wc.port); err != nil {
		return nil, false, fmt.Errorf("connection wmpci error: %v", err)
	} else {
		wc.clientConn = client
		return []byte("connection successfully"), true, nil
	}
}

func Client(host, port string) (*rpc.Client, error) {
	client_conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}

	client_of_wmpci := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(client_conn))
	return client_of_wmpci, nil
}
