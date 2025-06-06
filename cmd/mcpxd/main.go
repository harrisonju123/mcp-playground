package main

import (
	"context"
	"encoding/json"
	"github.com/harrisonju123/mcp-agent-poc/config"
	"github.com/harrisonju123/mcp-agent-poc/internal/router"
	"github.com/harrisonju123/mcp-agent-poc/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	r := router.NewRouter(nil)
	r.Replace([]router.Tool{{
		Name:        "echo",
		Description: "Return args unchanged",
		Handler: func(ctx context.Context, in []byte) ([]byte, error) {
			var v any
			if err := json.Unmarshal(in, &v); err != nil {
				return nil, err
			}
			return in, nil
		},
	}})

	//  watch tools
	go func() {
		if err := config.NewWatcher("../../config/tools.yaml", r).Run(ctx); err != nil {
			log.Fatalf("config watcher failed: %v", err)
		}
	}()

	cfg := config.Load()
	log.Printf("starting gRPC on port=%d (reflection=%v)", cfg.Port, cfg.EnableReflection)
	// start gRPC server
	if err := server.Start(ctx, r, cfg); err != nil {
		log.Fatalf("grpc server: %v", err)
	}
}
