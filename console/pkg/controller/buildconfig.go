package controller

import (
	"github.com/jinzhu/copier"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/mio/console/pkg/aggregate"
	"hidevops.io/mio/console/pkg/builder"
	"hidevops.io/mio/console/pkg/command"
	"hidevops.io/mio/console/pkg/service"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
)

type BuildConfigController struct {
	at.RestController
	buildConfigService   service.BuildConfigService
	buildConfigAggregate aggregate.BuildConfigAggregate
	buildAggregate       aggregate.BuildAggregate
	buildNode            builder.BuildNode
}

func init() {
	app.Register(newSourceConfigController)
}

func newSourceConfigController(buildConfigService service.BuildConfigService, buildConfigAggregate aggregate.BuildConfigAggregate, buildNode builder.BuildNode, buildAggregate aggregate.BuildAggregate) *BuildConfigController {
	return &BuildConfigController{
		buildConfigService:   buildConfigService,
		buildConfigAggregate: buildConfigAggregate,
		buildNode:            buildNode,
		buildAggregate:       buildAggregate,
	}
}

func (c *BuildConfigController) PostSourceCodePull(command *command.SourceCodePullCommand) (model.Response, error) {
	response := new(model.BaseResponse)
	err := c.buildConfigService.SourceCodePull(command.Host, command.Port, command)
	return response, err
}

func (c *BuildConfigController) PostCompile(command *command.CompileCommand) (model.Response, error) {
	response := new(model.BaseResponse)
	err := c.buildConfigService.Compile("localhost", "7578", command)
	return response, err
}

func (c *BuildConfigController) PostImageBuild(command *command.ImageBuildCommand) (model.Response, error) {
	response := new(model.BaseResponse)
	err := c.buildConfigService.ImageBuild("localhost", "7578", command)
	return response, err
}

func (c *BuildConfigController) PostImagePush(command *command.ImagePushCommand) (model.Response, error) {
	response := new(model.BaseResponse)
	err := c.buildConfigService.ImagePush("localhost", "7578", command)
	return response, err
}

func (c *BuildConfigController) PostCreateService(command *command.ServiceNode) (model.Response, error) {
	response := new(model.BaseResponse)
	err := c.buildNode.CreateServiceNode(command)
	return response, err
}

func (c *BuildConfigController) PostDeployNode(command *command.DeployNode) (model.Response, error) {
	response := new(model.BaseResponse)
	deploy, err := c.buildNode.Start(command)
	response.SetData(deploy)
	return response, err
}

func (c *BuildConfigController) PostWatch(command *command.PipelineStart) (model.Response, error) {
	response := new(model.BaseResponse)
	err := c.buildAggregate.WatchPod(command.Name, command.Namespace)
	return response, err
}

func (c *BuildConfigController) Post(template *command.BuildConfigTemplate) (model.Response, error) {
	pipeline := new(v1alpha1.Pipeline)
	copier.Copy(&pipeline, template)
	build, err := c.buildConfigAggregate.Create(template.Name, template.PipelineName, template.Namespace, template.SourceType, template.Version, template.Branch)
	base := new(model.BaseResponse)
	base.SetData(build)
	return base, err
}

func (c *BuildConfigController) DeleteByNameNamespace(name, namespace string) (model.Response, error) {
	err := c.buildConfigAggregate.Delete(name, namespace)
	response := new(model.BaseResponse)
	return response, err
}
