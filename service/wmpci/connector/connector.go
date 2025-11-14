package connector

import (
	"context"
	"fmt"
	"sync"

	"github.com/suitcase/butler/wmpci"
	connect_builtin "github.com/suitcase/butler/wmpci/connector/connect/built-in"
	connect_jsonrpc "github.com/suitcase/butler/wmpci/connector/connect/json-rpc"
)

var DefaultWMPConnectorSupports *WMPConnectorSupports

func init() {
	DefaultWMPConnectorSupports = NewWMPConnectorSupports()

	DefaultWMPConnectorSupports.AddConnectorGenerator("json-rpc", func() WMPConnector {
		return new(connect_jsonrpc.WMPRPCClient)
	})

	DefaultWMPConnectorSupports.AddConnectorGenerator("built-in", func() WMPConnector {
		return new(connect_builtin.BuiltinClient)
	})
}

type WMPConnector interface {
	WMPRegistrars() ([]byte, error)
	WMPCallHandle(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error)
	WMPConfig(ctx context.Context, conf map[string]wmpci.WMPCustom) (bool, error)
	WMPHealth(ctx context.Context, text ...string) (bool, error)
	WMPConnectStatus() (bool, error)
	WMPConnectorParameter() map[string]wmpci.WMPCustom
	WMPConnectionParameter(conn map[string]wmpci.WMPCustom) ([]byte, bool, error)
}

type WMPConnectorGenerator struct {
	Generator func() WMPConnector
}

type WMPApplication interface{}

type WMPConnectorSupports struct {
	mu                  sync.RWMutex
	connectorGenerators map[string]func() WMPConnector
}

func NewWMPConnectorSupports() *WMPConnectorSupports {

	return &WMPConnectorSupports{
		connectorGenerators: make(map[string]func() WMPConnector),
	}
}

func (wcs *WMPConnectorSupports) AddConnectorGenerator(name string, generator func() WMPConnector) (bool, error) {
	wcs.mu.Lock()
	defer wcs.mu.Unlock()

	if _, ok := wcs.connectorGenerators[name]; !ok {
		wcs.connectorGenerators[name] = generator
		return true, nil
	}

	return false, fmt.Errorf("this [%s] generator already exist", name)
}

func (wcs *WMPConnectorSupports) DelConnectorGenerator(name string) (bool, error) {
	wcs.mu.Lock()
	defer wcs.mu.Unlock()

	if _, ok := wcs.connectorGenerators[name]; ok {
		delete(wcs.connectorGenerators, name)
		return true, nil
	}

	return false, fmt.Errorf("this [%s] generator does not exist", name)
}

func (wcs *WMPConnectorSupports) GetConnectorGenerator(name string) (func() WMPConnector, error) {
	wcs.mu.RLock()
	defer wcs.mu.RUnlock()

	if generator, ok := wcs.connectorGenerators[name]; ok {
		return generator, nil
	}

	return nil, fmt.Errorf("this [%s] generator does not exist", name)
}

func (wcs *WMPConnectorSupports) GetConnectorGeneratorList() []string {
	wcs.mu.RLock()
	defer wcs.mu.RUnlock()

	var lists []string

	for name := range wcs.connectorGenerators {
		lists = append(lists, name)
	}

	return lists
}

func (wcs *WMPConnectorSupports) ConnectorGenerator(name string) (WMPConnector, error) {
	if generator, err := wcs.GetConnectorGenerator(name); err != nil {
		return nil, err
	} else {
		return generator(), nil
	}
}

func (wcs *WMPConnectorSupports) GetConnectorConnectCustomConfig(name string) (map[string]wmpci.WMPCustom, error) {
	if connector, err := wcs.ConnectorGenerator(name); err != nil {
		return nil, err
	} else {
		return connector.WMPConnectorParameter(), nil
	}
}

func (wcs *WMPConnectorSupports) GetConnectorWithCustomConfig(name string, conn map[string]wmpci.WMPCustom) (WMPConnector, []byte, error) {

	if connector, err := wcs.ConnectorGenerator(name); err != nil {
		return nil, nil, err
	} else {
		if data, ok, err := connector.WMPConnectionParameter(conn); err != nil {
			return nil, data, err
		} else if ok {
			return connector, data, nil
		} else {
			return connector, data, fmt.Errorf("this [%s] connector connection failed: %v", name, err)
		}
	}

}
