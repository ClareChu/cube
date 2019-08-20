package controller

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
)

type GatewayController struct {
	at.RestController
}

func init() {
	app.Register(newGatewayController)
}

func newGatewayController() *GatewayController {
	return &GatewayController{
	}
}



func (c *GatewayController) Post() (response model.Response, err error) {
	log.Debugf("create deployment config template: %v", "a")

	return response, err
}
