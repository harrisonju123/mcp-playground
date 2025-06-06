package config

import (
	"context"
	tp "github.com/Shopify/toxiproxy/v2/client"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/harrisonju123/mcp-agent-poc/internal/router"
)

var toxCli = tp.NewClient("http://toxiproxy:8474")

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
		tool, err := newToolFromConf(c)
		if err != nil {
			return err
		}
		tools = append(tools, tool)
	}

	w.router.Replace(tools)
	return nil
}

func newToolFromConf(tc ToolConf) (router.Tool, error) {
	target := tc.Command
	if strings.HasPrefix(target, "toxiproxy://") {
		up := strings.TrimPrefix(target, "toxiproxy://") // echo1:8080
		// proxy name is deterministic -> idempotent
		pxName := "px_" + strings.ReplaceAll(up, ":", "_")
		listen := "0.0.0.0:" + freePort(pxName)

		// idempotent "create or get"
		px, err := toxCli.Proxy(pxName)
		if err != nil {
			px, err = toxCli.CreateProxy(pxName, listen, up)
			if err != nil {
				return router.Tool{}, err
			}
		}
		target = px.Listen // 0.0.0.0:32001
	}
	return router.Tool{
		Name:        tc.Name,
		Description: tc.Description,
		Handler: func(ctx context.Context, in []byte) ([]byte, error) {
			//TODO: execute script/http/grpc
			return in, nil
		},
	}, nil
}

func freePort(seed string) string {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}
