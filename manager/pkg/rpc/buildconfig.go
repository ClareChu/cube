package rpc

import (
	"context"
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/protobuf"
	"hidevops.io/hiboot/pkg/starter/grpc"
)

type buildConfigServiceServerImpl struct {
	buildConfigAggregate aggregate.BuildConfigAggregate
}

func init() {
	grpc.Server(protobuf.RegisterBuildConfigServiceServer, newBuildConfigServiceServerImpl)
}

func newBuildConfigServiceServerImpl(buildConfigAggregate aggregate.BuildConfigAggregate) protobuf.BuildConfigServiceServer {
	return &buildConfigServiceServerImpl{
		buildConfigAggregate: buildConfigAggregate,
	}
}

func (s *buildConfigServiceServerImpl) SourceCodePull(ctx context.Context, request *protobuf.SourceCodePullRequest) (*protobuf.SourceCodePullResponse, error) {
	// response to client
	response := &protobuf.SourceCodePullResponse{}
	return response, nil
}

func (s *buildConfigServiceServerImpl) Compile(ctx context.Context, request *protobuf.CompileRequest) (*protobuf.CompileResponse, error) {
	// response to client
	response := &protobuf.CompileResponse{}
	return response, nil
}

func (s *buildConfigServiceServerImpl) ImageBuild(ctx context.Context, request *protobuf.ImageBuildRequest) (*protobuf.ImageBuildResponse, error) {
	// response to client
	response := &protobuf.ImageBuildResponse{}
	return response, nil
}

func (s *buildConfigServiceServerImpl) ImagePush(ctx context.Context, request *protobuf.ImagePushRequest) (*protobuf.ImagePushResponse, error) {
	// response to client
	response := &protobuf.ImagePushResponse{}
	return response, nil
}

func (s *buildConfigServiceServerImpl) Test(ctx context.Context, request *protobuf.TestsRequest) (*protobuf.TestsResponse, error) {
	// response to client
	response := &protobuf.TestsResponse{}
	return response, nil
}


func (s *buildConfigServiceServerImpl) Command(ctx context.Context, request *protobuf.CommandRequest) (*protobuf.CommandResponse, error) {
	// response to client
	response := &protobuf.CommandResponse{}
	return response, nil
}