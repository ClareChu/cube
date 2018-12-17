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

func TestDeploymentConfig(t *testing.T) {
	deploy := new(mocks.DeploymentConfigAggregate)
	appInfo := newDeploymentConfigController(deploy, nil)
	deploy.On("Create", "", "", "", "", "", "", "dev").Return(nil, nil)

	app := web.NewTestApp(t, appInfo).
		SetProperty("kube.serviceHost", "test").
		Run(t)
	log.SetLevel(log.DebugLevel)
	log.Println(io.GetWorkDir())
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/deploymentConfig/create").WithJSON(&command.DeployConfigType{}).Expect().Status(http.StatusOK)
	})
}
