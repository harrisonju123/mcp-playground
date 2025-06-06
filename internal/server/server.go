package server

import (
	"context"
	"fmt"
	pb "github.com/harrisonju123/mcp-agent-poc/api/gen"
	"github.com/harrisonju123/mcp-agent-poc/config"
	"github.com/harrisonju123/mcp-agent-poc/internal/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"net"
)

type GRPCServer struct {
	pb.UnimplementedAggregatorServer
	r *router.Router
}

func New(r *router.Router) *GRPCServer { return &GRPCServer{r: r} }

func (s *GRPCServer) ListTools(ctx context.Context, _ *pb.ListToolsRequest) (*pb.ListToolsResponse, error) {
	var infos []*pb.ToolInfo
	for _, t := range s.r.List() {
		infos = append(infos, &pb.ToolInfo{
			Name:        t.Name,
			Description: t.Description,
		})
	}
	return &pb.ListToolsResponse{Tools: infos}, nil
}

func (s *GRPCServer) CallTool(ctx context.Context, req *pb.CallToolRequest) (*pb.CallToolResponse, error) {
	argsByte, err := protojson.Marshal(req.ArgsJson)
	if err != nil {
		return nil, err
	}

	out, err := s.r.Call(ctx, req.Name, argsByte)
	if err != nil {
		return nil, err
	}

	var resultStruct structpb.Struct
	if err := protojson.Unmarshal(out, &resultStruct); err != nil {
		return nil, err
	}

	return &pb.CallToolResponse{
		ResultJson: &resultStruct,
	}, nil
}

func Start(ctx context.Context, router *router.Router, config config.Config) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAggregatorServer(grpcServer, New(router))
	// reflection
	if config.EnableReflection {
		reflection.Register(grpcServer)
	}

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- grpcServer.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		return nil
	case err := <-serveErr:
		return err
	}
}
