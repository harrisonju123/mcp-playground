package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/harrisonju123/mcp-agent-poc/api/gen"
	"github.com/harrisonju123/mcp-agent-poc/config"
	"github.com/harrisonju123/mcp-agent-poc/internal/router"
	"github.com/harrisonju123/mcp-agent-poc/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	r := router.NewRouter(nil)
	r.Replace([]router.Tool{{
		Name:        "echo",
		Description: "Return args unchanged",
		Handler: func(in []byte) ([]byte, error) {
			// validate if json
			var v any
			if err := json.Unmarshal(in, &v); err != nil {
				return nil, err
			}
			return in, nil
		},
	}})
	cfg := config.Load()
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	s := grpc.NewServer()
	pb.RegisterAggregatorServer(s, server.New(r))
	if cfg.EnableReflection {
		reflection.Register(s)
	}

	log.Printf("port=%d reflection=%v registry=%s",
		cfg.Port, cfg.EnableReflection, cfg.RegistryURL)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}

}
