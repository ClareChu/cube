package controller

import (
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
)

type AppController struct {
	at.RestController
	appService service.AppService
}

func init() {
	app.Register(newAppController)
}

func newAppController(appService service.AppService) *AppController {
	return &AppController{
		appService: appService,
	}
}

func (a *AppController) Post(cmd *command.PipelineStart) (response model.Response, err error) {
	response = new(model.BaseResponse)
	app, err := a.appService.Create(cmd)
	response.SetData(app)
	return response, err
}
