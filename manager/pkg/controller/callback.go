package controller

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
)

type CallbackController struct {
	at.RestController
}

func init() {
	app.Register(newCallbackController)
}

func newCallbackController() *CallbackController {
	return &CallbackController{
	}
}

type Rep struct {
	model.RequestBody
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (c *CallbackController) Get(r *Rep) {
	log.Infof("get info success")
}
