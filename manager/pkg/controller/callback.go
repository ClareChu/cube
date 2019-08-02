package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
)

type CallbackController struct {
	at.RestController
	callbackAggregate aggregate.CallbackAggregate
}

func init() {
	app.Register(newCallbackController)
}

func newCallbackController(callbackAggregate aggregate.CallbackAggregate) *CallbackController {
	return &CallbackController{
		callbackAggregate: callbackAggregate,
	}
}

type Rep struct {
	model.RequestBody
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (c *CallbackController) Get() {
	log.Infof("get info success")
	c.callbackAggregate.GetPod("ide-my-app22-37-v1-5b5d6d55cc-j9mzw", "test")
}
