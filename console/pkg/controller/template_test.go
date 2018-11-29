package controller

import (
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/io"
	"hidevops.io/mio/console/pkg/aggregate/mocks"
	"hidevops.io/mio/console/pkg/command"
	"net/http"
	"testing"
)

func TestTemplate(t *testing.T) {
	serviceConfigAggregate := new(mocks.ServiceConfigAggregate)
	pipelineConfigAggregate := new(mocks.PipelineConfigAggregate)
	buildConfigAggregate := new(mocks.BuildConfigAggregate)
	deploymentConfigAggregate := new(mocks.DeploymentConfigAggregate)
	gatewayAggregate := new(mocks.GatewayConfigAggregate)
	service := newTemplateController(pipelineConfigAggregate, buildConfigAggregate, deploymentConfigAggregate, serviceConfigAggregate, gatewayAggregate)

	app := web.NewTestApp(t, service).
		SetProperty("kube.serviceHost", "test").
		Run(t)
	log.SetLevel(log.DebugLevel)
	log.Println(io.GetWorkDir())
	pipelineConfigAggregate.On("NewPipelineConfigTemplate", &command.PipelineConfigTemplate{}).Return(nil, nil)
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/template/pipeline").WithJSON(&command.PipelineConfigTemplate{}).Expect().Status(http.StatusOK)
	})

	serviceConfigAggregate.On("Template", &command.ServiceConfig{}).Return(nil, nil)
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/template/serviceConfig").WithJSON(&command.ServiceConfig{}).Expect().Status(http.StatusOK)
	})

	gatewayAggregate.On("Template", &command.GatewayConfig{}).Return(nil, nil)
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/template/gatewayConfig").WithJSON(&command.GatewayConfig{}).Expect().Status(http.StatusOK)
	})

	buildConfigAggregate.On("Template", &command.BuildConfig{}).Return(nil, nil)
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/template/buildConfig").WithJSON(&command.BuildConfig{}).Expect().Status(http.StatusOK)
	})

	deploymentConfigAggregate.On("Template", &command.DeploymentConfig{}).Return(nil, nil)
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/template/deploymentConfig").WithJSON(&command.DeploymentConfig{}).Expect().Status(http.StatusOK)
	})
}
