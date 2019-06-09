package controller

import (
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
)

type AppController struct {
	at.RestController
}

func init() {
	app.Register(newAppController)
}

func newAppController() *AppController {
	return &AppController{
	}
}

func (a *AppController) Post(cmd *command.PipelineStart) (model.Response, error) {

	return response, err
}