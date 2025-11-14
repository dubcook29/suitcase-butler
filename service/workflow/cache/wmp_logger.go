package cache

import (
	"sync"
	"time"
)

type WMPCallLogManager struct {
	mu               sync.RWMutex
	SchedulerId      string
	SchedulerName    string
	StartTimeAt      int64
	DoneTimeAt       int64
	ExceptionStatus  bool
	ExceptionMessage string
	WMPCallHistory   []WMPCallLogger
}

func NewWMPCallLogManager(sid, sname string) *WMPCallLogManager {
	return &WMPCallLogManager{
		SchedulerId:   sid,
		SchedulerName: sname,
		StartTimeAt:   time.Now().Unix(),
	}
}

func (wcm *WMPCallLogManager) Logger(log WMPCallLogger) {
	wcm.mu.Lock()
	defer wcm.mu.Unlock()
	wcm.WMPCallHistory = append(wcm.WMPCallHistory, log)
}

func (wcm *WMPCallLogManager) Exit() {
	wcm.mu.Lock()
	defer wcm.mu.Unlock()
	wcm.DoneTimeAt = time.Now().Unix()
}

func (wcm *WMPCallLogManager) Exception(msg string) {
	wcm.mu.Lock()
	defer wcm.mu.Unlock()
	wcm.ExceptionStatus = true
	wcm.ExceptionMessage = msg
}

type WMPCallLogger struct {
	WMPID               string
	WMPName             string
	WMPCallAt           int64
	WMPDoneAt           int64
	WMPCurrentCustom    any
	WMPCurrentLimit     any
	WMPCurrentInData    any
	WMPCurrentOutData   any
	WMPCurrentException any
}

func NewWMPCallLogger(wid, wname string) WMPCallLogger {
	return WMPCallLogger{
		WMPID:     wid,
		WMPName:   wname,
		WMPCallAt: time.Now().Unix(),
	}
}
