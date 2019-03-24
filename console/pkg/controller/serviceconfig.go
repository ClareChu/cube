package controller

import (
	"github.com/prometheus/common/log"
	"hidevops.io/cube/console/pkg/aggregate"
	"hidevops.io/cube/console/pkg/command"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
)

type ServiceConfigController struct {
	at.RestController
	serviceConfigAggregate aggregate.ServiceConfigAggregate
}

func init() {
	app.Register(newServiceConfigController)
}

func newServiceConfigController(serviceConfigAggregate aggregate.ServiceConfigAggregate) *ServiceConfigController {
	return &ServiceConfigController{
		serviceConfigAggregate: serviceConfigAggregate,
	}
}

func (c *ServiceConfigController) PostCreate(cmd *command.DeployConfigType) (model.Response, error) {
	log.Debugf("create deployment config template: %v", cmd)
	deploy, err := c.serviceConfigAggregate.Create(cmd.Name, cmd.PipelineName, cmd.Namespace, cmd.SourceType, cmd.Version, cmd.Profile)
	response := new(model.BaseResponse)
	response.SetData(deploy)
	return response, err
}
