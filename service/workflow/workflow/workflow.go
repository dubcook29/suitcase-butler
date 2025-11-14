package workflow

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/suitcase/butler/data/meta"
	"github.com/suitcase/butler/grid"
	workflowcache "github.com/suitcase/butler/workflow/cache"
)

const (
	Created       = iota // = 0
	Waiting              // = 1
	Running              // = 2
	Pause                // = 3
	Stop                 // = 4
	Done                 // = 5
	Exit                 // = 6
	Exception            // = 7
	ExceptionQuit        // = 8
)

var WorkflowStatus map[uint]string = map[uint]string{
	Created:       "currently created",
	Waiting:       "currently waiting",
	Running:       "currently running",
	Pause:         "currently pause",
	Stop:          "workflow stop and exit",
	Done:          "workflow waiting close",
	Exit:          "workflow close and exit",
	Exception:     "workflow runtime exception and not is close and exit",
	ExceptionQuit: "workflow exception close and exit",
}

type Workflow struct {
	ctx       context.Context
	done      context.CancelFunc
	wg        sync.WaitGroup
	Status    uint
	createdAt time.Time
	startAt   time.Time
	closeAt   time.Time

	workflowDataBuffers *workflowcache.WorkflowDataBuffer
	schedulerGridMaps   *grid.Grid
}

func Workflower(ctx context.Context, grids *grid.Grid, asset meta.AssetMetaData, nextAsset chan string) *Workflow {
	var flow = new(Workflow)

	flow.ctx, flow.done = context.WithCancel(ctx)                                                           // with context
	flow.workflowDataBuffers = workflowcache.NewWorkflowDataBuffer(flow.ctx, asset, nextAsset)              // workflow data buffer created and initialize
	flow.schedulerGridMaps = grids.Initial(flow.ctx).InitialWorkflowDataProcessor(flow.workflowDataBuffers) // initialize grid handle and bind workflow data buffer

	return flow
}

func (f *Workflow) RefreshContext(ctx context.Context) {
	f.ctx, f.done = context.WithCancel(ctx)
	f.workflowDataBuffers.RefreshContext(f.ctx)
	f.schedulerGridMaps.RefreshContext(ctx)
}

// Starter this is the workflow startup entrance,
// waiting for the workflow execution to end and return to the final execution state
func (f *Workflow) Starter(ctx context.Context) bool {
	// prevent repeated startups if running
	if f.Status == Running {
		return false
	} else {
		f.Status = Running // update workflow status
	}

	f.RefreshContext(ctx) // refresh scheduler and buffer all context

	f.startAt = time.Now() // write now time

	// Add a pre-exit processing process, mainly processing the update of close time
	f.wg.Add(1)
	defer func() { defer f.wg.Done() }()

	logrus.Infof("[workflow] prepare to start workflow scheduling")

	f.schedulerGridMaps.Running()

	f.closeAt = time.Now()

	if f.Status == Running {
		f.Status = Done
		f.Status = f.CheckExceptionStatus()
	}

	logrus.Infof("[workflow] Goodbye, this Workflow is ready to be shut down, its status is [%s]", WorkflowStatus[f.Status])

	return true
}

func (f *Workflow) CheckExceptionStatus() uint {
	if f.Status != Done {
		return f.Status
	}

	// todo check exception in workflow cache

	return Exit
}

func (f *Workflow) Wait() {
	f.wg.Wait()
}

// Close This is a function that closes the workflow and blocks until the current workflow is closed.
func (f *Workflow) Close() {

	// is the workflow running
	if f.Status != Running {
		return
	} else {
		f.Status = Stop
	}

	// send context.CancelFunc() close context
	f.done()

	// Wait for all (scheduler)tasks to exit
	f.schedulerGridMaps.Wait()

	// Wait for workflow Starter to exit, or delete this workflow from the workflow manager
	f.Wait()
}
