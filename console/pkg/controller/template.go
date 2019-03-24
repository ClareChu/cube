package controller

import (
	"hidevops.io/cube/console/pkg/aggregate"
	"hidevops.io/cube/console/pkg/command"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
)

type TemplateController struct {
	at.RestController
	pipelineConfigAggregate   aggregate.PipelineConfigAggregate
	buildConfigAggregate      aggregate.BuildConfigAggregate
	deploymentConfigAggregate aggregate.DeploymentConfigAggregate
	serviceConfigAggregate    aggregate.ServiceConfigAggregate
	gatewayConfigAggregate    aggregate.GatewayConfigAggregate
}

func init() {
	app.Register(newTemplateController)
}

func newTemplateController(pipelineConfigAggregate aggregate.PipelineConfigAggregate, buildConfigAggregate aggregate.BuildConfigAggregate, deploymentConfigAggregate aggregate.DeploymentConfigAggregate,
	serviceConfigAggregate aggregate.ServiceConfigAggregate, gatewayConfigAggregate aggregate.GatewayConfigAggregate) *TemplateController {
	return &TemplateController{
		pipelineConfigAggregate:   pipelineConfigAggregate,
		buildConfigAggregate:      buildConfigAggregate,
		deploymentConfigAggregate: deploymentConfigAggregate,
		serviceConfigAggregate:    serviceConfigAggregate,
		gatewayConfigAggregate:    gatewayConfigAggregate,
	}
}

func (c *TemplateController) PostPipeline(cmd *command.PipelineConfigTemplate) (model.Response, error) {
	response := new(model.BaseResponse)
	pipeline, err := c.pipelineConfigAggregate.NewPipelineConfigTemplate(cmd)
	response.SetData(pipeline)
	return response, err
}

func (c *TemplateController) PostBuildConfig(cmd *command.BuildConfig) (model.Response, error) {
	response := new(model.BaseResponse)
	pipeline, err := c.buildConfigAggregate.Template(cmd)
	response.SetData(pipeline)
	return response, err
}

func (c *TemplateController) PostDeploymentConfig(cmd *command.DeploymentConfig) (model.Response, error) {
	log.Debugf("create deployment config template: %v", cmd)
	deploy, err := c.deploymentConfigAggregate.Template(cmd)
	response := new(model.BaseResponse)
	response.SetData(deploy)
	return response, err
}

func (c *TemplateController) PostServiceConfig(cmd *command.ServiceConfig) (model.Response, error) {
	log.Debugf("create deployment config template: %v", cmd)
	deploy, err := c.serviceConfigAggregate.Template(cmd)
	response := new(model.BaseResponse)
	response.SetData(deploy)
	return response, err
}

func (c *TemplateController) PostGatewayConfig(cmd *command.GatewayConfig) (model.Response, error) {
	log.Debugf("create deployment config template: %v", cmd)
	deploy, err := c.gatewayConfigAggregate.Template(cmd)
	response := new(model.BaseResponse)
	response.SetData(deploy)
	return response, err
}
