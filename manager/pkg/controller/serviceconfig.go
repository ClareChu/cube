package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
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
	param := &command.PipelineReqParams{}
	deploy, err := c.serviceConfigAggregate.Create(param)
	response := new(model.BaseResponse)
	response.SetData(deploy)
	return response, err
}
