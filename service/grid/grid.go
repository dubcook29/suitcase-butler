package grid

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	wmpcimanager "github.com/suitcase/butler/wmpci/manager"
	workflowcache "github.com/suitcase/butler/workflow/cache"
	"gopkg.in/yaml.v2"
)

type Grid struct {
	mu                    sync.RWMutex
	wg                    sync.WaitGroup
	ctx                   context.Context
	done                  context.CancelFunc
	workflowDataProcessor *workflowcache.WorkflowDataBuffer
	sessions              *wmpcimanager.WMPSessionManager
	taskStatusMachine     map[string]*TaskStatusMachine

	GridId          string         `json:"grid_id,omitempty" yaml:"grid_id,omitempty"`
	GridName        string         `json:"grid_name,omitempty" yaml:"grid_name,omitempty"`
	GridDescriptive string         `json:"grid_descriptive,omitempty" yaml:"grid_descriptive,omitempty"`
	GridTasks       []Task         `json:"grid_tasks,omitempty" yaml:"grid_tasks,omitempty"`
	GridDataBind    []GridDataBind `json:"grid_data_bind,omitempty" yaml:"grid_data_bind,omitempty"`
}

func (g *Grid) SelfCheck() error {
	return nil
}

func (g *Grid) Initial(ctx context.Context) *Grid {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.ctx, g.done = context.WithCancel(ctx)

	g.taskStatusMachine = make(map[string]*TaskStatusMachine)

	return g
}

func (g *Grid) RefreshContext(ctx context.Context) {
	g.ctx, g.done = context.WithCancel(ctx)
}

func (g *Grid) Health() (json.RawMessage, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	return json.Marshal(g.taskStatusMachine)
}

func (g *Grid) Close() {
	g.done()
	g.wg.Wait()
}

func (g *Grid) Wait() {
	g.wg.Wait()
}

// InitialWorkflowDataProcessor
// Initialize the workflow data processor to process all input and output data for grid running,
// but not grid data bind relationships
func (g *Grid) InitialWorkflowDataProcessor(proc *workflowcache.WorkflowDataBuffer) *Grid {
	g.mu.Lock()
	defer g.mu.Unlock()
	defer logrus.Debugf("[Grid] the grid is initialized successfully and the workflow data buffer is successfully bound")

	g.workflowDataProcessor = proc
	return g
}

func (g *Grid) InitialWMPSessionManager(sessions *wmpcimanager.WMPSessionManager) *Grid {
	g.mu.Lock()
	defer g.mu.Unlock()

	if false {
		// todo self-check workflow Data Processor function
		// return
	}

	g.sessions = sessions
	return g
}

func (g *Grid) ExpandGridTasks() {
	fmt.Printf(".\n")
	l := len(g.GridTasks) - 1
	for i, task := range g.GridTasks {
		if i == l {
			fmt.Printf("└── [%d] %s | %s | %s\n", i, task.TaskId, task.TaskName, task.WMPId)
		} else {
			fmt.Printf("├── [%d] %s | %s | %s\n", i, task.TaskId, task.TaskName, task.WMPId)
		}
		task.lookupTasks("│   ")
	}
}

func (g *Grid) Running() {
	var task_channel = make(chan Task, 10)
	go g.run(task_channel)
	for _, task := range g.GridTasks {
		task_channel <- task
	}
	close(task_channel)
	time.Sleep(2 * time.Second)
	g.wg.Wait()
}

func (g *Grid) run(tasks chan Task) {
	g.wg.Add(1)
	defer g.wg.Done()
	for task := range tasks {
		task.StatusMachine = g.status(task.TaskId)
		select {
		case <-g.ctx.Done():
			task.StatusMachine.UpdateTaskStatus(Abort, "")
		default:
			task.StatusMachine.UpdateTaskStatus(Started, fmt.Sprintf("[running] %s | %s | %s", task.TaskId, task.TaskName, task.WMPId))

			if request, err := task.CallWMPCIRequest(g.sessions); err != nil {
				logrus.Errorf("wmpci session connection request failed: %v", err)
				task.StatusMachine.UpdateTaskStatus(Exception, fmt.Sprintf("wmpci session connection request failed: %v", err))
			} else {

				if request, err = g.pull(task, request); err != nil {
					// request pull failed and error
					logrus.Errorf("request pull from wmpci session failed: %v", err)
					task.StatusMachine.UpdateTaskStatus(Exception, fmt.Sprintf("request pull from wmpci session failed: %v", err))
				} else {
					// request pull successfully
					logrus.Infof("request pull from wmpci session successfully")
					task.StatusMachine.UpdateTaskStatus(Inputed, "request pull from wmpci session successfully")
				}

				if response, err := task.CallWMPCIService(g.ctx, request, g.sessions); err != nil {
					logrus.Errorf("wmpci session connection service failed: %v", err)
					task.StatusMachine.UpdateTaskStatus(Exception, fmt.Sprintf("wmpci session connection service failed: %v", err))
				} else {
					task.StatusMachine.UpdateTaskStatus(Execute, "wmp application service called successfully")

					if err := g.push(task, response); err != nil {
						// response push failed and error
						logrus.Errorf("response push to workflow failed: %v", err)
						task.StatusMachine.UpdateTaskStatus(Exception, fmt.Sprintf("response push to workflow failed: %v", err))
					} else {
						// response push successfully
						logrus.Infof("response push to workflow successfully")
						task.StatusMachine.UpdateTaskStatus(Outputs, "response push to workflow successfully")

					}
				}

			}

			task.StatusMachine.UpdateTaskStatus(Exit, "")

			g.run(task.running())
		}
	}
}

func (g *Grid) status(task_id string) *TaskStatusMachine {
	g.mu.Lock()
	defer g.mu.Unlock()

	if status, ok := g.taskStatusMachine[task_id]; ok {
		return status
	} else {
		status = NewTaskStatusMachine()
		g.taskStatusMachine[task_id] = status
		return status
	}
}

func (g *Grid) pull(task Task, request map[string][]interface{}) (map[string][]interface{}, error) {
	if request, err := g.workflowDataProcessor.PullRequest(request, task.WMPId); err != nil {
		return nil, err
	} else {
		for index, value := range request {
			if data := g.bindReader(task.WMPId, index); len(data) > 0 {
				request[index] = append(value, data...)
			}
		}
		return request, nil
	}
}

func (g *Grid) push(task Task, response map[string][]interface{}) error {
	for index, value := range response {
		if err := g.bindWriter(task.WMPId, index, value); err != nil {
			return err
		}
	}

	return g.workflowDataProcessor.PushResponse(response, task.WMPId)
}

func (g *Grid) Serialization() ([]byte, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return yaml.Marshal(g)
}

func (g *Grid) Deserialization(in []byte) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	return yaml.Unmarshal(in, g)
}

func (g *Grid) Sync(path string, cover bool) (bool, error) {
	if !pathExists(path) {
		return false, fmt.Errorf("dir/file <%s> does not exist", path)
	}

	filepath := filepath.Join(path, g.GridId+".yaml")

	if ok, err := fileExists(filepath); err != nil {
		return false, err
	} else if ok && !cover {
		// if file is exist and cover source yaml file
		// read file and unmarshal to regist
		if data, err := os.ReadFile(filepath); err != nil {
			return false, err
		} else {
			return true, g.Deserialization(data)
		}
	}

	// write
	if data, err := g.Serialization(); err != nil {
		return false, err
	} else {
		return false, os.WriteFile(filepath, data, 0644)
	}
}
