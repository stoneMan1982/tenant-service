package config

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

// Manager is the config manager which is designed for hot reloading
type Manager struct {
	cfg     *AppConfig
	path    string
	cfgMu   sync.RWMutex
	watcher *fsnotify.Watcher

	// config listeners
	listeners  []ConfigListener
	listenerMu sync.RWMutex
}

// ConfigEvent config changes event
type ConfigEvent struct {
	OldConfig *AppConfig // old config
	NewConfig *AppConfig // new config
}

// ConfigListener config changes listener interface
type ConfigListener interface {
	OnConfigChange(event ConfigEvent) // config changes callback
}

// ConfigCallback config changes callback function type
type ConfigCallback func(event ConfigEvent)

// callbackListener implements ConfigListener interface
type callbackListener struct {
	callback ConfigCallback
}

func (c *callbackListener) OnConfigChange(event ConfigEvent) {
	c.callback(event)
}

// NewManager creates a new config manager
// path: config file path
func NewManager(path string) (*Manager, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	m := &Manager{
		path:    absPath,
		watcher: watcher,
	}

	if err := m.load(); err != nil {
		watcher.Close()
		return nil, err
	}

	go m.watch()

	return m, nil
}

// Load loads config file and notifies listeners
func (m *Manager) load() error {
	data, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}

	var newCfg AppConfig
	if err := yaml.Unmarshal(data, &newCfg); err != nil {
		return err
	}

	// save old config
	m.cfgMu.RLock()
	oldCfg := m.cfg
	m.cfgMu.RUnlock()

	// update new config
	m.cfgMu.Lock()
	m.cfg = &newCfg
	m.cfgMu.Unlock()

	// trigger config change event
	if oldCfg != nil { // initial load does not trigger event (oldCfg is nil)
		m.notifyListeners(ConfigEvent{
			OldConfig: oldCfg,
			NewConfig: &newCfg,
		})
	}

	return nil
}

// watch monitors config file changes and reloads config
// it will block until the manager is closed
func (m *Manager) watch() {
	if err := m.watcher.Add(m.path); err != nil {
		return
	}
	defer m.watcher.Close()

	for {
		select {
		case event, ok := <-m.watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {

				time.Sleep(100 * time.Millisecond)
				if err := m.load(); err != nil {
					continue
				}
			}

		case _, ok := <-m.watcher.Errors:
			if !ok {
				return
			}
		}
	}
}

// notify all listeners
func (m *Manager) notifyListeners(event ConfigEvent) {
	m.listenerMu.RLock()
	defer m.listenerMu.RUnlock()

	for _, listener := range m.listeners {
		go listener.OnConfigChange(event)
	}
}

// RegisterListener
func (m *Manager) RegisterListener(listener ConfigListener) {
	m.listenerMu.Lock()
	defer m.listenerMu.Unlock()
	m.listeners = append(m.listeners, listener)
}

// RegisterCallback
func (m *Manager) RegisterCallback(callback ConfigCallback) {
	m.RegisterListener(&callbackListener{callback: callback})
}

// GetConfig
func (m *Manager) GetConfig() *AppConfig {
	m.cfgMu.RLock()
	defer m.cfgMu.RUnlock()
	return m.cfg
}

// Close
func (m *Manager) Close() error {
	return m.watcher.Close()
}
