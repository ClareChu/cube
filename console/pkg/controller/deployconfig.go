package controller

import (
	"github.com/prometheus/common/log"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/mio/console/pkg/aggregate"
	"hidevops.io/mio/console/pkg/command"
)

type DeploymentConfigController struct {
	at.RestController
	deploymentConfigAggregate aggregate.DeploymentConfigAggregate
}

func init() {
	app.Register(newDeploymentConfigController)
}

func newDeploymentConfigController(deploymentConfigAggregate aggregate.DeploymentConfigAggregate) *DeploymentConfigController {
	return &DeploymentConfigController{
		deploymentConfigAggregate: deploymentConfigAggregate,
	}
}

func (c *DeploymentConfigController) PostCreate(cmd *command.DeployConfigType) (model.Response, error) {
	log.Debugf("create deployment config template: %v", cmd)
	deploy, err := c.deploymentConfigAggregate.Create(cmd.Name, cmd.PipelineName, cmd.Namespace, cmd.SourceType, cmd.Version, "", "dev")
	response := new(model.BaseResponse)
	response.SetData(deploy)
	return response, err
}
