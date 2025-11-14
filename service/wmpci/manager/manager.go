package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/wmpci"
	"github.com/suitcase/butler/wmpci/connector"
)

type WMPSessionManager struct {
	mu                sync.RWMutex
	ctx               context.Context
	wmpci_work_path   string
	sessions          map[string]*connector.WMPCISessions
	connectorSupports *connector.WMPConnectorSupports
}

func InitialSessionManager(ctx context.Context, workpath string) *WMPSessionManager {
	// initial wmp session manager
	defer logrus.Infof("[wmpci session] workflow sessions manager created successfully!")
	return &WMPSessionManager{
		ctx:               ctx,
		wmpci_work_path:   filepath.Join(workpath, ".wmpci/"),
		sessions:          make(map[string]*connector.WMPCISessions),
		connectorSupports: connector.DefaultWMPConnectorSupports,
	}
}

// shows all supported connectors for connecting to WMP
func (m *WMPSessionManager) ConnectorSupportLists() ([]string, error) {
	if lists := m.connectorSupports.GetConnectorGeneratorList(); lists != nil {
		return lists, nil
	} else {
		return nil, fmt.Errorf("[wmpci session] connector support list is null")
	}
}

// get connector configuration parameters
func (m *WMPSessionManager) SelectConnectorConnectionCumstom(name string) (map[string]wmpci.WMPCustom, error) {
	if custom, err := m.connectorSupports.GetConnectorConnectCustomConfig(name); err != nil {
		return nil, err
	} else {
		return custom, nil
	}
}

// write connector configuration parameters and created connector session, return connection session id
func (m *WMPSessionManager) ConnectorConnectionSession(name string, custom map[string]wmpci.WMPCustom) (string, error) {
	if connect, err := m.connectorSupports.ConnectorGenerator(name); err != nil {
		return "", err
	} else {
		session := connector.NewWMPCISession(m.wmpci_work_path)
		if data, ok, err := session.Connect(m.ctx, connect, custom); err != nil {
			return string(data), err
		} else if ok {
			return string(data), m.addSession(session)
		} else {
			return string(data), nil
		}
	}
}

// create connector session -> connector registration request -> connector session persistence
func (m *WMPSessionManager) addSession(session *connector.WMPCISessions) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[session.SessionId]; ok {
		return fmt.Errorf("[wmpci session] the current [%s] session already exists", session.SessionId)
	} else {
		m.sessions[session.SessionId] = session
	}

	return nil
}

func (m *WMPSessionManager) CloseSession(session_id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[session_id]; ok {
		delete(m.sessions, session_id)
		return nil
	} else {
		return fmt.Errorf("[wmpci session] the current [%s] session does not exists", session_id)
	}

}

func (m *WMPSessionManager) Session(session_id string) (*connector.WMPCISessions, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if session, ok := m.sessions[session_id]; ok {
		return session, nil
	} else {
		return nil, fmt.Errorf("[wmpci session] the current [%s] session does not exists", session_id)
	}
}

func (m *WMPSessionManager) SessionByWMPID(wmp_id string) (*connector.WMPCISessions, error) {
	if applications, err := m.SessionApplicationBasic(); err != nil {
		return nil, err
	} else {
		for session_id, basic := range applications {
			if basic.Id == wmp_id {
				return m.Session(session_id)
			}
		}
	}
	return nil, fmt.Errorf("[wmpci session] the current all session cannot find [%s]", wmp_id)
}

func (m *WMPSessionManager) SessionRegistration(session_id string) (*connector.WMPRegistration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if session, ok := m.sessions[session_id]; ok {
		return session.Registration, nil
	} else {
		return nil, fmt.Errorf("[wmpci session] the current [%s] session does not exists", session_id)
	}
}

func (m *WMPSessionManager) SessionApplicationBasic(session_id ...string) (map[string]wmpci.WMPBasic, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result = make(map[string]wmpci.WMPBasic)
	if session_id == nil {
		for sid, session := range m.sessions {
			result[sid] = session.Registration.RegistWMPBasic
		}
	} else {
		for _, sid := range session_id {
			if session, ok := m.sessions[sid]; ok {
				result[sid] = session.Registration.RegistWMPBasic
			} else {
				// TODO
			}
		}
	}

	return result, nil
}

func (m *WMPSessionManager) SessionMap() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return json.Marshal(m.sessions)
}

func (m *WMPSessionManager) Sessions() ([]connector.WMPCISessions, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var data []connector.WMPCISessions

	for _, session := range m.sessions {
		data = append(data, *session)
	}

	return data, nil
}

func (m *WMPSessionManager) ConnectorWMPRequest(session_id string) (wmpci.WMPRequest, error) {
	if session, err := m.Session(session_id); err != nil {
		return nil, err
	} else {
		return session.Registration.RegistWMPRequest, nil
	}
}

func (m *WMPSessionManager) ConnectorWMPResponse(session_id string) (wmpci.WMPResponse, error) {
	if session, err := m.Session(session_id); err != nil {
		return nil, err
	} else {
		return session.Registration.RegistWMPResponse, nil
	}
}

func (m *WMPSessionManager) ConnectorWMPCustom(session_id string) (map[string]wmpci.WMPCustom, error) {
	if session, err := m.Session(session_id); err != nil {
		return nil, err
	} else {
		return session.Registration.RegistWMPCustom, nil
	}
}

func (m *WMPSessionManager) ConnectorWMPCustomUpdated(session_id string, custom map[string]wmpci.WMPCustom) (bool, error) {
	if session, err := m.Session(session_id); err != nil {
		return false, err
	} else {
		return session.Config(m.ctx, custom)
	}
}
