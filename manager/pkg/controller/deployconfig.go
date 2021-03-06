package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentConfigController struct {
	at.RestController
	deploymentConfigAggregate aggregate.DeploymentConfigAggregate
	deploymentConfigBuilder   builder.DeploymentConfigBuilder
}

func init() {
	app.Register(newDeploymentConfigController)
}

func newDeploymentConfigController(deploymentConfigAggregate aggregate.DeploymentConfigAggregate, deploymentConfigBuilder builder.DeploymentConfigBuilder) *DeploymentConfigController {
	return &DeploymentConfigController{
		deploymentConfigAggregate: deploymentConfigAggregate,
		deploymentConfigBuilder:   deploymentConfigBuilder,
	}
}

func (c *DeploymentConfigController) PostCreate(cmd *command.DeployConfigType) (model.Response, error) {
	log.Debugf("create deployment config template: %v", cmd)
	param := &command.PipelineReqParams{}
	err := c.deploymentConfigAggregate.Create(param)
	response := new(model.BaseResponse)
	return response, err
}

type App struct {
	model.RequestBody
}

func (c *DeploymentConfigController) PostApp(app *App) (model.Response, error) {
	deploy := &v1alpha1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
		Spec: v1alpha1.DeploymentConfigSpec{
			Version: "v1",
			Labels: map[string]string{
				"app":     "hello-world",
				"version": "v1",
			},
			Container: corev1.Container{
				Ports: []corev1.ContainerPort{
					corev1.ContainerPort{
						Name: "http",
					},
				},
			},
		},
	}
	err := c.deploymentConfigBuilder.CreateApp(deploy)
	return nil, err
}
