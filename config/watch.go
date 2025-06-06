package config

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"

	"github.com/harrisonju123/mcp-agent-poc/internal/router"
)

type Watcher struct {
	path   string
	router *router.Router
}

func NewWatcher(path string, r *router.Router) *Watcher {
	return &Watcher{path: filepath.Clean(path), router: r}
}

const debounce = 250 * time.Millisecond

func (w *Watcher) Run(ctx context.Context) error {
	if err := w.loadAndSwap(); err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer watcher.Close()
	if err := watcher.Add(filepath.Dir(w.path)); err != nil {
		return err
	}

	timer := time.NewTimer(0) // already fired
	<-timer.C                 //drain

	for {
		select {
		case <-ctx.Done():
			return nil

		case ev := <-watcher.Events:
			if ev.Name == w.path && (ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename)) != 0 {
				timer.Reset(debounce)
			}
		case <-timer.C:
			_ = w.loadAndSwap() // ignore error keep the last good on fail
		case err := <-watcher.Errors:
			return err
		}
	}
}

func (w *Watcher) loadAndSwap() error {
	raw, err := os.ReadFile(w.path)
	if err != nil {
		return err
	}
	var cfg []ToolConf
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return err
	}
	tools := make([]router.Tool, 0, len(cfg))
	for _, c := range cfg {
		tools = append(tools, newToolFromConf(c))
	}

	w.router.Replace(tools)
	return nil
}

func newToolFromConf(tc ToolConf) router.Tool {
	return router.Tool{
		Name:        tc.Name,
		Description: tc.Description,
		Handler: func(ctx context.Context, in []byte) ([]byte, error) {
			//TODO: execute script/http/grpc
			return in, nil
		},
	}
}
