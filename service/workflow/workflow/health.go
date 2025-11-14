package workflow

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/data/meta"
)

type WorkflowRuntimeHealth struct {
	Assets          meta.AssetMetaData                  `json:"assets" bson:"assets"`
	SchedulerError  []string                            `json:"scheduler_error,omitempty" bson:"scheduler_error"`
	SchedulerStatus map[string]interface{}              `json:"scheduler_status,omitempty" bson:"scheduler_status"`
	WMPResponseData map[string]map[string][]interface{} `json:"wmp_response_data,omitempty" bson:"-"`
	WMPDMSData      map[string][]interface{}            `json:"wmp_dms_data,omitempty" bson:"-"`

	WorkflowStatus    string `json:"workflow_status,omitempty" bson:"workflow_status"`
	WorkflowCreatedAt string `json:"workflow_created_at,omitempty" bson:"workflow_created_at"`
	WorkflowStartAt   string `json:"workflow_start_at,omitempty" bson:"workflow_start_at"`
	WorkflowCloseAt   string `json:"workflow_close_at,omitempty" bson:"workflow_close_at"`
}

func (f *Workflow) RuntimeHealth() WorkflowRuntimeHealth {

	var health WorkflowRuntimeHealth

	health.WorkflowStatus = WorkflowStatus[f.Status]
	health.WorkflowStartAt = f.startAt.Format("2006-01-02 15:04:05")
	health.WorkflowCreatedAt = f.createdAt.Format("2006-01-02 15:04:05")
	health.WorkflowCloseAt = f.closeAt.Format("2006-01-02 15:04:05")
	health.InitialSchedulerAndFlowBufferHealth(func() map[string]json.RawMessage {

		var health map[string]json.RawMessage = make(map[string]json.RawMessage)
		var err error

		if health["scheudler"], err = f.schedulerGridMaps.Health(); err != nil {
			logrus.Errorf("[scheduler] scheduler health json marshal error: %v", err)
		}

		for k, v := range f.workflowDataBuffers.Health() {
			health[k] = v
		}

		return health
	}())

	return health
}

func (w *WorkflowRuntimeHealth) InitialSchedulerAndFlowBufferHealth(health map[string]json.RawMessage) *WorkflowRuntimeHealth {
	for k, v := range health {
		switch k {
		case "asset":
			if err := json.Unmarshal(v, &w.Assets); err != nil {
				logrus.Errorf("[workflow] %s health data unmarshal error: %v", k, err)
			}
		case "wmp_response_buffer":
			if err := json.Unmarshal(v, &w.WMPResponseData); err != nil {
				logrus.Errorf("[workflow] %s health data unmarshal error: %v", k, err)
			}
		case "wmp_dms_bind":
			if err := json.Unmarshal(v, &w.WMPDMSData); err != nil {
				logrus.Errorf("[workflow] %s health data unmarshal error: %v", k, err)
			}
		case "wmp_exception":
			if err := json.Unmarshal(v, &w.SchedulerError); err != nil {
				logrus.Errorf("[workflow] %s health data unmarshal error: %v", k, err)
			}
		case "scheudler":
			if err := json.Unmarshal(v, &w.SchedulerStatus); err != nil {
				logrus.Errorf("[workflow] %s health data unmarshal error: %v", k, err)
			}
		}

	}
	return w
}
