package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/data/meta"
	"github.com/suitcase/butler/db"
	"github.com/suitcase/butler/wmpci/data/model"
)

type WorkflowDataBuffer struct {
	mu                    sync.RWMutex
	ctx                   context.Context
	asset                 meta.AssetMetaData
	wmpResponseDataBuffer map[string]map[string][]interface{}
	exceptionDataBuffer   []string
	throwToEntryChannel   chan string
}

func (w *WorkflowDataBuffer) RefreshContext(ctx context.Context) {
	w.ctx = ctx
}

// nolint:unused
// func (w *WorkflowDataBuffer) WMPRequestDataReader(req map[string][]interface{}, wmp_id string) { 	panic("Reserved Functions Don't Call") }

func (w *WorkflowDataBuffer) WMPRequestDefault() map[string][]interface{} {
	var wmp_request = make(map[string][]interface{})

	var default_lists = map[string][]interface{}{
		"asset_id": nil,
	}

	for k, v := range default_lists {
		if val := w.asset.Reader(k); val != nil {
			wmp_request[k] = val
		} else {
			wmp_request[k] = v
		}
	}

	return wmp_request
}

func (w *WorkflowDataBuffer) PullRequest(request map[string][]interface{}, wmp_id string) (map[string][]interface{}, error) {

	var wmp_request = w.WMPRequestDefault()

	for k := range request {
		if val := w.asset.Reader(k); val != nil {
			wmp_request[k] = val
		}
	}

	return wmp_request, nil
}

func (w *WorkflowDataBuffer) PushResponse(response map[string][]interface{}, wmp_id string) error {

	for k, v := range response {
		if w.asset.IsPrivateKeys(k) {
			w.asset.Writer(k, v)
		} else if k == "^" {
			for _, item := range v {
				if v, ok := item.(string); ok {
					if vl := strings.Split(v, ":"); len(vl) >= 2 {
						go w.throwToEntry(v)

					}
				}
			}
			// todo
		} else if k == "#" {
			// todo
		} else {
			defer w.asset.AddWMPData(k)
			w.wmpResponseDataBufferWriter(v, wmp_id, k)
		}
	}

	return nil
}

func (w *WorkflowDataBuffer) wmpResponseDataBufferWriter(in []interface{}, wmp_id string, k string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.wmpResponseDataBuffer[wmp_id]; !ok {
		w.wmpResponseDataBuffer[wmp_id] = make(map[string][]interface{})
	}

	w.storage(wmp_id, k, in)

	// TODO need to consider the problem of overlapping and crossing results. How to deal with it in the future?
	if kv, ok := w.wmpResponseDataBuffer[wmp_id][k]; ok {
		w.wmpResponseDataBuffer[wmp_id][k] = append(kv, in)
	} else {
		w.wmpResponseDataBuffer[wmp_id][k] = in
	}
}

func (w *WorkflowDataBuffer) storage(wmp_id string, k string, v []interface{}) {

	var data []model.WMPDataModelInterface
	for _, vv := range v {
		if wmp_data_model, ok := vv.(model.WMPDataModelInterface); ok {
			data = append(data, wmp_data_model)
		} else {
			data = append(data, model.DefaultWMPDataModel{
				MetaData: model.NewWMPDataModelBasicStructure(w.asset.AssetId, k),
				Value:    vv,
				Formats:  fmt.Sprintf("%T", vv),
				From:     wmp_id,
			})
		}
	}

	if count, err := model.WMPDataModelMongoInsertFunc(context.TODO(), db.GetCurrentMongoClient(), "wmpdata", k, data); err != nil {
		w.AddExcetpion(fmt.Sprintf("[workflow buffer] wmp response data storage is failed: %v", err))
	} else {
		logrus.Debugf("[workflow buffer] wmp response data storage is successed, total (%d) data.", count)
	}

}

// AddExcetpion add a new eexception log to the `exceptDataBuffer`,
// any exception non-thrower in the buffer needs to be recorded in the `exceptDataBuffer`
func (w *WorkflowDataBuffer) AddExcetpion(e string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	logrus.Error(e)
	w.exceptionDataBuffer = append(w.exceptionDataBuffer, e)
}

func (w *WorkflowDataBuffer) throwToEntry(iv string) {
	if w.throwToEntryChannel == nil {
		return
	}
	w.throwToEntryChannel <- iv
}

// NewWorkflowDataBuffer: create a new workflow data buffer, input workflowManager
func NewWorkflowDataBuffer(ctx context.Context, asset meta.AssetMetaData, backCall chan string) *WorkflowDataBuffer {
	defer logrus.Debugf("[workflow buffer] workflow (%s) data buffer created successfully", asset.AssetId)
	return &WorkflowDataBuffer{
		asset:                 asset,
		ctx:                   ctx,
		mu:                    sync.RWMutex{},
		wmpResponseDataBuffer: make(map[string]map[string][]interface{}),
		throwToEntryChannel:   backCall,
	}
}

func (w *WorkflowDataBuffer) SelfCheck() error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return nil
}

// Health: returns wmpData, wmpDMS, exceptData in the current buffer
func (w *WorkflowDataBuffer) Health() map[string]json.RawMessage {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var health map[string]json.RawMessage = make(map[string]json.RawMessage)
	var err error

	if health["asset"], err = json.Marshal(w.asset); err != nil {
		logrus.Errorf("[workflow buffer] buffer health json marshal error: %v", err)
	}

	if health["wmp_response_buffer"], err = json.Marshal(w.wmpResponseDataBuffer); err != nil {
		logrus.Errorf("[workflow buffer] buffer health json marshal error: %v", err)
	}

	if health["wmp_exception"], err = json.Marshal(w.exceptionDataBuffer); err != nil {
		logrus.Errorf("[workflow buffer] buffer health json marshal error: %v", err)
	}

	return health
}
