package connector

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/wmpci"
	"gopkg.in/yaml.v2"
)

type WMPRegistration struct {
	mu                sync.RWMutex
	RegistWMPBasic    wmpci.WMPBasic             `json:"wmp_basic" yaml:"wmp_basic" bson:"wmp_basic"`
	RegistWMPRequest  map[string][]interface{}   `json:"wmp_request" yaml:"wmp_request" bson:"wmp_request"`
	RegistWMPResponse map[string][]interface{}   `json:"wmp_response" yaml:"wmp_response" bson:"wmp_response"`
	RegistWMPCustom   map[string]wmpci.WMPCustom `json:"wmp_custom" yaml:"wmp_custom" bson:"wmp_custom"`
}

/*
  - TODO
    json: unsupported type: map[interface {}]interface {}
*/

func (regist *WMPRegistration) PrintAllWmpregistration() {
	regist.mu.RLock()
	defer regist.mu.RUnlock()

	logrus.Infof("    [WMPID] > %s", regist.RegistWMPBasic.Id)
	logrus.Infof("    [NAMEs] > %s", regist.RegistWMPBasic.Name)
	logrus.Infof("  [Version] > %s", regist.RegistWMPBasic.Version)
	logrus.Infof("   [Custom] > %v", regist.RegistWMPCustom)
	logrus.Infof("  [Request] > %v", regist.RegistWMPRequest)
	logrus.Infof(" [Response] > %v", regist.RegistWMPResponse)
	fmt.Println("")
}

func (regist *WMPRegistration) JSON() ([]byte, error) {
	regist.mu.RLock()
	defer regist.mu.RUnlock()

	return json.Marshal(regist)
}

func (regist *WMPRegistration) LoadRegistrationFromJSON(in []byte) error {
	regist.mu.Lock()
	defer regist.mu.Unlock()

	return json.Unmarshal(in, regist)
}

func (regist *WMPRegistration) Serialization() ([]byte, error) {
	regist.mu.RLock()
	defer regist.mu.RUnlock()

	return yaml.Marshal(regist)
}

func (regist *WMPRegistration) Deserialization(in []byte) error {
	regist.mu.Lock()
	defer regist.mu.Unlock()

	return yaml.Unmarshal(in, regist)
}

func (regist *WMPRegistration) GetFullRequest() map[string][]interface{} {
	regist.mu.RLock()
	defer regist.mu.RUnlock()

	return deepCopyWMPMap(regist.RegistWMPRequest)
}

func deepCopyWMPMap(original map[string][]interface{}) map[string][]interface{} {
	CopyMaps := make(map[string][]interface{})

	for key, value := range original {
		newSlice := make([]interface{}, len(value))
		copy(newSlice, value)
		CopyMaps[key] = newSlice
	}

	return CopyMaps
}

func (regist *WMPRegistration) updateCustom(custom map[string]wmpci.WMPCustom) error {
	for k, v := range custom {
		if vv, ok := regist.RegistWMPCustom[k]; ok {
			if vv.Value != v.Value {
				regist.RegistWMPCustom[k] = v
			}
		} else {
			regist.RegistWMPCustom[k] = v
		}
	}
	return nil
}

func (regist *WMPRegistration) Sync(path string, cover bool) (bool, error) {
	if !pathExists(path) {
		return false, fmt.Errorf("dir/file <%s> does not exist", path)
	}

	filepath := filepath.Join(path, regist.RegistWMPBasic.Id+".yaml")

	if ok, err := fileExists(filepath); err != nil {
		return false, err
	} else if ok && !cover {
		// if file is exist and cover source yaml file
		// read file and unmarshal to regist
		if data, err := os.ReadFile(filepath); err != nil {
			return false, err
		} else {
			return true, regist.Deserialization(data)
		}
	}

	// write
	if data, err := regist.Serialization(); err != nil {
		return false, err
	} else {
		return false, os.WriteFile(filepath, data, 0644)
	}
}

func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// exists check if the directory exists
func pathExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// ensureDir check if the directory exists, and if it does not exist, create it
func ensureDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
