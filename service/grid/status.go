package grid

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	Created uint = iota
	Started
	Inputed
	Execute
	Outputs
	Exception
	Exit
	Abort
)

var StatusCode = map[uint]string{
	Created:   "created",
	Started:   "started",
	Inputed:   "input",
	Execute:   "execute",
	Outputs:   "output",
	Exception: "exception",
	Exit:      "exit",
	Abort:     "abort",
}

type StatusAttributes struct {
	StatusCodes    string `json:"status_codes"`
	StatusMessages string `json:"status_messages"`
	StatusAt       string `json:"status_at"`
}

type TaskStatusMachine struct {
	Time          int64              `json:"time"`
	CurrentStatus uint               `json:"current_status"`
	HistoryStatus []StatusAttributes `json:"history_status"`
}

func NewTaskStatusMachine() *TaskStatusMachine {
	tsrm := &TaskStatusMachine{Time: time.Now().Unix(), HistoryStatus: []StatusAttributes{}}
	tsrm.UpdateTaskStatus(Created, "created success")
	return tsrm
}

func (t *TaskStatusMachine) JSON() (json.RawMessage, error) {
	return json.Marshal(&t)
}

func (t *TaskStatusMachine) UpdateTaskStatus(types uint, msg string) {
	t.CurrentStatus = types
	t.Time = time.Now().Unix()
	logrus.Debugf("[Grid] %d - %s : %s", t.Time, StatusCode[types], msg)
	t.HistoryStatus = append(t.HistoryStatus, StatusAttributes{
		StatusCodes:    StatusCode[types],
		StatusMessages: msg,
		StatusAt:       time.Now().Format("2006-06-05 12:24:34"),
	})
}
