package server

import (
	"context"
	pb "github.com/harrisonju123/mcp-agent-poc/api/gen"
	"github.com/harrisonju123/mcp-agent-poc/router"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
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

	out, err := s.r.Call(req.Name, argsByte)
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
