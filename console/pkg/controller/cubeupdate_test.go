package controller

import (
	"hidevops.io/cube/console/pkg/aggregate/mocks"
	"hidevops.io/cube/console/pkg/command"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/io"
	"net/http"
	"testing"
)

func TestCubeUpdate(t *testing.T) {
	cubeUpdateAggregate := new(mocks.CubeUpdateAggregate)
	gateway := newCubeUpdateController(cubeUpdateAggregate)
	cubeUpdateAggregate.On("Add", "", &command.CubeUpdate{}).Return(nil)
	cubeUpdateAggregate.On("Delete", "ab").Return(nil)
	update := new(command.CubeUpdate)
	update.Version = "1"
	cubeUpdateAggregate.On("Get", "ab").Return(update, nil)
	app := web.NewTestApp(t, gateway).
		SetProperty("kube.serviceHost", "test").
		Run(t)
	log.SetLevel(log.DebugLevel)
	log.Println(io.GetWorkDir())
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/cubeUpdate").WithJSON(&command.CubeUpdate{}).Expect().Status(http.StatusOK)
	})

	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Delete("/cubeUpdate/type/a/arch/b").Expect().Status(http.StatusOK)
	})

	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Get("/cubeUpdate/type/a/arch/b/version/1").Expect().Status(http.StatusOK)
	})
}
