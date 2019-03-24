package controller

import (
	"hidevops.io/cube/console/pkg/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/starter/websocket"
)

type webSocketController struct {
	at.RestController
	register websocket.Register
}

func newWebSocketController(register websocket.Register) *webSocketController {
	c := &webSocketController{
		register: register,
	}
	return c
}

func init() {
	app.Register(newWebSocketController)
}

func (c *webSocketController) GetBuildLogs(handler *service.LogsHandler, conn *websocket.Connection) {
	c.register(handler, conn)
}

func (c *webSocketController) GetAppLogs(handler *service.AppLogsHandler, conn *websocket.Connection) {
	c.register(handler, conn)
}
