package model

import (
	"sync"

	"errors"
)

var DataModelManager *WMPDataModelInterfaceManager = NewWMPDataModelInterfaceManager()

type WMPDataModelInterfaceManager struct {
	mu         sync.RWMutex
	tableMap   map[string]WMPDataModelInterface
	privateKey []string
}

func NewWMPDataModelInterfaceManager() *WMPDataModelInterfaceManager {

	return &WMPDataModelInterfaceManager{}
}

func (wdmm *WMPDataModelInterfaceManager) ShowAllWMPDataModelInterface() map[string]WMPDataModelInterface {
	wdmm.mu.RLock()
	defer wdmm.mu.RUnlock()
	// TODO need deep copy
	return wdmm.tableMap
}

func (wdmm *WMPDataModelInterfaceManager) RegisterWMPDataModelInterface(wdm WMPDataModelInterface) error {
	wdmm.mu.Lock()
	defer wdmm.mu.Unlock()
	if _, ok := wdmm.tableMap[wdm.Key()]; ok {
		return errors.New("the same predefined reserved WMPDataModel index is already occupied")
	} else {
		wdmm.tableMap[wdm.Key()] = wdm
		return nil
	}
}

func (wdmm *WMPDataModelInterfaceManager) AddPrivateKey(key string) *WMPDataModelInterfaceManager {
	wdmm.mu.Lock()
	defer wdmm.mu.Unlock()

	wdmm.privateKey = removeDuplicates(append(wdmm.privateKey, key))

	return wdmm
}

func (wdmm *WMPDataModelInterfaceManager) CheckPrivateKey(key string) bool {
	wdmm.mu.RLock()
	defer wdmm.mu.RUnlock()

	for _, v := range wdmm.privateKey {
		if key == v {
			return true
		}
	}

	return false
}

func (wdmm *WMPDataModelInterfaceManager) AddPrivateKeys(keys []string) *WMPDataModelInterfaceManager {
	wdmm.mu.Lock()
	defer wdmm.mu.Unlock()

	wdmm.privateKey = removeDuplicates(append(wdmm.privateKey, keys...))

	return wdmm
}

func (wdmm *WMPDataModelInterfaceManager) GetAllPrivateKey() []string {
	wdmm.mu.RLock()
	defer wdmm.mu.RUnlock()

	var list []string = []string{
		"group_id", "GroupId", "groupid",
		"asset_id", "AssetId", "assetid",
		"org", "org_name",
		"asn", "as", "as_number",
		"ip", "ips", "address", "ip_address", "IpAddress", "ipaddress",
		"domain", "domain_name", "DomainName", "domainname",
		"input", "other", "cloud", "cdn",
		"scheduler_id", "task_id",
	}
	for _, v := range wdmm.tableMap {
		list = append(list, v.Key())
	}
	return list
}

func (wdmm *WMPDataModelInterfaceManager) GetWMPDataModelInterface(key string) WMPDataModelInterface {

	wdmm.mu.RLock()
	defer wdmm.mu.RUnlock()

	if v, ok := wdmm.tableMap[key]; ok {
		return v
	} else {
		return DefaultWMPDataModel{}
	}

}

func (wdmm *WMPDataModelInterfaceManager) CheckWMPDataModelInterfaceKey(key string) bool {
	wdmm.mu.RLock()
	defer wdmm.mu.RUnlock()

	if _, ok := wdmm.tableMap[key]; ok {
		return true
	} else {
		return false
	}
}
