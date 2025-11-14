package connector

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/suitcase/butler/wmpci"
)

/* **** wmp new registration  **** */

type WMPCISessions struct {
	SessionId           string                     `json:"session_id"`
	LastConnectionAt    int64                      `json:"last_connection"`
	FirstConnectionAt   int64                      `json:"first_connection"`
	ConnectionLogTables []string                   `json:"connection_logs"`
	Registration        *WMPRegistration           `json:"application_registration"`
	ConnectionCustom    map[string]wmpci.WMPCustom `json:"connection_custom"`

	ctx                  context.Context
	workHomePath         string
	wmp_connector        WMPConnector
	wmp_connector_status bool
	wmp_connector_errors string
}

func NewWMPCISession(workpath string) *WMPCISessions {
	ensureDir(workpath)
	return &WMPCISessions{
		SessionId:        uuid.NewString(),
		workHomePath:     workpath,
		ConnectionCustom: make(map[string]wmpci.WMPCustom),
	}
}

// WMPCISessions.Connect
func (session *WMPCISessions) Connect(ctx context.Context, connector WMPConnector, conn map[string]wmpci.WMPCustom) ([]byte, bool, error) {
	// 1. check connector is the status active ?
	if ok, err := connector.WMPConnectStatus(); err != nil {
		// not is the status active, and connection error
		return nil, false, err
	} else if !ok {
		if data, ok, err := connector.WMPConnectionParameter(conn); err != nil {
			return nil, false, err
		} else if !ok {
			return nil, false, fmt.Errorf("no activated connection")
		} else {
			if ok, err := session.registrars(connector); err != nil {
				return data, ok, err
			} else {
				return data, ok, nil
			}
		}
	} else {
		if ok, err := session.registrars(connector); err != nil {
			return []byte("activated connection"), ok, err
		} else {
			return []byte("activated connection"), ok, nil
		}
	}

}

func (session *WMPCISessions) Config(ctx context.Context, config map[string]wmpci.WMPCustom) (bool, error) {
	if ok, err := session.isActive(); !ok {
		return false, err
	}

	// TODO update config , will sync local yaml file updated

	if err := session.Registration.updateCustom(config); err != nil {
		return false, err
	} else {
		if _, err := session.Registration.Sync(session.workHomePath, true); err != nil {
			return false, err
		}
	}

	return session.wmp_connector.WMPConfig(ctx, session.Registration.RegistWMPCustom)
}

func (session *WMPCISessions) Health(ctx context.Context, txt ...string) (bool, error) {
	if ok, err := session.isActive(); !ok {
		return false, err
	}

	return session.wmp_connector.WMPHealth(ctx, txt...)
}

func (session *WMPCISessions) Service(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error) {
	if ok, err := session.isActive(); !ok {
		return nil, err
	}

	return session.wmp_connector.WMPCallHandle(ctx, request)
}

func (session *WMPCISessions) DebugService(ctx context.Context, request wmpci.WMPRequest) (wmpci.WMPResponse, error) {
	return nil, nil
}

func (session *WMPCISessions) isActive() (bool, error) {
	if !session.wmp_connector_status || session.wmp_connector == nil {
		return false, fmt.Errorf("the current wmpci connector is inactive")
	} else {
		return true, nil
	}
}

func (session *WMPCISessions) registrars(connector WMPConnector) (bool, error) {
	session.wmp_connector = connector
	if data, err := session.wmp_connector.WMPRegistrars(); err != nil {
		session.wmp_connector_status = false
		return false, err
	} else {
		if data != nil {
			if err := json.Unmarshal(data, &session.Registration); err != nil {
				return session.wmp_connector_status, err
			}
		}
		session.wmp_connector_status = true
		if cache, err := session.Registration.Sync(session.workHomePath, false); err != nil {
			return false, err
		} else if cache {
			return session.wmp_connector.WMPConfig(session.ctx, session.Registration.RegistWMPCustom)
		} else {
			return true, nil
		}
	}
}

func (session *WMPCISessions) WMPRequest() map[string][]interface{} {
	return session.Registration.GetFullRequest()
}
