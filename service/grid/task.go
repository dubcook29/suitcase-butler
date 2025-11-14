package grid

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/wmpci/manager"
)

type Task struct {
	TaskId        string             `json:"task_id,omitempty" yaml:"task_id,omitempty"`
	TaskName      string             `json:"task_name,omitempty" yaml:"task_name,omitempty"`
	WMPId         string             `json:"wmp_id" yaml:"wmp_id"`
	SessionId     string             `json:"session_id" yaml:"session_id"`
	PrevTaskId    []string           `json:"prev_task_id,omitempty" yaml:"prev_task,omitempty"`
	NextTask      []Task             `json:"next_task,omitempty" yaml:"next_task,omitempty"`
	Concurrency   bool               `json:"concurrency" yaml:"concurrency"`
	StatusMachine *TaskStatusMachine `json:"status_machine,omitempty" yaml:"-"`
}

func (t Task) CallWMPCIService(ctx context.Context, wmpRequest map[string][]interface{}, sessions *manager.WMPSessionManager) (map[string][]interface{}, error) {
	// session_id is generated randomly and has no hard connection with wmp, so it is impossible to call the session interface through WMP-id.
	logrus.Infof("=== current call %s wmpci service\n wmp request: \n%+v ", t.WMPId, wmpRequest)
	if session, err := sessions.SessionByWMPID(t.WMPId); err != nil {
		return nil, err
	} else {
		if response, err := session.Service(context.TODO(), wmpRequest); err != nil {
			logrus.Errorf("## ERROR ## current call to [%s]: request: %+v\n response: (failed), %v", session.Registration.RegistWMPBasic.Name, wmpRequest, err)
			return nil, err
		} else {
			logrus.Debugf("## DEBUG ## current call to [%s]: request: %+v\n response:%+v", session.Registration.RegistWMPBasic.Name, wmpRequest, response)
			return response, nil
		}
	}
}

func (t Task) CallWMPCIRequest(sessions *manager.WMPSessionManager) (map[string][]interface{}, error) {
	logrus.Infof("=== current call %s wmpci request", t.WMPId)

	if session, err := sessions.SessionByWMPID(t.WMPId); err != nil {
		return nil, err
	} else {
		return session.Registration.GetFullRequest(), nil
	}
}

func (t Task) lookupTasks(v ...string) {
	l := len(t.NextTask) - 1
	for i, task := range t.NextTask {
		if i == l {
			fmt.Printf("%s└── [%d] %s | %s | %s\n", strings.Join(v, ""), i, task.TaskId, task.TaskName, task.WMPId)
		} else {
			fmt.Printf("%s├── [%d] %s | %s | %s\n", strings.Join(v, ""), i, task.TaskId, task.TaskName, task.WMPId)
		}
		task.lookupTasks(strings.Join(v, ""), "│   ")
	}
}

func (t Task) running() chan Task {

	var taskChannel = make(chan Task, 10)

	go func() {
		for _, task := range t.NextTask {
			taskChannel <- task
		}
		close(taskChannel)
	}()

	return taskChannel
}
