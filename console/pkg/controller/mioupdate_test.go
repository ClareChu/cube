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

func TestMioUpdate(t *testing.T) {
	mioUpdateAggregate := new(mocks.MioUpdateAggregate)
	gateway := newMioUpdateController(mioUpdateAggregate)
	mioUpdateAggregate.On("Add", "", &command.MioUpdate{}).Return(nil)
	mioUpdateAggregate.On("Delete", "ab").Return(nil)
	update := new(command.MioUpdate)
	update.Version = "1"
	mioUpdateAggregate.On("Get", "ab").Return(update, nil)
	app := web.NewTestApp(t, gateway).
		SetProperty("kube.serviceHost", "test").
		Run(t)
	log.SetLevel(log.DebugLevel)
	log.Println(io.GetWorkDir())
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/mioUpdate").WithJSON(&command.MioUpdate{}).Expect().Status(http.StatusOK)
	})

	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Delete("/mioUpdate/type/a/arch/b").Expect().Status(http.StatusOK)
	})

	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Get("/mioUpdate/type/a/arch/b/version/1").Expect().Status(http.StatusOK)
	})
}
