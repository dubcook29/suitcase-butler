package grid

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	wmpcimanager "github.com/suitcase/butler/wmpci/manager"
)

type GridManager struct {
	mu sync.RWMutex

	path      string
	filenames map[string]string
	sessions  *wmpcimanager.WMPSessionManager
}

func (g *GridManager) Initial(workpath string, sessions *wmpcimanager.WMPSessionManager) *GridManager {
	g.path = filepath.Join(workpath, ".grid")
	ensureDir(g.path)
	g.filenames = make(map[string]string)
	g.sessions = sessions
	return g
}

func (g *GridManager) loadingAllGridFiles() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	return filepath.Walk(g.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml")) {
			var key string = info.Name()
			key = strings.TrimSuffix(key, ".yaml")
			g.filenames[key] = path
		}
		return nil
	})
}

func (g *GridManager) NewGrid(grid_id string) (*Grid, error) {

	g.mu.RLock()
	defer g.mu.RUnlock()
	if _, ok := g.filenames[grid_id]; ok {
		var grid *Grid = new(Grid)
		grid.GridId = grid_id
		grid.sessions = g.sessions
		if ok, err := grid.Sync(g.path, false); err != nil {
			return nil, err
		} else if ok {
			return grid, nil
		}
	}

	return nil, fmt.Errorf("grid does not exist")
}

func (g *GridManager) AddGrid(grid *Grid) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.filenames[grid.GridId]; !ok {
		g.filenames[grid.GridId] = filepath.Join(g.path, grid.GridId+".yaml")
		if ok, err := grid.Sync(g.path, false); err != nil {
			return false, err
		} else {
			return ok, nil
		}
	} else {
		return false, fmt.Errorf("grid already exists")
	}

}

func (g *GridManager) ModGrid(grid *Grid) (bool, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, ok := g.filenames[grid.GridId]; ok {
		if _, err := grid.Sync(g.path, true); err != nil {
			return false, err
		} else {
			return true, nil
		}
	}

	return false, fmt.Errorf("grid does not exist")
}

func (g *GridManager) DelGrid(grid_id string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.filenames[grid_id]; ok {
		delete(g.filenames, grid_id)
		// delete file from local files
		return ok, nil
	}

	return false, fmt.Errorf("grid does not exist")
}

func (g *GridManager) GetAllOnlineGrid() ([]*Grid, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var data []*Grid

	for id := range g.filenames {
		if grider, err := g.NewGrid(id); err != nil {
			return nil, err
		} else {
			data = append(data, grider)
		}
	}

	return data, nil
}

func (g *GridManager) Sessions() *wmpcimanager.WMPSessionManager {
	return g.sessions
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
