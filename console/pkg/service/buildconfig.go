package service

import (
	"context"
	"github.com/jinzhu/copier"
	"hidevops.io/cube/console/pkg/command"
	"hidevops.io/cube/console/pkg/protobuf"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/grpc"
	"hidevops.io/hiboot/pkg/utils/idgen"
	"time"
)

const ConnectorTime = 60 * time.Minute

type BuildConfigService interface {
	SourceCodePull(host, port string, command *command.SourceCodePullCommand) error
	Compile(host, port string, cmd *command.CompileCommand) error
	ImageBuild(host, port string, cmd *command.ImageBuildCommand) error
	ImagePush(host, port string, cmd *command.ImagePushCommand) error
}

type BuildConfigServiceImpl struct {
	BuildConfigService
	clientConnector grpc.ClientConnector
}

func init() {
	app.Register(newBuildConfigCommand)
}

func newBuildConfigCommand(clientConnector grpc.ClientConnector) BuildConfigService {
	return &BuildConfigServiceImpl{
		clientConnector: clientConnector,
	}
}

func (s *BuildConfigServiceImpl) SourceCodePull(host, port string, cmd *command.SourceCodePullCommand) error {
	log.Debug("build config create")
	prop := &grpc.ClientProperties{
		Host: host,
		Port: port,
	}
	id, _ := idgen.NextString()
	gRpcCli, err := s.clientConnector.Connect(id, protobuf.NewBuildConfigServiceClient, prop)
	if err != nil {
		log.Errorf("client connect err :%v", err)
		return err
	}
	request := &protobuf.SourceCodePullRequest{}
	copier.Copy(request, cmd)
	ctx, cancel := context.WithTimeout(context.Background(), ConnectorTime)
	defer cancel()
	buildConfigServiceClient := gRpcCli.(protobuf.BuildConfigServiceClient)
	response, err := buildConfigServiceClient.SourceCodePull(ctx, request)
	log.Debug(response)
	if err != nil {
		return err
	}
	return err
}

func (s *BuildConfigServiceImpl) Compile(host, port string, cmd *command.CompileCommand) error {
	log.Debug("build config create")
	prop := &grpc.ClientProperties{
		Host: host,
		Port: port,
	}
	id, _ := idgen.NextString()
	gRpcCli, err := s.clientConnector.Connect(id, protobuf.NewBuildConfigServiceClient, prop)
	if err != nil {
		log.Errorf("client connect err :%v", err)
		return err
	}
	request := &protobuf.CompileRequest{}
	copier.Copy(request, cmd)
	ctx, cancel := context.WithTimeout(context.Background(), ConnectorTime)
	defer cancel()
	buildConfigServiceClient := gRpcCli.(protobuf.BuildConfigServiceClient)
	response, err := buildConfigServiceClient.Compile(ctx, request)
	log.Debug(response)
	if err != nil {
		return err
	}
	return err
}

func (s *BuildConfigServiceImpl) ImageBuild(host, port string, cmd *command.ImageBuildCommand) error {
	log.Debug("build config create")
	prop := &grpc.ClientProperties{
		Host: host,
		Port: port,
	}
	id, _ := idgen.NextString()
	gRpcCli, err := s.clientConnector.Connect(id, protobuf.NewBuildConfigServiceClient, prop)
	if err != nil {
		log.Errorf("client connect err :%v", err)
		return err
	}
	request := &protobuf.ImageBuildRequest{}
	copier.Copy(request, cmd)
	ctx, cancel := context.WithTimeout(context.Background(), ConnectorTime)
	defer cancel()
	buildConfigServiceClient := gRpcCli.(protobuf.BuildConfigServiceClient)
	response, err := buildConfigServiceClient.ImageBuild(ctx, request)
	log.Debug(response)
	if err != nil {
		return err
	}
	return err
}

func (s *BuildConfigServiceImpl) ImagePush(host, port string, cmd *command.ImagePushCommand) error {
	log.Debug("build config create")
	prop := &grpc.ClientProperties{
		Host: host,
		Port: port,
	}
	id, _ := idgen.NextString()
	gRpcCli, err := s.clientConnector.Connect(id, protobuf.NewBuildConfigServiceClient, prop)
	if err != nil {
		log.Errorf("client connect err :%v", err)
		return err
	}
	request := &protobuf.ImagePushRequest{}
	copier.Copy(request, cmd)
	ctx, cancel := context.WithTimeout(context.Background(), ConnectorTime)
	defer cancel()
	buildConfigServiceClient := gRpcCli.(protobuf.BuildConfigServiceClient)
	response, err := buildConfigServiceClient.ImagePush(ctx, request)
	log.Debug(response)
	if err != nil {
		return err
	}
	return err
}
