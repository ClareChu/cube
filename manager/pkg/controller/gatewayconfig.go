package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
)

type GatewayConfigController struct {
	at.RestController
	gatewayConfigAggregate aggregate.GatewayConfigAggregate
}

func init() {
	app.Register(newGatewayConfigController)
}

func newGatewayConfigController(gatewayConfigAggregate aggregate.GatewayConfigAggregate) *GatewayConfigController {
	return &GatewayConfigController{
		gatewayConfigAggregate: gatewayConfigAggregate,
	}
}

func (c *GatewayConfigController) PostCreate(cmd *command.DeployConfigType) (model.Response, error) {
	log.Debugf("create deployment config template: %v", cmd)
	deploy, err := c.gatewayConfigAggregate.Create(cmd.Name, cmd.PipelineName, cmd.Namespace, cmd.SourceType, cmd.Version, cmd.Profile)
	response := new(model.BaseResponse)
	response.SetData(deploy)
	return response, err
}
