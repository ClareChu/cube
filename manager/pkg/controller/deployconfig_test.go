package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate/mocks"
	builder "hidevops.io/cube/manager/pkg/builder/mocks"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"testing"
)

func TestDeploymentConfig(t *testing.T) {
	deploy := new(mocks.DeploymentConfigAggregate)
	deployment := new(builder.DeploymentConfigBuilder)
	appInfo := newDeploymentConfigController(deploy, deployment)
	deploy.On("Create", "", "", "", "", "", "", "dev").Return(nil, nil)

	app := web.NewTestApp(t, appInfo).
		SetProperty("kube.serviceHost", "test").
		Run(t)
	log.SetLevel(log.DebugLevel)
	log.Println(io.GetWorkDir())
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/deploymentConfig/create").WithJSON(&command.DeployConfigType{}).Expect().Status(http.StatusOK)
	})
	ds := &v1alpha1.Deployment{
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
			Port: []corev1.ContainerPort{
				corev1.ContainerPort{
					ContainerPort: 8080,
					Protocol:      "TCP",
				},
				corev1.ContainerPort{
					ContainerPort: 7575,
					Protocol:      "TCP",
				},
			},
		},
	}
	deployment.On("CreateApp", ds).Return(nil)
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/deploymentConfig/app").WithJSON(&App{}).Expect().Status(http.StatusOK)
	})
}
