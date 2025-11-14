package connect_builtin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/suitcase/butler/wmpci"
)

type BuiltinClient struct {
	clientConn BuiltinServer
}

func (wc *BuiltinClient) WMPCallHandle(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error) {
	return wc.clientConn.WMPService(ctx, request)
}

func (wc *BuiltinClient) WMPConfig(ctx context.Context, conf map[string]wmpci.WMPCustom) (bool, error) {
	return wc.clientConn.WMPConfig(conf)
}

func (wc *BuiltinClient) WMPRegistrars() ([]byte, error) {
	wmpci, registr := wc.clientConn.WMPRegist()
	wc.clientConn = wmpci
	return json.Marshal(registr)
}

func (wc *BuiltinClient) WMPHealth(ctx context.Context, text ...string) (bool, error) {
	return true, nil
}

func (wc *BuiltinClient) WMPConnectStatus() (bool, error) {
	if wc.clientConn != nil {
		return true, nil
	} else {
		return false, nil
	}
}

func (wc *BuiltinClient) WMPConnection(conn ...any) ([]byte, bool, error) {
	if len(conn) >= 1 {
		for _, v := range conn {
			if vv, ok := v.(BuiltinServer); ok {
				wc.clientConn = vv
				return []byte("connection successfully"), true, nil
			}
		}
	}
	return nil, false, fmt.Errorf("connection Parameter parsing error: %v", conn...)
}

func (wc *BuiltinClient) WMPConnectorParameter() map[string]wmpci.WMPCustom {
	return map[string]wmpci.WMPCustom{
		"WMPCI": {
			Name:        "WMPCI",
			Required:    false,
			Description: "BuiltinServer interface",
		},
	}
}

func (wc *BuiltinClient) WMPConnectionParameter(conn map[string]wmpci.WMPCustom) ([]byte, bool, error) {

	if v, ok := conn["WMPCI"]; ok {
		if vv, ok := v.Value.(BuiltinServer); ok {
			wc.clientConn = vv
			return []byte("connection successfully"), true, nil
		}
	}

	return nil, false, fmt.Errorf("connection Parameter parsing error: %v", conn)
}
